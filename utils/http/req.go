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

// ParseUintQuery 解析 query 参数为 uint，返回 (value, ok)
func ParseUintQuery(ctx *gin.Context, name string) (uint, bool) {
	valStr := ctx.Query(name)
	if valStr == "" {
		return 0, false
	}
	v, err := strconv.ParseUint(valStr, 10, 32)
	if err != nil {
		return 0, false
	}
	return uint(v), true
}

// ParseUint 解析字符串为 uint
func ParseUint(s string, result *uint) bool {
	if s == "" {
		return false
	}
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return false
	}
	*result = uint(v)
	return true
}

// BindJSON 统一 JSON 绑定与错误处理，返回是否继续
func BindJSON(ctx *gin.Context, dst interface{}) bool {
	if err := ctx.ShouldBindJSON(dst); err != nil {
		Error(ctx, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}

// RequireAdmin 校验是否管理员，不是则返回 403，返回是否通过
func RequireAdmin(ctx *gin.Context) bool {
	roleCodeVal, exists := ctx.Get("roleCode")
	if !exists {
		Error(ctx, http.StatusForbidden, "forbidden")
		return false
	}
	roleCode, _ := roleCodeVal.(string)
	// 允许多种管理员角色
	if roleCode != "admin" && roleCode != "super_admin" && roleCode != "store_admin" {
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
