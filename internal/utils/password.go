package utils

import (
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

const DefaultCost = 10

// HashPassword 生成密码哈希
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		slog.Error("Password hash failed: ", slog.Any("error", err))
		return ""
	}
	return string(bytes)
}

// CheckPassword 验证密码是否匹配
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
