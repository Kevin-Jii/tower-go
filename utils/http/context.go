package http

import (
	"errors"
	"github.com/gin-gonic/gin"
)

// GetStoreID 从 Gin Context 中安全地提取 StoreID
func GetStoreID(ctx *gin.Context) (uint, error) {
	// 检查 StoreID 是否存在并转换为 uint 类型
	val, exists := ctx.Get("StoreID")
	if !exists {
		return 0, errors.New("StoreID not found in context (Middleware missing?)")
	}
	storeID, ok := val.(uint)
	if !ok || storeID == 0 {
		return 0, errors.New("Invalid StoreID in context")
	}
	return storeID, nil
}

// 示例：如果还需要 Gorm Scope 函数，可以在这里定义
// func StoreScope(storeID uint) func(db *gorm.DB) *gorm.DB {
//     return func(db *gorm.DB) *gorm.DB {
//         return db.Where("store_id = ?", storeID)
//     }
// }
