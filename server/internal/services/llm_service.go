// Package services 承载业务逻辑，对 handler 层提供无 HTTP 上下文的纯函数接口
package services

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/zhitu/server/internal/config"
)

// LLMService 封装 OpenAI Chat Completions 与 Anthropic Messages 接口的调用
// 直接用 net/http，不引入第三方 SDK，便于切换 provider
type LLMService struct {
	cfg    *config.LLMConfig
	client *http.Client
}

// ChatMessage OpenAI Chat Completions 消息格式
type ChatMessage struct {
	Role    string `json:"role"` // system / user / assistant
	Content string `json:"content"`
}

// chatRequest Chat Completions 请求体
type chatRequest struct {
	Model          string        `json:"model"`
	Messages       []ChatMessage `json:"messages"`
	Temperature    float64       `json:"temperature,omitempty"`
	MaxTokens      int           `json:"max_tokens,omitempty"`
	Stream         bool          `json:"stream,omitempty"`
	ResponseFormat *respFormat   `json:"response_format,omitempty"`
}

// respFormat 强制 JSON 输出
type respFormat struct {
	Type string `json:"type"` // "json_object"
}

type chatResponseMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponseChoice struct {
	Message      chatResponseMessage `json:"message"`
	FinishReason string              `json:"finish_reason"`
}

// chatResponse 非流式响应体
type chatResponse struct {
	Choices []chatResponseChoice `json:"choices"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// streamChunk 流式响应的单个 SSE chunk
type streamChunk struct {
	Choices []struct {
		Delta struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

// anthropicRequest Anthropic Messages 请求体。
// Anthropic 将 system 提示词放在顶层，其余消息只允许 user / assistant。
type anthropicRequest struct {
	Model       string        `json:"model"`
	System      string        `json:"system,omitempty"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens"`
	Stream      bool          `json:"stream,omitempty"`
}

type anthropicResponse struct {
	Type    string `json:"type"`
	Role    string `json:"role"`
	Model   string `json:"model"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	StopReason string `json:"stop_reason"`
	Usage      struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

type anthropicStreamEvent struct {
	Type  string `json:"type"`
	Delta struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"delta"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

// === MiMo ASR (mimo-v2.5-asr) 请求/响应类型 ===
// mimo ASR 不走 /v1/audio/transcriptions，而是通过 /v1/chat/completions
// 用 input_audio 类型的 content 传入 base64 编码的音频

// asrInputAudio ASR 输入音频
type asrInputAudio struct {
	Data string `json:"data"` // "data:{MIME};base64,{BASE64_AUDIO}"
}

// asrContentItem ASR 消息 content 数组项
type asrContentItem struct {
	Type       string        `json:"type"` // "input_audio"
	InputAudio asrInputAudio `json:"input_audio"`
}

// asrMessage ASR 请求消息（content 是数组而非字符串）
type asrMessage struct {
	Role    string           `json:"role"` // "user"
	Content []asrContentItem `json:"content"`
}

// asrRequest mimo-v2.5-asr 请求体
type asrRequest struct {
	Model      string       `json:"model"`
	Messages   []asrMessage `json:"messages"`
	ASROptions struct {
		Language string `json:"language"` // auto / zh / en
	} `json:"asr_options"`
}

// === MiMo TTS (mimo-v2.5-tts) 请求/响应类型 ===
// mimo TTS 不走 /v1/audio/speech，而是通过 /v1/chat/completions
// 目标文本放在 role=assistant 的 content 中，audio 参数指定格式和音色

// ttsMessage TTS 请求消息
type ttsMessage struct {
	Role    string `json:"role"`    // "assistant"
	Content string `json:"content"` // 要合成的文本
}

// ttsAudioConfig TTS 音频配置
type ttsAudioConfig struct {
	Format string `json:"format"` // "wav" / "pcm16"
	Voice  string `json:"voice"`  // "Chloe" / "冰糖" 等预置音色
}

// ttsRequestMimo mimo-v2.5-tts 请求体
type ttsRequestMimo struct {
	Model    string         `json:"model"`
	Messages []ttsMessage   `json:"messages"`
	Audio    ttsAudioConfig `json:"audio"`
}

// ttsResponseMimo TTS 响应（响应中 message.audio.data 是 base64 编码的音频）
type ttsResponseMimo struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
			Audio   struct {
				ID   string `json:"id"`
				Data string `json:"data"` // base64 编码的音频字节
			} `json:"audio"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

// LLM 错误
var (
	ErrLLMNotConfigured = errors.New("llm not configured: api_key is empty")
	ErrLLMEmptyResponse = errors.New("llm returned empty content")
	ErrLLMStreamFailed  = errors.New("llm stream interrupted")
)

// NewLLMService 构造 LLMService
func NewLLMService(cfg *config.LLMConfig) *LLMService {
	timeout := time.Duration(cfg.TimeoutSec) * time.Second
	if timeout == 0 {
		timeout = 60 * time.Second
	}
	return &LLMService{
		cfg: cfg,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// IsConfigured 是否已配置可用（api_key 非空且非占位符）
func (s *LLMService) IsConfigured() bool {
	return s.cfg.APIKey != "" && !strings.Contains(s.cfg.APIKey, "please-replace")
}

// Chat 非流式对话，返回完整文本
//
// 模型：cfg.ChatModel —— 面试功能对接 mimo 时应配置为 "mimo-v2.5-pro"
// 用途：面试文字分析（提问生成、评分、复盘报告等所有非流式 LLM 调用）
func (s *LLMService) Chat(ctx context.Context, messages []ChatMessage) (string, error) {
	if !s.IsConfigured() {
		return "", ErrLLMNotConfigured
	}

	reqBody := chatRequest{
		Model:       s.cfg.ChatModel,
		Messages:    messages,
		Temperature: s.cfg.Temperature,
		MaxTokens:   s.cfg.MaxTokens,
	}

	respBody, err := s.doChat(ctx, reqBody)
	if err != nil {
		return "", err
	}
	if len(respBody.Choices) == 0 || respBody.Choices[0].Message.Content == "" {
		return "", ErrLLMEmptyResponse
	}
	return respBody.Choices[0].Message.Content, nil
}

// ChatStream 流式对话，通过 onDelta 回调逐 token 推送
// 调用方负责在 HTTP handler 中设置 SSE 头并写入响应流
//
// 模型：cfg.ChatModel —— 面试功能对接 mimo 时应配置为 "mimo-v2.5-pro"
// 用途：面试文字分析（AI 面试官流式提问）
func (s *LLMService) ChatStream(ctx context.Context, messages []ChatMessage, onDelta func(delta string)) error {
	if !s.IsConfigured() {
		return ErrLLMNotConfigured
	}

	reqBody := chatRequest{
		Model:       s.cfg.ChatModel,
		Messages:    messages,
		Temperature: s.cfg.Temperature,
		MaxTokens:   s.cfg.MaxTokens,
		Stream:      true,
	}

	if s.isAnthropic() {
		return s.chatStreamAnthropic(ctx, reqBody, onDelta)
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.BaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.cfg.APIKey)
	req.Header.Set("Accept", "text/event-stream")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return s.wrapAPIError(resp)
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024) // 扩大 buffer 防止长行截断
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}
		var chunk streamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue // 跳过无法解析的 chunk
		}
		if len(chunk.Choices) > 0 {
			delta := chunk.Choices[0].Delta.Content
			if delta != "" && onDelta != nil {
				onDelta(delta)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("%w: %v", ErrLLMStreamFailed, err)
	}
	return nil
}

// ChatJSON 对话并要求返回 JSON，自动解析到 out
// out 必须是指针。prompt 中应明确要求 JSON 格式与字段。
//
// 模型：cfg.ChatModel —— 面试功能对接 mimo 时应配置为 "mimo-v2.5-pro"
// 用途：面试文字分析（评分、复盘报告等结构化输出）
func (s *LLMService) ChatJSON(ctx context.Context, messages []ChatMessage, out interface{}) error {
	if !s.IsConfigured() {
		return ErrLLMNotConfigured
	}

	reqBody := chatRequest{
		Model:          s.cfg.ChatModel,
		Messages:       messages,
		Temperature:    s.cfg.Temperature,
		MaxTokens:      s.cfg.MaxTokens,
		ResponseFormat: &respFormat{Type: "json_object"},
	}

	respBody, err := s.doChat(ctx, reqBody)
	if err != nil {
		return err
	}
	if len(respBody.Choices) == 0 || respBody.Choices[0].Message.Content == "" {
		return ErrLLMEmptyResponse
	}

	content := strings.TrimSpace(respBody.Choices[0].Message.Content)
	if err := json.Unmarshal([]byte(content), out); err != nil {
		// 兜底：尝试从 markdown json 代码块中提取
		extracted := extractJSONBlock(content)
		if extracted != "" {
			if err2 := json.Unmarshal([]byte(extracted), out); err2 == nil {
				return nil
			}
		}
		return fmt.Errorf("unmarshal llm json: %w (raw: %s)", err, truncate(content, 500))
	}
	return nil
}

// Transcribe 调用 MiMo ASR 将音频转写为文本
//
// mimo ASR 不走 OpenAI 的 /v1/audio/transcriptions，而是通过 /v1/chat/completions
// 将音频 base64 编码后放在 input_audio 类型的 content 中传入。
//
// 模型：cfg.WhisperModel → "mimo-v2.5-asr"（语音识别）
// 用途：面试语音识别（用户语音回答转文字）
func (s *LLMService) Transcribe(ctx context.Context, audio io.Reader, filename string) (string, error) {
	if !s.IsConfigured() {
		return "", ErrLLMNotConfigured
	}

	// 1. 读取音频并 base64 编码
	audioBytes, err := io.ReadAll(audio)
	if err != nil {
		return "", fmt.Errorf("read audio: %w", err)
	}
	audioBase64 := base64.StdEncoding.EncodeToString(audioBytes)
	mime := audioMIME(filename)
	dataURL := fmt.Sprintf("data:%s;base64,%s", mime, audioBase64)

	// 2. 构造 ASR 请求
	reqBody := asrRequest{
		Model: s.cfg.WhisperModel,
		Messages: []asrMessage{{
			Role: "user",
			Content: []asrContentItem{{
				Type:       "input_audio",
				InputAudio: asrInputAudio{Data: dataURL},
			}},
		}},
		ASROptions: struct {
			Language string `json:"language"`
		}{Language: "auto"},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	// 3. POST /v1/chat/completions
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.openAIBaseURL()+"/chat/completions",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.cfg.APIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", s.wrapAPIError(resp)
	}

	// 4. 解析响应：choices[0].message.content 是识别出的文字
	var result chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode asr response: %w", err)
	}
	if len(result.Choices) == 0 || result.Choices[0].Message.Content == "" {
		return "", ErrLLMEmptyResponse
	}
	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}

// Synthesize 调用 MiMo TTS 将文本合成为 WAV 音频字节
//
// mimo TTS 不走 OpenAI 的 /v1/audio/speech，而是通过 /v1/chat/completions
// 目标文本放在 role=assistant 的 content 中，audio 参数指定格式和音色。
// 响应中 message.audio.data 是 base64 编码的音频字节。
//
// 模型：cfg.TTSModel → "mimo-v2.5-tts"（语音生成）
// 用途：面试语音生成（AI 面试官朗读，前端用户可选开关）
func (s *LLMService) Synthesize(ctx context.Context, text string) ([]byte, error) {
	if !s.IsConfigured() {
		return nil, ErrLLMNotConfigured
	}
	if text == "" {
		return nil, errors.New("tts input text is empty")
	}

	// 音色兜底：未配置或为 OpenAI 默认值时用 mimo 预置音色
	voice := s.cfg.TTSVoice
	if voice == "" || voice == "alloy" {
		voice = "mimo_default"
	}

	// 1. 构造 TTS 请求
	reqBody := ttsRequestMimo{
		Model: s.cfg.TTSModel,
		Messages: []ttsMessage{{
			Role:    "assistant",
			Content: text,
		}},
		Audio: ttsAudioConfig{
			Format: "wav",
			Voice:  voice,
		},
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// 2. POST /v1/chat/completions
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.openAIBaseURL()+"/chat/completions",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.cfg.APIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, s.wrapAPIError(resp)
	}

	// 3. 解析响应：message.audio.data 是 base64 编码的 WAV 音频
	var result ttsResponseMimo
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode tts response: %w", err)
	}
	if len(result.Choices) == 0 || result.Choices[0].Message.Audio.Data == "" {
		return nil, ErrLLMEmptyResponse
	}

	// 4. base64 解码为音频字节
	audio, err := base64.StdEncoding.DecodeString(result.Choices[0].Message.Audio.Data)
	if err != nil {
		return nil, fmt.Errorf("decode base64 audio: %w", err)
	}
	return audio, nil
}

// audioMIME 根据音频文件扩展名返回 MIME 类型
func audioMIME(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".wav":
		return "audio/wav"
	case ".mp3":
		return "audio/mpeg"
	case ".m4a":
		return "audio/mp4"
	case ".ogg":
		return "audio/ogg"
	case ".flac":
		return "audio/flac"
	default:
		return "audio/wav"
	}
}

// doChat 执行非流式 Chat 请求
func (s *LLMService) doChat(ctx context.Context, reqBody chatRequest) (*chatResponse, error) {
	if s.isAnthropic() {
		return s.doChatAnthropic(ctx, reqBody)
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.BaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.cfg.APIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, s.wrapAPIError(resp)
	}

	var result chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}

func (s *LLMService) isAnthropic() bool {
	provider := strings.ToLower(strings.TrimSpace(s.cfg.Provider))
	return provider == "anthropic" || provider == "mimo-anthropic"
}

func (s *LLMService) anthropicEndpoint() string {
	base := strings.TrimRight(strings.TrimSpace(s.cfg.BaseURL), "/")
	if strings.HasSuffix(base, "/messages") {
		return base
	}
	return base + "/messages"
}

// openAIBaseURL 返回语音接口使用的 OpenAI 兼容 API 根地址。
// 当聊天使用 MiMo Anthropic 入口时，语音能力仍位于同域名的 /v1 下。
func (s *LLMService) openAIBaseURL() string {
	base := strings.TrimRight(strings.TrimSpace(s.cfg.BaseURL), "/")
	if !s.isAnthropic() {
		return base
	}
	if marker := strings.Index(base, "/anthropic/"); marker >= 0 {
		return base[:marker] + "/v1"
	}
	if strings.HasSuffix(base, "/messages") {
		return strings.TrimSuffix(base, "/messages")
	}
	return base
}

func (s *LLMService) newAnthropicRequest(ctx context.Context, payload interface{}) (*http.Request, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.anthropicEndpoint(),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", s.cfg.APIKey)
	req.Header.Set("Anthropic-Version", "2023-06-01")
	return req, nil
}

func toAnthropicRequest(reqBody chatRequest) anthropicRequest {
	var systemParts []string
	messages := make([]ChatMessage, 0, len(reqBody.Messages))
	for _, message := range reqBody.Messages {
		switch strings.ToLower(message.Role) {
		case "system":
			if strings.TrimSpace(message.Content) != "" {
				systemParts = append(systemParts, message.Content)
			}
		case "assistant":
			messages = append(messages, ChatMessage{Role: "assistant", Content: message.Content})
		default:
			messages = append(messages, ChatMessage{Role: "user", Content: message.Content})
		}
	}
	return anthropicRequest{
		Model:       reqBody.Model,
		System:      strings.Join(systemParts, "\n\n"),
		Messages:    messages,
		Temperature: reqBody.Temperature,
		MaxTokens:   reqBody.MaxTokens,
		Stream:      reqBody.Stream,
	}
}

func (s *LLMService) doChatAnthropic(ctx context.Context, reqBody chatRequest) (*chatResponse, error) {
	req, err := s.newAnthropicRequest(ctx, toAnthropicRequest(reqBody))
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, s.wrapAPIError(resp)
	}

	var anthropicResult anthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&anthropicResult); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	var content strings.Builder
	for _, block := range anthropicResult.Content {
		if block.Type == "text" {
			content.WriteString(block.Text)
		}
	}

	return &chatResponse{
		Choices: []chatResponseChoice{{
			Message: chatResponseMessage{
				Role:    "assistant",
				Content: content.String(),
			},
			FinishReason: anthropicResult.StopReason,
		}},
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     anthropicResult.Usage.InputTokens,
			CompletionTokens: anthropicResult.Usage.OutputTokens,
			TotalTokens:      anthropicResult.Usage.InputTokens + anthropicResult.Usage.OutputTokens,
		},
	}, nil
}

func (s *LLMService) chatStreamAnthropic(
	ctx context.Context,
	reqBody chatRequest,
	onDelta func(delta string),
) error {
	anthropicBody := toAnthropicRequest(reqBody)
	anthropicBody.Stream = true
	req, err := s.newAnthropicRequest(ctx, anthropicBody)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/event-stream")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return s.wrapAPIError(resp)
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	sawText := false
	sawStop := false
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		var event anthropicStreamEvent
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}
		if event.Type == "error" {
			return fmt.Errorf(
				"%w: anthropic %s: %s",
				ErrLLMStreamFailed,
				event.Error.Type,
				truncate(event.Error.Message, 500),
			)
		}
		if event.Type == "message_stop" {
			sawStop = true
			break
		}
		if event.Type == "content_block_delta" &&
			event.Delta.Type == "text_delta" &&
			event.Delta.Text != "" {
			sawText = true
			if onDelta != nil {
				onDelta(event.Delta.Text)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("%w: %v", ErrLLMStreamFailed, err)
	}
	if !sawStop {
		return fmt.Errorf("%w: missing message_stop", ErrLLMStreamFailed)
	}
	if !sawText {
		return fmt.Errorf("%w: empty text stream", ErrLLMStreamFailed)
	}
	return nil
}

// wrapAPIError 读取错误响应体并包装为 error
func (s *LLMService) wrapAPIError(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	return fmt.Errorf("llm api error: status=%d body=%s", resp.StatusCode, truncate(string(body), 500))
}

// extractJSONBlock 从 markdown ```json ... ``` 代码块中提取 JSON
func extractJSONBlock(s string) string {
	start := strings.Index(s, "```json")
	if start < 0 {
		start = strings.Index(s, "```")
		if start < 0 {
			return ""
		}
		start += 3
	} else {
		start += 7
	}
	end := strings.Index(s[start:], "```")
	if end < 0 {
		return ""
	}
	return strings.TrimSpace(s[start : start+end])
}

// truncate 截断字符串到指定长度
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
