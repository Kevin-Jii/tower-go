package http

import (
	"net/http"
	"reflect"
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

// GetPage 获取分页页码，默认为 1
func GetPage(ctx *gin.Context) int {
	pageStr := ctx.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 1
	}
	return page
}

// GetPageSize 获取分页大小，默认为 10
func GetPageSize(ctx *gin.Context) int {
	pageSizeStr := ctx.DefaultQuery("page_size", "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		return 10
	}
	return pageSize
}

// StructToMap 将结构体转换为 map[string]interface{}
func StructToMap(obj interface{}) map[string]interface{} {
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

		// 只处理可访问的字段
		if !value.CanInterface() {
			continue
		}

		// 获取 json 标签
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" || jsonTag == "" {
			jsonTag = field.Name
		}

		// 跳过忽略的字段
		if jsonTag == "-" {
			continue
		}

		// 处理嵌套结构体
		if value.Kind() == reflect.Struct {
			nested := StructToMap(value.Interface())
			for k, v := range nested {
				result[k] = v
			}
		} else {
			result[jsonTag] = value.Interface()
		}
	}

	return result
}
