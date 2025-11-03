package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	defaultPage     = 1
	defaultPageSize = 10
	maxPageSize     = 100
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Meta 分页或额外信息元数据
type Meta struct {
	Total     int64 `json:"total"`
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	PageCount int   `json:"page_count"`
	HasMore   bool  `json:"has_more"`
}

// StandardResponse 统一响应（支持可选 meta）
type StandardResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta"` // 始终返回（为空则为 null），方便前端统一解析
}

// PaginationResponse (已兼容旧结构，保留以避免前端立刻修改) — 标记为 deprecated
// Deprecated: 请使用 StandardResponse + Meta
type PaginationResponse struct {
	Response
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, StandardResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
	LogDebug("Success Response",
		zap.String("path", ctx.Request.URL.Path),
		zap.String("method", ctx.Request.Method),
	)
}

// SuccessWithPagination 旧调用入口，内部转换为新结构
func SuccessWithPagination(ctx *gin.Context, data interface{}, total int64, page, pageSize int) {
	pageCount := int((total + int64(pageSize) - 1) / int64(pageSize))
	hasMore := page < pageCount
	ctx.JSON(http.StatusOK, StandardResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
		Meta: &Meta{
			Total: total, Page: page, PageSize: pageSize, PageCount: pageCount, HasMore: hasMore,
		},
	})
}

// SuccessWithMeta 新的显式调用（如果未来有非分页 meta 也可复用）
func SuccessWithMeta(ctx *gin.Context, data interface{}, meta *Meta) {
	ctx.JSON(http.StatusOK, StandardResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
		Meta:    meta,
	})
}

func Error(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, StandardResponse{
		Code:    code,
		Message: message,
	})
	LogWarn("Error Response",
		zap.Int("code", code),
		zap.String("message", message),
		zap.String("path", ctx.Request.URL.Path),
		zap.String("method", ctx.Request.Method),
		zap.String("ip", ctx.ClientIP()),
	)
}

// ErrorWithCode 使用错误码响应
func ErrorWithCode(ctx *gin.Context, errCode ErrorCode) {
	httpCode := http.StatusBadRequest
	if errCode.Code >= 2000 && errCode.Code < 3000 {
		httpCode = http.StatusUnauthorized
	} else if errCode.Code >= 1000 && errCode.Code < 2000 {
		httpCode = http.StatusInternalServerError
	}

	ctx.JSON(httpCode, StandardResponse{
		Code:    errCode.Code,
		Message: errCode.Message,
	})

	LogBusinessError(errCode, nil,
		zap.String("path", ctx.Request.URL.Path),
		zap.String("method", ctx.Request.Method),
		zap.String("ip", ctx.ClientIP()),
	)
}

// ErrorWithCodeAndData 使用错误码响应（带额外数据）
func ErrorWithCodeAndData(ctx *gin.Context, errCode ErrorCode, data interface{}) {
	httpCode := http.StatusBadRequest
	if errCode.Code >= 2000 && errCode.Code < 3000 {
		httpCode = http.StatusUnauthorized
	} else if errCode.Code >= 1000 && errCode.Code < 2000 {
		httpCode = http.StatusInternalServerError
	}

	ctx.JSON(httpCode, StandardResponse{
		Code:    errCode.Code,
		Message: errCode.Message,
		Data:    data,
	})

	LogBusinessError(errCode, nil,
		zap.String("path", ctx.Request.URL.Path),
		zap.String("method", ctx.Request.Method),
		zap.String("ip", ctx.ClientIP()),
		zap.Any("data", data),
	)
}

func GetPage(ctx *gin.Context) int {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if page < 1 {
		return defaultPage
	}
	return page
}

func GetPageSize(ctx *gin.Context) int {
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if pageSize < 1 {
		return defaultPageSize
	}
	if pageSize > maxPageSize {
		return maxPageSize
	}
	return pageSize
}
