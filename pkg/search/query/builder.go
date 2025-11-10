package query

import (
	"fmt"
	"strings"
)

// Builder 搜索查询构建器
type Builder struct {
	conditions []string
	args       []interface{}
	orderBy    []string
	limit      *int
	offset     *int
}

// New 创建新的查询构建器
func New() *Builder {
	return &Builder{
		conditions: make([]string, 0),
		args:       make([]interface{}, 0),
		orderBy:    make([]string, 0),
	}
}

// Where 添加WHERE条件
func (b *Builder) Where(condition string, args ...interface{}) *Builder {
	if condition == "" {
		return b
	}
	b.conditions = append(b.conditions, condition)
	b.args = append(b.args, args...)
	return b
}

// WhereOr 添加OR条件组
func (b *Builder) WhereOr(conditions []string, args ...interface{}) *Builder {
	if len(conditions) == 0 {
		return b
	}

	orCondition := "(" + strings.Join(conditions, " OR ") + ")"
	b.conditions = append(b.conditions, orCondition)
	b.args = append(b.args, args...)
	return b
}

// WhereIn 添加IN条件
func (b *Builder) WhereIn(field string, values []interface{}) *Builder {
	if len(values) == 0 {
		return b
	}

	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
	}

	condition := fmt.Sprintf("%s IN (%s)", field, strings.Join(placeholders, ","))
	b.conditions = append(b.conditions, condition)
	b.args = append(b.args, values...)
	return b
}

// WhereNotIn 添加NOT IN条件
func (b *Builder) WhereNotIn(field string, values []interface{}) *Builder {
	if len(values) == 0 {
		return b
	}

	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
	}

	condition := fmt.Sprintf("%s NOT IN (%s)", field, strings.Join(placeholders, ","))
	b.conditions = append(b.conditions, condition)
	b.args = append(b.args, values...)
	return b
}

// WhereLike 添加LIKE条件
func (b *Builder) WhereLike(field string, pattern string) *Builder {
	condition := fmt.Sprintf("%s LIKE ?", field)
	b.conditions = append(b.conditions, condition)
	b.args = append(b.args, pattern)
	return b
}

// WhereILike 添加不区分大小写的LIKE条件（MySQL的LIKE默认不区分大小写）
func (b *Builder) WhereILike(field string, pattern string) *Builder {
	condition := fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", field)
	b.conditions = append(b.conditions, condition)
	b.args = append(b.args, pattern)
	return b
}

// WhereBetween 添加BETWEEN条件
func (b *Builder) WhereBetween(field string, start, end interface{}) *Builder {
	condition := fmt.Sprintf("%s BETWEEN ? AND ?", field)
	b.conditions = append(b.conditions, condition)
	b.args = append(b.args, start, end)
	return b
}

// WhereNull 添加IS NULL条件
func (b *Builder) WhereNull(field string) *Builder {
	condition := fmt.Sprintf("%s IS NULL", field)
	b.conditions = append(b.conditions, condition)
	return b
}

// WhereNotNull 添加IS NOT NULL条件
func (b *Builder) WhereNotNull(field string) *Builder {
	condition := fmt.Sprintf("%s IS NOT NULL", field)
	b.conditions = append(b.conditions, condition)
	return b
}

// WhereDate 添加日期条件
func (b *Builder) WhereDate(field string, date string) *Builder {
	condition := fmt.Sprintf("DATE(%s) = ?", field)
	b.conditions = append(b.conditions, condition)
	b.args = append(b.args, date)
	return b
}

// WhereDateBetween 添加日期范围条件
func (b *Builder) WhereDateBetween(field string, startDate, endDate string) *Builder {
	condition := fmt.Sprintf("DATE(%s) BETWEEN ? AND ?", field)
	b.conditions = append(b.conditions, condition)
	b.args = append(b.args, startDate, endDate)
	return b
}

// OrderBy 添加排序
func (b *Builder) OrderBy(field string, direction ...string) *Builder {
	dir := "ASC"
	if len(direction) > 0 {
		dir = strings.ToUpper(direction[0])
		if dir != "ASC" && dir != "DESC" {
			dir = "ASC"
		}
	}

	b.orderBy = append(b.orderBy, fmt.Sprintf("%s %s", field, dir))
	return b
}

// OrderByDesc 添加降序排序
func (b *Builder) OrderByDesc(field string) *Builder {
	return b.OrderBy(field, "DESC")
}

// Limit 添加限制
func (b *Builder) Limit(limit int) *Builder {
	b.limit = &limit
	return b
}

// Offset 添加偏移量
func (b *Builder) Offset(offset int) *Builder {
	b.offset = &offset
	return b
}

// Build 构建SQL查询和参数
func (b *Builder) Build() (string, []interface{}) {
	var sqlParts []string

	// 构建WHERE子句
	if len(b.conditions) > 0 {
		sqlParts = append(sqlParts, "WHERE "+strings.Join(b.conditions, " AND "))
	}

	// 构建ORDER BY子句
	if len(b.orderBy) > 0 {
		sqlParts = append(sqlParts, "ORDER BY "+strings.Join(b.orderBy, ", "))
	}

	// 构建LIMIT和OFFSET子句
	if b.limit != nil {
		sqlParts = append(sqlParts, fmt.Sprintf("LIMIT %d", *b.limit))
		if b.offset != nil {
			sqlParts = append(sqlParts, fmt.Sprintf("OFFSET %d", *b.offset))
		}
	}

	sql := strings.Join(sqlParts, " ")
	return sql, b.args
}

// BuildCount 构建计数查询（排除LIMIT和OFFSET）
func (b *Builder) BuildCount() (string, []interface{}) {
	var sqlParts []string

	// 构建WHERE子句
	if len(b.conditions) > 0 {
		sqlParts = append(sqlParts, "WHERE "+strings.Join(b.conditions, " AND "))
	}

	sql := strings.Join(sqlParts, " ")
	return sql, b.args
}

// Clone 克隆构建器
func (b *Builder) Clone() *Builder {
	clone := New()
	clone.conditions = make([]string, len(b.conditions))
	copy(clone.conditions, b.conditions)
	clone.args = make([]interface{}, len(b.args))
	copy(clone.args, b.args)
	clone.orderBy = make([]string, len(b.orderBy))
	copy(clone.orderBy, b.orderBy)

	if b.limit != nil {
		limit := *b.limit
		clone.limit = &limit
	}

	if b.offset != nil {
		offset := *b.offset
		clone.offset = &offset
	}

	return clone
}

// Reset 重置构建器
func (b *Builder) Reset() *Builder {
	b.conditions = make([]string, 0)
	b.args = make([]interface{}, 0)
	b.orderBy = make([]string, 0)
	b.limit = nil
	b.offset = nil
	return b
}

// HasConditions 检查是否有条件
func (b *Builder) HasConditions() bool {
	return len(b.conditions) > 0
}

// GetArgs 获取参数
func (b *Builder) GetArgs() []interface{} {
	return b.args
}

// GetConditions 获取条件
func (b *Builder) GetConditions() []string {
	return b.conditions
}