package search

import (
	"regexp"
	"strings"
)

// SearchType 搜索类型
type SearchType int

const (
	SearchTypeExact  SearchType = iota // 精确匹配
	SearchTypePrefix                   // 前缀匹配（username LIKE 'keyword%'）
	SearchTypeSuffix                   // 后缀匹配（phone LIKE '%keyword'）
	SearchTypeFull                     // 全文匹配（LIKE '%keyword%'）
)

// SearchCondition 搜索条件
type SearchCondition struct {
	Field      string
	Value      string
	SearchType SearchType
}

// OptimizeSearchKeyword 智能优化搜索关键字
// 根据关键字特征自动选择最优查询方式
func OptimizeSearchKeyword(keyword string) []SearchCondition {
	if keyword == "" {
		return []SearchCondition{}
	}

	keyword = strings.TrimSpace(keyword)
	conditions := []SearchCondition{}

	// 判断是否是纯数字（可能是手机号）
	if isNumeric(keyword) {
		// 数字关键字：优先用于手机号后缀匹配
		if len(keyword) >= 4 {
			// 手机号后缀匹配（如输入 "1234" 匹配 "***1234"）
			conditions = append(conditions, SearchCondition{
				Field:      "phone",
				Value:      keyword,
				SearchType: SearchTypeSuffix,
			})
		}
		// 同时支持精确ID查询（如果是短数字）
		if len(keyword) <= 10 {
			conditions = append(conditions, SearchCondition{
				Field:      "id",
				Value:      keyword,
				SearchType: SearchTypeExact,
			})
		}
	}

	// 中文或字母：用于用户名前缀匹配
	if containsChineseOrLetter(keyword) {
		conditions = append(conditions, SearchCondition{
			Field:      "username",
			Value:      keyword,
			SearchType: SearchTypePrefix,
		})
	}

	// 如果包含 @，可能是邮箱
	if strings.Contains(keyword, "@") {
		conditions = append(conditions, SearchCondition{
			Field:      "email",
			Value:      keyword,
			SearchType: SearchTypePrefix,
		})
	}

	// 如果没有匹配任何优化规则，使用全文搜索（最慢但最全面）
	if len(conditions) == 0 {
		conditions = append(conditions, SearchCondition{
			Field:      "username",
			Value:      keyword,
			SearchType: SearchTypeFull,
		})
		conditions = append(conditions, SearchCondition{
			Field:      "phone",
			Value:      keyword,
			SearchType: SearchTypeFull,
		})
	}

	return conditions
}

// BuildSearchSQL 构建搜索 SQL 条件
func BuildSearchSQL(conditions []SearchCondition) (string, []interface{}) {
	if len(conditions) == 0 {
		return "", nil
	}

	var sqlParts []string
	var args []interface{}

	for _, cond := range conditions {
		switch cond.SearchType {
		case SearchTypeExact:
			sqlParts = append(sqlParts, cond.Field+" = ?")
			args = append(args, cond.Value)
		case SearchTypePrefix:
			sqlParts = append(sqlParts, cond.Field+" LIKE ?")
			args = append(args, cond.Value+"%")
		case SearchTypeSuffix:
			sqlParts = append(sqlParts, cond.Field+" LIKE ?")
			args = append(args, "%"+cond.Value)
		case SearchTypeFull:
			sqlParts = append(sqlParts, cond.Field+" LIKE ?")
			args = append(args, "%"+cond.Value+"%")
		}
	}

	sql := strings.Join(sqlParts, " OR ")
	return sql, args
}

// isNumeric 判断字符串是否为纯数字
func isNumeric(s string) bool {
	matched, _ := regexp.MatchString(`^\d+$`, s)
	return matched
}

// containsChineseOrLetter 判断是否包含中文或字母
func containsChineseOrLetter(s string) bool {
	// 匹配中文字符或字母
	matched, _ := regexp.MatchString(`[\p{Han}a-zA-Z]`, s)
	return matched
}

// OptimizePhoneSearch 优化手机号搜索
// 手机号通常查后4位或后6位，使用后缀匹配性能更好
func OptimizePhoneSearch(keyword string) string {
	keyword = strings.TrimSpace(keyword)
	// 如果是纯数字且长度合理，使用后缀匹配
	if isNumeric(keyword) && len(keyword) >= 4 && len(keyword) <= 11 {
		return "%" + keyword // 后缀匹配：phone LIKE '%1234'
	}
	// 否则使用全文匹配
	return "%" + keyword + "%"
}

// OptimizeUsernameSearch 优化用户名搜索
// 用户名通常从开头输入，使用前缀匹配性能更好
func OptimizeUsernameSearch(keyword string) string {
	keyword = strings.TrimSpace(keyword)
	// 如果包含中文或字母，使用前缀匹配
	if containsChineseOrLetter(keyword) {
		return keyword + "%" // 前缀匹配：username LIKE '张%'
	}
	// 否则使用全文匹配
	return "%" + keyword + "%"
}
