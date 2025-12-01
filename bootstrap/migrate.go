package bootstrap

import (
	"fmt"
	"os"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils/database"
	"github.com/Kevin-Jii/tower-go/utils/logging"

	"go.uber.org/zap"
)

func AutoMigrateAndSeeds() {
	if os.Getenv("SKIP_AUTO_MIGRATE") == "1" {
		logging.LogInfo("跳过数据库迁移（SKIP_AUTO_MIGRATE=1）")
		return
	}
	// 可通过环境变量 SKIP_MIGRATION=1 跳过迁移和种子数据初始化（加快启动速度）
	if skipMigration := fmt.Sprintf("%v", database.DB.Migrator().HasTable(&model.User{})); skipMigration == "true" {
		// 表已存在，跳过详细检查
		logging.LogInfo("数据表已存在，跳过迁移检查（如需强制迁移请删除表）")
	}

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
		&model.RoleMenu{},
		&model.StoreRoleMenu{},
		&model.DingTalkBot{},
		&model.Supplier{},
		&model.SupplierCategory{},
		&model.SupplierProduct{},
		&model.StoreSupplierProduct{},
		&model.PurchaseOrder{},
		&model.PurchaseOrderItem{},
	}
	for _, m := range migrateModels {
		if err := database.DB.AutoMigrate(m); err != nil {
			logging.LogError("数据表迁移失败", zap.String("model", fmt.Sprintf("%T", m)), zap.Error(err))
			logging.LogWarn("迁移失败，后续种子数据将跳过")
			return
		}
	}
	logging.LogInfo("数据表迁移完成")

	if err := database.CreateOptimizedIndexes(database.DB); err != nil {
		logging.LogError("创建优化索引失败", zap.Error(err))
	}

	// 种子数据初始化已移至 SQL 文件: migrations/init_seed_data.sql
	// 首次部署请手动执行: mysql -u用户名 -p密码 数据库名 < migrations/init_seed_data.sql
	logging.LogInfo("数据表迁移完成，种子数据请执行 migrations/init_seed_data.sql")
}
