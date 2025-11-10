package database
	
import (
	"database/sql"
	"fmt"
	"strings"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils/logging"

	"go.uber.org/zap"
)

// MigrateDishCategoryData 将旧 dishes.category 文本迁移到 dish_categories 并填充 category_id
// 安全且幂等：若旧列不存在或已无数据则直接跳过。执行完毕后尝试删除旧列。
func MigrateDishCategoryData() error {
	db := DB
	if db == nil {
		return fmt.Errorf("db not initialized")
	}

	// 检查旧列是否存在
	var exists int
	if err := db.Raw("SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='dishes' AND COLUMN_NAME='category' LIMIT 1").Scan(&exists).Error; err != nil {
		return err
	}
	if exists == 0 {
		return nil // 无需迁移
	}

	type Row struct {
		StoreID  uint
		Category sql.NullString
	}
	var rows []Row
	if err := db.Raw("SELECT DISTINCT store_id, category FROM dishes WHERE category IS NOT NULL AND category <> '' AND category_id IS NULL").Scan(&rows).Error; err != nil {
		return err
	}

	for _, r := range rows {
		name := strings.TrimSpace(r.Category.String)
		if name == "" {
			continue
		}
		// 找或建分类
		var cat model.DishCategory
		if err := db.Where("store_id = ? AND name = ?", r.StoreID, name).First(&cat).Error; err != nil {
			// 创建
			cat = model.DishCategory{StoreID: r.StoreID, Name: name, Status: 1}
			if err2 := db.Create(&cat).Error; err2 != nil {
				logging.LogWarn("创建分类失败", zap.Uint("store", r.StoreID), zap.String("name", name), zap.Error(err2))
				continue
			}
		}
		// 更新 dishes 绑定 category_id
		if err := db.Model(&model.Dish{}).Where("store_id = ? AND category = ? AND category_id IS NULL", r.StoreID, name).Update("category_id", cat.ID).Error; err != nil {
			logging.LogWarn("更新菜品分类ID失败", zap.Uint("store", r.StoreID), zap.String("name", name), zap.Error(err))
		}
	}

	// 尝试删除旧列
	if err := db.Exec("ALTER TABLE dishes DROP COLUMN category").Error; err != nil {
		logging.LogWarn("删除旧category列失败，可手动处理", zap.Error(err))
	} else {
		logging.LogInfo("旧category列已删除")
	}
	return nil
}
