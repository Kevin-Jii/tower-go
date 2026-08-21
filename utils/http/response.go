package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/pkg/apicode"
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
		Data:    normalizeResponseData(data),
	}

	// 打印响应数据到控制台（已禁用）
	// printResponse(c, resp)

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

// ErrorApp 使用统一错误码结构响应（推荐新业务与鉴权链路使用）
func ErrorApp(c *gin.Context, co apicode.Code) {
	Error(c, co.Num, co.Msg)
}

// ErrorFrom 将 Service/Module 返回的统一业务错误转换为 API 响应。
// 未包装的底层错误不会直接暴露给客户端，只返回统一内部错误提示。
func ErrorFrom(c *gin.Context, err error) {
	if err == nil {
		ErrorApp(c, apicode.InternalError)
		return
	}
	if code, ok := apicode.Resolve(err); ok {
		ErrorApp(c, code)
		return
	}
	logging.LogError("API Error", zap.Error(err))
	ErrorApp(c, apicode.InternalError)
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	resp := Response{
		Code:    code,
		Message: message,
	}

	// 打印错误响应到控制台（已禁用）
	// printErrorResponse(c, resp)

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

	// 打印错误响应到控制台（已禁用）
	// printErrorResponse(c, resp)

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
		Data:    normalizeResponseData(data),
	})
}

const (
	responseDateFormat     = "2006-01-02"
	responseDateTimeFormat = "2006-01-02 15:04:05"
)

var (
	timeType          = reflect.TypeOf(time.Time{})
	jsonMarshalerType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
)

// normalizeResponseData converts time values before Gin serializes the response.
// This keeps date formatting consistent for nested structs, slices and maps.
func normalizeResponseData(data interface{}) interface{} {
	if data == nil {
		return nil
	}
	return normalizeResponseValue(reflect.ValueOf(data), "")
}

func normalizeResponseValue(value reflect.Value, fieldName string) interface{} {
	if !value.IsValid() {
		return nil
	}
	if value.Kind() == reflect.Interface {
		if value.IsNil() {
			return nil
		}
		return normalizeResponseValue(value.Elem(), fieldName)
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		return normalizeResponseValue(value.Elem(), fieldName)
	}
	if value.Type() == timeType {
		date := value.Interface().(time.Time)
		if isResponseDateOnlyField(fieldName) {
			return date.Format(responseDateFormat)
		}
		return date.Format(responseDateTimeFormat)
	}
	if value.CanInterface() && (value.Type().Implements(jsonMarshalerType) ||
		(value.CanAddr() && value.Addr().Type().Implements(jsonMarshalerType))) {
		return value.Interface()
	}

	switch value.Kind() {
	case reflect.Struct:
		return normalizeResponseStruct(value)
	case reflect.Slice, reflect.Array:
		if value.Type().Elem().Kind() == reflect.Uint8 {
			if value.CanInterface() {
				return value.Interface()
			}
			return nil
		}
		items := make([]interface{}, value.Len())
		for i := 0; i < value.Len(); i++ {
			items[i] = normalizeResponseValue(value.Index(i), fieldName)
		}
		return items
	case reflect.Map:
		if value.IsNil() {
			return nil
		}
		if value.Type().Key().Kind() == reflect.String {
			result := make(map[string]interface{}, value.Len())
			for _, key := range value.MapKeys() {
				result[key.String()] = normalizeResponseValue(value.MapIndex(key), key.String())
			}
			return result
		}
		result := make(map[interface{}]interface{}, value.Len())
		for _, key := range value.MapKeys() {
			result[key.Interface()] = normalizeResponseValue(value.MapIndex(key), fmt.Sprint(key.Interface()))
		}
		return result
	default:
		if value.CanInterface() {
			return value.Interface()
		}
		return nil
	}
}

func normalizeResponseStruct(value reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})
	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			continue
		}
		tag := field.Tag.Get("json")
		parts := strings.Split(tag, ",")
		name := parts[0]
		if name == "-" {
			continue
		}
		if name == "" {
			name = field.Name
		}
		fieldValue := value.Field(i)
		if hasJSONOption(parts[1:], "omitempty") && isResponseEmptyValue(fieldValue) {
			continue
		}
		if field.Anonymous && name == field.Name {
			if embedded, ok := normalizeResponseValue(fieldValue, fieldNameFromType(fieldValue)).(map[string]interface{}); ok {
				for key, item := range embedded {
					result[key] = item
				}
			}
			continue
		}
		result[name] = normalizeResponseValue(fieldValue, name)
	}
	return result
}

func hasJSONOption(options []string, target string) bool {
	for _, option := range options {
		if option == target {
			return true
		}
	}
	return false
}

func isResponseEmptyValue(value reflect.Value) bool {
	if !value.IsValid() {
		return true
	}
	return value.IsZero()
}

func fieldNameFromType(value reflect.Value) string {
	if value.IsValid() {
		return value.Type().Name()
	}
	return ""
}

func isResponseDateOnlyField(fieldName string) bool {
	name := strings.ToLower(fieldName)
	return name == "date" || strings.HasSuffix(name, "_date")
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
