package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Kevin-Jii/tower-go/utils/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Response API响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	resp := Response{
		Code:    200,
		Message: "success",
		Data:    data,
	}

	// 打印响应数据到控制台
	printResponse(c, resp)

	c.JSON(http.StatusOK, resp)
}

// printResponse 打印响应数据到控制台
func printResponse(c *gin.Context, resp Response) {
	jsonData, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("📤 API Response [%s %s]\n", c.Request.Method, c.Request.URL.Path)
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println(string(jsonData))
	fmt.Println(strings.Repeat("=", 60))
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	resp := Response{
		Code:    code,
		Message: message,
	}

	// 打印错误响应到控制台
	printErrorResponse(c, resp)

	c.JSON(http.StatusOK, resp)
}

// WithError 带错误信息的响应
func WithError(c *gin.Context, code int, message string, err error) {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
		logging.LogError("API Error", zap.String("message", message), zap.Error(err))
	}

	resp := Response{
		Code:    code,
		Message: message,
		Error:   errorMsg,
	}

	// 打印错误响应到控制台
	printErrorResponse(c, resp)

	c.JSON(http.StatusOK, resp)
}

// printErrorResponse 打印错误响应到控制台
func printErrorResponse(c *gin.Context, resp Response) {
	jsonData, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("❌ API Error [%s %s]\n", c.Request.Method, c.Request.URL.Path)
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println(string(jsonData))
	fmt.Println(strings.Repeat("=", 60))
}

// BadRequest 400 响应
func BadRequest(c *gin.Context, message string) {
	Error(c, 400, message)
}

// Unauthorized 401 响应
func Unauthorized(c *gin.Context, message string) {
	Error(c, 401, message)
}

// Forbidden 403 响应
func Forbidden(c *gin.Context, message string) {
	Error(c, 403, message)
}

// NotFound 404 响应
func NotFound(c *gin.Context, message string) {
	Error(c, 404, message)
}

// InternalServerError 500 响应
func InternalServerError(c *gin.Context, message string) {
	Error(c, 500, message)
}

// ValidationFailed 验证失败响应
func ValidationFailed(c *gin.Context, message string) {
	Error(c, 422, message)
}

// TooManyRequests 429 响应
func TooManyRequests(c *gin.Context, message string) {
	Error(c, 429, message)
}

// Custom 自定义响应
func Custom(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// Paginated 分页响应
type PaginatedResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	PageNum  int         `json:"page_num"`
}

// SuccessWithPagination 分页成功响应
func SuccessWithPagination(c *gin.Context, items interface{}, total int64, page, pageSize int) {
	pageNum := int(total / int64(pageSize))
	if total%int64(pageSize) != 0 {
		pageNum++
	}

	Custom(c, 200, "success", PaginatedResponse{
		List:     items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		PageNum:  pageNum,
	})
}

// File 文件下载响应
func File(c *gin.Context, data []byte, filename string) {
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.Data(http.StatusOK, "application/octet-stream", data)
}

// Stream 流式响应
func Stream(c *gin.Context, data []byte) {
	c.Data(http.StatusOK, "application/octet-stream", data)
}
