package api

import (
	"github.com/Kevin-Jii/tower-go/controller"
	userModulePkg "github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/database"
)

// Controllers 应用控制器容器
type Controllers struct {
	User              *controller.UserController
	Store             *controller.StoreController
	Menu              *controller.MenuController
	DingTalkBot       *controller.DingTalkBotController
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

	// 初始化服务层
	userService := service.NewUserService(userModule)
	storeService := service.NewStoreService(storeModule)
	dingTalkService := service.NewDingTalkService(dingTalkBotModule)
	menuService := service.NewMenuService(menuModule, roleMenuModule, storeRoleMenuModule)
	supplierService := service.NewSupplierService(supplierModule)
	supplierProductService := service.NewSupplierProductService(supplierProductModule, supplierCategoryModule, supplierModule)
	storeSupplierService := service.NewStoreSupplierService(storeSupplierModule)
	purchaseOrderService := service.NewPurchaseOrderService(purchaseOrderModule, supplierProductModule, storeSupplierModule)

	return &Controllers{
		User:              controller.NewUserController(userService),
		Store:             controller.NewStoreController(storeService),
		Menu:              controller.NewMenuController(menuService),
		DingTalkBot:       controller.NewDingTalkBotController(dingTalkService),
		Supplier:          controller.NewSupplierController(supplierService),
		SupplierProduct:   controller.NewSupplierProductController(supplierProductService),
		StoreSupplier:     controller.NewStoreSupplierController(storeSupplierService),
		PurchaseOrder:     controller.NewPurchaseOrderController(purchaseOrderService),
		DingTalkBotModule: dingTalkBotModule,
	}
}
