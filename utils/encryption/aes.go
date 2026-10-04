// Package encryption 提供 AES-256-GCM 加解密工具，用于敏感字段（如阿里云 SMS AccessKey Secret）
// 在数据库中加密落表。密钥来源于 APP_ENCRYPTION_KEY 环境变量（推荐 32 字节 base64 / hex / 任意 32+ 字节），
// 若未设置则退化为从 JWT_SECRET 派生（仅用于本地开发，生产环境务必显式设置 APP_ENCRYPTION_KEY）。
package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"strings"
)

var (
	ErrEmptyPlaintext    = errors.New("encryption: plaintext is empty")
	ErrEmptyCiphertext   = errors.New("encryption: ciphertext is empty")
	ErrCiphertextCorrupt = errors.New("encryption: ciphertext is corrupted")
)

const (
	nonceSize = 12
	// 加密结果 base64(nonce + ciphertext + tag)。空字符串表示明文为空字符串。
	emptyMarker = "\x00"
)

// keyMaterial 从环境变量派生 32 字节密钥。
func keyMaterial() []byte {
	raw := strings.TrimSpace(os.Getenv("APP_ENCRYPTION_KEY"))
	if raw == "" {
		raw = strings.TrimSpace(os.Getenv("JWT_SECRET"))
	}
	sum := sha256.Sum256([]byte(raw))
	return sum[:]
}

// Encrypt 用 AES-256-GCM 加密明文，返回 base64(nonce || ciphertext || tag)。
// 明文为空时返回空字符串。
func Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	block, err := aes.NewCipher(keyMaterial())
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	cipherText := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	buf := make([]byte, 0, len(nonce)+len(cipherText))
	buf = append(buf, nonce...)
	buf = append(buf, cipherText...)
	return base64.StdEncoding.EncodeToString(buf), nil
}

// Decrypt 解密 Encrypt 返回的 base64。
// 密文为空时返回空字符串。
func Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}
	raw, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", ErrCiphertextCorrupt
	}
	if len(raw) < nonceSize {
		return "", ErrCiphertextCorrupt
	}
	block, err := aes.NewCipher(keyMaterial())
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := raw[:nonceSize]
	sealed := raw[nonceSize:]
	plain, err := gcm.Open(nil, nonce, sealed, nil)
	if err != nil {
		return "", ErrCiphertextCorrupt
	}
	return string(plain), nil
}

// MaskSecret 仅展示 AccessKey Secret 前后各 4 位，中间省略。
func MaskSecret(s string) string {
	if len(s) <= 8 {
		return strings.Repeat("*", len(s))
	}
	return s[:4] + strings.Repeat("*", len(s)-8) + s[len(s)-4:]
}
