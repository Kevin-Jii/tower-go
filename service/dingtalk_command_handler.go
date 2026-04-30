package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
)

// DingTalkCommandHandler 钉钉命令处理器
type DingTalkCommandHandler struct {
	inventoryModule    *module.InventoryModule
	storeAccountModule *module.StoreAccountModule
	storeModule        *module.StoreModule
	userModule         *module.UserModule
	templateService    *MessageTemplateService
}

var globalCommandHandler *DingTalkCommandHandler

// InitCommandHandler 初始化命令处理器
func InitCommandHandler(
	inventoryModule *module.InventoryModule,
	storeAccountModule *module.StoreAccountModule,
	storeModule *module.StoreModule,
	userModule *module.UserModule,
	templateService *MessageTemplateService,
) {
	globalCommandHandler = &DingTalkCommandHandler{
		inventoryModule:    inventoryModule,
		storeAccountModule: storeAccountModule,
		storeModule:        storeModule,
		userModule:         userModule,
		templateService:    templateService,
	}
}

// GetCommandHandler 获取命令处理器
func GetCommandHandler() *DingTalkCommandHandler {
	return globalCommandHandler
}

// HandleCommand 处理用户命令
func (h *DingTalkCommandHandler) HandleCommand(ctx context.Context, data *chatbot.BotCallbackDataModel) (string, string) {
	if h == nil {
		return "系统提示", "命令处理器未初始化"
	}

	content := strings.TrimSpace(data.Text.Content)

	// 获取用户所属门店
	storeID := h.getStoreIDByStaffId(data.SenderStaffId)

	switch {
	case content == "帮助" || content == "菜单" || content == "help":
		return h.handleHelp()

	case content == "库存查询" || content == "查库存":
		return h.handleInventoryQuery(storeID)

	case content == "今日记账" || content == "记账查询":
		return h.handleTodayAccount(storeID)

	case content == "今日入库" || content == "入库查询":
		return h.handleTodayInventoryIn(storeID)

	case strings.HasPrefix(content, "查询库存 "):
		keyword := strings.TrimPrefix(content, "查询库存 ")
		return h.handleInventorySearch(storeID, keyword)

	default:
		return h.handleUnknown(content)
	}
}

// getStoreIDByStaffId 根据钉钉用户ID获取门店ID
func (h *DingTalkCommandHandler) getStoreIDByStaffId(staffId string) uint {
	if h.userModule == nil {
		return 999 // 默认总部
	}
	// 通过钉钉用户ID查找用户
	user, err := h.userModule.GetByDingTalkID(staffId)
	if err != nil || user == nil {
		return 999 // 默认总部
	}
	return user.StoreID
}

// handleHelp 帮助菜单
func (h *DingTalkCommandHandler) handleHelp() (string, string) {
	if h.templateService != nil {
		title, text, err := h.templateService.RenderTemplate(model.TemplateBotHelp, nil)
		if err == nil && text != "" {
			return title, text
		}
	}
	// 默认回复
	return "📋 功能菜单", "发送 库存查询、今日记账、今日入库 等命令"
}

// handleInventoryQuery 库存查询
func (h *DingTalkCommandHandler) handleInventoryQuery(storeID uint) (string, string) {
	if h.inventoryModule == nil {
		return "库存查询", "库存模块未初始化"
	}

	req := &model.ListInventoryReq{
		StoreID:  storeID,
		Page:     1,
		PageSize: 10,
	}
	list, total, err := h.inventoryModule.List(req)
	if err != nil {
		return "库存查询", fmt.Sprintf("查询失败: %v", err)
	}

	if total == 0 {
		return "库存查询", "暂无库存数据"
	}

	var lines []string
	for i, item := range list {
		lines = append(lines, fmt.Sprintf("%d. %s: %.2f%s", i+1, item.ProductName, item.Quantity, item.Unit))
	}

	data := map[string]interface{}{
		"Total":      total,
		"ItemList":   strings.Join(lines, "\n\n"),
		"CreateTime": time.Now().Format("2006-01-02 15:04:05"),
	}

	if h.templateService != nil {
		title, text, err := h.templateService.RenderTemplate(model.TemplateBotInventoryQuery, data)
		if err == nil && text != "" {
			return title, text
		}
	}
	return "📦 库存查询", fmt.Sprintf("共%d项\n\n%s", total, strings.Join(lines, "\n"))
}

// handleTodayAccount 今日记账查询
func (h *DingTalkCommandHandler) handleTodayAccount(storeID uint) (string, string) {
	if h.storeAccountModule == nil {
		return "记账查询", "记账模块未初始化"
	}

	today := time.Now().Format("2006-01-02")
	totalAmount, netIncomeAmount, count, err := h.storeAccountModule.GetStatsByDateRange(storeID, today, today)
	if err != nil {
		return "记账查询", fmt.Sprintf("查询失败: %v", err)
	}

	data := map[string]interface{}{
		"Date":        today,
		"Count":       count,
		"TotalAmount": fmt.Sprintf("%.2f", totalAmount),
		"NetIncome":   fmt.Sprintf("%.2f", netIncomeAmount),
		"CreateTime":  time.Now().Format("2006-01-02 15:04:05"),
	}

	if h.templateService != nil {
		title, text, err := h.templateService.RenderTemplate(model.TemplateBotTodayAccount, data)
		if err == nil && text != "" {
			return title, text
		}
	}
	return "📝 今日记账", fmt.Sprintf("日期: %s\n笔数: %d\n总销售额: ¥%.2f\n净利润: ¥%.2f", today, count, totalAmount, netIncomeAmount)
}

// handleTodayInventoryIn 今日入库查询
func (h *DingTalkCommandHandler) handleTodayInventoryIn(storeID uint) (string, string) {
	if h.inventoryModule == nil {
		return "入库查询", "库存模块未初始化"
	}

	today := time.Now().Format("2006-01-02")
	orderType := int8(1) // 入库
	req := &model.ListInventoryOrderReq{
		StoreID:  storeID,
		Type:     &orderType,
		Date:     today,
		Page:     1,
		PageSize: 10,
	}
	list, total, err := h.inventoryModule.ListOrders(req)
	if err != nil {
		return "入库查询", fmt.Sprintf("查询失败: %v", err)
	}

	if total == 0 {
		return "今日入库", fmt.Sprintf("日期: %s\n暂无入库记录", today)
	}

	var lines []string
	var totalQty float64
	for _, order := range list {
		lines = append(lines, fmt.Sprintf("- %s: %.2f件 (%s)", order.OrderNo, order.TotalQuantity, order.Reason))
		totalQty += order.TotalQuantity
	}

	data := map[string]interface{}{
		"Date":          today,
		"Count":         total,
		"TotalQuantity": fmt.Sprintf("%.2f", totalQty),
		"ItemList":      strings.Join(lines, "\n\n"),
		"CreateTime":    time.Now().Format("2006-01-02 15:04:05"),
	}

	if h.templateService != nil {
		title, text, err := h.templateService.RenderTemplate(model.TemplateBotTodayInventory, data)
		if err == nil && text != "" {
			return title, text
		}
	}
	return "📦 今日入库", fmt.Sprintf("日期: %s\n单数: %d\n总量: %.2f", today, total, totalQty)
}

// handleInventorySearch 搜索库存
func (h *DingTalkCommandHandler) handleInventorySearch(storeID uint, keyword string) (string, string) {
	if h.inventoryModule == nil {
		return "库存搜索", "库存模块未初始化"
	}

	req := &model.ListInventoryReq{
		StoreID:  storeID,
		Keyword:  keyword,
		Page:     1,
		PageSize: 10,
	}
	list, total, err := h.inventoryModule.List(req)
	if err != nil {
		return "库存搜索", fmt.Sprintf("查询失败: %v", err)
	}

	if total == 0 {
		return "库存搜索", fmt.Sprintf("未找到包含「%s」的商品", keyword)
	}

	var lines []string
	for i, item := range list {
		lines = append(lines, fmt.Sprintf("%d. %s: %.2f%s", i+1, item.ProductName, item.Quantity, item.Unit))
	}

	data := map[string]interface{}{
		"Keyword":    keyword,
		"Total":      total,
		"ItemList":   strings.Join(lines, "\n\n"),
		"CreateTime": time.Now().Format("2006-01-02 15:04:05"),
	}

	if h.templateService != nil {
		title, text, err := h.templateService.RenderTemplate(model.TemplateBotSearchResult, data)
		if err == nil && text != "" {
			return title, text
		}
	}
	return "🔍 库存搜索", fmt.Sprintf("关键词: %s\n共%d项\n\n%s", keyword, total, strings.Join(lines, "\n"))
}

// handleUnknown 未知命令
func (h *DingTalkCommandHandler) handleUnknown(content string) (string, string) {
	data := map[string]interface{}{
		"Content": content,
	}

	if h.templateService != nil {
		title, text, err := h.templateService.RenderTemplate(model.TemplateBotUnknown, data)
		if err == nil && text != "" {
			return title, text
		}
	}
	return "🤖 智能助手", fmt.Sprintf("您发送: %s\n\n发送 帮助 查看可用功能", content)
}
