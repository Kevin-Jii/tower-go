package performance

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"gorm.io/gorm"
)

// QueryOptimizer 查询优化器
type QueryOptimizer struct {
	db             *gorm.DB
	indexAnalyzer  *IndexAnalyzer
	joinDeduplicator *JoinDeduplicator
	mu             sync.RWMutex
}

// QueryAnalysisResult 查询分析结果
type QueryAnalysisResult struct {
	Query            string
	Issues           []QueryIssue
	Recommendations  []string
	EstimatedCost    float64
	IndexUsage       IndexUsage
	JoinCount        int
	WhereClauses     []string
}

// QueryIssue 查询问题
type QueryIssue struct {
	Type       IssueType
	Severity   Severity
	Message    string
	Location   string
	Suggestion string
}

// IssueType 问题类型
type IssueType string

const (
	IssueTypeMissingIndex     IssueType = "missing_index"
	IssueTypeNoWhereClause   IssueType = "no_where_clause"
	IssueTypeSelectAll       IssueType = "select_all"
	IssueTypeNPlusOne        IssueType = "n_plus_one"
	IssueTypeDuplicateJoin   IssueType = "duplicate_join"
	IssueTypeOffsetPagination IssueType = "offset_pagination"
	IssueTypeFullTableScan   IssueType = "full_table_scan"
	IssueTypeLikePrefix      IssueType = "like_prefix"
)

// Severity 严重程度
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh      Severity = "high"
	SeverityMedium    Severity = "medium"
	SeverityLow       Severity = "low"
)

// IndexUsage 索引使用情况
type IndexUsage struct {
	UsedIndexes      []string
	MissingIndexes   []string
	PotentialIndexes []string
	TableScan        bool
}

// JoinDeduplicator JOIN去重器
type JoinDeduplicator struct {
	joins map[string]bool
}

// NewJoinDeduplicator 创建JOIN去重器
func NewJoinDeduplicator() *JoinDeduplicator {
	return &JoinDeduplicator{joins: make(map[string]bool)}
}

// Add 添加JOIN
func (jd *JoinDeduplicator) Add(table string) bool {
	key := strings.ToLower(table)
	if jd.joins[key] {
		return false
	}
	jd.joins[key] = true
	return true
}

// Reset 重置
func (jd *JoinDeduplicator) Reset() {
	jd.joins = make(map[string]bool)
}

// IndexAnalyzer 索引分析器
type IndexAnalyzer struct {
	db *gorm.DB
}

// NewIndexAnalyzer 创建索引分析器
func NewIndexAnalyzer(db *gorm.DB) *IndexAnalyzer {
	return &IndexAnalyzer{db: db}
}

// AnalyzeIndexUsage 分析索引使用情况
func (ia *IndexAnalyzer) AnalyzeIndexUsage(table string, whereColumns []string) IndexUsage {
	result := IndexUsage{
		UsedIndexes:      []string{},
		MissingIndexes:   []string{},
		PotentialIndexes: []string{},
		TableScan:        false,
	}

	commonPatterns := map[string][]string{
		"store_id":    {"idx_store_id", "store_id"},
		"created_at":  {"idx_created_at", "created_at"},
		"account_date": {"idx_account_date", "account_date"},
		"product_id":  {"idx_product_id", "product_id"},
	}

	for _, col := range whereColumns {
		if patterns, ok := commonPatterns[col]; ok {
			result.PotentialIndexes = append(result.PotentialIndexes, patterns...)
		}
		if col == "" || col == "*" {
			result.TableScan = true
		}
	}

	return result
}

// NewQueryOptimizer 创建查询优化器
func NewQueryOptimizer(db *gorm.DB) *QueryOptimizer {
	return &QueryOptimizer{
		db:              db,
		indexAnalyzer:   NewIndexAnalyzer(db),
		joinDeduplicator: NewJoinDeduplicator(),
	}
}

// Analyze 分析SQL查询
func (qo *QueryOptimizer) Analyze(sql string) (*QueryAnalysisResult, error) {
	result := &QueryAnalysisResult{
		Query:           sql,
		Issues:          []QueryIssue{},
		Recommendations: []string{},
	}

	result.WhereClauses = qo.extractWhereClauses(sql)
	result.JoinCount = qo.countJoins(sql)

	qo.checkMissingIndex(sql, result)
	qo.checkOffsetPagination(sql, result)
	qo.checkFullTableScan(sql, result)
	qo.checkLikePrefix(sql, result)
	qo.checkDuplicateJoin(sql, result)
	qo.checkSelectAll(sql, result)

	qo.generateRecommendations(result)

	table := qo.extractTable(sql)
	if table != "" {
		result.IndexUsage = qo.indexAnalyzer.AnalyzeIndexUsage(table, result.WhereClauses)
	}

	return result, nil
}

// extractWhereClauses 提取WHERE条件
func (qo *QueryOptimizer) extractWhereClauses(sql string) []string {
	whereRe := regexp.MustCompile(`(?i)WHERE\s+(\w+)\s*=\s*\?`)
	matches := whereRe.FindAllStringSubmatch(sql, -1)

	var clauses []string
	for _, match := range matches {
		if len(match) > 1 {
			clauses = append(clauses, match[1])
		}
	}
	return clauses
}

// countJoins 统计JOIN数量
func (qo *QueryOptimizer) countJoins(sql string) int {
	joinRe := regexp.MustCompile(`(?i)\bJOIN\b`)
	return len(joinRe.FindAllStringIndex(sql, -1))
}

// extractTable 提取表名
func (qo *QueryOptimizer) extractTable(sql string) string {
	fromRe := regexp.MustCompile(`(?i)FROM\s+(\w+)`)
	matches := fromRe.FindAllStringSubmatch(sql, -1)
	if len(matches) > 0 && len(matches[0]) > 1 {
		return matches[0][1]
	}
	return ""
}

// checkMissingIndex 检查缺失索引
func (qo *QueryOptimizer) checkMissingIndex(sql string, result *QueryAnalysisResult) {
	whereColumns := result.WhereClauses
	if len(whereColumns) == 0 {
		return
	}

	unindexedFields := []string{"remark", "tag_name", "operator_name"}
	for _, col := range whereColumns {
		for _, unindexed := range unindexedFields {
			if col == unindexed {
				result.Issues = append(result.Issues, QueryIssue{
					Type:       IssueTypeMissingIndex,
					Severity:   SeverityMedium,
					Message:    fmt.Sprintf("字段 '%s' 可能缺少索引", col),
					Location:   "WHERE clause",
					Suggestion: fmt.Sprintf("考虑为字段 '%s' 添加索引", col),
				})
			}
		}
	}
}

// checkOffsetPagination 检查OFFSET分页
func (qo *QueryOptimizer) checkOffsetPagination(sql string, result *QueryAnalysisResult) {
	offsetRe := regexp.MustCompile(`(?i)OFFSET\s+(\d+)`)
	if offsetRe.MatchString(sql) {
		result.Issues = append(result.Issues, QueryIssue{
			Type:       IssueTypeOffsetPagination,
			Severity:   SeverityHigh,
			Message:    "使用OFFSET分页在大数据量时性能较差",
			Location:   "LIMIT/OFFSET",
			Suggestion: "考虑使用游标分页(CURSOR)替代OFFSET",
		})
	}
}

// checkFullTableScan 检查全表扫描
func (qo *QueryOptimizer) checkFullTableScan(sql string, result *QueryAnalysisResult) {
	whereRe := regexp.MustCompile(`(?i)\bWHERE\b`)
	if !whereRe.MatchString(sql) {
		result.Issues = append(result.Issues, QueryIssue{
			Type:       IssueTypeFullTableScan,
			Severity:   SeverityCritical,
			Message:    "查询没有WHERE条件，可能导致全表扫描",
			Location:   "Query",
			Suggestion: "添加WHERE条件或使用LIMIT限制结果集",
		})
		result.IndexUsage.TableScan = true
	}
}

// checkLikePrefix 检查LIKE前缀
func (qo *QueryOptimizer) checkLikePrefix(sql string, result *QueryAnalysisResult) {
	likeRe := regexp.MustCompile(`(?i)LIKE\s+'%([^%]+)'`)
	if likeRe.MatchString(sql) {
		result.Issues = append(result.Issues, QueryIssue{
			Type:       IssueTypeLikePrefix,
			Severity:   SeverityMedium,
			Message:    "LIKE查询以通配符开头，无法使用索引",
			Location:   "WHERE clause",
			Suggestion: "避免使用 '%xxx' 模式，使用 'xxx%' 前缀匹配",
		})
	}
}

// checkDuplicateJoin 检查重复JOIN
func (qo *QueryOptimizer) checkDuplicateJoin(sql string, result *QueryAnalysisResult) {
	qo.joinDeduplicator.Reset()

	joinTableRe := regexp.MustCompile(`(?i)JOIN\s+(\w+)`)
	matches := joinTableRe.FindAllStringSubmatch(sql, -1)

	for _, match := range matches {
		if len(match) > 1 {
			table := match[1]
			if !qo.joinDeduplicator.Add(table) {
				result.Issues = append(result.Issues, QueryIssue{
					Type:       IssueTypeDuplicateJoin,
					Severity:   SeverityMedium,
					Message:    fmt.Sprintf("表 '%s' 被JOIN了多次", table),
					Location:   "JOIN clause",
					Suggestion: "合并重复的JOIN",
				})
			}
		}
	}
}

// checkSelectAll 检查SELECT *
func (qo *QueryOptimizer) checkSelectAll(sql string, result *QueryAnalysisResult) {
	selectRe := regexp.MustCompile(`(?i)SELECT\s+\*`)
	if selectRe.MatchString(sql) {
		result.Issues = append(result.Issues, QueryIssue{
			Type:       IssueTypeSelectAll,
			Severity:   SeverityLow,
			Message:    "使用SELECT *会读取不必要的字段",
			Location:   "SELECT clause",
			Suggestion: "只选择需要的字段以减少数据传输",
		})
	}
}

// generateRecommendations 生成优化建议
func (qo *QueryOptimizer) generateRecommendations(result *QueryAnalysisResult) {
	if result.JoinCount > 3 {
		result.Recommendations = append(result.Recommendations, "查询包含较多JOIN，考虑拆分为多个简单查询")
	}

	if len(result.WhereClauses) > 5 {
		result.Recommendations = append(result.Recommendations, "WHERE条件较多，考虑使用复合索引")
	}

	if result.IndexUsage.TableScan {
		result.Recommendations = append(result.Recommendations, "添加适当的索引避免全表扫描")
	}

	for _, issue := range result.Issues {
		if issue.Suggestion != "" {
			result.Recommendations = append(result.Recommendations, issue.Suggestion)
		}
	}
}

// OptimizeQuery 优化查询
func (qo *QueryOptimizer) OptimizeQuery(sql string) (string, error) {
	result, err := qo.Analyze(sql)
	if err != nil {
		return sql, err
	}

	for _, issue := range result.Issues {
		if issue.Severity == SeverityCritical {
			return sql, fmt.Errorf("查询存在严重问题: %s", issue.Message)
		}
	}

	return sql, nil
}