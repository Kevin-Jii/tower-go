package utils

import (
	"regexp"
	"strings"
)

// ApplyMultiTermFuzzy 将原始关键字按空格拆分为多个 term；每个 term 生成一组 (字段 LIKE %term% OR ... OR numericField = term[数字]) 条件，
// 多个 term 之间 AND 连接，实现“多词 AND + 多字段 OR + 数字精确匹配”。
// fields: 需要 OR 模糊匹配的字段列表
// numericField: 若 term 是纯数字，额外增加 numericField = term 精确匹配（可为空字符串禁用）
// raw: 原始搜索关键字
// 返回：构造完条件的 QueryBuilder（可继续追加其它条件）
func ApplyMultiTermFuzzy(qb *QueryBuilder, fields []string, raw string, numericField string) *QueryBuilder {
	raw = strings.TrimSpace(raw)
	if raw == "" || len(fields) == 0 {
		return qb
	}
	terms := strings.Fields(raw)
	for _, term := range terms {
		term = strings.TrimSpace(term)
		if term == "" {
			continue
		}
		likeConds := make([]string, 0, len(fields)+1)
		args := make([]interface{}, 0, len(fields)+1)
		for _, f := range fields {
			likeConds = append(likeConds, f+" LIKE ?")
			args = append(args, "%"+term+"%")
		}
		if numericField != "" && isDigits(term) {
			likeConds = append(likeConds, numericField+" = ?")
			args = append(args, term)
		}
		// 单个 term 的 OR 组整体作为一个 AND 条件加入
		qb.Where("("+strings.Join(likeConds, " OR ")+")", args...)
	}
	return qb
}

func isDigits(s string) bool {
	matched, _ := regexp.MatchString(`^\d+$`, s)
	return matched
}
