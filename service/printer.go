package service

import (
	"errors"
	"fmt"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/pkg/xpyun"
	updatesPkg "github.com/Kevin-Jii/tower-go/utils/updates"
)

type PrinterService struct {
	printerModule *module.PrinterModule
	storeModule   *module.StoreModule
	xpyunClient   *xpyun.Client
}

func NewPrinterService(printerModule *module.PrinterModule, storeModule *module.StoreModule) *PrinterService {
	return &PrinterService{
		printerModule: printerModule,
		storeModule:   storeModule,
	}
}

// InitXpyunClient 初始化芯烨云客户端（需要在应用启动时调用）
func (s *PrinterService) InitXpyunClient(user, userKey string) {
	s.xpyunClient = xpyun.NewClient(user, userKey)
}

// BindPrinter 绑定打印机到门店
func (s *PrinterService) BindPrinter(req *model.BindPrinterReq) error {
	// 验证门店存在
	_, err := s.storeModule.GetByID(req.StoreID)
	if err != nil {
		return errors.New("store not found")
	}

	// 检查SN是否已被绑定
	existing, err := s.printerModule.GetBySn(req.Sn)
	if err == nil && existing != nil {
		return fmt.Errorf("printer sn %s already bound to store %d", req.Sn, existing.StoreID)
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

	// 推送到芯烨云
	if s.xpyunClient != nil {
		resp := s.xpyunClient.AddPrinter(req.Sn, name)
		if resp.Content != nil && !resp.Content.IsSuccess() {
			// 推送失败，但本地已绑定，记录日志即可
			fmt.Printf("push printer to xpyun failed: %s\n", resp.Content.Msg)
		}
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
	if s.xpyunClient == nil {
		return "", errors.New("xpyun client not initialized")
	}

	resp := s.xpyunClient.PrintReceipt(sn, content, copies)
	if resp.Content == nil {
		return "", errors.New("print failed")
	}

	if !resp.Content.IsSuccess() {
		return "", fmt.Errorf("print error: %s", resp.Content.Msg)
	}

	return resp.Content.OrderId, nil
}

// GetPrinterWithStatus 获取打印机信息及在线状态
func (s *PrinterService) GetPrinterWithStatus(id uint) (*model.PrinterResp, error) {
	printer, err := s.printerModule.GetByID(id)
	if err != nil {
		return nil, errors.New("printer not found")
	}

	resp := &model.PrinterResp{
		ID:         printer.ID,
		StoreID:    printer.StoreID,
		Sn:         printer.Sn,
		Name:       printer.Name,
		Type:       int(printer.Type),
		Status:     printer.Status,
		IsDefault:  printer.IsDefault,
		Remark:     printer.Remark,
		CreatedAt:  printer.CreatedAt,
		Online:     0, // 默认离线
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

	// 查询在线状态
	if s.xpyunClient != nil {
		statusResp := s.xpyunClient.QueryPrinterStatus(printer.Sn)
		if statusResp.Content != nil && statusResp.Content.IsSuccess() {
			if data, ok := statusResp.Content.Data.(map[string]interface{}); ok {
				if status, ok := data["status"].(float64); ok {
					resp.Online = int(status)
				}
			}
		}
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