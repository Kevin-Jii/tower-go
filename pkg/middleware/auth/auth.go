package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Kevin-Jii/tower-go/pkg/auth/jwt"
	"github.com/Kevin-Jii/tower-go/pkg/http/response"
)

// Config 认证中间件配置
type Config struct {
	SecretKey              string        `yaml:"secret_key" json:"secret_key"`
	TokenLookup            string        `yaml:"token_lookup" json:"token_lookup"`             // "header:Authorization", "query:token", "cookie:token"
	AuthScheme             string        `yaml:"auth_scheme" json:"auth_scheme"`               // "Bearer"
	TimeFunc               func() time.Time `yaml:"-" json:"-"`                              // 时间函数，便于测试
	Timeout                time.Duration `yaml:"timeout" json:"timeout"`                     // 请求超时时间
	MaxRefresh             time.Duration `yaml:"max_refresh" json:"max_refresh"`             // 最大刷新时间
	IdentityKey            string        `yaml:"identity_key" json:"identity_key"`           // 身份键名
	IdentityHeaders        []string      `yaml:"identity_headers" json:"identity_headers"`   // 从HTTP头中提取的身份字段
	SendUnauthorizedHeader bool          `yaml:"send_unauthorized_header" json:"send_unauthorized_header"`
	UnauthorizedStatusCode int           `yaml:"unauthorized_status_code" json:"unauthorized_status_code"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		TokenLookup:            "header:Authorization",
		AuthScheme:             "Bearer",
		TimeFunc:               time.Now,
		Timeout:                time.Hour,
		MaxRefresh:             time.Hour * 24,
		IdentityKey:            "id",
		IdentityHeaders:        []string{},
		SendUnauthorizedHeader: true,
		UnauthorizedStatusCode: http.StatusUnauthorized,
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.SecretKey == "" {
		return fmt.Errorf("secret_key is required")
	}
	if c.TokenLookup == "" {
		return fmt.Errorf("token_lookup is required")
	}
	if c.AuthScheme == "" {
		return fmt.Errorf("auth_scheme is required")
	}
	if c.IdentityKey == "" {
		return fmt.Errorf("identity_key is required")
	}
	if c.UnauthorizedStatusCode <= 0 {
		c.UnauthorizedStatusCode = http.StatusUnauthorized
	}
	return nil
}

// JWTMiddleware JWT认证中间件
type JWTMiddleware struct {
	config      *Config
	jwtManager  *jwt.JWT
	extractor   TokenExtractor
}

// New 创建JWT认证中间件
func New(config *Config) (*JWTMiddleware, error) {
	if config == nil {
		config = DefaultConfig()
	}
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// 创建JWT管理器
	jwtConfig := jwt.DefaultConfig()
	jwtConfig.Secret = config.SecretKey
	jwtConfig.Expiration = config.Timeout
	jwtManager, err := jwt.New(jwtConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT manager: %w", err)
	}

	// 创建token提取器
	extractor, err := NewExtractor(config.TokenLookup, config.AuthScheme)
	if err != nil {
		return nil, fmt.Errorf("failed to create token extractor: %w", err)
	}

	return &JWTMiddleware{
		config:     config,
		jwtManager: jwtManager,
		extractor:  extractor,
	}, nil
}

// Middleware 返回Gin中间件函数
func (m *JWTMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 提取token
		token, err := m.extractor.ExtractToken(c)
		if err != nil {
			m.handleUnauthorized(c, err)
			return
		}

		// 验证token
		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			m.handleUnauthorized(c, fmt.Errorf("invalid token: %w", err))
			return
		}

		// 设置用户信息到上下文
		m.setIdentity(c, claims)

		// 继续处理请求
		c.Next()
	}
}

// RefreshMiddleware 支持token刷新的中间件
func (m *JWTMiddleware) RefreshMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 提取token
		token, err := m.extractor.ExtractToken(c)
		if err != nil {
			m.handleUnauthorized(c, err)
			return
		}

		// 验证refresh token
		claims, err := m.jwtManager.ValidateRefreshToken(token)
		if err != nil {
			m.handleUnauthorized(c, fmt.Errorf("invalid refresh token: %w", err))
			return
		}

		// 生成新的token对
		newTokenPair, err := m.jwtManager.GenerateTokenPair(claims)
		if err != nil {
			m.handleError(c, http.StatusInternalServerError, "failed to refresh token")
			return
		}

		// 设置用户信息到上下文
		m.setIdentity(c, claims)

		// 将新token对存储到上下文中
		c.Set("token_pair", newTokenPair)

		c.Next()
	}
}

// OptionalMiddleware 可选认证中间件（token存在时验证，不存在时跳过）
func (m *JWTMiddleware) OptionalMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试提取token
		token, err := m.extractor.ExtractToken(c)
		if err != nil {
			// token不存在，继续处理请求
			c.Next()
			return
		}

		// 验证token
		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			// token无效，继续处理请求
			c.Next()
			return
		}

		// 设置用户信息到上下文
		m.setIdentity(c, claims)

		c.Next()
	}
}

// RoleMiddleware 角色验证中间件
func (m *JWTMiddleware) RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleCode, exists := c.Get("role_code")
		if !exists {
			m.handleUnauthorized(c, fmt.Errorf("user role not found"))
			return
		}

		userRole, ok := roleCode.(string)
		if !ok {
			m.handleUnauthorized(c, fmt.Errorf("invalid user role format"))
			return
		}

		// 检查用户角色是否在允许列表中
		allowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				allowed = true
				break
			}
		}

		if !allowed {
			m.handleForbidden(c, fmt.Errorf("insufficient permissions: required roles %v", allowedRoles))
			return
		}

		c.Next()
	}
}

// AdminMiddleware 管理员验证中间件
func (m *JWTMiddleware) AdminMiddleware() gin.HandlerFunc {
	return m.RoleMiddleware("admin")
}

// StoreMiddleware 门店验证中间件
func (m *JWTMiddleware) StoreMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		storeID, exists := c.Get("store_id")
		if !exists {
			m.handleUnauthorized(c, fmt.Errorf("user store not found"))
			return
		}

		// 检查storeID是否有效
		if storeIDVal, ok := storeID.(uint); !ok || storeIDVal == 0 {
			m.handleUnauthorized(c, fmt.Errorf("invalid store"))
			return
		}

		c.Next()
	}
}

// setIdentity 设置用户身份信息到上下文
func (m *JWTMiddleware) setIdentity(c *gin.Context, claims *jwt.Claims) {
	c.Set("user_id", claims.UserID)
	c.Set("username", claims.Username)
	c.Set("store_id", claims.StoreID)
	c.Set("role_code", claims.RoleCode)
	c.Set("role_id", claims.RoleID)
	c.Set("token_type", claims.TokenType)

	// 设置自定义字段
	if claims.Custom != nil {
		for k, v := range claims.Custom {
			c.Set(k, v)
		}
	}
}

// handleUnauthorized 处理未授权错误
func (m *JWTMiddleware) handleUnauthorized(c *gin.Context, err error) {
	if m.config.SendUnauthorizedHeader {
		c.Header("WWW-Authenticate", m.config.AuthScheme+" realm=\"Restricted\"")
	}
	response.Error(c, m.config.UnauthorizedStatusCode, "Unauthorized")
	c.Abort()
}

// handleForbidden 处理禁止访问错误
func (m *JWTMiddleware) handleForbidden(c *gin.Context, err error) {
	response.Error(c, http.StatusForbidden, "Forbidden")
	c.Abort()
}

// handleError 处理一般错误
func (m *JWTMiddleware) handleError(c *gin.Context, statusCode int, message string) {
	response.Error(c, statusCode, message)
	c.Abort()
}

// GetConfig 获取配置
func (m *JWTMiddleware) GetConfig() *Config {
	return m.config
}

// GetJWTManager 获取JWT管理器
func (m *JWTMiddleware) GetJWTManager() *jwt.JWT {
	return m.jwtManager
}

// SetConfig 更新配置
func (m *JWTMiddleware) SetConfig(config *Config) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	m.config = config
	return nil
}

// ========== 工具函数 ==========

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) uint {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			return id
		}
	}
	return 0
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		if name, ok := username.(string); ok {
			return name
		}
	}
	return ""
}

// GetStoreID 从上下文获取门店ID
func GetStoreID(c *gin.Context) uint {
	if storeID, exists := c.Get("store_id"); exists {
		if id, ok := storeID.(uint); ok {
			return id
		}
	}
	return 0
}

// GetRoleCode 从上下文获取角色代码
func GetRoleCode(c *gin.Context) string {
	if roleCode, exists := c.Get("role_code"); exists {
		if code, ok := roleCode.(string); ok {
			return code
		}
	}
	return ""
}

// GetRoleID 从上下文获取角色ID
func GetRoleID(c *gin.Context) uint {
	if roleID, exists := c.Get("role_id"); exists {
		if id, ok := roleID.(uint); ok {
			return id
		}
	}
	return 0
}

// IsAdmin 检查是否为管理员
func IsAdmin(c *gin.Context) bool {
	return GetRoleCode(c) == "admin"
}

// IsAuthenticated 检查是否已认证
func IsAuthenticated(c *gin.Context) bool {
	return GetUserID(c) > 0
}

// HasRole 检查是否具有指定角色
func HasRole(c *gin.Context, role string) bool {
	return GetRoleCode(c) == role
}

// HasAnyRole 检查是否具有任意指定角色
func HasAnyRole(c *gin.Context, roles ...string) bool {
	userRole := GetRoleCode(c)
	for _, role := range roles {
		if userRole == role {
			return true
		}
	}
	return false
}

// RequireAuth 要求认证的中间件
func RequireAuth() gin.HandlerFunc {
	middleware, _ := New(DefaultConfig())
	return middleware.Middleware()
}

// RequireAdmin 要求管理员权限的中间件
func RequireAdmin() gin.HandlerFunc {
	middleware, _ := New(DefaultConfig())
	return middleware.AdminMiddleware()
}

// RequireRole 要求指定角色的中间件
func RequireRole(roles ...string) gin.HandlerFunc {
	middleware, _ := New(DefaultConfig())
	return middleware.RoleMiddleware(roles...)
}

// RequireStore 要求门店权限的中间件
func RequireStore() gin.HandlerFunc {
	middleware, _ := New(DefaultConfig())
	return middleware.StoreMiddleware()
}

// OptionalAuth 可选认证中间件
func OptionalAuth() gin.HandlerFunc {
	middleware, _ := New(DefaultConfig())
	return middleware.OptionalMiddleware()
}