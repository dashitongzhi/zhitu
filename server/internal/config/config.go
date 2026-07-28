// Package config 负责加载和管理 YAML 配置文件
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config 全局配置根结构
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Admin    AdminConfig    `yaml:"admin"`
	LLM      LLMConfig      `yaml:"llm"`
	Storage  StorageConfig  `yaml:"storage"`
}

// ServerConfig HTTP 服务相关配置
type ServerConfig struct {
	Port         int      `yaml:"port"`
	Mode         string   `yaml:"mode"`
	AllowOrigins []string `yaml:"allow_origins"`
}

// DatabaseConfig 数据库相关配置
type DatabaseConfig struct {
	Path     string `yaml:"path"`
	LogLevel string `yaml:"log_level"`
}

// JWTConfig JWT 相关配置
type JWTConfig struct {
	Secret      string `yaml:"secret"`
	AdminSecret string `yaml:"admin_secret"` // 管理员 token 独立签名密钥，与用户隔离
	Issuer      string `yaml:"issuer"`
	ExpireHours int    `yaml:"expire_hours"`
}

// AdminConfig 后台管理员凭据（从 config 读取，不入库）
type AdminConfig struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

// LLMConfig LLM 接口配置
//
// 面试功能（InterviewRoom）模型选型约定：
//   - ChatModel      → "mimo-v2.5-pro"   文字分析（面试提问、评分、复盘报告）
//   - WhisperModel   → "mimo-v2.5-asr"   语音识别（用户语音回答转文字）
//   - TTSModel       → "mimo-v2.5-tts"   语音生成（AI 面试官朗读，用户可选开关）
//
// 部署时在 config.yaml 中把三者改为上述 mimo 模型即可对接 mimo 服务。
type LLMConfig struct {
	Provider     string  `yaml:"provider"` // openai / anthropic / mimo
	BaseURL      string  `yaml:"base_url"`
	APIKey       string  `yaml:"api_key"`
	ChatModel    string  `yaml:"chat_model"`    // 文字分析模型，推荐 mimo-v2.5-pro
	WhisperModel string  `yaml:"whisper_model"` // 语音识别 (ASR) 模型，推荐 mimo-v2.5-asr
	TTSModel     string  `yaml:"tts_model"`     // 语音生成 (TTS) 模型，推荐 mimo-v2.5-tts（用户可选）
	TTSVoice     string  `yaml:"tts_voice"`     // TTS 音色
	Temperature  float64 `yaml:"temperature"`
	MaxTokens    int     `yaml:"max_tokens"`
	TimeoutSec   int     `yaml:"timeout_sec"`
}

// StorageConfig 文件存储配置（本地磁盘）
type StorageConfig struct {
	BaseDir     string `yaml:"base_dir"`
	AudioDir    string `yaml:"audio_dir"`
	TTSDir      string `yaml:"tts_dir"`
	ResumeDir   string `yaml:"resume_dir"`
	MaxUploadMB int    `yaml:"max_upload_mb"`
}

// Load 从指定路径加载配置文件
// path 为空时按优先级查找：configs/config.local.yaml > configs/config.yaml
func Load(path string) (*Config, error) {
	candidates := []string{}
	if path != "" {
		candidates = append(candidates, path)
	} else {
		candidates = append(candidates, "configs/config.local.yaml", "configs/config.yaml")
	}

	var resolved string
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			resolved = p
			break
		}
	}
	if resolved == "" {
		return nil, fmt.Errorf("config file not found in: %s", strings.Join(candidates, ", "))
	}

	data, err := os.ReadFile(resolved)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", resolved, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", resolved, err)
	}

	if err := cfg.applyDefaults(); err != nil {
		return nil, err
	}

	// 确保数据库所在目录存在
	if cfg.Database.Path != "" {
		if dir := filepath.Dir(cfg.Database.Path); dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return nil, fmt.Errorf("create db dir %s: %w", dir, err)
			}
		}
	}

	return &cfg, nil
}

// applyDefaults 填充缺失字段的默认值
func (c *Config) applyDefaults() error {
	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	if c.Server.Mode == "" {
		c.Server.Mode = "debug"
	}
	if c.Database.Path == "" {
		c.Database.Path = "./data/zhitu.db"
	}
	if c.Database.LogLevel == "" {
		c.Database.LogLevel = "warn"
	}
	if c.JWT.Secret == "" {
		c.JWT.Secret = "zhitu-default-secret-please-change"
	}
	// 管理员 token 独立密钥，默认不回退到用户密钥，保证签名隔离
	if c.JWT.AdminSecret == "" {
		c.JWT.AdminSecret = "zhitu-admin-secret-please-change"
	}
	if c.JWT.Issuer == "" {
		c.JWT.Issuer = "zhitu"
	}
	if c.JWT.ExpireHours <= 0 {
		c.JWT.ExpireHours = 168
	}
	if c.Admin.Email == "" {
		c.Admin.Email = "admin@zhitu.com"
	}
	if c.Admin.Password == "" {
		c.Admin.Password = "admin123456"
	}
	// LLM 默认值
	if c.LLM.Provider == "" {
		c.LLM.Provider = "openai"
	}
	if c.LLM.BaseURL == "" {
		c.LLM.BaseURL = "https://api.openai.com/v1"
	}
	if c.LLM.ChatModel == "" {
		c.LLM.ChatModel = "gpt-4o"
	}
	if c.LLM.WhisperModel == "" {
		c.LLM.WhisperModel = "whisper-1"
	}
	if c.LLM.TTSModel == "" {
		c.LLM.TTSModel = "tts-1"
	}
	if c.LLM.TTSVoice == "" {
		c.LLM.TTSVoice = "alloy"
	}
	if c.LLM.Temperature == 0 {
		c.LLM.Temperature = 0.7
	}
	if c.LLM.MaxTokens == 0 {
		c.LLM.MaxTokens = 4096
	}
	if c.LLM.TimeoutSec == 0 {
		c.LLM.TimeoutSec = 60
	}
	// Storage 默认值
	if c.Storage.BaseDir == "" {
		c.Storage.BaseDir = "./data/uploads"
	}
	if c.Storage.AudioDir == "" {
		c.Storage.AudioDir = "./data/uploads/audio"
	}
	if c.Storage.TTSDir == "" {
		c.Storage.TTSDir = "./data/uploads/tts"
	}
	if c.Storage.ResumeDir == "" {
		c.Storage.ResumeDir = "./data/uploads/resume"
	}
	if c.Storage.MaxUploadMB == 0 {
		c.Storage.MaxUploadMB = 20
	}
	// 确保存储目录存在
	for _, dir := range []string{c.Storage.BaseDir, c.Storage.AudioDir, c.Storage.TTSDir, c.Storage.ResumeDir} {
		if dir != "" {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("create storage dir %s: %w", dir, err)
			}
		}
	}
	return nil
}
