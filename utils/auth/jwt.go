package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// getJWTSecret 从环境变量获取JWT密钥，如果不存在则报错
func getJWTSecret() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET environment variable is not set")
	}
	if len(secret) < 32 {
		return nil, errors.New("JWT_SECRET must be at least 32 characters long")
	}
	return []byte(secret), nil
}

// Claims JWT载荷结构体，用于存储用户ID和门店ID
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	StoreID  uint   `json:"store_id"`  // 门店 ID
	RoleCode string `json:"role_code"` // 角色代码
	RoleID   uint   `json:"role_id"`   // 角色 ID
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint, username string, storeID uint, roleCode string, roleID uint) (string, int64, error) {
	expiration := time.Now().Add(24 * time.Hour)
	claims := Claims{
		UserID:   userID,
		Username: username,
		StoreID:  storeID,
		RoleCode: roleCode,
		RoleID:   roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	secret, err := getJWTSecret()
	if err != nil {
		return "", 0, fmt.Errorf("failed to get JWT secret: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", 0, err
	}
	return signed, int64(time.Until(expiration).Seconds()), nil
}

// ValidateToken 验证并解析 Token
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	// 获取JWT密钥
	secret, err := getJWTSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to get JWT secret: %w", err)
	}

	// 使用 secret 作为密钥进行解析
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 校验签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// 返回密钥
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	// 这里返回的 claims 已经被 ParseWithClaims 填充
	return claims, nil
}

// ParseToken 解析JWT token (与 ValidateToken 功能类似，可选择保留一个)
func ParseToken(tokenString string) (*Claims, error) {
	// 为了简化和避免冗余，建议只保留 ValidateToken
	return ValidateToken(tokenString)
}
