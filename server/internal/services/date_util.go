package services

import (
	"fmt"
	"strings"
	"time"
)

// parseDate 解析日期字符串，支持 YYYY-MM-DD 与 YYYY-MM
func parseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, fmt.Errorf("empty date string")
	}
	// 尝试 YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	// 尝试 YYYY-MM
	if t, err := time.Parse("2006-01", s); err == nil {
		return t, nil
	}
	// 尝试 YYYY/MM/DD
	if t, err := time.Parse("2006/01/02", s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("unrecognized date format: %s", s)
}
