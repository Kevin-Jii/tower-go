package bootstrap

import (
	"tower-go/config"
	"tower-go/utils"

	"go.uber.org/zap"
)

func InitDatabase() {
	dbConfig := config.GetDatabaseConfig()
	if err := utils.InitDB(dbConfig); err != nil {
		utils.LogFatal("数据库连接失败", zap.Error(err))
	}
	utils.LogInfo("数据库连接成功")
}
