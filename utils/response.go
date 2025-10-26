package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

type PaginationResponse struct {
	Response
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

func SuccessWithPagination(ctx *gin.Context, data interface{}, total int64, page, pageSize int) {
	ctx.JSON(http.StatusOK, PaginationResponse{
		Response: Response{
			Code:    http.StatusOK,
			Message: "success",
			Data:    data,
		},
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func Error(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, Response{
		Code:    code,
		Message: message,
	})
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
