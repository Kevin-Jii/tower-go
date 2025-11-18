package service

import (
	"io"
	"time"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/utils"
	"github.com/Kevin-Jii/tower-go/utils/events"
)

type MenuReportService struct {
	menuReportModule *module.MenuReportModule
	dishModule       *module.DishModule
	storeModule      *module.StoreModule
	userModule       *module.UserModule
	botModule        *module.DingTalkBotModule
	eventBus         *events.EventBus
}

func NewMenuReportService(
	menuReportModule *module.MenuReportModule,
	dishModule *module.DishModule,
	storeModule *module.StoreModule,
	userModule *module.UserModule,
	botModule *module.DingTalkBotModule,
	eventBus *events.EventBus,
) *MenuReportService {
	return &MenuReportService{
		menuReportModule: menuReportModule,
		dishModule:       dishModule,
		storeModule:      storeModule,
		userModule:       userModule,
		botModule:        botModule,
		eventBus:         eventBus,
	}
}

// CreateMenuReportOrder 创建报菜记录单（包含详情）
func (s *MenuReportService) CreateMenuReportOrder(storeID uint, userID uint, req *model.CreateMenuReportOrderReq) (*model.MenuReportOrder, error) {
	// 验证所有菜品是否存在且属于当前门店
	for _, item := range req.Items {
		if _, err := s.dishModule.GetByID(item.DishID, storeID); err != nil {
			return nil, err
		}
	}

	// 创建报菜记录单和详情
	order := &model.MenuReportOrder{
		StoreID: storeID,
		UserID:  userID,
		Remark:  req.Remark,
	}

	// 创建详情列表
	items := make([]*model.MenuReportItem, len(req.Items))
	for i, reqItem := range req.Items {
		items[i] = &model.MenuReportItem{
			DishID:   reqItem.DishID,
			Quantity: reqItem.Quantity,
			Remark:   reqItem.Remark,
		}
	}
	order.Items = items

	if err := s.menuReportModule.CreateOrder(order); err != nil {
		return nil, err
	}

	// 发布报菜创建事件（异步）
	if s.eventBus != nil {
		// 获取门店和用户信息
		storeName := ""
		userName := ""
		storePhone := ""
		storeAddress := ""

		if store, err := s.storeModule.GetByID(storeID); err == nil && store != nil {
			storeName = store.Name
			storePhone = store.Phone
			storeAddress = store.Address
		}

		if user, err := s.userModule.GetByID(userID); err == nil && user != nil {
			userName = user.Username
		}

		// 获取门店的第一个启用的机器人
		var botID uint = 0
		if bots, err := s.botModule.ListEnabledByStoreID(storeID); err == nil && len(bots) > 0 {
			botID = bots[0].ID
		}

		event := MenuReportOrderCreatedEvent{
			Order:        order,
			StoreName:    storeName,
			UserName:     userName,
			StorePhone:   storePhone,
			StoreAddress: storeAddress,
			BotID:        botID, // 使用门店的第一个启用的机器人
		}

		s.eventBus.PublishAsync(event)
	}

	return order, nil
}

// GetMenuReportOrder 获取报菜记录单详情
func (s *MenuReportService) GetMenuReportOrder(id uint, storeID uint) (*model.MenuReportOrder, error) {
	return s.menuReportModule.GetOrderByID(id, storeID)
}

// ListMenuReportOrders 获取报菜记录单列表
func (s *MenuReportService) ListMenuReportOrders(storeID uint, page, pageSize int) ([]*model.MenuReportOrder, int64, error) {
	return s.menuReportModule.ListOrders(storeID, page, pageSize)
}

// ListMenuReportOrdersByDateRange 根据日期范围查询
func (s *MenuReportService) ListMenuReportOrdersByDateRange(storeID uint, startDate, endDate time.Time) ([]*model.MenuReportOrder, error) {
	return s.menuReportModule.ListOrdersByDateRange(storeID, startDate, endDate)
}

// UpdateMenuReportOrder 更新报菜记录单
func (s *MenuReportService) UpdateMenuReportOrder(id uint, storeID uint, req *model.UpdateMenuReportOrderReq) error {
	// 验证记录是否存在
	_, err := s.menuReportModule.GetOrderByID(id, storeID)
	if err != nil {
		return err
	}

	return s.menuReportModule.UpdateOrder(id, storeID, req.Remark)
}

// DeleteMenuReportOrder 删除报菜记录单
func (s *MenuReportService) DeleteMenuReportOrder(id uint, storeID uint) error {
	return s.menuReportModule.DeleteOrder(id, storeID)
}

// AddMenuReportItem 添加报菜详情项
func (s *MenuReportService) AddMenuReportItem(orderID uint, storeID uint, req *model.MenuReportItemReq) (*model.MenuReportItem, error) {
	// 验证订单存在且属于当前门店
	_, err := s.menuReportModule.GetOrderByID(orderID, storeID)
	if err != nil {
		return nil, err
	}

	// 验证菜品是否存在且属于当前门店
	if _, err := s.dishModule.GetByID(req.DishID, storeID); err != nil {
		return nil, err
	}

	item := &model.MenuReportItem{
		ReportOrderID: orderID,
		DishID:        req.DishID,
		Quantity:      req.Quantity,
		Remark:        req.Remark,
	}

	if err := s.menuReportModule.AddItem(item); err != nil {
		return nil, err
	}

	return item, nil
}

// UpdateMenuReportItem 更新报菜详情项
func (s *MenuReportService) UpdateMenuReportItem(id uint, orderID uint, storeID uint, dishID *uint, quantity *int, remark *string) error {
	// 验证订单存在且属于当前门店
	_, err := s.menuReportModule.GetOrderByID(orderID, storeID)
	if err != nil {
		return err
	}

	// 如果更新了菜品，验证新菜品是否存在
	if dishID != nil {
		if _, err := s.dishModule.GetByID(*dishID, storeID); err != nil {
			return err
		}
	}

	updates := make(map[string]interface{})
	if dishID != nil {
		updates["dish_id"] = *dishID
	}
	if quantity != nil {
		updates["quantity"] = *quantity
	}
	if remark != nil {
		updates["remark"] = *remark
	}

	return s.menuReportModule.UpdateItem(id, updates)
}

// DeleteMenuReportItem 删除报菜详情项
func (s *MenuReportService) DeleteMenuReportItem(id uint, orderID uint, storeID uint) error {
	// 验证订单存在且属于当前门店
	_, err := s.menuReportModule.GetOrderByID(orderID, storeID)
	if err != nil {
		return err
	}

	return s.menuReportModule.DeleteItem(id)
}

// GetStatsByDateRange 获取统计数据
func (s *MenuReportService) GetStatsByDateRange(storeID uint, startDate, endDate time.Time) (*model.MenuReportStats, error) {
	return s.menuReportModule.GetStatsByDateRange(storeID, startDate, endDate)
}

// GetStatsByDateRangeAllStores 获取所有门店统计数据（仅总部）
func (s *MenuReportService) GetStatsByDateRangeAllStores(startDate, endDate time.Time) (*model.MenuReportStats, error) {
	return s.menuReportModule.GetStatsByDateRangeAllStores(startDate, endDate)
}

// GenerateExcel 生成报菜记录单Excel文件
func (s *MenuReportService) GenerateExcel(order *model.MenuReportOrder) (io.Reader, error) {
	// 获取门店和用户信息
	storeName := ""
	userName := ""
	storePhone := ""
	storeAddress := ""

	// 优先使用预加载的数据，避免额外查询
	if order.Store != nil {
		storeName = order.Store.Name
		storePhone = order.Store.Phone
		storeAddress = order.Store.Address
	} else if order.StoreID > 0 {
		// 如果没有预加载，则查询
		if store, err := s.storeModule.GetByID(order.StoreID); err == nil && store != nil {
			storeName = store.Name
			storePhone = store.Phone
			storeAddress = store.Address
		}
	}

	if order.User != nil {
		userName = order.User.Username
	} else if order.UserID > 0 {
		if user, err := s.userModule.GetByID(order.UserID); err == nil && user != nil {
			userName = user.Username
		}
	}

	// 调用工具函数生成Excel
	return utils.GenerateMenuReportExcel(order, storeName, userName, storePhone, storeAddress)
}
