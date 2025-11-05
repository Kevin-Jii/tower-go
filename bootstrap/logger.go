package bootstrap

import (
	"fmt"
	"os"
	"tower-go/utils"
)

func InitLogger() func() {
	logConfig := &utils.LogConfig{
		Level:      "info",
		FilePath:   "logs/app.log",
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
		Console:    true,
	}
	if err := utils.InitLogger(logConfig); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	utils.LogInfo("日志系统初始化完成")
	return utils.CloseLogger
}
