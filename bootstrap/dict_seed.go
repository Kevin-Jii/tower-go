package bootstrap

import (
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils/database"
	"github.com/Kevin-Jii/tower-go/utils/logging"
)

// InitDefaultDicts 初始化默认字典数据
func InitDefaultDicts() {
	initSalesChannel()
	initOrderSource()
	initInventoryReason()
	logging.LogInfo("字典数据初始化完成")
}

// initSalesChannel 初始化销售渠道字典
func initSalesChannel() {
	typeCode := "sales_channel"
	typeName := "销售渠道"

	// 创建或获取字典类型
	typeID := ensureDictType(typeCode, typeName, "门店记账-销售渠道")
	if typeID == 0 {
		return
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
		{"其他", "other", 99},
	}

	for _, item := range items {
		ensureDictData(typeID, typeCode, item.Label, item.Value, item.Sort)
	}
}

// initOrderSource 初始化订单来源字典
func initOrderSource() {
	typeCode := "order_source"
	typeName := "订单来源"

	typeID := ensureDictType(typeCode, typeName, "门店记账-订单来源")
	if typeID == 0 {
		return
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

	for _, item := range items {
		ensureDictData(typeID, typeCode, item.Label, item.Value, item.Sort)
	}
}

// initInventoryReason 初始化出入库原因字典
func initInventoryReason() {
	typeCode := "inventory_reason"
	typeName := "出入库原因"

	typeID := ensureDictType(typeCode, typeName, "库存管理-出入库原因")
	if typeID == 0 {
		return
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

	for _, item := range items {
		ensureDictData(typeID, typeCode, item.Label, item.Value, item.Sort)
	}
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
func ensureDictData(typeID uint, typeCode, label, value string, sort int) {
	var count int64
	database.GetDB().Model(&model.DictData{}).Where("type_code = ? AND value = ?", typeCode, value).Count(&count)
	if count > 0 {
		return // 已存在
	}

	data := model.DictData{
		TypeID:   typeID,
		TypeCode: typeCode,
		Label:    label,
		Value:    value,
		Sort:     sort,
		Status:   1,
	}
	database.GetDB().Create(&data)
}
