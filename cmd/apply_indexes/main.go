package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/bootstrap"
	"github.com/Kevin-Jii/tower-go/config"
	"github.com/Kevin-Jii/tower-go/utils/database"

	"gorm.io/gorm"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化数据库连接
	bootstrap.InitDatabase()

	// 获取数据库实例
	db := database.GetDB()
	if db == nil {
		log.Fatal("数据库连接失败")
	}

	fmt.Println("==============================================")
	fmt.Println("数据库索引优化工具")
	fmt.Println("==============================================")
	fmt.Println()

	// 读取SQL脚本
	sqlContent, err := os.ReadFile("migrations/add_performance_indexes.sql")
	if err != nil {
		log.Fatalf("读取SQL文件失败: %v", err)
	}

	// 解析SQL语句
	statements := parseSQL(string(sqlContent))

	fmt.Printf("共找到 %d 条索引创建语句\n\n", len(statements))

	// 执行每条SQL语句
	successCount := 0
	failCount := 0
	skippedCount := 0

	for i, stmt := range statements {
		indexName := extractIndexName(stmt)
		fmt.Printf("[%d/%d] 创建索引: %s\n", i+1, len(statements), indexName)

		// 检查索引是否已存在
		exists := indexExists(db, stmt)
		if exists {
			fmt.Printf("  ⏭️  索引已存在，跳过\n\n")
			skippedCount++
			continue
		}

		// 执行SQL
		startTime := time.Now()
		err := db.Exec(stmt).Error
		duration := time.Since(startTime)

		if err != nil {
			fmt.Printf("  ❌ 失败: %v\n\n", err)
			failCount++
		} else {
			fmt.Printf("  ✅ 成功 (耗时: %v)\n\n", duration)
			successCount++
		}

		// 短暂延迟，避免数据库压力过大
		time.Sleep(100 * time.Millisecond)
	}

	// 输出统计
	fmt.Println("==============================================")
	fmt.Println("执行结果统计:")
	fmt.Printf("  ✅ 成功: %d\n", successCount)
	fmt.Printf("  ❌ 失败: %d\n", failCount)
	fmt.Printf("  ⏭️  跳过: %d\n", skippedCount)
	fmt.Println("==============================================")

	if failCount > 0 {
		os.Exit(1)
	}
}

// parseSQL 解析SQL文件，提取CREATE INDEX语句
func parseSQL(content string) []string {
	var statements []string
	lines := strings.Split(content, "\n")

	var currentStmt strings.Builder
	inStatement := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 跳过注释和空行
		if trimmed == "" || strings.HasPrefix(trimmed, "--") {
			continue
		}

		// 检测CREATE INDEX语句开始
		if strings.HasPrefix(strings.ToUpper(trimmed), "CREATE") {
			inStatement = true
			currentStmt.Reset()
		}

		if inStatement {
			currentStmt.WriteString(line)
			currentStmt.WriteString("\n")

			// 检测语句结束（分号）
			if strings.HasSuffix(trimmed, ";") {
				stmt := strings.TrimSpace(currentStmt.String())
				if stmt != "" {
					statements = append(statements, stmt)
				}
				inStatement = false
			}
		}
	}

	return statements
}

// extractIndexName 从CREATE INDEX语句中提取索引名称
func extractIndexName(stmt string) string {
	// CREATE INDEX idx_name ON table_name(...)
	parts := strings.Fields(stmt)
	for i, part := range parts {
		if strings.ToUpper(part) == "INDEX" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return "unknown"
}

// indexExists 检查索引是否已存在
func indexExists(db *gorm.DB, stmt string) bool {
	indexName := extractIndexName(stmt)
	tableName := extractTableName(stmt)

	if tableName == "" {
		return false
	}

	var count int64
	query := `
		SELECT COUNT(*) 
		FROM information_schema.statistics 
		WHERE table_schema = DATABASE() 
		AND table_name = ? 
		AND index_name = ?
	`

	err := db.Raw(query, tableName, indexName).Scan(&count).Error
	if err != nil {
		return false
	}

	return count > 0
}

// extractTableName 从CREATE INDEX语句中提取表名
func extractTableName(stmt string) string {
	// CREATE INDEX idx_name ON table_name(...)
	parts := strings.Fields(stmt)
	for i, part := range parts {
		if strings.ToUpper(part) == "ON" && i+1 < len(parts) {
			tableName := parts[i+1]
			// 移除括号
			if idx := strings.Index(tableName, "("); idx != -1 {
				tableName = tableName[:idx]
			}
			return tableName
		}
	}
	return ""
}
