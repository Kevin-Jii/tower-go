package bootstrap

import (
	"fmt"
	"os"
	"github.com/Kevin-Jii/tower-go/utils/logging"
)

func InitLogger() func() {
	logConfig := &logging.LogConfig{
		Level:      "info",
		FilePath:   "logs/app.log",
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
		Console:    true,
	}
	if err := logging.InitLogger(logConfig); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	logging.LogInfo("日志系统初始化完成")
	return logging.CloseLogger
}
