package api

import (
	"fmt"

	"github.com/Kevin-Jii/tower-go/config"
	"github.com/Kevin-Jii/tower-go/controller"
	userModulePkg "github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/database"
	"github.com/Kevin-Jii/tower-go/utils/logging"
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
	File              *controller.FileController
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

	// 初始化MinIO文件服务（可选）
	var fileController *controller.FileController
	minioConfig := config.GetMinIOConfig()
	fmt.Printf("📁 MinIO配置: enabled=%v, endpoint=%s, bucket=%s\n", minioConfig.Enabled, minioConfig.Endpoint, minioConfig.Bucket)
	if minioConfig.Enabled {
		fmt.Println("📁 正在连接MinIO服务...")
		minioService, err := service.NewMinIOService(
			minioConfig.Endpoint,
			minioConfig.AccessKey,
			minioConfig.SecretKey,
			minioConfig.Bucket,
			minioConfig.UseSSL,
		)
		if err != nil {
			fmt.Printf("❌ MinIO服务连接失败: %v\n", err)
			logging.LogWarn("MinIO服务连接失败，文件服务不可用: " + err.Error())
		} else {
			fileController = controller.NewFileController(minioService)
			fmt.Println("✅ MinIO文件服务已启用")
			logging.LogInfo("MinIO文件服务已启用")
		}
	} else {
		fmt.Println("⚠️  MinIO文件服务未启用 (MINIO_ENABLED=false)")
	}

	return &Controllers{
		User:              controller.NewUserController(userService),
		Store:             controller.NewStoreController(storeService),
		Menu:              controller.NewMenuController(menuService),
		DingTalkBot:       controller.NewDingTalkBotController(dingTalkService),
		Supplier:          controller.NewSupplierController(supplierService),
		SupplierProduct:   controller.NewSupplierProductController(supplierProductService),
		StoreSupplier:     controller.NewStoreSupplierController(storeSupplierService),
		PurchaseOrder:     controller.NewPurchaseOrderController(purchaseOrderService),
		File:              fileController,
		DingTalkBotModule: dingTalkBotModule,
	}
}
