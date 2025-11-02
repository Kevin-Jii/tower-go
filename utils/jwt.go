package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// NOTE: Ensure this secret is consistent across your application, or load it from config.
var jwtSecret = []byte("32fdsfdsgfdsbvzxvasdfdsaf3qr2343@#@312r32/*-+dd2s3")

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
func GenerateToken(userID uint, username string, storeID uint, roleCode string, roleID uint) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		StoreID:  storeID,
		RoleCode: roleCode,
		RoleID:   roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken 验证并解析 Token
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	// 使用 jwtSecret 作为密钥进行解析
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 校验签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// 返回密钥
		return jwtSecret, nil
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
