package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// 密码长度限制（与 PRD 字段规则对齐：至少 8 位）
const (
	PasswordMinLength = 8
	PasswordMaxLength = 64
)

// ErrPasswordTooShort 密码过短
var ErrPasswordTooShort = errors.New("password must be at least 8 characters")

// ErrPasswordTooLong 密码过长
var ErrPasswordTooLong = errors.New("password must be at most 64 characters")

// HashPassword 使用 bcrypt 哈希密码
func HashPassword(plain string) (string, error) {
	if len(plain) < PasswordMinLength {
		return "", ErrPasswordTooShort
	}
	if len(plain) > PasswordMaxLength {
		return "", ErrPasswordTooLong
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePassword 校验明文密码与哈希是否匹配
func ComparePassword(hashed, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}

// ValidatePasswordStrength 校验密码强度（长度 + 字母 + 数字）
// PRD 要求：至少 8 位，含字母和数字
func ValidatePasswordStrength(plain string) error {
	if len(plain) < PasswordMinLength {
		return ErrPasswordTooShort
	}
	if len(plain) > PasswordMaxLength {
		return ErrPasswordTooLong
	}
	hasLetter := false
	hasDigit := false
	for _, r := range plain {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z'):
			hasLetter = true
		case r >= '0' && r <= '9':
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return errors.New("password must contain both letters and digits")
	}
	return nil
}
