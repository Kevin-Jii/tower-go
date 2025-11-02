// @title Tower Go API
// @version 1.0
// @description Tower Go 用户管理系统 API 文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

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
	"fmt"
	"log"
	"os"

	"tower-go/config"
	"tower-go/controller"
	_ "tower-go/docs"
	"tower-go/middleware"
	"tower-go/model"
	userModulePkg "tower-go/module"
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

	// 如果需要可以通过环境变量覆盖端口（如 PORT=10024）
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		// 简单转换
		var p int
		_, err := fmt.Sscanf(portEnv, "%d", &p)
		if err == nil && p > 0 {
			cfg := config.GetConfig()
			cfg.App.Port = p
		}
	}

	// 初始化数据库连接
	dbConfig := config.GetDatabaseConfig()
	if err := utils.InitDB(dbConfig); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	fmt.Printf("\n数据库连接成功！\n")

	// 外键前置数据完整性检查（users.store_id 但 stores 中缺失）
	var invalidUserCount int64
	utils.DB.Raw("SELECT COUNT(*) FROM users u LEFT JOIN stores s ON u.store_id = s.id WHERE u.store_id <> 0 AND s.id IS NULL").Scan(&invalidUserCount)
	if invalidUserCount > 0 {
		log.Printf("警告: 发现 %d 条用户记录的 store_id 无效，请修复 stores 表或重置这些用户的 store_id=0", invalidUserCount)
	}

	// 自动迁移数据表（按外键依赖顺序）
	// 顺序：Store -> Role -> Menu -> User -> Dish -> MenuReport -> RoleMenu -> StoreRoleMenu
	migrateModels := []interface{}{&model.Store{}, &model.Role{}, &model.Menu{}, &model.User{}, &model.Dish{}, &model.MenuReport{}, &model.RoleMenu{}, &model.StoreRoleMenu{}}
	for _, m := range migrateModels {
		if err := utils.DB.AutoMigrate(m); err != nil {
			log.Printf("AutoMigrate model %T failed: %v", m, err)
			log.Printf("迁移失败，后续种子数据将跳过。")
			goto SKIP_SEED
		}
	}
	log.Println("数据表迁移完成")

	// 初始化种子数据（仅在迁移成功后）
	if err := utils.InitRoleSeeds(utils.DB); err != nil {
		log.Printf("InitRoleSeeds failed: %v", err)
	} else {
		fmt.Println("角色基础数据初始化成功")
	}
	if err := utils.InitMenuSeeds(utils.DB); err != nil {
		log.Printf("InitMenuSeeds failed: %v", err)
	} else {
		fmt.Println("菜单种子数据初始化成功")
	}

	if err := utils.InitRoleMenuSeeds(utils.DB); err != nil {
		log.Printf("InitRoleMenuSeeds failed: %v", err)
	} else {
		fmt.Println("角色菜单权限初始化成功")
	}

SKIP_SEED:

	// 初始化模块/服务/控制器
	userModule := userModulePkg.NewUserModule(utils.DB)
	userService := service.NewUserService(userModule)
	userController := controller.NewUserController(userService)

	storeModule := userModulePkg.NewStoreModule(utils.DB)
	storeService := service.NewStoreService(storeModule)
	storeController := controller.NewStoreController(storeService)

	dishModule := userModulePkg.NewDishModule(utils.DB)
	dishService := service.NewDishService(dishModule)
	dishController := controller.NewDishController(dishService)

	menuReportModule := userModulePkg.NewMenuReportModule(utils.DB)
	menuReportService := service.NewMenuReportService(menuReportModule, dishModule)
	menuReportController := controller.NewMenuReportController(menuReportService)

	// 菜单权限模块
	menuModule := userModulePkg.NewMenuModule(utils.DB)
	roleMenuModule := userModulePkg.NewRoleMenuModule(utils.DB)
	storeRoleMenuModule := userModulePkg.NewStoreRoleMenuModule(utils.DB)
	menuService := service.NewMenuService(menuModule, roleMenuModule, storeRoleMenuModule)
	menuController := controller.NewMenuController(menuService)

	// 初始化 WebSocket 会话管理：策略可配置，这里先写死 single (单点登录)
	utils.InitSessionManager("single", 3)

	// 启动 HTTP 服务
	r := gin.Default()
	r.Use(middleware.RequestLoggerMiddleware(4096))
	v1 := r.Group("/api/v1")

	// 公开路由 (无需认证)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", userController.Register)
		auth.POST("/login", userController.Login)
	}

	// 需要认证和门店隔离的路由
	users := v1.Group("/users")
	// 强制门店认证和 StoreID 提取
	users.Use(middleware.StoreAuthMiddleware())
	{
		// 个人信息相关接口 (使用 UserID 访问自己的信息)
		users.GET("/profile", userController.GetProfile)
		users.PUT("/profile", userController.UpdateProfile)

		// 用户管理接口 (强制使用 StoreID 过滤数据)
		users.POST("", userController.CreateUser) // 新增创建用户路由
		users.GET("", userController.ListUsers)
		users.GET("/:id", userController.GetUser)
		users.PUT("/:id", userController.UpdateUser)
		users.DELETE("/:id", userController.DeleteUser)
		users.POST(":id/reset-password", userController.ResetUserPassword)
	}

	// 门店管理路由 (需要认证)
	stores := v1.Group("/stores")
	stores.Use(middleware.AuthMiddleware())
	{
		stores.POST("", storeController.CreateStore)       // 仅管理员
		stores.GET("", storeController.ListStores)         // 所有认证用户
		stores.GET("/:id", storeController.GetStore)       // 所有认证用户
		stores.PUT("/:id", storeController.UpdateStore)    // 仅管理员
		stores.DELETE("/:id", storeController.DeleteStore) // 仅管理员
	}

	// 菜品管理路由 (需要认证和门店隔离)
	dishes := v1.Group("/dishes")
	dishes.Use(middleware.StoreAuthMiddleware())
	{
		dishes.POST("", dishController.CreateDish)
		dishes.GET("", dishController.ListDishes)
		dishes.GET("/:id", dishController.GetDish)
		dishes.PUT("/:id", dishController.UpdateDish)
		dishes.DELETE("/:id", dishController.DeleteDish)
	}

	// 报菜管理路由 (需要认证和门店隔离)
	menuReports := v1.Group("/menu-reports")
	menuReports.Use(middleware.StoreAuthMiddleware())
	{
		menuReports.POST("", menuReportController.CreateMenuReport)
		menuReports.GET("", menuReportController.ListMenuReports)
		menuReports.GET("/statistics", menuReportController.GetStatistics) // 统计接口
		menuReports.GET("/:id", menuReportController.GetMenuReport)
		menuReports.PUT("/:id", menuReportController.UpdateMenuReport)
		menuReports.DELETE("/:id", menuReportController.DeleteMenuReport)
	}

	// 菜单权限管理路由 (需要认证)
	menus := v1.Group("/menus")
	menus.Use(middleware.AuthMiddleware())
	{
		// 菜单CRUD（仅总部管理员）
		menus.POST("", menuController.CreateMenu)
		menus.GET("", menuController.ListMenus)
		menus.GET("/tree", menuController.GetMenuTree)
		menus.GET("/:id", menuController.GetMenu)
		menus.PUT("/:id", menuController.UpdateMenu)
		menus.DELETE("/:id", menuController.DeleteMenu)

		// 角色菜单权限管理（仅总部管理员）
		menus.POST("/assign-role", menuController.AssignMenusToRole)
		menus.GET("/role", menuController.GetRoleMenus)
		menus.GET("/role-ids", menuController.GetRoleMenuIDs)

		// 门店角色菜单权限管理（总部管理员或门店管理员）
		menus.POST("/assign-store-role", menuController.AssignMenusToStoreRole)
		menus.GET("/store-role", menuController.GetStoreRoleMenus)
		menus.GET("/store-role-ids", menuController.GetStoreRoleMenuIDs)
		menus.POST("/copy-store", menuController.CopyStoreMenus)

		// 获取当前用户的菜单和权限（所有认证用户）
		menus.GET("/user-menus", menuController.GetUserMenus)
		menus.GET("/user-permissions", menuController.GetUserPermissions)
	}

	// WebSocket 路由 (基于 JWT 的连接认证)
	r.GET("/ws", controller.WebSocketHandler)

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
