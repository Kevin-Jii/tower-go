package utils

import (
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	"math/big"
)

// HashPassword 使用 bcrypt 加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash 验证密码
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateStrongPassword 生成强随机密码
func GenerateStrongPassword(length int) (string, error) {
	if length < 8 {
		length = 8
	}

	// 定义字符集
	const (
		lowercase = "abcdefghijklmnopqrstuvwxyz"
		uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numbers   = "0123456789"
		special   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
		allChars  = lowercase + uppercase + numbers + special
	)

	password := make([]byte, length)

	// 确保至少包含一个大写字母、小写字母、数字和特殊字符
	charsets := []string{lowercase, uppercase, numbers, special}
	for i, charset := range charsets {
		if i >= length {
			break
		}
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[n.Int64()]
	}

	// 填充剩余位置
	for i := len(charsets); i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			return "", err
		}
		password[i] = allChars[n.Int64()]
	}

	// 随机打乱密码字符顺序
	for i := range password {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(length)))
		if err != nil {
			return "", err
		}
		password[i], password[j.Int64()] = password[j.Int64()], password[i]
	}

	return string(password), nil
}
