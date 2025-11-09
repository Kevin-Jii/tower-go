package batch

import (
	"fmt"
	"reflect"
	"tower-go/utils/logging"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SortItem 通用排序更新条目
type SortItem struct {
	ID   uint
	Sort int
}

// BatchUpdateSort 通用批量排序更新（按 store_id + id 条件）
// modelParam: 传入 &ModelStruct{}，用于 gorm.Model 解析表
func BatchUpdateSort(tx *gorm.DB, modelParam interface{}, storeID uint, items []SortItem) error {
	return tx.Transaction(func(t *gorm.DB) error {
		for _, it := range items {
			if err := t.Model(modelParam).Where("id = ? AND store_id = ?", it.ID, storeID).Update("sort", it.Sort).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BatchSize 默认批量操作大小
const DefaultBatchSize = 100

// BatchInsert 批量插入（自动分批）
func BatchInsert(db *gorm.DB, records interface{}, batchSize ...int) error {
	size := DefaultBatchSize
	if len(batchSize) > 0 && batchSize[0] > 0 {
		size = batchSize[0]
	}

	err := db.CreateInBatches(records, size).Error
	if err != nil {
		logging.LogError("批量插入失败", zap.Error(err))
		return err
	}

	logging.LogInfo("批量插入成功", zap.Int("batch_size", size))
	return nil
}

// BatchUpdate 批量更新（通过主键）
func BatchUpdate(db *gorm.DB, model interface{}, ids []uint, updates map[string]interface{}) error {
	if len(ids) == 0 {
		return fmt.Errorf("ids cannot be empty")
	}

	err := db.Model(model).Where("id IN ?", ids).Updates(updates).Error
	if err != nil {
		logging.LogError("批量更新失败", zap.Error(err), zap.Int("count", len(ids)))
		return err
	}

	logging.LogInfo("批量更新成功", zap.Int("count", len(ids)))
	return nil
}

// BatchDelete 批量删除（软删除）
func BatchDelete(db *gorm.DB, model interface{}, ids []uint) error {
	if len(ids) == 0 {
		return fmt.Errorf("ids cannot be empty")
	}

	err := db.Where("id IN ?", ids).Delete(model).Error
	if err != nil {
		logging.LogError("批量删除失败", zap.Error(err), zap.Int("count", len(ids)))
		return err
	}

	logging.LogInfo("批量删除成功", zap.Int("count", len(ids)))
	return nil
}

// BatchHardDelete 批量硬删除
func BatchHardDelete(db *gorm.DB, model interface{}, ids []uint) error {
	if len(ids) == 0 {
		return fmt.Errorf("ids cannot be empty")
	}

	err := db.Unscoped().Where("id IN ?", ids).Delete(model).Error
	if err != nil {
		logging.LogError("批量硬删除失败", zap.Error(err), zap.Int("count", len(ids)))
		return err
	}

	logging.LogInfo("批量硬删除成功", zap.Int("count", len(ids)))
	return nil
}

// BatchUpdateField 批量更新单个字段
func BatchUpdateField(db *gorm.DB, model interface{}, ids []uint, field string, value interface{}) error {
	if len(ids) == 0 {
		return fmt.Errorf("ids cannot be empty")
	}

	updates := map[string]interface{}{field: value}
	return BatchUpdate(db, model, ids, updates)
}

// BatchUpdateStatus 批量更新状态
func BatchUpdateStatus(db *gorm.DB, model interface{}, ids []uint, status int) error {
	return BatchUpdateField(db, model, ids, "status", status)
}

// BatchOperation 批量操作辅助函数（带事务）
func BatchOperation(db *gorm.DB, operation func(tx *gorm.DB) error) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := operation(tx); err != nil {
			logging.LogError("批量操作失败，事务已回滚", zap.Error(err))
			return err
		}
		logging.LogInfo("批量操作成功，事务已提交")
		return nil
	})
}

// BatchInsertWithTransaction 批量插入（带事务）
func BatchInsertWithTransaction(db *gorm.DB, records interface{}, batchSize ...int) error {
	return BatchOperation(db, func(tx *gorm.DB) error {
		return BatchInsert(tx, records, batchSize...)
	})
}

// ChunkProcess 分块处理大量数据（避免内存溢出）
func ChunkProcess(db *gorm.DB, model interface{}, chunkSize int, processor func(records interface{}) error) error {
	offset := 0

	for {
		// 使用反射创建切片
		modelType := reflect.TypeOf(model)
		if modelType.Kind() == reflect.Ptr {
			modelType = modelType.Elem()
		}
		sliceType := reflect.SliceOf(reflect.PointerTo(modelType))
		records := reflect.New(sliceType).Interface()

		err := db.Offset(offset).Limit(chunkSize).Find(records).Error
		if err != nil {
			logging.LogError("分块查询失败", zap.Error(err), zap.Int("offset", offset))
			return err
		}

		// 检查是否还有数据
		slice := reflect.ValueOf(records).Elem()
		count := slice.Len()
		if count == 0 {
			break
		}

		// 处理这批数据
		if err := processor(records); err != nil {
			logging.LogError("分块处理失败", zap.Error(err), zap.Int("offset", offset))
			return err
		}

		// 如果返回的数据少于 chunkSize，说明已经是最后一批
		if count < chunkSize {
			break
		}

		offset += chunkSize
		logging.LogDebug("分块处理进度", zap.Int("processed", offset))
	}

	logging.LogInfo("分块处理完成", zap.Int("total_processed", offset))
	return nil
}

// BulkUpsert 批量 Upsert（插入或更新）
func BulkUpsert(db *gorm.DB, records interface{}, conflictColumns []string, updateColumns []string) error {
	// 构建更新字段
	updateMap := make(map[string]clause.Column)
	for _, col := range updateColumns {
		updateMap[col] = clause.Column{Name: col}
	}

	err := db.Clauses(clause.OnConflict{
		Columns:   convertToColumns(conflictColumns),
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(records).Error

	if err != nil {
		logging.LogError("批量 Upsert 失败", zap.Error(err))
		return err
	}

	logging.LogInfo("批量 Upsert 成功")
	return nil
}

// convertToColumns 转换字符串切片为 Column 切片
func convertToColumns(cols []string) []clause.Column {
	columns := make([]clause.Column, len(cols))
	for i, col := range cols {
		columns[i] = clause.Column{Name: col}
	}
	return columns
}
