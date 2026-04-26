package bootstrap

import (
	"fmt"
	"os"
	"strings"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils/database"
	"github.com/Kevin-Jii/tower-go/utils/logging"

	"go.uber.org/zap"
)

const migrationVersionFile = ".migration_version"
const currentMigrationVersion = "1"

// autoMigrateModels 与下方 AutoMigrate 顺序一致；shouldSkipMigration 会校验每张表均存在后才允许跳过。
var autoMigrateModels = []interface{}{
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
	&model.ProductUnitSpec{},
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
	&model.StoreAccountItem{},
	&model.DingTalkUser{},
	&model.MessageTemplate{},
	&model.Member{},
	&model.WalletLog{},
	&model.RechargeOrder{},
}

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
	for _, m := range autoMigrateModels {
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
	// 如果数据库关键表不存在，则不允许跳过（避免本地标记文件误导导致线上缺表）
	db := database.GetDB()
	if db == nil {
		return false
	}
	migrator := db.Migrator()
	// 任意一张 AutoMigrate 目标表缺失，都必须执行迁移（避免仅有部分表 + .migration_version 时误跳过）
	for _, m := range autoMigrateModels {
		if !migrator.HasTable(m) {
			return false
		}
	}

	// 记账表新增字段后，若线上库未加列则仍需迁移（避免 .migration_version 导致永远不 AutoMigrate）
	if migrator.HasTable(&model.StoreAccount{}) {
		if !migrator.HasColumn(&model.StoreAccount{}, "other_expense_amount") ||
			!migrator.HasColumn(&model.StoreAccount{}, "net_income_amount") {
			return false
		}
	}

	// 供应商商品新增双价格字段后，若线上库未加列则仍需迁移
	if migrator.HasTable(&model.SupplierProduct{}) {
		if !migrator.HasColumn(&model.SupplierProduct{}, "bottle_price") ||
			!migrator.HasColumn(&model.SupplierProduct{}, "case_price") ||
			!migrator.HasColumn(&model.SupplierProduct{}, "bottles_per_case") {
			return false
		}
	}

	// roles 表 data_scope（数据权限）
	if migrator.HasTable(&model.Role{}) && !migrator.HasColumn(&model.Role{}, "data_scope") {
		return false
	}

	// 检查标记文件
	if _, err := os.Stat(migrationVersionFile); err == nil {
		// 文件存在，读取版本
		data, err := os.ReadFile(migrationVersionFile)
		if err == nil && strings.TrimSpace(string(data)) == currentMigrationVersion {
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
