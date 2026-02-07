package bootstrap

import (
	"fmt"
	"os"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils/database"
	"github.com/Kevin-Jii/tower-go/utils/logging"

	"go.uber.org/zap"
)

const migrationVersionFile = ".migration_version"
const currentMigrationVersion = "1"

func AutoMigrateAndSeeds() {
	if os.Getenv("SKIP_AUTO_MIGRATE") == "1" {
		logging.LogInfo("跳过数据库迁移（SKIP_AUTO_MIGRATE=1）")
		return
	}

	// 检查迁移版本，避免重复迁移
	if shouldSkipMigration() {
		logging.LogInfo("数据表已迁移，跳过迁移检查")
		return
	}

	// 执行迁移
	migrateModels := []interface{}{
		&model.Store{},
		&model.Role{},
		&model.Menu{},
		&model.User{},
		&model.RoleMenu{},
		&model.StoreRoleMenu{},
		&model.DingTalkBot{},
		&model.Supplier{},
		&model.SupplierCategory{},
		&model.SupplierProduct{},
		&model.StoreSupplier{},
		&model.PurchaseOrder{},
		&model.PurchaseOrderItem{},
		&model.DictType{},
		&model.DictData{},
		&model.Inventory{},
		&model.InventoryOrder{},
		&model.InventoryOrderItem{},
		&model.Gallery{},
		&model.StoreAccount{},
		&model.DingTalkUser{},
		&model.MessageTemplate{},
		&model.Member{},
		&model.WalletLog{},
		&model.RechargeOrder{},
	}

	for _, m := range migrateModels {
		if err := database.GetDB().AutoMigrate(m); err != nil {
			logging.LogError("数据表迁移失败", zap.String("model", fmt.Sprintf("%T", m)), zap.Error(err))
			logging.LogWarn("迁移失败，后续种子数据将跳过")
			return
		}
	}
	logging.LogInfo("数据表迁移完成")

	if err := database.CreateOptimizedIndexes(database.GetDB()); err != nil {
		logging.LogError("创建优化索引失败", zap.Error(err))
	}

	// 标记迁移已完成
	markMigrationComplete()

	logging.LogInfo("数据表迁移完成，种子数据请执行 migrations/init_seed_data.sql")
}

// shouldSkipMigration 检查是否应该跳过迁移
func shouldSkipMigration() bool {
	// 检查标记文件
	if _, err := os.Stat(migrationVersionFile); err == nil {
		// 文件存在，读取版本
		data, err := os.ReadFile(migrationVersionFile)
		if err == nil && string(data) == currentMigrationVersion {
			return true
		}
	}
	return false
}

// markMigrationComplete 标记迁移已完成
func markMigrationComplete() {
	if err := os.WriteFile(migrationVersionFile, []byte(currentMigrationVersion), 0644); err != nil {
		logging.LogWarn("无法写入迁移标记文件", zap.Error(err))
	}
}
