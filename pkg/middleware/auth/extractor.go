package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// TokenExtractor token提取器接口
type TokenExtractor interface {
	ExtractToken(c *gin.Context) (string, error)
}

// NewExtractor 创建token提取器
func NewExtractor(lookup, authScheme string) (TokenExtractor, error) {
	if lookup == "" {
		lookup = "header:Authorization"
	}
	if authScheme == "" {
		authScheme = "Bearer"
	}

	parts := strings.Split(lookup, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid token lookup format: %s", lookup)
	}

	source := strings.TrimSpace(parts[0])
	name := strings.TrimSpace(parts[1])

	switch source {
	case "header":
		return &HeaderExtractor{
			Name:       name,
			AuthScheme: authScheme,
		}, nil
	case "query":
		return &QueryExtractor{
			Name: name,
		}, nil
	case "cookie":
		return &CookieExtractor{
			Name: name,
		}, nil
	case "form":
		return &FormExtractor{
			Name: name,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported token source: %s", source)
	}
}

// HeaderExtractor 从HTTP头提取token
type HeaderExtractor struct {
	Name       string
	AuthScheme string
}

// ExtractToken 提取token
func (e *HeaderExtractor) ExtractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader(e.Name)
	if authHeader == "" {
		return "", fmt.Errorf("missing %s header", e.Name)
	}

	// 如果不需要验证认证方案，直接返回
	if e.AuthScheme == "" {
		return authHeader, nil
	}

	// 验证认证方案
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && strings.EqualFold(parts[0], e.AuthScheme)) {
		return "", fmt.Errorf("invalid authorization header format, expected %s scheme", e.AuthScheme)
	}

	return strings.TrimSpace(parts[1]), nil
}

// QueryExtractor 从查询参数提取token
type QueryExtractor struct {
	Name string
}

// ExtractToken 提取token
func (e *QueryExtractor) ExtractToken(c *gin.Context) (string, error) {
	token := c.Query(e.Name)
	if token == "" {
		return "", fmt.Errorf("missing %s query parameter", e.Name)
	}
	return token, nil
}

// CookieExtractor 从Cookie提取token
type CookieExtractor struct {
	Name string
}

// ExtractToken 提取token
func (e *CookieExtractor) ExtractToken(c *gin.Context) (string, error) {
	cookie, err := c.Cookie(e.Name)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", fmt.Errorf("missing %s cookie", e.Name)
		}
		return "", fmt.Errorf("failed to read %s cookie: %w", e.Name, err)
	}
	return cookie, nil
}

// FormExtractor 从表单数据提取token
type FormExtractor struct {
	Name string
}

// ExtractToken 提取token
func (e *FormExtractor) ExtractToken(c *gin.Context) (string, error) {
	// 检查Content-Type
	contentType := c.GetHeader("Content-Type")
	if !strings.Contains(contentType, "application/x-www-form-urlencoded") &&
		!strings.Contains(contentType, "multipart/form-data") {
		return "", fmt.Errorf("invalid content type for form extraction: %s", contentType)
	}

	// 解析表单
	if err := c.Request.ParseForm(); err != nil {
		return "", fmt.Errorf("failed to parse form: %w", err)
	}

	token := c.PostForm(e.Name)
	if token == "" {
		return "", fmt.Errorf("missing %s form field", e.Name)
	}

	return token, nil
}

// MultiExtractor 多源token提取器
type MultiExtractor struct {
	extractors []TokenExtractor
}

// NewMultiExtractor 创建多源token提取器
func NewMultiExtractor(extractors ...TokenExtractor) *MultiExtractor {
	return &MultiExtractor{
		extractors: extractors,
	}
}

// ExtractToken 按顺序尝试从多个源提取token
func (e *MultiExtractor) ExtractToken(c *gin.Context) (string, error) {
	for _, extractor := range e.extractors {
		token, err := extractor.ExtractToken(c)
		if err == nil {
			return token, nil
		}
	}
	return "", fmt.Errorf("failed to extract token from all sources")
}

// ChainExtractor 链式token提取器
type ChainExtractor struct {
	extractors []TokenExtractor
}

// NewChainExtractor 创建链式token提取器
func NewChainExtractor(extractors ...TokenExtractor) *ChainExtractor {
	return &ChainExtractor{
		extractors: extractors,
	}
}

// ExtractToken 链式提取token，按优先级尝试
func (e *ChainExtractor) ExtractToken(c *gin.Context) (string, error) {
	var lastErr error

	for _, extractor := range e.extractors {
		token, err := extractor.ExtractToken(c)
		if err == nil {
			return token, nil
		}
		lastErr = err
	}

	return "", lastErr
}

// ConditionalExtractor 条件token提取器
type ConditionalExtractor struct {
	condition func(*gin.Context) bool
	extractor TokenExtractor
}

// NewConditionalExtractor 创建条件token提取器
func NewConditionalExtractor(condition func(*gin.Context) bool, extractor TokenExtractor) *ConditionalExtractor {
	return &ConditionalExtractor{
		condition: condition,
		extractor: extractor,
	}
}

// ExtractToken 条件提取token
func (e *ConditionalExtractor) ExtractToken(c *gin.Context) (string, error) {
	if e.condition != nil && e.condition(c) {
		return e.extractor.ExtractToken(c)
	}
	return "", fmt.Errorf("condition not met for token extraction")
}

// FallbackExtractor 回退token提取器
type FallbackExtractor struct {
	primary    TokenExtractor
	secondary  TokenExtractor
	fallback   TokenExtractor
}

// NewFallbackExtractor 创建回退token提取器
func NewFallbackExtractor(primary, secondary, fallback TokenExtractor) *FallbackExtractor {
	return &FallbackExtractor{
		primary:   primary,
		secondary: secondary,
		fallback:  fallback,
	}
}

// ExtractToken 回退提取token
func (e *FallbackExtractor) ExtractToken(c *gin.Context) (string, error) {
	// 尝试主要提取器
	if token, err := e.primary.ExtractToken(c); err == nil {
		return token, nil
	}

	// 尝试次要提取器
	if e.secondary != nil {
		if token, err := e.secondary.ExtractToken(c); err == nil {
			return token, nil
		}
	}

	// 使用回退提取器
	if e.fallback != nil {
		return e.fallback.ExtractToken(c)
	}

	return "", fmt.Errorf("all token extractors failed")
}