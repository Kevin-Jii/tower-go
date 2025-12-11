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
	Dict              *controller.DictController
	Inventory         *controller.InventoryController
	File              *controller.FileController
	Gallery           *controller.GalleryController
	StoreAccount      *controller.StoreAccountController
	Statistics        *controller.StatisticsController
	MessageTemplate   *controller.MessageTemplateController
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
	dingTalkUserModule := userModulePkg.NewDingTalkUserModule(database.DB)
	supplierModule := userModulePkg.NewSupplierModule(database.DB)
	supplierCategoryModule := userModulePkg.NewSupplierCategoryModule(database.DB)
	supplierProductModule := userModulePkg.NewSupplierProductModule(database.DB)
	storeSupplierModule := userModulePkg.NewStoreSupplierModule(database.DB)
	purchaseOrderModule := userModulePkg.NewPurchaseOrderModule(database.DB)
	dictModule := userModulePkg.NewDictModule(database.DB)
	inventoryModule := userModulePkg.NewInventoryModule(database.DB)
	galleryModule := userModulePkg.NewGalleryModule(database.DB)
	storeAccountModule := userModulePkg.NewStoreAccountModule(database.DB)
	statisticsModule := userModulePkg.NewStatisticsModule(database.DB)
	messageTemplateModule := userModulePkg.NewMessageTemplateModule(database.DB)

	userModulePkg.SetDB(database.DB)

	// 初始化RustFS文件服务（可选）- 提前初始化以便其他服务使用
	var rustfsService *service.RustFSService
	var imageGeneratorService *service.ImageGeneratorService

	rustfsConfig := config.GetRustFSConfig()
	fmt.Printf("📁 RustFS配置: enabled=%v, endpoint=%s, bucket=%s, notifyBucket=%s\n", rustfsConfig.Enabled, rustfsConfig.Endpoint, rustfsConfig.Bucket, rustfsConfig.NotifyBucket)
	if rustfsConfig.Enabled {
		fmt.Println("📁 正在连接RustFS服务...")
		var err error
		rustfsService, err = service.NewRustFSServiceWithNotify(
			rustfsConfig.Endpoint,
			rustfsConfig.AccessKey,
			rustfsConfig.SecretKey,
			rustfsConfig.Bucket,
			rustfsConfig.NotifyBucket,
			rustfsConfig.UseSSL,
		)
		if err != nil {
			fmt.Printf("❌ RustFS服务连接失败: %v\n", err)
			logging.LogWarn("RustFS服务连接失败，文件服务不可用: " + err.Error())
		} else {
			fmt.Println("✅ RustFS文件服务已启用")
			logging.LogInfo("RustFS文件服务已启用")
			// 初始化图片生成服务
			imageGeneratorService = service.NewImageGeneratorService(rustfsService)
		}
	} else {
		fmt.Println("⚠️  RustFS文件服务未启用 (RUSTFS_ENABLED=false)")
	}

	// 初始化服务层
	userService := service.NewUserService(userModule)
	storeService := service.NewStoreService(storeModule)
	dingTalkService := service.NewDingTalkService(dingTalkBotModule, dingTalkUserModule)
	menuService := service.NewMenuService(menuModule, roleMenuModule, storeRoleMenuModule)
	supplierService := service.NewSupplierService(supplierModule)
	supplierProductService := service.NewSupplierProductService(supplierProductModule, supplierCategoryModule, supplierModule)
	storeSupplierService := service.NewStoreSupplierService(storeSupplierModule)
	purchaseOrderService := service.NewPurchaseOrderService(purchaseOrderModule, supplierProductModule, storeSupplierModule)
	dictService := service.NewDictService(dictModule)
	messageTemplateService := service.NewMessageTemplateService(messageTemplateModule)
	inventoryService := service.NewInventoryService(inventoryModule, userModule, storeModule, supplierProductModule, dingTalkService, dingTalkBotModule, messageTemplateService)
	storeAccountService := service.NewStoreAccountService(storeAccountModule, supplierProductModule, storeModule, userModule, dictModule, dingTalkService, dingTalkBotModule, messageTemplateService, imageGeneratorService)
	statisticsService := service.NewStatisticsService(statisticsModule)

	// 初始化默认消息模板
	if err := messageTemplateService.InitDefaultTemplates(); err != nil {
		logging.LogWarn("初始化消息模板失败: " + err.Error())
	}

	// 初始化钉钉命令处理器
	service.InitCommandHandler(inventoryModule, storeAccountModule, storeModule, userModule, messageTemplateService)

	// 初始化文件和图库控制器（依赖RustFS）
	var fileController *controller.FileController
	var galleryController *controller.GalleryController
	if rustfsService != nil {
		fileController = controller.NewFileController(rustfsService)
		galleryService := service.NewGalleryService(galleryModule, rustfsService)
		galleryController = controller.NewGalleryController(galleryService, rustfsService)
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
		Dict:              controller.NewDictController(dictService),
		Inventory:         controller.NewInventoryController(inventoryService),
		File:              fileController,
		Gallery:           galleryController,
		StoreAccount:      controller.NewStoreAccountController(storeAccountService),
		Statistics:        controller.NewStatisticsController(statisticsService),
		MessageTemplate:   controller.NewMessageTemplateController(messageTemplateService),
		DingTalkBotModule: dingTalkBotModule,
	}
}
