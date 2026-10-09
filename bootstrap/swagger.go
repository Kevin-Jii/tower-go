package bootstrap

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/Kevin-Jii/tower-go/apidocs"
	"github.com/Kevin-Jii/tower-go/utils/logging"
	"go.uber.org/zap"
)

// GenerateSwaggerDocs 在应用启动时自动执行 swag init
// 可通过环境变量 SWAG_AUTO=0 禁用；SWAG_ARGS 追加自定义参数
func GenerateSwaggerDocs() {
	if v := os.Getenv("SWAG_AUTO"); v == "0" || strings.ToLower(v) == "false" {
		return
	}
	root := applicationRoot()
	args := []string{"init", "-d", root, "-g", "cmd/main.go"}
	if extra := os.Getenv("SWAG_ARGS"); extra != "" {
		// 简单拆分追加
		parts := strings.Fields(extra)
		args = append(args, parts...)
	}
	cmd := exec.Command("swag", args...)
	cmd.Dir = root
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		logging.LogWarn("自动生成 API 文档失败", zap.Error(err), zap.String("output", out.String()))
		return
	}
	if err := apidocs.ConvertFile(applicationFile("docs/swagger.json"), applicationFile("docs/openapi.json")); err != nil {
		logging.LogWarn("转换 OpenAPI 3 文档失败", zap.Error(err))
		return
	}
	logging.LogInfo("OpenAPI 3 文档已自动生成并校验")
}
