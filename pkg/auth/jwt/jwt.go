package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Config JWT配置
type Config struct {
	Secret           string        `yaml:"secret" json:"secret"`
	Expiration       time.Duration `yaml:"expiration" json:"expiration"`
	Issuer           string        `yaml:"issuer" json:"issuer"`
	RefreshTokenExp  time.Duration `yaml:"refresh_token_exp" json:"refresh_token_exp"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Secret:          "", // 必须设置
		Expiration:      24 * time.Hour,
		Issuer:          "tower-go",
		RefreshTokenExp: 7 * 24 * time.Hour, // 7天
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Secret == "" {
		return errors.New("JWT secret is required")
	}
	if len(c.Secret) < 32 {
		return errors.New("JWT secret must be at least 32 characters long")
	}
	if c.Expiration <= 0 {
		return errors.New("JWT expiration must be positive")
	}
	return nil
}

// Claims JWT载荷结构体
type Claims struct {
	UserID    uint                   `json:"user_id"`
	Username  string                 `json:"username"`
	StoreID   uint                   `json:"store_id,omitempty"`
	RoleCode  string                 `json:"role_code,omitempty"`
	RoleID    uint                   `json:"role_id,omitempty"`
	TokenType string                 `json:"token_type,omitempty"` // access, refresh
	Custom    map[string]interface{} `json:"custom,omitempty"`      // 自定义字段
	jwt.RegisteredClaims
}

// TokenPair 令牌对
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// JWT JWT管理器
type JWT struct {
	config *Config
}

// New 创建JWT管理器
func New(config *Config) (*JWT, error) {
	if config == nil {
		config = DefaultConfig()
	}
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid JWT config: %w", err)
	}
	return &JWT{config: config}, nil
}

// NewWithSecret 使用密钥创建JWT管理器
func NewWithSecret(secret string) (*JWT, error) {
	config := DefaultConfig()
	config.Secret = secret
	return New(config)
}

// GenerateToken 生成访问令牌
func (j *JWT) GenerateToken(claims *Claims) (string, int64, error) {
	return j.generateTokenWithExpiration(claims, j.config.Expiration, "access")
}

// GenerateRefreshToken 生成刷新令牌
func (j *JWT) GenerateRefreshToken(claims *Claims) (string, int64, error) {
	return j.generateTokenWithExpiration(claims, j.config.RefreshTokenExp, "refresh")
}

// GenerateTokenPair 生成令牌对
func (j *JWT) GenerateTokenPair(claims *Claims) (*TokenPair, error) {
	// 生成访问令牌
	accessToken, expiresIn, err := j.GenerateToken(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// 生成刷新令牌
	refreshToken, _, err := j.GenerateRefreshToken(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
	}, nil
}

// generateTokenWithExpiration 生成指定过期时间的令牌
func (j *JWT) generateTokenWithExpiration(claims *Claims, expiration time.Duration, tokenType string) (string, int64, error) {
	now := time.Now()
	exp := now.Add(expiration)

	// 复制claims以避免修改原始数据
	tokenClaims := *claims
	tokenClaims.ExpiresAt = jwt.NewNumericDate(exp)
	tokenClaims.IssuedAt = jwt.NewNumericDate(now)
	tokenClaims.NotBefore = jwt.NewNumericDate(now)
	tokenClaims.Issuer = j.config.Issuer
	tokenClaims.TokenType = tokenType

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims)
	signed, err := token.SignedString([]byte(j.config.Secret))
	if err != nil {
		return "", 0, err
	}

	return signed, int64(expiration.Seconds()), nil
}

// ValidateToken 验证令牌
func (j *JWT) ValidateToken(tokenString string) (*Claims, error) {
	return j.parseToken(tokenString)
}

// ValidateRefreshToken 验证刷新令牌
func (j *JWT) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := j.parseToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid refresh token")
	}
	return claims, nil
}

// parseToken 解析令牌
func (j *JWT) parseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.config.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// 检查令牌类型
	if claims.TokenType == "" {
		// 兼容旧版本令牌
		claims.TokenType = "access"
	}

	return claims, nil
}

// RefreshToken 刷新访问令牌
func (j *JWT) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	// 验证刷新令牌
	claims, err := j.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// 生成新的令牌对
	newClaims := *claims
	newClaims.TokenType = "" // 清除令牌类型，让生成器重新设置
	newClaims.Custom = nil   // 清除自定义字段

	return j.GenerateTokenPair(&newClaims)
}

// ExtendToken 延长令牌有效期
func (j *JWT) ExtendToken(tokenString string, additionalDuration time.Duration) (string, int64, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", 0, fmt.Errorf("invalid token: %w", err)
	}

	// 计算新的过期时间
	newExp := time.Now().Add(additionalDuration)
	if claims.ExpiresAt != nil && newExp.After(claims.ExpiresAt.Time) {
		// 新的过期时间不能超过原有过期时间的某个限制
		maxExp := claims.ExpiresAt.Time.Add(j.config.Expiration)
		if newExp.After(maxExp) {
			newExp = maxExp
		}
	}

	claims.ExpiresAt = jwt.NewNumericDate(newExp)
	claims.IssuedAt = jwt.NewNumericDate(time.Now())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(j.config.Secret))
	if err != nil {
		return "", 0, err
	}

	return signed, int64(newExp.Sub(time.Now()).Seconds()), nil
}

// GetConfig 获取配置
func (j *JWT) GetConfig() *Config {
	return j.config
}

// SetConfig 更新配置
func (j *JWT) SetConfig(config *Config) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	j.config = config
	return nil
}