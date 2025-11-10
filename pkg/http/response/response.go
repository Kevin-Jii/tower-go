package response

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// StatusCode 标准状态码
const (
	StatusSuccess           = 200  // 成功
	StatusCreated           = 201  // 创建成功
	StatusAccepted          = 202  // 已接受
	StatusNoContent         = 204  // 无内容
	StatusBadRequest        = 400  // 请求错误
	StatusUnauthorized      = 401  // 未授权
	StatusForbidden         = 403  // 禁止访问
	StatusNotFound          = 404  // 未找到
	StatusMethodNotAllowed  = 405  // 方法不允许
	StatusConflict          = 409  // 冲突
	StatusValidationFailed  = 422  // 验证失败
	StatusTooManyRequests   = 429  // 请求过多
	StatusInternalServerError = 500 // 内部服务器错误
	StatusNotImplemented    = 501  // 未实现
	StatusBadGateway        = 502  // 网关错误
	StatusServiceUnavailable = 503 // 服务不可用
)

// Response 标准API响应结构
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp int64       `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// Pagination 分页信息
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasMore    bool  `json:"has_more"`
}

// PaginatedResponse 分页响应结构
type PaginatedResponse struct {
	Response
	Pagination Pagination `json:"pagination"`
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Response
	Details interface{} `json:"details,omitempty"`
}

// Config 响应配置
type Config struct {
	EnableRequestID bool    `yaml:"enable_request_id" json:"enable_request_id"`
	DefaultMessage  string  `yaml:"default_message" json:"default_message"`
	StackTrace      bool    `yaml:"stack_trace" json:"stack_trace"`
	TimeFormat      string  `yaml:"time_format" json:"time_format"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		EnableRequestID: true,
		DefaultMessage:  "success",
		StackTrace:      false,
		TimeFormat:      time.RFC3339,
	}
}

// Responder 响应器
type Responder struct {
	config *Config
}

// New 创建响应器
func New(config *Config) *Responder {
	if config == nil {
		config = DefaultConfig()
	}
	return &Responder{config: config}
}

// NewWithDefaults 使用默认配置创建响应器
func NewWithDefaults() *Responder {
	return New(DefaultConfig())
}

// buildResponse 构建基础响应
func (r *Responder) buildResponse(ctx *gin.Context, code int, message string, data interface{}, err string) Response {
	if message == "" {
		message = r.config.DefaultMessage
	}

	resp := Response{
		Code:      code,
		Message:   message,
		Data:      data,
		Error:     err,
		Timestamp: time.Now().Unix(),
	}

	// 添加请求ID（如果启用且从context中获取）
	if r.config.EnableRequestID && ctx != nil {
		requestID := ctx.Value("request_id")
		if requestID != nil {
			if id, ok := requestID.(string); ok {
				resp.RequestID = id
			}
		}
	}

	return resp
}

// Success 成功响应
func (r *Responder) Success(ginContext *gin.Context, data interface{}) {
	resp := r.buildResponse(ginContext, StatusSuccess, "", data, "")
	ginContext.JSON(http.StatusOK, resp)
}

// SuccessWithMessage 带消息的成功响应
func (r *Responder) SuccessWithMessage(ginContext *gin.Context, message string, data interface{}) {
	resp := r.buildResponse(ginContext, StatusSuccess, message, data, "")
	ginContext.JSON(http.StatusOK, resp)
}

// Error 错误响应
func (r *Responder) Error(ginContext *gin.Context, code int, message string) {
	resp := r.buildResponse(ginContext, code, message, nil, "")
	statusCode := r.getHTTPStatus(code)
	ginContext.JSON(statusCode, resp)
}

// ErrorWithDetails 带详细信息的错误响应
func (r *Responder) ErrorWithDetails(ginContext *gin.Context, code int, message string, details interface{}) {
	resp := ErrorResponse{
		Response: r.buildResponse(ginContext, code, message, nil, ""),
		Details:  details,
	}
	statusCode := r.getHTTPStatus(code)
	ginContext.JSON(statusCode, resp)
}

// ErrorWithException 带异常的错误响应
func (r *Responder) ErrorWithException(ginContext *gin.Context, code int, message string, err error) {
	errorMsg := ""
	if err != nil && r.config.StackTrace {
		errorMsg = err.Error()
	}
	resp := r.buildResponse(ginContext, code, message, nil, errorMsg)
	statusCode := r.getHTTPStatus(code)
	ginContext.JSON(statusCode, resp)
}

// Custom 自定义HTTP状态码响应
func (r *Responder) Custom(ginContext *gin.Context, httpCode int, code int, message string, data interface{}) {
	resp := r.buildResponse(ginContext, code, message, data, "")
	ginContext.JSON(httpCode, resp)
}

// Paginated 分页响应
func (r *Responder) Paginated(ginContext *gin.Context, items interface{}, total int64, page, pageSize int) {
	totalPages := int(total / int64(pageSize))
	if total%int64(pageSize) != 0 {
		totalPages++
	}

	hasMore := page < totalPages

	pagination := Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		HasMore:    hasMore,
	}

	resp := PaginatedResponse{
		Response:   r.buildResponse(ginContext, StatusSuccess, "", items, ""),
		Pagination: pagination,
	}

	ginContext.JSON(http.StatusOK, resp)
}

// Created 创建成功响应
func (r *Responder) Created(ginContext *gin.Context, data interface{}) {
	resp := r.buildResponse(ginContext, StatusCreated, "Created successfully", data, "")
	ginContext.JSON(http.StatusCreated, resp)
}

// Accepted 已接受响应
func (r *Responder) Accepted(ginContext *gin.Context, message string) {
	if message == "" {
		message = "Request accepted"
	}
	resp := r.buildResponse(ginContext, StatusAccepted, message, nil, "")
	ginContext.JSON(http.StatusAccepted, resp)
}

// NoContent 无内容响应
func (r *Responder) NoContent(ginContext *gin.Context) {
	ginContext.Status(http.StatusNoContent)
}

// File 文件下载响应
func (r *Responder) File(ginContext *gin.Context, data []byte, filename string) {
	ginContext.Header("Content-Description", "File Transfer")
	ginContext.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ginContext.Header("Content-Type", "application/octet-stream")
	ginContext.Data(http.StatusOK, "application/octet-stream", data)
}

// Stream 流式响应
func (r *Responder) Stream(ginContext *gin.Context, data []byte, contentType string) {
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	ginContext.Data(http.StatusOK, contentType, data)
}

// getHTTPStatus 根据业务状态码获取HTTP状态码
func (r *Responder) getHTTPStatus(code int) int {
	switch {
	case code >= 200 && code < 300:
		return http.StatusOK
	case code == StatusCreated:
		return http.StatusCreated
	case code == StatusAccepted:
		return http.StatusAccepted
	case code == StatusNoContent:
		return http.StatusNoContent
	case code == StatusBadRequest:
		return http.StatusBadRequest
	case code == StatusUnauthorized:
		return http.StatusUnauthorized
	case code == StatusForbidden:
		return http.StatusForbidden
	case code == StatusNotFound:
		return http.StatusNotFound
	case code == StatusMethodNotAllowed:
		return http.StatusMethodNotAllowed
	case code == StatusConflict:
		return http.StatusConflict
	case code == StatusValidationFailed:
		return http.StatusUnprocessableEntity
	case code == StatusTooManyRequests:
		return http.StatusTooManyRequests
	case code == StatusInternalServerError:
		return http.StatusInternalServerError
	case code == StatusNotImplemented:
		return http.StatusNotImplemented
	case code == StatusBadGateway:
		return http.StatusBadGateway
	case code == StatusServiceUnavailable:
		return http.StatusServiceUnavailable
	default:
		// 对于未知状态码，返回4xx或5xx对应的HTTP状态码
		if code >= 400 && code < 500 {
			return http.StatusBadRequest
		} else if code >= 500 {
			return http.StatusInternalServerError
		}
		return http.StatusOK
	}
}

// GetConfig 获取配置
func (r *Responder) GetConfig() *Config {
	return r.config
}

// SetConfig 更新配置
func (r *Responder) SetConfig(config *Config) {
	r.config = config
}

// 全局默认响应器
var defaultResponder = NewWithDefaults()

// 便利函数（使用默认响应器）

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	defaultResponder.Success(c, data)
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	defaultResponder.Error(c, code, message)
}

// ErrorWithException 带异常的错误响应
func ErrorWithException(c *gin.Context, code int, message string, err error) {
	defaultResponder.ErrorWithException(c, code, message, err)
}

// BadRequest 400响应
func BadRequest(c *gin.Context, message string) {
	defaultResponder.Error(c, StatusBadRequest, message)
}

// Unauthorized 401响应
func Unauthorized(c *gin.Context, message string) {
	defaultResponder.Error(c, StatusUnauthorized, message)
}

// Forbidden 403响应
func Forbidden(c *gin.Context, message string) {
	defaultResponder.Error(c, StatusForbidden, message)
}

// NotFound 404响应
func NotFound(c *gin.Context, message string) {
	defaultResponder.Error(c, StatusNotFound, message)
}

// InternalServerError 500响应
func InternalServerError(c *gin.Context, message string) {
	defaultResponder.Error(c, StatusInternalServerError, message)
}

// ValidationFailed 验证失败响应
func ValidationFailed(c *gin.Context, message string) {
	defaultResponder.Error(c, StatusValidationFailed, message)
}

// TooManyRequests 429响应
func TooManyRequests(c *gin.Context, message string) {
	defaultResponder.Error(c, StatusTooManyRequests, message)
}

// SuccessWithPagination 分页成功响应
func SuccessWithPagination(c *gin.Context, items interface{}, total int64, page, pageSize int) {
	defaultResponder.Paginated(c, items, total, page, pageSize)
}

// Created 创建成功响应
func Created(c *gin.Context, data interface{}) {
	defaultResponder.Created(c, data)
}

// Accepted 已接受响应
func Accepted(c *gin.Context, message string) {
	defaultResponder.Accepted(c, message)
}

// Custom 自定义响应
func Custom(c *gin.Context, code int, message string, data interface{}) {
	defaultResponder.SuccessWithMessage(c, message, data)
}

// File 文件下载响应
func File(c *gin.Context, data []byte, filename string) {
	defaultResponder.File(c, data, filename)
}

// Stream 流式响应
func Stream(c *gin.Context, data []byte, contentType string) {
	defaultResponder.Stream(c, data, contentType)
}