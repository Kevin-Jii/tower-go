package bootstrap

import (
	"tower-go/config"
	"tower-go/utils/database"
	"tower-go/utils/logging"

	"go.uber.org/zap"
)

func InitDatabase() {
	dbConfig := config.GetDatabaseConfig()
	if err := database.InitDB(dbConfig); err != nil {
		logging.LogFatal("数据库连接失败", zap.Error(err))
	}
	logging.LogInfo("数据库连接成功")
}
