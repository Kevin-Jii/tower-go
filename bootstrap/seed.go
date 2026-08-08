package bootstrap

import (
	"os"
	"strings"

	"github.com/Kevin-Jii/tower-go/utils/database"
	"github.com/Kevin-Jii/tower-go/utils/logging"

	"go.uber.org/zap"
)

const seedDataVersion = "2"

// RunSeedSQL 执行种子数据 SQL 文件
func RunSeedSQL() {
	// 检查是否跳过（与 InitDefaultDicts、InitDefaultTemplates 一致，便于清空库后仅手工执行 SQL）
	if os.Getenv("SKIP_SEED_DATA") == "1" {
		logging.LogInfo("跳过种子数据 SQL（SKIP_SEED_DATA=1），不执行 migrations/init_seed_data.sql")
		return
	}
	markerPath := initializationMarkerPath(seedMarkerFile)
	if initializationComplete(markerPath, seedDataVersion) {
		logging.LogInfo("跳过种子数据 SQL（初始化版本已完成）", zap.String("version", seedDataVersion))
		return
	}

	// 读取 SQL 文件
	sqlFile := applicationFile("migrations/init_seed_data.sql")
	content, err := os.ReadFile(sqlFile)
	if err != nil {
		logging.LogWarn("读取种子数据文件失败", zap.String("file", sqlFile), zap.Error(err))
		return
	}

	// 按分号分割 SQL 语句
	statements := splitSQLStatements(string(content))
	successCount := 0
	skipCount := 0
	failureCount := 0

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
			failureCount++
			logging.LogWarn("执行种子数据语句失败", zap.Error(err), zap.String("stmt", stmt[:min(100, len(stmt))]))
		} else {
			successCount++
		}
	}

	logging.LogInfo("种子数据初始化完成",
		zap.Int("success", successCount),
		zap.Int("skipped", skipCount),
		zap.Int("failed", failureCount))
	if failureCount > 0 {
		logging.LogWarn("种子数据存在失败语句，不写入初始化标记，下次启动将重试")
		return
	}
	if err := markInitializationComplete(markerPath, seedDataVersion); err != nil {
		logging.LogWarn("无法写入种子数据初始化标记", zap.Error(err))
	}
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
