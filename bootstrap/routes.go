package bootstrap

import (
	"fmt"
	"tower-go/config"
	"tower-go/controller"
	"tower-go/middleware"
	userModulePkg "tower-go/module"
	"tower-go/service"
	"tower-go/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type AppControllers struct {
	User         *controller.UserController
	Store        *controller.StoreController
	Dish         *controller.DishController
	DishCategory *controller.DishCategoryController
	MenuReport   *controller.MenuReportController
	Menu         *controller.MenuController
}

func BuildControllers() *AppControllers {
	userModule := userModulePkg.NewUserModule(utils.DB)
	storeModule := userModulePkg.NewStoreModule(utils.DB)
	dishModule := userModulePkg.NewDishModule(utils.DB)
	menuReportModule := userModulePkg.NewMenuReportModule(utils.DB)
	dishCategoryModule := userModulePkg.NewDishCategoryModule(utils.DB)
	menuModule := userModulePkg.NewMenuModule(utils.DB)
	roleMenuModule := userModulePkg.NewRoleMenuModule(utils.DB)
	storeRoleMenuModule := userModulePkg.NewStoreRoleMenuModule(utils.DB)

	// 初始化角色模块全局 DB（旧实现依赖 SetDB）
	userModulePkg.SetDB(utils.DB)

	userService := service.NewUserService(userModule)
	storeService := service.NewStoreService(storeModule)
	dishService := service.NewDishService(dishModule)
	menuReportService := service.NewMenuReportService(menuReportModule, dishModule)
	dishCategoryService := service.NewDishCategoryService(dishCategoryModule)
	menuService := service.NewMenuService(menuModule, roleMenuModule, storeRoleMenuModule)

	return &AppControllers{
		User:         controller.NewUserController(userService),
		Store:        controller.NewStoreController(storeService),
		Dish:         controller.NewDishController(dishService),
		DishCategory: controller.NewDishCategoryController(dishCategoryService),
		MenuReport:   controller.NewMenuReportController(menuReportService),
		Menu:         controller.NewMenuController(menuService),
	}
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
		cats.POST("", c.DishCategory.CreateCategory)
		cats.GET("", c.DishCategory.ListCategories)
		cats.PUT("/:id", c.DishCategory.UpdateCategory)
		cats.DELETE("/:id", c.DishCategory.DeleteCategory)
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

	r.GET("/ws", controller.WebSocketHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf(":%d", config.GetConfig().App.Port)
	swaggerURL := fmt.Sprintf("http://localhost%s/swagger/index.html", addr)
	fmt.Printf("📚 Swagger UI: %s\n\n", swaggerURL)
}
