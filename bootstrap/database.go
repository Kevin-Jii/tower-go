package bootstrap

import (
	"github.com/Kevin-Jii/tower-go/config"
	"github.com/Kevin-Jii/tower-go/utils/database"
	"github.com/Kevin-Jii/tower-go/utils/logging"

	"go.uber.org/zap"
)

func InitDatabase() {
	dbConfig := config.GetDatabaseConfig()
	if err := database.InitDB(dbConfig); err != nil {
		logging.LogFatal("数据库连接失败", zap.Error(err))
	}
	logging.LogInfo("数据库连接成功")
}
