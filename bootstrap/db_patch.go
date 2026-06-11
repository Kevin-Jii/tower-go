package bootstrap

import (
	"github.com/Kevin-Jii/tower-go/utils/logging"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// applyUserStorePatches 超级管理员 store_id=0 不引用 stores 表，需移除 GORM 自动创建的外键并归一化历史数据。
func applyUserStorePatches(db *gorm.DB) {
	if db == nil {
		return
	}

	var fkName *string
	if err := db.Raw(`
		SELECT CONSTRAINT_NAME
		FROM information_schema.TABLE_CONSTRAINTS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'users'
		  AND CONSTRAINT_TYPE = 'FOREIGN KEY'
		  AND CONSTRAINT_NAME = 'fk_users_store'
		LIMIT 1
	`).Scan(&fkName).Error; err != nil {
		logging.LogWarn("检查 users 外键失败", zap.Error(err))
		return
	}
	if fkName != nil && *fkName != "" {
		if err := db.Exec("ALTER TABLE `users` DROP FOREIGN KEY `fk_users_store`").Error; err != nil {
			logging.LogWarn("删除 users.store_id 外键失败", zap.Error(err))
		} else {
			logging.LogInfo("已删除 users.store_id 外键 fk_users_store")
		}
	}

	if err := db.Exec(`
		UPDATE users u
		INNER JOIN roles r ON u.role_id = r.id
		SET u.store_id = 0
		WHERE r.code = 'super_admin' AND u.store_id <> 0
	`).Error; err != nil {
		logging.LogWarn("归一化超级管理员 store_id 失败", zap.Error(err))
	}
}
