package helper

import (
	"golang.org/x/crypto/bcrypt"
)

// VerifyPassword 验证登录密码
func VerifyPassword(password []byte, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)

	return err != nil
}

// GeneratePassword 加密登录密码
func GeneratePassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
