package bootstrap

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"tower-go/utils/logging"

	"go.uber.org/zap"
)

// GenerateSwaggerDocs 在应用启动时自动执行 swag init
// 可通过环境变量 SWAG_AUTO=0 禁用；SWAG_ARGS 追加自定义参数
func GenerateSwaggerDocs() {
	if v := os.Getenv("SWAG_AUTO"); v == "0" || strings.ToLower(v) == "false" {
		return
	}
	args := []string{"init", "-g", "cmd/main.go"}
	if extra := os.Getenv("SWAG_ARGS"); extra != "" {
		// 简单拆分追加
		parts := strings.Fields(extra)
		args = append(args, parts...)
	}
	cmd := exec.Command("swag", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		logging.LogWarn("自动生成 swagger 失败", zap.Error(err))
		return
	}
	logging.LogInfo("swagger 文档已自动生成")
}
