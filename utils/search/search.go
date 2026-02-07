package search

import (
	"strings"
)

// OptimizeSearchKeyword 优化搜索关键字
// 将用户输入转换为数据库查询条件
func OptimizeSearchKeyword(keyword string) []string {
	if keyword == "" {
		return nil
	}

	// 移除首尾空格
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil
	}

	// 支持多个关键字（空格分隔）
	keywords := strings.Fields(keyword)
	if len(keywords) == 0 {
		return nil
	}

	return keywords
}

// BuildSearchSQL 构建搜索 SQL
// 返回构建好的 SQL 条件和参数
func BuildSearchSQL(conditions []string) (string, []interface{}) {
	if len(conditions) == 0 {
		return "", nil
	}

	// 构建 OR 条件，支持多字段搜索
	placeholders := make([]string, len(conditions))
	args := make([]interface{}, 0, len(conditions)*3)

	for i, keyword := range conditions {
		placeholders[i] = "(username LIKE ? OR phone LIKE ? OR nickname LIKE ?)"
		likeKeyword := "%" + keyword + "%"
		args = append(args, likeKeyword, likeKeyword, likeKeyword)
	}

	// 用 OR 连接所有条件
	sql := strings.Join(placeholders, " OR ")
	return sql, args
}

// BuildAdvancedSearchSQL 构建高级搜索 SQL
// 支持多个字段和条件的组合搜索
func BuildAdvancedSearchSQL(fieldConditions map[string][]string) (string, []interface{}) {
	if len(fieldConditions) == 0 {
		return "", nil
	}

	var conditions []string
	var args []interface{}

	for field, keywords := range fieldConditions {
		for _, keyword := range keywords {
			conditions = append(conditions, field+" LIKE ?")
			args = append(args, "%"+keyword+"%")
		}
	}

	if len(conditions) == 0 {
		return "", nil
	}

	return strings.Join(conditions, " AND "), args
}
