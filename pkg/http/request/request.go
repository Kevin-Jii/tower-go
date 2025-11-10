package request

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Validator 验证器实例
var validate = validator.New()

// PaginationConfig 分页配置
type PaginationConfig struct {
	DefaultPage int `yaml:"default_page" json:"default_page"`
	MinPage     int `yaml:"min_page" json:"min_page"`
	MaxPage     int `yaml:"max_page" json:"max_page"`

	DefaultPageSize int `yaml:"default_page_size" json:"default_page_size"`
	MinPageSize     int `yaml:"min_page_size" json:"min_page_size"`
	MaxPageSize     int `yaml:"max_page_size" json:"max_page_size"`
}

// DefaultPaginationConfig 返回默认分页配置
func DefaultPaginationConfig() *PaginationConfig {
	return &PaginationConfig{
		DefaultPage:     1,
		MinPage:         1,
		MaxPage:         1000,
		DefaultPageSize: 10,
		MinPageSize:     1,
		MaxPageSize:     100,
	}
}

// Request 请求处理器
type Request struct {
	paginationConfig *PaginationConfig
}

// New 创建请求处理器
func New(config *PaginationConfig) *Request {
	if config == nil {
		config = DefaultPaginationConfig()
	}
	return &Request{paginationConfig: config}
}

// NewWithDefaults 使用默认配置创建请求处理器
func NewWithDefaults() *Request {
	return New(DefaultPaginationConfig())
}

// ParseUintParam 解析路径参数为uint
func (r *Request) ParseUintParam(ctx *gin.Context, name string) (uint, error) {
	valStr := ctx.Param(name)
	if valStr == "" {
		return 0, fmt.Errorf("missing param: %s", name)
	}

	val, err := strconv.ParseUint(valStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: must be a positive integer", name)
	}

	return uint(val), nil
}

// ParseUintParamWithDefault 解析路径参数为uint，支持默认值
func (r *Request) ParseUintParamWithDefault(ctx *gin.Context, name string, defaultValue uint) (uint, error) {
	valStr := ctx.Param(name)
	if valStr == "" {
		return defaultValue, nil
	}

	val, err := strconv.ParseUint(valStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: must be a positive integer", name)
	}

	return uint(val), nil
}

// ParseIntParam 解析路径参数为int
func (r *Request) ParseIntParam(ctx *gin.Context, name string) (int, error) {
	valStr := ctx.Param(name)
	if valStr == "" {
		return 0, fmt.Errorf("missing param: %s", name)
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: must be an integer", name)
	}

	return val, nil
}

// ParseIntParamWithDefault 解析路径参数为int，支持默认值
func (r *Request) ParseIntParamWithDefault(ctx *gin.Context, name string, defaultValue int) (int, error) {
	valStr := ctx.Param(name)
	if valStr == "" {
		return defaultValue, nil
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: must be an integer", name)
	}

	return val, nil
}

// BindJSON 绑定JSON请求体并进行验证
func (r *Request) BindJSON(ctx *gin.Context, dst interface{}) error {
	if err := ctx.ShouldBindJSON(dst); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	if err := validate.Struct(dst); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return r.formatValidationError(validationErrors)
		}
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

// BindQuery 绑定查询参数并进行验证
func (r *Request) BindQuery(ctx *gin.Context, dst interface{}) error {
	if err := ctx.ShouldBindQuery(dst); err != nil {
		return fmt.Errorf("invalid query parameters: %w", err)
	}

	if err := validate.Struct(dst); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return r.formatValidationError(validationErrors)
		}
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

// BindURI 绑定URI参数并进行验证
func (r *Request) BindURI(ctx *gin.Context, dst interface{}) error {
	if err := ctx.ShouldBindUri(dst); err != nil {
		return fmt.Errorf("invalid URI parameters: %w", err)
	}

	if err := validate.Struct(dst); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return r.formatValidationError(validationErrors)
		}
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

// GetPage 获取分页页码
func (r *Request) GetPage(ctx *gin.Context) int {
	pageStr := ctx.DefaultQuery("page", strconv.Itoa(r.paginationConfig.DefaultPage))
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return r.paginationConfig.DefaultPage
	}

	if page < r.paginationConfig.MinPage {
		return r.paginationConfig.MinPage
	}

	if r.paginationConfig.MaxPage > 0 && page > r.paginationConfig.MaxPage {
		return r.paginationConfig.MaxPage
	}

	return page
}

// GetPageSize 获取分页大小
func (r *Request) GetPageSize(ctx *gin.Context) int {
	pageSizeStr := ctx.DefaultQuery("page_size", strconv.Itoa(r.paginationConfig.DefaultPageSize))
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		return r.paginationConfig.DefaultPageSize
	}

	if pageSize < r.paginationConfig.MinPageSize {
		return r.paginationConfig.MinPageSize
	}

	if r.paginationConfig.MaxPageSize > 0 && pageSize > r.paginationConfig.MaxPageSize {
		return r.paginationConfig.MaxPageSize
	}

	return pageSize
}

// GetOffset 计算分页偏移量
func (r *Request) GetOffset(ctx *gin.Context) int {
	page := r.GetPage(ctx)
	pageSize := r.GetPageSize(ctx)
	return (page - 1) * pageSize
}

// GetKeyword 获取搜索关键字
func (r *Request) GetKeyword(ctx *gin.Context) string {
	keyword := strings.TrimSpace(ctx.Query("keyword"))
	return keyword
}

// GetSort 获取排序字段和方向
func (r *Request) GetSort(ctx *gin.Context) (field, direction string) {
	sort := ctx.Query("sort")
	if sort == "" {
		return "", ""
	}

	// 解析排序字段和方向（例如：name:asc, created_at:desc）
	parts := strings.Split(sort, ":")
	if len(parts) >= 2 {
		field = strings.TrimSpace(parts[0])
		direction = strings.ToLower(strings.TrimSpace(parts[1]))
		if direction != "asc" && direction != "desc" {
			direction = "asc"
		}
	} else {
		field = strings.TrimSpace(parts[0])
		direction = "asc"
	}

	return field, direction
}

// GetFilter 获取过滤条件
func (r *Request) GetFilter(ctx *gin.Context) map[string]string {
	filters := make(map[string]string)

	for key, values := range ctx.Request.URL.Query() {
		if strings.HasPrefix(key, "filter_") {
			filterKey := strings.TrimPrefix(key, "filter_")
			if len(values) > 0 {
				filters[filterKey] = values[0]
			}
		}
	}

	return filters
}

// GetClientIP 获取客户端IP
func (r *Request) GetClientIP(ctx *gin.Context) string {
	// 优先检查X-Forwarded-For头
	xForwardedFor := ctx.GetHeader("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	// 检查X-Real-IP头
	xRealIP := ctx.GetHeader("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// 使用RemoteAddr
	return ctx.ClientIP()
}

// GetUserAgent 获取用户代理
func (r *Request) GetUserAgent(ctx *gin.Context) string {
	return ctx.GetHeader("User-Agent")
}

// GetBearerToken 获取Bearer令牌
func (r *Request) GetBearerToken(ctx *gin.Context) string {
	auth := ctx.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}

// RequireRole 检查用户角色
func (r *Request) RequireRole(ctx *gin.Context, requiredRole string) error {
	role, exists := ctx.Get("roleCode")
	if !exists {
		return fmt.Errorf("user role not found")
	}

	userRole, ok := role.(string)
	if !ok {
		return fmt.Errorf("invalid user role format")
	}

	if userRole != requiredRole {
		return fmt.Errorf("insufficient permissions: required role %s", requiredRole)
	}

	return nil
}

// RequireAnyRole 检查用户是否具有任意一个指定角色
func (r *Request) RequireAnyRole(ctx *gin.Context, roles ...string) error {
	if len(roles) == 0 {
		return nil
	}

	userRole, exists := ctx.Get("roleCode")
	if !exists {
		return fmt.Errorf("user role not found")
	}

	roleStr, ok := userRole.(string)
	if !ok {
		return fmt.Errorf("invalid user role format")
	}

	for _, role := range roles {
		if roleStr == role {
			return nil
		}
	}

	return fmt.Errorf("insufficient permissions: required one of roles %v", roles)
}

// IsAdmin 检查是否为管理员
func (r *Request) IsAdmin(ctx *gin.Context) bool {
	err := r.RequireRole(ctx, "admin")
	return err == nil
}

// formatValidationError 格式化验证错误信息
func (r *Request) formatValidationError(errs validator.ValidationErrors) error {
	var messages []string
	for _, err := range errs {
		message := ""
		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required", err.Field())
		case "min":
			message = fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
		case "max":
			message = fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param())
		case "len":
			message = fmt.Sprintf("%s must be %s characters", err.Field(), err.Param())
		case "email":
			message = fmt.Sprintf("%s must be a valid email", err.Field())
		case "numeric":
			message = fmt.Sprintf("%s must be numeric", err.Field())
		case "alpha":
			message = fmt.Sprintf("%s must contain only letters", err.Field())
		case "alphanum":
			message = fmt.Sprintf("%s must contain only letters and numbers", err.Field())
		default:
			message = fmt.Sprintf("%s is invalid", err.Field())
		}
		messages = append(messages, message)
	}
	return fmt.Errorf(strings.Join(messages, "; "))
}

// StructToMap 将结构体转换为map
func (r *Request) StructToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return result
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// 跳过非导出字段
		if !field.IsExported() {
			continue
		}

		// 获取json标签
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" || jsonTag == "" {
			jsonTag = field.Name
		} else {
			// 处理json标签中的omitempty等选项
			if parts := strings.Split(jsonTag, ","); len(parts) > 0 {
				jsonTag = parts[0]
			}
		}

		// 处理嵌套结构体
		if value.Kind() == reflect.Struct {
			nested := r.StructToMap(value.Interface())
			for k, v := range nested {
				result[k] = v
			}
		} else {
			result[jsonTag] = value.Interface()
		}
	}

	return result
}

// GetPaginationConfig 获取分页配置
func (r *Request) GetPaginationConfig() *PaginationConfig {
	return r.paginationConfig
}

// SetPaginationConfig 设置分页配置
func (r *Request) SetPaginationConfig(config *PaginationConfig) {
	r.paginationConfig = config
}

// 全局默认请求处理器
var defaultRequest = NewWithDefaults()

// 便利函数（使用默认请求处理器）

// ParseUintParam 解析路径参数为uint
func ParseUintParam(ctx *gin.Context, name string) (uint, error) {
	return defaultRequest.ParseUintParam(ctx, name)
}

// ParseUintParamWithDefault 解析路径参数为uint，支持默认值
func ParseUintParamWithDefault(ctx *gin.Context, name string, defaultValue uint) (uint, error) {
	return defaultRequest.ParseUintParamWithDefault(ctx, name, defaultValue)
}

// ParseIntParam 解析路径参数为int
func ParseIntParam(ctx *gin.Context, name string) (int, error) {
	return defaultRequest.ParseIntParam(ctx, name)
}

// BindJSON 绑定JSON请求体
func BindJSON(ctx *gin.Context, dst interface{}) error {
	return defaultRequest.BindJSON(ctx, dst)
}

// BindQuery 绑定查询参数
func BindQuery(ctx *gin.Context, dst interface{}) error {
	return defaultRequest.BindQuery(ctx, dst)
}

// GetPage 获取分页页码
func GetPage(ctx *gin.Context) int {
	return defaultRequest.GetPage(ctx)
}

// GetPageSize 获取分页大小
func GetPageSize(ctx *gin.Context) int {
	return defaultRequest.GetPageSize(ctx)
}

// GetOffset 计算分页偏移量
func GetOffset(ctx *gin.Context) int {
	return defaultRequest.GetOffset(ctx)
}

// GetKeyword 获取搜索关键字
func GetKeyword(ctx *gin.Context) string {
	return defaultRequest.GetKeyword(ctx)
}

// RequireRole 检查用户角色
func RequireRole(ctx *gin.Context, requiredRole string) error {
	return defaultRequest.RequireRole(ctx, requiredRole)
}

// RequireAdmin 检查是否为管理员
func RequireAdmin(ctx *gin.Context) bool {
	return defaultRequest.IsAdmin(ctx)
}

// StructToMap 将结构体转换为map
func StructToMap(obj interface{}) map[string]interface{} {
	return defaultRequest.StructToMap(obj)
}