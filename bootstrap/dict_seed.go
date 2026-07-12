package bootstrap

import (
	"os"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils/database"
	"github.com/Kevin-Jii/tower-go/utils/logging"
	"go.uber.org/zap"
)

const dictDataVersion = "1"

// InitDefaultDicts 初始化默认字典数据（与 RunSeedSQL 相同开关：SKIP_SEED_DATA=1 时跳过）
func InitDefaultDicts() {
	if os.Getenv("SKIP_SEED_DATA") == "1" {
		logging.LogInfo("跳过内存字典种子（SKIP_SEED_DATA=1）")
		return
	}
	markerPath := initializationMarkerPath(dictMarkerFile)
	if initializationComplete(markerPath, dictDataVersion) {
		logging.LogInfo("跳过默认字典初始化（初始化版本已完成）", zap.String("version", dictDataVersion))
		return
	}
	if !initSalesChannel() || !initOrderSource() || !initInventoryReason() {
		logging.LogWarn("默认字典初始化未完成，不写入初始化标记，下次启动将重试")
		return
	}

	logging.LogInfo("字典数据初始化完成")
	if err := markInitializationComplete(markerPath, dictDataVersion); err != nil {
		logging.LogWarn("无法写入字典初始化标记", zap.Error(err))
	}
}

// initSalesChannel 初始化销售渠道字典
func initSalesChannel() bool {
	typeCode := "sales_channel"
	typeName := "销售渠道"

	// 创建或获取字典类型
	typeID := ensureDictType(typeCode, typeName, "门店记账-销售渠道")
	if typeID == 0 {
		return false
	}

	// 字典数据
	items := []struct {
		Label string
		Value string
		Sort  int
	}{
		{"线下门店", "offline", 1},
		{"美团外卖", "meituan", 2},
		{"饿了么", "eleme", 3},
		{"抖音", "douyin", 4},
		{"小红书", "xiaohongshu", 5},
		{"微信小程序", "wechat_mini", 6},
		{"团购", "group_buy", 7},
		{"商城", "mall", 8},
		{"其他", "other", 99},
	}

	// 批量检查已存在的数据
	existingValues, ok := getExistingDictValues(typeCode)
	if !ok {
		return false
	}
	for _, item := range items {
		if _, exists := existingValues[item.Value]; !exists {
			if !ensureDictData(typeID, typeCode, item.Label, item.Value, item.Sort) {
				return false
			}
		}
	}
	return true
}

// initOrderSource 初始化订单来源字典
func initOrderSource() bool {
	typeCode := "order_source"
	typeName := "订单来源"

	typeID := ensureDictType(typeCode, typeName, "门店记账-订单来源")
	if typeID == 0 {
		return false
	}

	items := []struct {
		Label string
		Value string
		Sort  int
	}{
		{"堂食", "dine_in", 1},
		{"外卖", "takeout", 2},
		{"自提", "pickup", 3},
		{"团购", "group_buy", 4},
		{"预订", "reservation", 5},
		{"其他", "other", 99},
	}

	// 批量检查已存在的数据
	existingValues, ok := getExistingDictValues(typeCode)
	if !ok {
		return false
	}
	for _, item := range items {
		if _, exists := existingValues[item.Value]; !exists {
			if !ensureDictData(typeID, typeCode, item.Label, item.Value, item.Sort) {
				return false
			}
		}
	}
	return true
}

// initInventoryReason 初始化出入库原因字典
func initInventoryReason() bool {
	typeCode := "inventory_reason"
	typeName := "出入库原因"

	typeID := ensureDictType(typeCode, typeName, "库存管理-出入库原因")
	if typeID == 0 {
		return false
	}

	items := []struct {
		Label string
		Value string
		Sort  int
	}{
		{"采购入库", "purchase_in", 1},
		{"退货入库", "return_in", 2},
		{"调拨入库", "transfer_in", 3},
		{"盘盈入库", "inventory_in", 4},
		{"销售出库", "sale_out", 10},
		{"报损出库", "loss_out", 11},
		{"调拨出库", "transfer_out", 12},
		{"盘亏出库", "inventory_out", 13},
		{"其他", "other", 99},
	}

	// 批量检查已存在的数据
	existingValues, ok := getExistingDictValues(typeCode)
	if !ok {
		return false
	}
	for _, item := range items {
		if _, exists := existingValues[item.Value]; !exists {
			if !ensureDictData(typeID, typeCode, item.Label, item.Value, item.Sort) {
				return false
			}
		}
	}
	return true
}

// ensureDictType 确保字典类型存在，返回类型ID
func ensureDictType(code, name, remark string) uint {
	var dictType model.DictType
	err := database.GetDB().Where("code = ?", code).First(&dictType).Error
	if err == nil {
		return dictType.ID
	}

	// 创建新类型
	dictType = model.DictType{
		Code:   code,
		Name:   name,
		Remark: remark,
		Status: 1,
	}
	if err := database.GetDB().Create(&dictType).Error; err != nil {
		logging.LogWarn("创建字典类型失败: " + code)
		return 0
	}
	return dictType.ID
}

// ensureDictData 确保字典数据存在
func ensureDictData(typeID uint, typeCode, label, value string, sort int) bool {
	var existing model.DictData
	err := database.GetDB().Where("type_code = ? AND value = ?", typeCode, value).First(&existing).Error
	if err == nil {
		return true // 已存在
	}

	data := model.DictData{
		TypeID:   typeID,
		TypeCode: typeCode,
		Label:    label,
		Value:    value,
		Sort:     sort,
		Status:   1,
	}
	if err := database.GetDB().Create(&data).Error; err != nil {
		logging.LogWarn("创建字典数据失败: " + typeCode + "." + value)
		return false
	}
	return true
}

// getExistingDictValues 批量获取已存在的字典值
func getExistingDictValues(typeCode string) (map[string]bool, bool) {
	var existingData []model.DictData
	if err := database.GetDB().Select("value").Where("type_code = ?", typeCode).Find(&existingData).Error; err != nil {
		logging.LogWarn("读取字典数据失败: " + typeCode)
		return nil, false
	}

	result := make(map[string]bool)
	for _, data := range existingData {
		result[data.Value] = true
	}
	return result, true
}
