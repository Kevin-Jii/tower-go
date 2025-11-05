package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ParseUintParam 解析 path 参数为 uint，失败统一返回 400，并给出消息。返回 (value, ok)
func ParseUintParam(ctx *gin.Context, name string) (uint, bool) {
	valStr := ctx.Param(name)
	if valStr == "" {
		Error(ctx, http.StatusBadRequest, "missing param: "+name)
		return 0, false
	}
	v, err := strconv.ParseUint(valStr, 10, 32)
	if err != nil || v == 0 {
		Error(ctx, http.StatusBadRequest, "invalid param: "+name)
		return 0, false
	}
	return uint(v), true
}

// BindJSON 统一 JSON 绑定与错误处理，返回是否继续
func BindJSON(ctx *gin.Context, dst interface{}) bool {
	if err := ctx.ShouldBindJSON(dst); err != nil {
		Error(ctx, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}

// RequireAdmin 校验是否 admin，不是则返回 403，返回是否通过
func RequireAdmin(ctx *gin.Context) bool {
	roleCodeVal, exists := ctx.Get("roleCode")
	if !exists {
		Error(ctx, http.StatusForbidden, "forbidden")
		return false
	}
	roleCode, _ := roleCodeVal.(string)
	if roleCode != "admin" { // 与 model.RoleCodeAdmin 保持一致，避免循环依赖
		Error(ctx, http.StatusForbidden, "Only admin allowed")
		return false
	}
	return true
}
