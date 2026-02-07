package bootstrap

import (
	"os"
	"strings"

	"github.com/Kevin-Jii/tower-go/utils/database"
	"github.com/Kevin-Jii/tower-go/utils/logging"

	"go.uber.org/zap"
)

// RunSeedSQL 执行种子数据 SQL 文件
func RunSeedSQL() {
	// 检查是否跳过
	if os.Getenv("SKIP_SEED_DATA") == "1" {
		logging.LogInfo("跳过种子数据初始化（SKIP_SEED_DATA=1）")
		return
	}

	// 检查数据库是否有数据（用户表为空则需要初始化）
	var userCount int64
	database.GetDB().Raw("SELECT COUNT(*) FROM users").Scan(&userCount)
	
	seedMarkerFile := ".seed_executed"
	if userCount > 0 {
		// 数据库有数据，确保标记文件存在
		os.WriteFile(seedMarkerFile, []byte("executed"), 0644)
		logging.LogInfo("数据库已有数据，跳过种子数据初始化")
		return
	}
	
	// 数据库为空，删除旧标记文件（如果存在）
	os.Remove(seedMarkerFile)

	// 读取 SQL 文件
	sqlFile := "migrations/init_seed_data.sql"
	content, err := os.ReadFile(sqlFile)
	if err != nil {
		logging.LogWarn("读取种子数据文件失败", zap.String("file", sqlFile), zap.Error(err))
		return
	}

	// 按分号分割 SQL 语句
	statements := splitSQLStatements(string(content))
	successCount := 0
	skipCount := 0

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}

		if err := database.GetDB().Exec(stmt).Error; err != nil {
			// 忽略重复键错误，继续执行
			if strings.Contains(err.Error(), "Duplicate") {
				skipCount++
				continue
			}
			logging.LogWarn("执行种子数据语句失败", zap.Error(err), zap.String("stmt", stmt[:min(100, len(stmt))]))
		} else {
			successCount++
		}
	}

	logging.LogInfo("种子数据初始化完成",
		zap.Int("success", successCount),
		zap.Int("skipped", skipCount))

	// 创建标记文件
	os.WriteFile(seedMarkerFile, []byte("executed"), 0644)
}

// splitSQLStatements 分割 SQL 语句（处理多行语句）
func splitSQLStatements(content string) []string {
	var statements []string
	var current strings.Builder
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// 跳过注释行
		if strings.HasPrefix(trimmed, "--") {
			continue
		}
		current.WriteString(line)
		current.WriteString("\n")

		// 检查是否以分号结尾
		if strings.HasSuffix(trimmed, ";") {
			statements = append(statements, current.String())
			current.Reset()
		}
	}

	// 处理最后一条没有分号的语句
	if current.Len() > 0 {
		statements = append(statements, current.String())
	}

	return statements
}


func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
