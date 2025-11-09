package bootstrap

import (
	"fmt"
	"tower-go/config"
	"tower-go/controller"
	"tower-go/middleware"
	userModulePkg "tower-go/module"
	"tower-go/service"
	"tower-go/utils/database"
	"tower-go/utils/events"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type AppControllers struct {
	User              *controller.UserController
	Store             *controller.StoreController
	Dish              *controller.DishController
	DishCategory      *controller.DishCategoryController
	MenuReport        *controller.MenuReportController
	Menu              *controller.MenuController
	DingTalkBot       *controller.DingTalkBotController
	DingTalkBotModule *userModulePkg.DingTalkBotModule // 用于 Stream 初始化
}

func BuildControllers() *AppControllers {
	// 初始化模块层
	userModule := userModulePkg.NewUserModule(database.DB)
	storeModule := userModulePkg.NewStoreModule(database.DB)
	dishModule := userModulePkg.NewDishModule(database.DB)
	menuReportModule := userModulePkg.NewMenuReportModule(database.DB)
	dishCategoryModule := userModulePkg.NewDishCategoryModule(database.DB)
	menuModule := userModulePkg.NewMenuModule(database.DB)
	roleMenuModule := userModulePkg.NewRoleMenuModule(database.DB)
	storeRoleMenuModule := userModulePkg.NewStoreRoleMenuModule(database.DB)
	dingTalkBotModule := userModulePkg.NewDingTalkBotModule(database.DB)

	// 初始化角色模块全局 DB（旧实现依赖 SetDB）
	userModulePkg.SetDB(database.DB)

	// 初始化事件总线
	eventBus := events.GetEventBus()

	// 初始化服务层
	userService := service.NewUserService(userModule)
	storeService := service.NewStoreService(storeModule)
	dishService := service.NewDishService(dishModule)
	dingTalkService := service.NewDingTalkService(dingTalkBotModule)
	menuReportService := service.NewMenuReportService(
		menuReportModule,
		dishModule,
		storeModule,
		userModule,
		eventBus,
	)
	dishCategoryService := service.NewDishCategoryService(dishCategoryModule, dishModule)
	menuService := service.NewMenuService(menuModule, roleMenuModule, storeRoleMenuModule)

	// 注册事件监听器
	menuReportListener := service.NewMenuReportEventListener(dingTalkService)
	service.RegisterMenuReportEventListeners(eventBus, menuReportListener)

	return &AppControllers{
		User:              controller.NewUserController(userService),
		Store:             controller.NewStoreController(storeService),
		Dish:              controller.NewDishController(dishService),
		DishCategory:      controller.NewDishCategoryController(dishCategoryService),
		MenuReport:        controller.NewMenuReportController(menuReportService),
		Menu:              controller.NewMenuController(menuService),
		DingTalkBot:       controller.NewDingTalkBotController(dingTalkService),
		DingTalkBotModule: dingTalkBotModule, // 暴露给 Stream 初始化使用
	}
}

// CRUDController 通用 CRUD 适配接口
type CRUDController interface {
	Create(*gin.Context)
	List(*gin.Context)
	Get(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}

// dishCategoryCRUDAdapter 适配已有 DishCategoryController 命名到标准 CRUD 接口
type dishCategoryCRUDAdapter struct {
	inner *controller.DishCategoryController
}

func (a dishCategoryCRUDAdapter) Create(ctx *gin.Context) { a.inner.CreateCategory(ctx) }
func (a dishCategoryCRUDAdapter) List(ctx *gin.Context)   { a.inner.ListCategories(ctx) }

// 暂无单分类详情接口，这里复用跨门店列表逻辑作为占位；后续可实现 GetCategory(id)
func (a dishCategoryCRUDAdapter) Get(ctx *gin.Context)    { a.inner.ListCategoriesForStore(ctx) }
func (a dishCategoryCRUDAdapter) Update(ctx *gin.Context) { a.inner.UpdateCategory(ctx) }
func (a dishCategoryCRUDAdapter) Delete(ctx *gin.Context) { a.inner.DeleteCategory(ctx) }

// registerCRUD 注册标准 CRUD 路由集合（不包含 GET /:id 逻辑占位示例可扩展）
func registerCRUD(group *gin.RouterGroup, ctrl CRUDController) {
	group.POST("", ctrl.Create)
	group.GET("", ctrl.List)
	group.PUT(":id", ctrl.Update)
	group.DELETE(":id", ctrl.Delete)
}

func RegisterRoutes(r *gin.Engine, c *AppControllers) {
	v1 := r.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/register", c.User.Register)
		auth.POST("/login", c.User.Login)
	}

	users := v1.Group("/users")
	users.Use(middleware.StoreAuthMiddleware())
	{
		users.GET("/profile", c.User.GetProfile)
		users.PUT("/profile", c.User.UpdateProfile)
		users.POST("", c.User.CreateUser)
		users.GET("", c.User.ListUsers)
		users.GET("/:id", c.User.GetUser)
		users.PUT("/:id", c.User.UpdateUser)
		users.DELETE("/:id", c.User.DeleteUser)
		users.POST(":id/reset-password", c.User.ResetUserPassword)
	}

	stores := v1.Group("/stores")
	stores.Use(middleware.AuthMiddleware())
	{
		stores.POST("", c.Store.CreateStore)
		stores.GET("", c.Store.ListStores)
		stores.GET("/all", c.Store.ListAllStores)
		// 门店分类：创建 / 列表 / 分类菜品 / 删除（在单门店详情前声明避免冲突）
		stores.POST("/:id/dish-categories", c.DishCategory.CreateCategoryForStore)
		stores.GET("/:id/dish-categories", c.DishCategory.ListCategoriesForStore)
		stores.POST("/:id/dish-categories/:cid/dishes", c.DishCategory.CreateDishForStoreCategory)
		stores.GET("/:id/dish-categories/:cid/dishes", c.DishCategory.ListDishesForStoreCategory)
		stores.PUT("/:id/dish-categories/:cid/dishes/:did", c.DishCategory.UpdateDishForStoreCategory)
		stores.DELETE("/:id/dish-categories/:cid/dishes/:did", c.DishCategory.DeleteDishForStoreCategory)
		stores.PUT("/:id/dish-categories/:cid", c.DishCategory.UpdateCategoryForStore)
		stores.DELETE("/:id/dish-categories/:cid", c.DishCategory.DeleteCategoryForStore)
		stores.PUT("/:id/dishes/:did", c.Dish.UpdateDishForStore)
		stores.DELETE("/:id/dishes/:did", c.Dish.DeleteDishForStore)
		stores.GET("/:id", c.Store.GetStore)
		stores.PUT("/:id", c.Store.UpdateStore)
		stores.DELETE("/:id", c.Store.DeleteStore)
	}

	dishes := v1.Group("/dishes")
	dishes.Use(middleware.StoreAuthMiddleware())
	{
		dishes.POST("", c.Dish.CreateDish)
		dishes.GET("", c.Dish.ListDishes)
		dishes.GET("/:id", c.Dish.GetDish)
		dishes.PUT("/:id", c.Dish.UpdateDish)
		dishes.DELETE("/:id", c.Dish.DeleteDish)
		dishes.GET("/by-category", c.DishCategory.ListCategoriesWithDishes)
	}

	cats := v1.Group("/dish-categories")
	cats.Use(middleware.StoreAuthMiddleware())
	{
		registerCRUD(cats, dishCategoryCRUDAdapter{inner: c.DishCategory})
		cats.POST("/reorder", c.DishCategory.ReorderCategories)
	}

	reports := v1.Group("/menu-reports")
	reports.Use(middleware.StoreAuthMiddleware())
	{
		reports.POST("", c.MenuReport.CreateMenuReport)
		reports.GET("", c.MenuReport.ListMenuReports)
		reports.GET("/statistics", c.MenuReport.GetStatistics)
		reports.GET("/:id", c.MenuReport.GetMenuReport)
		reports.PUT("/:id", c.MenuReport.UpdateMenuReport)
		reports.DELETE("/:id", c.MenuReport.DeleteMenuReport)
	}

	menus := v1.Group("/menus")
	menus.Use(middleware.AuthMiddleware())
	{
		menus.POST("", c.Menu.CreateMenu)
		menus.GET("", c.Menu.ListMenus)
		menus.GET("/tree", c.Menu.GetMenuTree)
		menus.GET("/:id", c.Menu.GetMenu)
		menus.PUT("/:id", c.Menu.UpdateMenu)
		menus.DELETE("/:id", c.Menu.DeleteMenu)
		menus.POST("/assign-role", c.Menu.AssignMenusToRole)
		menus.GET("/role", c.Menu.GetRoleMenus)
		menus.GET("/role-ids", c.Menu.GetRoleMenuIDs)
		menus.POST("/assign-store-role", c.Menu.AssignMenusToStoreRole)
		menus.GET("/store-role", c.Menu.GetStoreRoleMenus)
		menus.GET("/store-role-ids", c.Menu.GetStoreRoleMenuIDs)
		menus.POST("/copy-store", c.Menu.CopyStoreMenus)
		menus.GET("/user-menus", c.Menu.GetUserMenus)
		menus.GET("/user-permissions", c.Menu.GetUserPermissions)
	}

	roles := v1.Group("/roles")
	roles.Use(middleware.AuthMiddleware())
	{
		roles.POST("", controller.CreateRole)
		roles.GET("", controller.ListRoles)
		roles.GET("/:id", controller.GetRole)
		roles.PUT("/:id", controller.UpdateRole)
		roles.DELETE("/:id", controller.DeleteRole)
	}

	// 钉钉管理
	dingtalk := v1.Group("/dingtalk")
	dingtalk.Use(middleware.AuthMiddleware())
	{
		// 机器人配置
		robots := dingtalk.Group("/robots")
		{
			robots.POST("", c.DingTalkBot.CreateBot)
			robots.GET("", c.DingTalkBot.ListBots)
			robots.GET("/:id", c.DingTalkBot.GetBot)
			robots.PUT("/:id", c.DingTalkBot.UpdateBot)
			robots.DELETE("/:id", c.DingTalkBot.DeleteBot)
			robots.POST("/:id/test", c.DingTalkBot.TestBot)
		}
	}

	r.GET("/ws", controller.WebSocketHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf(":%d", config.GetConfig().App.Port)
	swaggerURL := fmt.Sprintf("http://localhost%s/swagger/index.html", addr)
	fmt.Printf("📚 Swagger UI: %s\n\n", swaggerURL)
}
