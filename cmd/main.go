// @title Tower Go API
// @version 1.0
// @description Tower Go 用户管理系统 API 文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email jeckwell@126.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:10024
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"tower-go/bootstrap"
	_ "tower-go/docs"
)

func main() { bootstrap.Run() }
