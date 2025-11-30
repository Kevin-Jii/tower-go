package api

import (
	"github.com/Kevin-Jii/tower-go/controller"
	userModulePkg "github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/database"
	"github.com/Kevin-Jii/tower-go/utils/events"
)

// Controllers 应用控制器容器
type Controllers struct {
	User              *controller.UserController
	Store             *controller.StoreController
	Dish              *controller.DishController
	DishCategory      *controller.DishCategoryController
	MenuReport        *controller.MenuReportController
	Menu              *controller.MenuController
	DingTalkBot       *controller.DingTalkBotController
	ReportBot         *controller.ReportBotController
	Supplier          *controller.SupplierController
	SupplierProduct   *controller.SupplierProductController
	StoreSupplier     *controller.StoreSupplierController
	PurchaseOrder     *controller.PurchaseOrderController
	DingTalkBotModule *userModulePkg.DingTalkBotModule
}

// BuildControllers 构建所有控制器及其依赖
func BuildControllers() *Controllers {
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
	supplierModule := userModulePkg.NewSupplierModule(database.DB)
	supplierCategoryModule := userModulePkg.NewSupplierCategoryModule(database.DB)
	supplierProductModule := userModulePkg.NewSupplierProductModule(database.DB)
	storeSupplierModule := userModulePkg.NewStoreSupplierModule(database.DB)
	purchaseOrderModule := userModulePkg.NewPurchaseOrderModule(database.DB)

	userModulePkg.SetDB(database.DB)

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
		dingTalkBotModule,
		eventBus,
	)
	dishCategoryService := service.NewDishCategoryService(dishCategoryModule, dishModule)
	menuService := service.NewMenuService(menuModule, roleMenuModule, storeRoleMenuModule)
	supplierService := service.NewSupplierService(supplierModule)
	supplierProductService := service.NewSupplierProductService(supplierProductModule, supplierCategoryModule, supplierModule)
	storeSupplierService := service.NewStoreSupplierService(storeSupplierModule)
	purchaseOrderService := service.NewPurchaseOrderService(purchaseOrderModule, supplierProductModule, storeSupplierModule)

	// 注册事件监听器
	menuReportListener := service.NewMenuReportEventListener(dingTalkService)
	service.RegisterMenuReportEventListeners(eventBus, menuReportListener)

	return &Controllers{
		User:              controller.NewUserController(userService),
		Store:             controller.NewStoreController(storeService),
		Dish:              controller.NewDishController(dishService),
		DishCategory:      controller.NewDishCategoryController(dishCategoryService),
		MenuReport:        controller.NewMenuReportController(menuReportService),
		Menu:              controller.NewMenuController(menuService),
		DingTalkBot:       controller.NewDingTalkBotController(dingTalkService),
		ReportBot:         controller.NewReportBotController(menuReportModule),
		Supplier:          controller.NewSupplierController(supplierService),
		SupplierProduct:   controller.NewSupplierProductController(supplierProductService),
		StoreSupplier:     controller.NewStoreSupplierController(storeSupplierService),
		PurchaseOrder:     controller.NewPurchaseOrderController(purchaseOrderService),
		DingTalkBotModule: dingTalkBotModule,
	}
}
