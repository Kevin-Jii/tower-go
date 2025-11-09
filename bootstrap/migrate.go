package bootstrap

import (
	"fmt"
	"tower-go/model"
	"tower-go/utils/database"
	"tower-go/utils/logging"
	"tower-go/utils/seeding"

	"go.uber.org/zap"
)

func AutoMigrateAndSeeds() {
	// 外键前置数据检查
	var invalidUserCount int64
	database.DB.Raw("SELECT COUNT(*) FROM users u LEFT JOIN stores s ON u.store_id = s.id WHERE u.store_id <> 0 AND s.id IS NULL").Scan(&invalidUserCount)
	if invalidUserCount > 0 {
		logging.LogWarn("发现无效用户记录", zap.Int64("count", invalidUserCount))
	}

	migrateModels := []interface{}{
		&model.Store{},
		&model.Role{},
		&model.Menu{},
		&model.User{},
		&model.DishCategory{},
		&model.Dish{},
		&model.MenuReport{},
		&model.RoleMenu{},
		&model.StoreRoleMenu{},
		&model.DingTalkBot{},
	}
	for _, m := range migrateModels {
		if err := database.DB.AutoMigrate(m); err != nil {
			logging.LogError("数据表迁移失败", zap.String("model", fmt.Sprintf("%T", m)), zap.Error(err))
			logging.LogWarn("迁移失败，后续种子数据将跳过")
			return
		}
	}
	logging.LogInfo("数据表迁移完成")

	if err := database.MigrateDishCategoryData(); err != nil {
		logging.LogWarn("菜品分类数据迁移失败", zap.Error(err))
	} else {
		logging.LogInfo("菜品分类数据迁移完成")
	}

	if err := database.CreateOptimizedIndexes(database.DB); err != nil {
		logging.LogError("创建优化索引失败", zap.Error(err))
	}

	if err := seeding.InitRoleSeeds(database.DB); err != nil {
		logging.LogError("角色基础数据初始化失败", zap.Error(err))
	}
	if err := seeding.InitMenuSeeds(database.DB); err != nil {
		logging.LogError("菜单种子数据初始化失败", zap.Error(err))
	}
	if err := seeding.InitRoleMenuSeeds(database.DB); err != nil {
		logging.LogError("角色菜单权限初始化失败", zap.Error(err))
	}
	if err := seeding.InitSuperAdmin(database.DB); err != nil {
		logging.LogError("超级管理员初始化失败", zap.Error(err))
	}
	if err := seeding.EnsureStoreCodes(database.DB); err != nil {
		logging.LogError("门店编码补全失败", zap.Error(err))
	}

	// 初始化钉钉管理菜单
	if err := seeding.InitDingTalkMenuSeeds(database.DB); err != nil {
		logging.LogError("钉钉管理菜单初始化失败", zap.Error(err))
	}
}
