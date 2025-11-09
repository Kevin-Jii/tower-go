package database

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// QueryBuilder 查询构造器
type QueryBuilder struct {
	db         *gorm.DB
	conditions []string
	args       []interface{}
	orders     []string
	limit      int
	offset     int
}

// NewQueryBuilder 创建查询构造器
func NewQueryBuilder(db *gorm.DB) *QueryBuilder {
	return &QueryBuilder{
		db:         db,
		conditions: make([]string, 0),
		args:       make([]interface{}, 0),
		orders:     make([]string, 0),
	}
}

// ForStore 统一门店过滤
func (qb *QueryBuilder) ForStore(storeID uint) *QueryBuilder {
	if storeID > 0 {
		qb.conditions = append(qb.conditions, "store_id = ?")
		qb.args = append(qb.args, storeID)
	}
	return qb
}

// WhereStatusEnabled 简化启用状态过滤（适用于 status=1 约定）
func (qb *QueryBuilder) WhereStatusEnabled() *QueryBuilder {
	qb.conditions = append(qb.conditions, "status = 1")
	return qb
}

// Where 添加 WHERE 条件
func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	qb.conditions = append(qb.conditions, condition)
	qb.args = append(qb.args, args...)
	return qb
}

// WhereIf 条件性添加 WHERE（当 condition 为 true 时）
func (qb *QueryBuilder) WhereIf(condition bool, query string, args ...interface{}) *QueryBuilder {
	if condition {
		return qb.Where(query, args...)
	}
	return qb
}

// WhereIn 添加 IN 条件
func (qb *QueryBuilder) WhereIn(column string, values interface{}) *QueryBuilder {
	if values == nil {
		return qb
	}
	qb.conditions = append(qb.conditions, fmt.Sprintf("%s IN ?", column))
	qb.args = append(qb.args, values)
	return qb
}

// WhereLike 添加 LIKE 条件（模糊查询）
func (qb *QueryBuilder) WhereLike(column string, value string) *QueryBuilder {
	if value == "" {
		return qb
	}
	qb.conditions = append(qb.conditions, fmt.Sprintf("%s LIKE ?", column))
	qb.args = append(qb.args, "%"+value+"%")
	return qb
}

// WhereMultiLike 多字段 LIKE 查询（OR 连接）
func (qb *QueryBuilder) WhereMultiLike(columns []string, value string) *QueryBuilder {
	if value == "" || len(columns) == 0 {
		return qb
	}

	var conditions []string
	for _, col := range columns {
		conditions = append(conditions, fmt.Sprintf("%s LIKE ?", col))
		qb.args = append(qb.args, "%"+value+"%")
	}

	qb.conditions = append(qb.conditions, fmt.Sprintf("(%s)", strings.Join(conditions, " OR ")))
	return qb
}

// WhereBetween 添加 BETWEEN 条件
func (qb *QueryBuilder) WhereBetween(column string, start, end interface{}) *QueryBuilder {
	qb.conditions = append(qb.conditions, fmt.Sprintf("%s BETWEEN ? AND ?", column))
	qb.args = append(qb.args, start, end)
	return qb
}

// WhereNull 添加 IS NULL 条件
func (qb *QueryBuilder) WhereNull(column string) *QueryBuilder {
	qb.conditions = append(qb.conditions, fmt.Sprintf("%s IS NULL", column))
	return qb
}

// WhereNotNull 添加 IS NOT NULL 条件
func (qb *QueryBuilder) WhereNotNull(column string) *QueryBuilder {
	qb.conditions = append(qb.conditions, fmt.Sprintf("%s IS NOT NULL", column))
	return qb
}

// OrderBy 添加排序
func (qb *QueryBuilder) OrderBy(column string, direction ...string) *QueryBuilder {
	dir := "ASC"
	if len(direction) > 0 {
		dir = strings.ToUpper(direction[0])
	}
	qb.orders = append(qb.orders, fmt.Sprintf("%s %s", column, dir))
	return qb
}

// OrderByDesc 降序排序
func (qb *QueryBuilder) OrderByDesc(column string) *QueryBuilder {
	return qb.OrderBy(column, "DESC")
}

// Limit 设置查询数量限制
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

// Offset 设置查询偏移量
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

// Page 设置分页（page 从 1 开始）
func (qb *QueryBuilder) Page(page, pageSize int) *QueryBuilder {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	qb.limit = pageSize
	qb.offset = (page - 1) * pageSize
	return qb
}

// Build 构建最终查询
func (qb *QueryBuilder) Build() *gorm.DB {
	query := qb.db

	// 添加 WHERE 条件
	if len(qb.conditions) > 0 {
		whereClause := strings.Join(qb.conditions, " AND ")
		query = query.Where(whereClause, qb.args...)
	}

	// 添加排序
	if len(qb.orders) > 0 {
		for _, order := range qb.orders {
			query = query.Order(order)
		}
	}

	// 添加 LIMIT
	if qb.limit > 0 {
		query = query.Limit(qb.limit)
	}

	// 添加 OFFSET
	if qb.offset > 0 {
		query = query.Offset(qb.offset)
	}

	return query
}

// Find 执行查询并返回结果
func (qb *QueryBuilder) Find(dest interface{}) error {
	return qb.Build().Find(dest).Error
}

// First 查询第一条记录
func (qb *QueryBuilder) First(dest interface{}) error {
	return qb.Build().First(dest).Error
}

// Count 统计记录数
func (qb *QueryBuilder) Count(count *int64) error {
	query := qb.db
	if len(qb.conditions) > 0 {
		whereClause := strings.Join(qb.conditions, " AND ")
		query = query.Where(whereClause, qb.args...)
	}
	return query.Count(count).Error
}

// Paginate 分页查询（返回数据和总数）
func (qb *QueryBuilder) Paginate(dest interface{}, page, pageSize int) (int64, error) {
	var total int64

	// 先统计总数（不包含 limit/offset）
	if err := qb.Count(&total); err != nil {
		return 0, err
	}

	// 如果没有数据，直接返回
	if total == 0 {
		return 0, nil
	}

	// 执行分页查询
	qb.Page(page, pageSize)
	if err := qb.Find(dest); err != nil {
		return 0, err
	}

	return total, nil
}

// PaginateResult 分页查询结果
type PaginateResult struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// PaginateWithResult 分页查询并返回结构化结果
func (qb *QueryBuilder) PaginateWithResult(dest interface{}, page, pageSize int) (*PaginateResult, error) {
	total, err := qb.Paginate(dest, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &PaginateResult{
		Data:       dest,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// Clone 克隆查询构造器（用于复用条件）
func (qb *QueryBuilder) Clone() *QueryBuilder {
	return &QueryBuilder{
		db:         qb.db,
		conditions: append([]string{}, qb.conditions...),
		args:       append([]interface{}{}, qb.args...),
		orders:     append([]string{}, qb.orders...),
		limit:      qb.limit,
		offset:     qb.offset,
	}
}

// Reset 重置查询构造器
func (qb *QueryBuilder) Reset() *QueryBuilder {
	qb.conditions = make([]string, 0)
	qb.args = make([]interface{}, 0)
	qb.orders = make([]string, 0)
	qb.limit = 0
	qb.offset = 0
	return qb
}
