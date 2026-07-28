package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// SaveFile 保存文件到指定目录，返回相对路径（用于 URL）与绝对路径
// filename 用于推断扩展名，实际文件名用随机 hex 防止冲突
func SaveFile(reader io.Reader, destDir, originalName string) (relPath, absPath string, err error) {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", "", fmt.Errorf("create dir: %w", err)
	}

	ext := filepath.Ext(originalName)
	randName := randomHex(16) + ext
	abs := filepath.Join(destDir, randName)

	f, err := os.Create(abs)
	if err != nil {
		return "", "", fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return "", "", fmt.Errorf("write file: %w", err)
	}

	// 返回相对路径用 / 分隔，便于前端拼接 URL
	cwd, _ := os.Getwd()
	rel := strings.TrimPrefix(abs, cwd+string(filepath.Separator))
	rel = filepath.ToSlash(rel)
	return rel, abs, nil
}

// SaveBytes 保存字节数据到文件
func SaveBytes(data []byte, destDir, filename string) (relPath, absPath string, err error) {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", "", fmt.Errorf("create dir: %w", err)
	}

	abs := filepath.Join(destDir, filename)
	if err := os.WriteFile(abs, data, 0o644); err != nil {
		return "", "", fmt.Errorf("write file: %w", err)
	}

	cwd, _ := os.Getwd()
	rel := strings.TrimPrefix(abs, cwd+string(filepath.Separator))
	rel = filepath.ToSlash(rel)
	return rel, abs, nil
}

// ValidateFileSize 校验文件大小是否超过限制（MB）
func ValidateFileSize(size int64, maxMB int) bool {
	return size <= int64(maxMB)*1024*1024
}

// IsAudioExt 判断是否为支持的音频扩展名
func IsAudioExt(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".mp3", ".wav", ".m4a", ".webm", ".ogg", ".flac":
		return true
	}
	return false
}

// IsMimoASRExt 判断文件是否为 MiMo-V2.5-ASR 当前支持的音频格式。
func IsMimoASRExt(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".mp3" || ext == ".wav"
}

// IsResumeExt 判断是否为支持的简历文件扩展名
func IsResumeExt(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".pdf", ".docx", ".doc", ".txt":
		return true
	}
	return false
}

// randomHex 生成 n 字节的随机 hex 字符串
func randomHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		// 极端情况退化为固定占位（实际不会发生）
		return "fallback00000000000000000000000000"
	}
	return hex.EncodeToString(b)
}
