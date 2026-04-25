package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/pkg/xpyun"
	"github.com/Kevin-Jii/tower-go/utils/logging"
	updatesPkg "github.com/Kevin-Jii/tower-go/utils/updates"
)

type PrinterService struct {
	printerModule       *module.PrinterModule
	storeModule         *module.StoreModule
	purchaseOrderModule *module.PurchaseOrderModule
	xpyunClient         *xpyun.Client
}

func NewPrinterService(printerModule *module.PrinterModule, storeModule *module.StoreModule, purchaseOrderModule *module.PurchaseOrderModule) *PrinterService {
	return &PrinterService{
		printerModule:       printerModule,
		storeModule:         storeModule,
		purchaseOrderModule: purchaseOrderModule,
	}
}

// InitXpyunClient 初始化芯烨云客户端（需要在应用启动时调用）
func (s *PrinterService) InitXpyunClient(user, userKey, baseURL string) {
	s.xpyunClient = xpyun.NewClientWithBaseURL(user, userKey, baseURL)
	fmt.Printf(">>>>>> InitXpyunClient called: user=%s, baseURL=%s, xpyunClient=%v\n", user, baseURL, s.xpyunClient)
}

// BindPrinter 绑定打印机到门店
func (s *PrinterService) BindPrinter(req *model.BindPrinterReq) error {
	// 验证门店存在
	_, err := s.storeModule.GetByID(req.StoreID)
	if err != nil {
		return errors.New("store not found")
	}

	// 检查SN是否已被本地绑定
	existing, err := s.printerModule.GetBySn(req.Sn)
	if err == nil && existing != nil {
		return fmt.Errorf("printer sn %s already bound to store %d", req.Sn, existing.StoreID)
	}

	// 推送到芯烨云（如果未存在会自动添加）
	if s.xpyunClient != nil {
		resp := s.xpyunClient.AddPrinter(req.Sn, req.Name)
		if resp.Content != nil && !resp.Content.IsSuccess() {
			return fmt.Errorf("push to xpyun failed: %s", resp.Content.Msg)
		}
	}

	// 构建打印机名称
	name := req.Name
	if name == "" {
		name = "芯烨云打印机"
	}

	printer := &model.Printer{
		StoreID:   req.StoreID,
		Sn:        req.Sn,
		Name:      name,
		Type:      model.PrinterType(req.Type),
		IsDefault: req.IsDefault,
		Remark:    req.Remark,
		Status:    1,
	}

	// 绑定到数据库
	if err := s.printerModule.BindStore(printer); err != nil {
		return err
	}

	return nil
}

// UnbindPrinter 解绑打印机
func (s *PrinterService) UnbindPrinter(id uint) error {
	printer, err := s.printerModule.GetByID(id)
	if err != nil {
		return errors.New("printer not found")
	}

	// 从芯烨云删除
	if s.xpyunClient != nil {
		resp := s.xpyunClient.DelPrinter([]string{printer.Sn})
		if resp.Content != nil && !resp.Content.IsSuccess() {
			fmt.Printf("remove printer from xpyun failed: %s\n", resp.Content.Msg)
		}
	}

	return s.printerModule.Delete(id)
}

// UpdatePrinter 更新打印机信息
func (s *PrinterService) UpdatePrinter(id uint, req *model.UpdatePrinterReq) error {
	_, err := s.printerModule.GetByID(id)
	if err != nil {
		return errors.New("printer not found")
	}

	updateMap := updatesPkg.BuildUpdatesFromReq(req)
	if len(updateMap) == 0 {
		return nil
	}

	// 如果设置默认，先清除其他默认
	if req.IsDefault != nil && *req.IsDefault == 1 {
		printer, _ := s.printerModule.GetByID(id)
		if printer != nil {
			s.printerModule.ClearDefault(printer.StoreID)
		}
	}

	return s.printerModule.UpdateByID(id, updateMap)
}

// ListPrintersByStore 获取门店下的打印机列表
func (s *PrinterService) ListPrintersByStore(storeID uint) ([]*model.Printer, error) {
	return s.printerModule.ListByStoreID(storeID)
}

// ListAllPrinters 获取所有打印机
func (s *PrinterService) ListAllPrinters() ([]*model.Printer, int64, error) {
	return s.printerModule.ListAll()
}

// GetPrinter 获取打印机详情
func (s *PrinterService) GetPrinter(id uint) (*model.Printer, error) {
	return s.printerModule.GetByID(id)
}

// GetDefaultPrinter 获取门店默认打印机
func (s *PrinterService) GetDefaultPrinter(storeID uint) (*model.Printer, error) {
	return s.printerModule.GetDefaultByStoreID(storeID)
}

// QueryPrinterStatus 查询打印机在线状态
func (s *PrinterService) QueryPrinterStatus(sn string) (int, error) {
	if s.xpyunClient == nil {
		return 0, errors.New("xpyun client not initialized")
	}

	resp := s.xpyunClient.QueryPrinterStatus(sn)
	if resp.Content == nil {
		return 0, errors.New("query failed")
	}

	// 状态码: 0-离线, 1-在线正常, 2-在线异常
	data, ok := resp.Content.Data.(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid response: %v", resp.Content.Data)
	}

	status, _ := data["status"].(float64)
	return int(status), nil
}

// PrintReceipt 打印小票
func (s *PrinterService) PrintReceipt(sn, content string, copies int) (string, error) {
	fmt.Printf(">>>>>> PrintReceipt called: xpyunClient=%v, sn=%s\n", s.xpyunClient, sn)
	if s.xpyunClient == nil {
		return "", errors.New("xpyun client not initialized")
	}

	resp := s.xpyunClient.PrintReceipt(sn, content, copies)
	fmt.Printf(">>>>>> PrintReceipt resp: httpStatus=%d, code=%d, msg=%s, orderId=%s\n",
		resp.HttpStatusCode, resp.Content.Code, resp.Content.Msg, resp.Content.OrderId)

	if resp.Content == nil {
		return "", errors.New("print failed")
	}

	if !resp.Content.IsSuccess() {
		return "", fmt.Errorf("print error: %s", resp.Content.Msg)
	}

	return resp.Content.OrderId, nil
}

// TestPrint 测试打印
func (s *PrinterService) TestPrint(printerID uint, content string, copies int) (string, error) {
	logging.LogInfo(fmt.Sprintf("[TestPrint] printerID=%d, copies=%d", printerID, copies))

	printer, err := s.printerModule.GetByID(printerID)
	if err != nil {
		return "", errors.New("printer not found")
	}

	// 如果没有提供内容，使用默认测试内容
	if content == "" {
		content = "<C>测试打印</C><BR>----------------<BR>这是一张测试小票<BR>打印机: " + printer.Name + "<BR>SN: " + printer.Sn + "<BR>时间: " + time.Now().Format("2006-01-02 15:04:05") + "<BR>----------------<BR>"
	}

	logging.LogInfo(fmt.Sprintf("[TestPrint] printer found, sn=%s, calling PrintReceipt...", printer.Sn))

	return s.PrintReceipt(printer.Sn, content, copies)
}

// GetPrinterWithStatus 获取打印机信息及在线状态
func (s *PrinterService) GetPrinterWithStatus(id uint) (*model.PrinterResp, error) {
	printer, err := s.printerModule.GetByID(id)
	if err != nil {
		return nil, errors.New("printer not found")
	}

	resp := &model.PrinterResp{
		ID:            printer.ID,
		StoreID:       printer.StoreID,
		Sn:            printer.Sn,
		Name:          printer.Name,
		Type:          int(printer.Type),
		Status:        printer.Status,
		IsDefault:     printer.IsDefault,
		Online:        printer.Online,
		LastHeartbeat: printer.LastHeartbeat,
		Remark:        printer.Remark,
		CreatedAt:     printer.CreatedAt,
	}

	// 设置类型名称
	if printer.Type == model.PrinterTypeReceipt {
		resp.TypeName = "小票打印机"
	} else {
		resp.TypeName = "标签打印机"
	}

	// 设置状态名称
	if printer.Status == 1 {
		resp.StatusName = "正常"
	} else {
		resp.StatusName = "停用"
	}

	// 获取门店名称
	store, err := s.storeModule.GetByID(printer.StoreID)
	if err == nil {
		resp.StoreName = store.Name
	}

	return resp, nil
}

// BatchQueryStatus 批量查询打印机状态
func (s *PrinterService) BatchQueryStatus(storeID uint) ([]*model.PrinterStatus, error) {
	printers, err := s.printerModule.ListByStoreID(storeID)
	if err != nil {
		return nil, err
	}

	var results []*model.PrinterStatus
	for _, p := range printers {
		status := &model.PrinterStatus{
			Sn:     p.Sn,
			Online: 0,
		}

		if s.xpyunClient != nil {
			resp := s.xpyunClient.QueryPrinterStatus(p.Sn)
			if resp.Content != nil && resp.Content.IsSuccess() {
				if data, ok := resp.Content.Data.(map[string]interface{}); ok {
					if s, ok := data["status"].(float64); ok {
						status.Online = int(s)
					}
				}
			}
		}

		results = append(results, status)
	}

	return results, nil
}

// SyncAllPrinterStatus 同步所有打印机状态（定时任务调用）
func (s *PrinterService) SyncAllPrinterStatus() error {
	if s.xpyunClient == nil {
		return errors.New("xpyun client not initialized")
	}

	total, err := s.printerModule.CountAll()
	if err != nil {
		return err
	}
	if total == 0 {
		// 没有任何打印机，跳过同步（避免无意义的定时任务执行）
		return nil
	}

	sns, err := s.printerModule.ListAllSn()
	if err != nil {
		return err
	}

	statuses := make(map[string]int)
	for _, sn := range sns {
		resp := s.xpyunClient.QueryPrinterStatus(sn)
		if resp.Content != nil && resp.Content.IsSuccess() {
			if data, ok := resp.Content.Data.(map[string]interface{}); ok {
				if status, ok := data["status"].(float64); ok {
					statuses[sn] = int(status)
				}
			}
		} else {
			// 查询失败默认为离线
			statuses[sn] = 0
		}
	}

	return s.printerModule.BatchUpdateOnlineStatus(statuses)
}

// PrintPurchaseOrder 打印采购单
func (s *PrinterService) PrintPurchaseOrder(printerID uint, orderID uint) (string, error) {
	logging.LogInfo(fmt.Sprintf("[PrintPurchaseOrder] printerID=%d, orderID=%d", printerID, orderID))

	// 获取打印机信息
	printer, err := s.printerModule.GetByID(printerID)
	if err != nil {
		return "", errors.New("printer not found")
	}

	// 获取采购单完整信息
	order, err := s.purchaseOrderModule.GetByIDWithDetails(orderID)
	if err != nil {
		return "", fmt.Errorf("purchase order not found: %v", err)
	}

	// 构建打印内容
	content := s.buildPurchaseOrderContent(order)

	logging.LogInfo(fmt.Sprintf("[PrintPurchaseOrder] printer found, sn=%s, calling PrintReceipt...", printer.Sn))

	return s.PrintReceipt(printer.Sn, content, 1)
}

// buildPurchaseOrderContent 构建采购单打印内容
func (s *PrinterService) buildPurchaseOrderContent(order *model.PurchaseOrder) string {
	var content string

	// 标题
	content += "<C><B>采购单</B></C><BR>"
	content += "================================<BR>"

	// 基本信息
	content += fmt.Sprintf("单号: %s<BR>", order.OrderNo)
	if order.Store != nil {
		content += fmt.Sprintf("门店: %s<BR>", order.Store.Name)
	}
	content += fmt.Sprintf("日期: %s<BR>", order.OrderDate.Format("2006-01-02"))
	if order.Creator != nil {
		content += fmt.Sprintf("制单人: %s<BR>", order.Creator.Username)
	}
	content += fmt.Sprintf("打印时间: %s<BR>", time.Now().Format("2006-01-02 15:04:05"))
	content += "================================<BR>"
	content += "<BR>"

	// 按供应商分组打印明细
	supplierItems := make(map[uint][]model.PurchaseOrderItem)
	for _, item := range order.Items {
		supplierItems[item.SupplierID] = append(supplierItems[item.SupplierID], item)
	}

	// 打印每个供应商的商品
	for _, items := range supplierItems {
		if len(items) > 0 && items[0].Supplier != nil {
			content += fmt.Sprintf("<B>%s</B><BR>", items[0].Supplier.SupplierName)
		}
		content += "--------------------------------<BR>"
		content += padRight("商品名称", 12) + padCenter("规格", 10) + "   单价<BR>"
		content += "--------------------------------<BR>"

		supplierTotal := 0.0
		for _, item := range items {
			productName := "未知商品"
			if item.Product != nil {
				productName = item.Product.Name
			}

			// 格式化商品名称（最多12个字符，中文算2个字符）
			displayName := truncateString(productName, 12)

			unit := ""
			if item.Product != nil {
				unit = item.Product.Unit
			}
			spec := fmt.Sprintf("%s*%s", strconv.FormatFloat(item.Quantity, 'f', -1, 64), unit)
			content += fmt.Sprintf("%s%s%7.2f<BR>",
				padRight(displayName, 12),
				padCenter(spec, 10),
				item.UnitPrice)

			if item.Remark != "" {
				content += fmt.Sprintf("  备注: %s<BR>", item.Remark)
			}

			supplierTotal += item.Amount
		}

		content += "--------------------------------<BR>"
		content += fmt.Sprintf("<R>小计: %.2f元</R><BR>", supplierTotal)
		content += "<BR>"
	}

	// 总计
	content += "================================<BR>"
	content += fmt.Sprintf("<R><B>合计: %.2f元</B></R><BR>", order.TotalAmount)
	content += "================================<BR>"

	// 备注
	if order.Remark != "" {
		content += "<BR>"
		content += fmt.Sprintf("备注: %s<BR>", order.Remark)
	}

	// 切纸
	content += "<CUT>"

	return content
}

// truncateString 截断字符串到指定长度（中文算2个字符）
func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	length := 0
	result := []rune{}

	for _, r := range runes {
		charLen := 1
		if r > 127 { // 中文字符
			charLen = 2
		}

		if length+charLen > maxLen {
			break
		}

		result = append(result, r)
		length += charLen
	}

	return string(result)
}

// displayWidth 计算字符串显示宽度（中文算2个字符）
func displayWidth(s string) int {
	w := 0
	for _, r := range s {
		if r > 127 {
			w += 2
		} else {
			w++
		}
	}
	return w
}

// padRight 按显示宽度右填充空格
func padRight(s string, width int) string {
	w := displayWidth(s)
	if w >= width {
		return s
	}
	return s + strings.Repeat(" ", width-w)
}

// padCenter 按显示宽度居中填充空格
func padCenter(s string, width int) string {
	w := displayWidth(s)
	if w >= width {
		return s
	}
	total := width - w
	left := total / 2
	right := total - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}
