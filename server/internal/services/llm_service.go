// Package services 承载业务逻辑，对 handler 层提供无 HTTP 上下文的纯函数接口
package services

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/zhitu/server/internal/config"
)

// LLMService 封装 OpenAI 兼容接口的调用
// 直接用 net/http，不引入第三方 SDK，便于切换 provider
type LLMService struct {
	cfg    *config.LLMConfig
	client *http.Client
}

// ChatMessage OpenAI Chat Completions 消息格式
type ChatMessage struct {
	Role    string `json:"role"`    // system / user / assistant
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

// chatResponse 非流式响应体
type chatResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
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

// whisperResponse Whisper 转写响应
type whisperResponse struct {
	Text string `json:"text"`
}

// ttsRequest TTS 请求体
type ttsRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
	Voice string `json:"voice"`
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

// Transcribe 调用 Whisper 将音频转写为文本
// audio 为音频文件内容（MP3/WAV/M4A 等），filename 用于 multipart 文件名
func (s *LLMService) Transcribe(ctx context.Context, audio io.Reader, filename string) (string, error) {
	if !s.IsConfigured() {
		return "", ErrLLMNotConfigured
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 音频文件字段
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("create form file: %w", err)
	}
	if _, err := io.Copy(part, audio); err != nil {
		return "", fmt.Errorf("copy audio: %w", err)
	}
	// model 字段
	if err := writer.WriteField("model", s.cfg.WhisperModel); err != nil {
		return "", fmt.Errorf("write model field: %w", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.BaseURL+"/audio/transcriptions", &buf)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.cfg.APIKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", s.wrapAPIError(resp)
	}

	var result whisperResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode whisper response: %w", err)
	}
	return strings.TrimSpace(result.Text), nil
}

// Synthesize 调用 TTS 将文本合成为 MP3 音频字节
func (s *LLMService) Synthesize(ctx context.Context, text string) ([]byte, error) {
	if !s.IsConfigured() {
		return nil, ErrLLMNotConfigured
	}
	if text == "" {
		return nil, errors.New("tts input text is empty")
	}

	reqBody := ttsRequest{
		Model: s.cfg.TTSModel,
		Input: text,
		Voice: s.cfg.TTSVoice,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.BaseURL+"/audio/speech", bytes.NewReader(body))
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

	audio, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read tts audio: %w", err)
	}
	return audio, nil
}

// doChat 执行非流式 Chat 请求
func (s *LLMService) doChat(ctx context.Context, reqBody chatRequest) (*chatResponse, error) {
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
