// @title Tower Go API
// @version 1.0
// @description Tower Go 用户管理系统 API 文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"fmt"
	"log"

	"tower-go/config"
	"tower-go/controller"
	_ "tower-go/docs"
	"tower-go/middleware"
	"tower-go/model"
	"tower-go/module"
	"tower-go/service"
	"tower-go/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// 加载配置文件
	if err := config.LoadConfig("config/config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库连接
	dbConfig := config.GetDatabaseConfig()
	if err := utils.InitDB(dbConfig); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	fmt.Printf("\n数据库连接成功！\n")

	// 自动迁移 user 表
	if err := utils.DB.AutoMigrate(&model.User{}); err != nil {
		log.Printf("AutoMigrate user failed: %v", err)
	}

	// 初始化模块/服务/控制器
	userModule := module.NewUserModule(utils.DB)
	userService := service.NewUserService(userModule)
	userController := controller.NewUserController(userService)

	// 启动 HTTP 服务
	r := gin.Default()
	v1 := r.Group("/api/v1")

	// 公开路由
	auth := v1.Group("/auth")
	{
		auth.POST("/register", userController.Register)
		auth.POST("/login", userController.Login)
	}

	// 需要认证的路由
	users := v1.Group("/users")
	users.Use(middleware.AuthMiddleware())
	{
		// 个人信息相关接口
		users.GET("/profile", userController.GetProfile)
		users.PUT("/profile", userController.UpdateProfile)

		// 用户管理接口
		users.GET("", userController.ListUsers)
		users.GET("/:id", userController.GetUser)
		users.PUT("/:id", userController.UpdateUser)
		users.DELETE("/:id", userController.DeleteUser)
	}

	// swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf(":%d", config.GetConfig().App.Port)
	fmt.Printf("starting server at %s\n", addr)

	// 打印 Swagger UI 路径到控制台
	swaggerURL := fmt.Sprintf("http://localhost%s/swagger/index.html", addr)
	fmt.Printf("Swagger UI: %s\n", swaggerURL)

	if err := r.Run(addr); err != nil {
		log.Fatalf("server exit: %v", err)
	}
}
