package services

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
)

// 简历解析相关错误
var (
	ErrUnsupportedResumeFormat = errors.New("unsupported resume format, only .pdf/.docx/.txt supported")
	ErrDocxParseFailed         = errors.New("failed to parse docx file")
)

// extractResumeText 根据文件扩展名提取简历纯文本
func extractResumeText(reader io.Reader, filename string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".txt":
		return extractTXT(reader)
	case ".pdf":
		return extractPDF(reader)
	case ".docx":
		return extractDOCX(reader)
	case ".doc":
		return "", errors.New("legacy .doc format not supported, please convert to .docx")
	default:
		return "", ErrUnsupportedResumeFormat
	}
}

// extractTXT 直接读取文本
func extractTXT(reader io.Reader) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// extractPDF 提取 PDF 文本
// 使用 github.com/ledongthuc/pdf（纯 Go，无 CGO）
func extractPDF(reader io.Reader) (string, error) {
	// ledongthuc/pdf 需要文件路径或 bytes.Reader
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("read pdf: %w", err)
	}
	if len(data) == 0 {
		return "", errors.New("empty pdf content")
	}

	r, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("open pdf: %w", err)
	}

	var buf strings.Builder
	numPages := r.NumPage()
	for i := 1; i <= numPages; i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		text, err := page.GetPlainText(nil)
		if err != nil {
			// 单页失败不影响整体
			continue
		}
		buf.WriteString(text)
		buf.WriteString("\n")
	}
	return buf.String(), nil
}

// extractDOCX 提取 .docx 文本
// .docx 本质是 zip，word/document.xml 内是正文
func extractDOCX(reader io.Reader) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("read docx: %w", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrDocxParseFailed, err)
	}

	var docXML []byte
	for _, f := range zipReader.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return "", fmt.Errorf("open document.xml: %w", err)
			}
			docXML, err = io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return "", fmt.Errorf("read document.xml: %w", err)
			}
			break
		}
	}
	if docXML == nil {
		return "", errors.New("document.xml not found in docx")
	}

	return extractTextFromDocXML(docXML), nil
}

// extractTextFromDocXML 从 document.xml 提取纯文本
// 简单实现：提取 <w:t> 标签内容，遇到 <w:p>（段落）添加换行
var (
	reWP = regexp.MustCompile(`<w:p[ >]`)
	reWT = regexp.MustCompile(`<w:t[^>]*>([^<]*)</w:t>`)
)

func extractTextFromDocXML(xml []byte) string {
	// 按段落分割
	paragraphs := reWP.Split(string(xml), -1)
	var buf strings.Builder
	for _, p := range paragraphs {
		// 提取该段所有 <w:t> 文本
		matches := reWT.FindAllStringSubmatch(p, -1)
		for _, m := range matches {
			if len(m) >= 2 {
				buf.WriteString(m[1])
			}
		}
		buf.WriteString("\n")
	}
	return strings.TrimSpace(buf.String())
}
