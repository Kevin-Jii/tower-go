package utils

import (
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// GenerateEmployeeNo 生成唯一的6位数工号
// 使用递增方式，从100000开始
func GenerateEmployeeNo(db *gorm.DB) (string, error) {
	var maxEmployeeNo string

	// 查询当前最大工号
	err := db.Raw("SELECT employee_no FROM users WHERE employee_no IS NOT NULL AND employee_no != '' ORDER BY employee_no DESC LIMIT 1").Scan(&maxEmployeeNo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		LogDatabaseError("查询最大工号失败", err)
		return "", fmt.Errorf("查询最大工号失败: %v", err)
	}

	var nextNo int
	if maxEmployeeNo == "" {
		// 没有找到任何工号，从100000开始
		nextNo = 100000
	} else {
		// 解析最大工号并+1
		_, err := fmt.Sscanf(maxEmployeeNo, "%d", &nextNo)
		if err != nil {
			LogError("解析工号失败")
			return "", fmt.Errorf("解析工号失败: %v", err)
		}
		nextNo++
	}

	// 检查是否超过6位数的最大值
	if nextNo > 999999 {
		LogError("工号已达到最大值")
		return "", fmt.Errorf("工号已达到最大值999999")
	}

	employeeNo := fmt.Sprintf("%06d", nextNo)

	LogInfo("生成新工号")
	return employeeNo, nil
}

// GenerateEmployeeNoRandom 使用随机方式生成6位数工号（备用方案）
// 在100000-999999范围内随机生成，检查唯一性，最多重试10次
func GenerateEmployeeNoRandom(db *gorm.DB) (string, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	maxRetries := 10

	for i := 0; i < maxRetries; i++ {
		// 生成6位随机数字 (100000-999999)
		randomNo := r.Intn(900000) + 100000
		employeeNo := fmt.Sprintf("%06d", randomNo)

		// 检查是否已存在
		var count int64
		err := db.Raw("SELECT COUNT(*) FROM users WHERE employee_no = ?", employeeNo).Scan(&count).Error
		if err != nil {
			LogDatabaseError("检查工号唯一性失败", err)
			continue
		}

		if count == 0 {
			LogInfo("随机生成新工号")
			return employeeNo, nil
		}
	}

	LogError("生成唯一工号失败")
	return "", fmt.Errorf("生成唯一工号失败，已重试%d次", maxRetries)
}

// GenerateStoreCode 生成唯一的门店编码 JWXXXX
func GenerateStoreCode(db *gorm.DB) (string, error) {
	if db == nil {
		return "", fmt.Errorf("数据库连接未初始化")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	maxRetries := 20

	for i := 0; i < maxRetries; i++ {
		code := fmt.Sprintf("JW%04d", r.Intn(10000))

		var count int64
		err := db.Raw("SELECT COUNT(*) FROM stores WHERE store_code = ?", code).Scan(&count).Error
		if err != nil {
			LogDatabaseError("检查门店编码唯一性失败", err)
			return "", fmt.Errorf("检查门店编码唯一性失败: %v", err)
		}

		if count == 0 {
			LogInfo("生成门店编码")
			return code, nil
		}
	}

	LogError("生成唯一门店编码失败")
	return "", fmt.Errorf("生成唯一门店编码失败，已重试%d次", maxRetries)
}

// GenerateDishCategoryCode 生成唯一的菜品分类编码 格式：DC{storeID}-{YYYYMMDD}{随机3位}
// 若需要更短，可改成 DC{storeID}{随机4位}
func GenerateDishCategoryCode(db *gorm.DB, storeID uint) (string, error) {
	if db == nil {
		return "", fmt.Errorf("数据库连接未初始化")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	maxRetries := 20
	datePart := time.Now().Format("20060102")
	for i := 0; i < maxRetries; i++ {
		suffix := r.Intn(1000) // 000-999
		code := fmt.Sprintf("DC%d-%s%03d", storeID, datePart, suffix)
		var count int64
		err := db.Raw("SELECT COUNT(*) FROM dish_categories WHERE code = ?", code).Scan(&count).Error
		if err != nil {
			LogDatabaseError("检查分类编码唯一性失败", err)
			return "", fmt.Errorf("检查分类编码唯一性失败: %v", err)
		}
		if count == 0 {
			LogInfo("生成菜品分类编码")
			return code, nil
		}
	}
	LogError("生成唯一菜品分类编码失败")
	return "", fmt.Errorf("生成唯一菜品分类编码失败，已重试%d次", maxRetries)
}
