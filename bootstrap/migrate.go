package bootstrap

import (
	"fmt"
	"tower-go/model"
	"tower-go/utils"

	"go.uber.org/zap"
)

func AutoMigrateAndSeeds() {
	// 外键前置数据检查
	var invalidUserCount int64
	utils.DB.Raw("SELECT COUNT(*) FROM users u LEFT JOIN stores s ON u.store_id = s.id WHERE u.store_id <> 0 AND s.id IS NULL").Scan(&invalidUserCount)
	if invalidUserCount > 0 {
		utils.LogWarn("发现无效用户记录", zap.Int64("count", invalidUserCount))
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
		if err := utils.DB.AutoMigrate(m); err != nil {
			utils.LogError("数据表迁移失败", zap.String("model", fmt.Sprintf("%T", m)), zap.Error(err))
			utils.LogWarn("迁移失败，后续种子数据将跳过")
			return
		}
	}
	utils.LogInfo("数据表迁移完成")

	if err := utils.MigrateDishCategoryData(); err != nil {
		utils.LogWarn("菜品分类数据迁移失败", zap.Error(err))
	} else {
		utils.LogInfo("菜品分类数据迁移完成")
	}

	if err := utils.CreateOptimizedIndexes(utils.DB); err != nil {
		utils.LogError("创建优化索引失败", zap.Error(err))
	}

	if err := utils.InitRoleSeeds(utils.DB); err != nil {
		utils.LogError("角色基础数据初始化失败", zap.Error(err))
	}
	if err := utils.InitMenuSeeds(utils.DB); err != nil {
		utils.LogError("菜单种子数据初始化失败", zap.Error(err))
	}
	if err := utils.InitRoleMenuSeeds(utils.DB); err != nil {
		utils.LogError("角色菜单权限初始化失败", zap.Error(err))
	}
	if err := utils.InitSuperAdmin(utils.DB); err != nil {
		utils.LogError("超级管理员初始化失败", zap.Error(err))
	}
	if err := utils.EnsureStoreCodes(utils.DB); err != nil {
		utils.LogError("门店编码补全失败", zap.Error(err))
	}

	// 初始化钉钉管理菜单
	if err := utils.InitDingTalkMenuSeeds(utils.DB); err != nil {
		utils.LogError("钉钉管理菜单初始化失败", zap.Error(err))
	}
}
