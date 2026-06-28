package api

import (
	"fmt"
	"os"

	"github.com/Kevin-Jii/tower-go/config"
	"github.com/Kevin-Jii/tower-go/controller"
	"github.com/Kevin-Jii/tower-go/cron"
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
	InventoryLoss     *controller.InventoryLossController
	File              *controller.FileController
	Gallery           *controller.GalleryController
	StoreAccount      *controller.StoreAccountController
	StoreExpense      *controller.StoreExpenseController
	StoreReturn       *controller.StoreReturnController
	MeituanAI         *controller.MeituanAIController
	Statistics        *controller.StatisticsController
	MessageTemplate   *controller.MessageTemplateController
	Member            *controller.MemberController
	Printer           *controller.PrinterController
	PriceList         *controller.PriceListController
	B2B               *controller.B2BController
	ThirdPartyAccount *controller.ThirdPartyAccountController
	ThirdPartyRoute   *controller.ThirdPartyRouteController
	AuditLog          *controller.AuditLogController
	DingTalkBotModule *userModulePkg.DingTalkBotModule
	PrinterService    *service.PrinterService
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
	productUnitSpecModule := userModulePkg.NewProductUnitSpecModule(database.DB)
	storeSupplierModule := userModulePkg.NewStoreSupplierModule(database.DB)
	purchaseOrderModule := userModulePkg.NewPurchaseOrderModule(database.DB)
	dictModule := userModulePkg.NewDictModule(database.DB)
	inventoryModule := userModulePkg.NewInventoryModule(database.DB)
	inventoryLossModule := userModulePkg.NewInventoryLossModule(database.DB)
	galleryModule := userModulePkg.NewGalleryModule(database.DB)
	storeAccountModule := userModulePkg.NewStoreAccountModule(database.DB)
	storeExpenseModule := userModulePkg.NewStoreExpenseModule(database.DB)
	storeReturnModule := userModulePkg.NewStoreReturnModule(database.DB)
	meituanAIModule := userModulePkg.NewMeituanAIModule(database.DB)
	statisticsModule := userModulePkg.NewStatisticsModule(database.DB)
	messageTemplateModule := userModulePkg.NewMessageTemplateModule(database.DB)
	memberModule := userModulePkg.NewMemberModule(database.DB)
	priceListModule := userModulePkg.NewPriceListModule(database.DB)
	b2bModule := userModulePkg.NewB2BModule(database.DB)
	thirdPartyAccountModule := userModulePkg.NewThirdPartyAccountModule(database.DB)
	thirdPartyOrderModule := userModulePkg.NewThirdPartyOrderModule(database.DB)
	thirdPartyRouteModule := userModulePkg.NewThirdPartyRouteModule(database.DB)
	auditLogModule := userModulePkg.NewAuditLogModule(database.DB)

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
			rustfsConfig.PublicBaseURL,
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
	userService := service.NewUserService(userModule, storeModule)
	storeService := service.NewStoreService(storeModule, thirdPartyAccountModule)
	dingTalkService := service.NewDingTalkService(dingTalkBotModule, dingTalkUserModule)
	menuService := service.NewMenuService(menuModule, roleMenuModule, storeRoleMenuModule)
	supplierService := service.NewSupplierService(supplierModule, storeSupplierModule)
	supplierProductService := service.NewSupplierProductService(supplierProductModule, productUnitSpecModule, dictModule, supplierCategoryModule, supplierModule)
	storeSupplierService := service.NewStoreSupplierService(storeSupplierModule, productUnitSpecModule)
	purchaseOrderService := service.NewPurchaseOrderService(purchaseOrderModule, supplierProductModule, storeSupplierModule, storeModule, dingTalkBotModule, dingTalkService)
	dictService := service.NewDictService(dictModule)
	messageTemplateService := service.NewMessageTemplateService(messageTemplateModule)
	inventoryService := service.NewInventoryService(inventoryModule, productUnitSpecModule, userModule, storeModule, supplierProductModule, dingTalkService, dingTalkBotModule, messageTemplateService)
	inventoryLossService := service.NewInventoryLossService(inventoryLossModule, supplierProductModule, productUnitSpecModule, memberModule, userModule, dictModule)
	storeAccountService := service.NewStoreAccountService(storeAccountModule, inventoryModule, supplierProductModule, productUnitSpecModule, storeModule, memberModule, userModule, dictModule, b2bModule, dingTalkService, dingTalkBotModule, messageTemplateService, imageGeneratorService)
	storeExpenseService := service.NewStoreExpenseService(storeExpenseModule, dictModule, userModule)
	storeReturnService := service.NewStoreReturnService(storeReturnModule, userModule)
	meituanAIService := service.NewMeituanAIService(meituanAIModule)
	statisticsService := service.NewStatisticsService(statisticsModule)
	memberService := service.NewMemberService(memberModule)
	memberService.SetDependencies(storeModule, dingTalkBotModule, dictModule, userModule, dingTalkService)
	priceListService := service.NewPriceListService(priceListModule, storeModule, supplierProductModule)
	b2bService := service.NewB2BService(b2bModule, storeModule, supplierProductModule, productUnitSpecModule, userModule)
	thirdPartyAccountService := service.NewThirdPartyAccountService(thirdPartyAccountModule, thirdPartyOrderModule)
	thirdPartyRouteService := service.NewThirdPartyRouteService(thirdPartyRouteModule, storeModule, thirdPartyOrderModule)
	auditLogService := service.NewAuditLogService(auditLogModule)

	// 初始化打印机模块
	printerModule := userModulePkg.NewPrinterModule(database.DB)
	printerService := service.NewPrinterService(printerModule, storeModule, purchaseOrderModule)

	// 从配置初始化芯烨云客户端（如果配置了）
	xpyunConfig := config.GetConfig().Xpyun
	if xpyunConfig.User != "" && xpyunConfig.UserKey != "" {
		printerService.InitXpyunClient(xpyunConfig.User, xpyunConfig.UserKey, xpyunConfig.BaseURL)
		fmt.Printf(">>>>>> 芯烨云配置: BaseURL=%s, User=%s\n", xpyunConfig.BaseURL, xpyunConfig.User)
		logging.LogInfo("芯烨云打印机客户端已初始化")
	}

	// 初始化默认消息模板（SKIP_SEED_DATA=1 时跳过，由 SQL 手工导入）
	if os.Getenv("SKIP_SEED_DATA") != "1" {
		if err := messageTemplateService.InitDefaultTemplates(); err != nil {
			logging.LogWarn("初始化消息模板失败: " + err.Error())
		}
	} else {
		logging.LogInfo("跳过空库默认消息模板写入（SKIP_SEED_DATA=1）")
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
		SupplierProduct:   controller.NewSupplierProductController(supplierProductService, storeSupplierService),
		StoreSupplier:     controller.NewStoreSupplierController(storeSupplierService),
		PurchaseOrder:     controller.NewPurchaseOrderController(purchaseOrderService),
		Dict:              controller.NewDictController(dictService),
		Inventory:         controller.NewInventoryController(inventoryService),
		InventoryLoss:     controller.NewInventoryLossController(inventoryLossService),
		File:              fileController,
		Gallery:           galleryController,
		StoreAccount:      controller.NewStoreAccountController(storeAccountService),
		StoreExpense:      controller.NewStoreExpenseController(storeExpenseService),
		StoreReturn:       controller.NewStoreReturnController(storeReturnService),
		MeituanAI:         controller.NewMeituanAIController(meituanAIService),
		Statistics:        controller.NewStatisticsController(statisticsService),
		MessageTemplate:   controller.NewMessageTemplateController(messageTemplateService),
		Member:            controller.NewMemberController(memberService),
		Printer:           controller.NewPrinterController(printerService),
		PriceList:         controller.NewPriceListController(priceListService),
		B2B:               controller.NewB2BController(b2bService),
		ThirdPartyAccount: controller.NewThirdPartyAccountController(thirdPartyAccountService),
		ThirdPartyRoute:   controller.NewThirdPartyRouteController(thirdPartyRouteService),
		AuditLog:          controller.NewAuditLogController(auditLogService),
		DingTalkBotModule: dingTalkBotModule,
		PrinterService:    printerService,
	}
}

// StartCronJobs 启动定时任务
func (c *Controllers) StartCronJobs() error {
	// 启动打印机状态同步任务
	_, err := cron.StartPrinterSync(c.PrinterService)
	return err
}
