package service

import (
	"time"
	"tower-go/model"
	"tower-go/module"
	"tower-go/utils"
)

type MenuReportService struct {
	menuReportModule *module.MenuReportModule
	dishModule       *module.DishModule
	storeModule      *module.StoreModule
	userModule       *module.UserModule
	eventBus         *utils.EventBus
}

func NewMenuReportService(
	menuReportModule *module.MenuReportModule,
	dishModule *module.DishModule,
	storeModule *module.StoreModule,
	userModule *module.UserModule,
	eventBus *utils.EventBus,
) *MenuReportService {
	return &MenuReportService{
		menuReportModule: menuReportModule,
		dishModule:       dishModule,
		storeModule:      storeModule,
		userModule:       userModule,
		eventBus:         eventBus,
	}
}

// CreateMenuReport 创建报菜记录
func (s *MenuReportService) CreateMenuReport(storeID uint, userID uint, req *model.CreateMenuReportReq) (*model.MenuReport, error) {
	// 验证菜品是否存在且属于当前门店
	dish, err := s.dishModule.GetByID(req.DishID, storeID)
	if err != nil {
		return nil, err
	}

	report := &model.MenuReport{
		StoreID:  storeID,
		DishID:   req.DishID,
		UserID:   userID,
		Quantity: req.Quantity,
		Remark:   req.Remark,
	}

	if err := s.menuReportModule.Create(report); err != nil {
		return nil, err
	}

	// 发布报菜创建事件（异步）
	if s.eventBus != nil {
		// 获取门店和用户信息
		storeName := ""
		userName := ""

		if store, err := s.storeModule.GetByID(storeID); err == nil && store != nil {
			storeName = store.Name
		}

		if user, err := s.userModule.GetByID(userID); err == nil && user != nil {
			userName = user.Username
		}

		event := MenuReportCreatedEvent{
			Report:    report,
			StoreName: storeName,
			DishName:  dish.Name,
			UserName:  userName,
		}

		s.eventBus.PublishAsync(event)
	}

	return report, nil
}

// GetMenuReport 获取报菜记录详情
func (s *MenuReportService) GetMenuReport(id uint, storeID uint) (*model.MenuReport, error) {
	return s.menuReportModule.GetByID(id, storeID)
}

// ListMenuReports 获取报菜记录列表
func (s *MenuReportService) ListMenuReports(storeID uint, page, pageSize int) ([]*model.MenuReport, int64, error) {
	return s.menuReportModule.List(storeID, page, pageSize)
}

// ListMenuReportsByDateRange 根据日期范围查询
func (s *MenuReportService) ListMenuReportsByDateRange(storeID uint, startDate, endDate time.Time) ([]*model.MenuReport, error) {
	return s.menuReportModule.ListByDateRange(storeID, startDate, endDate)
}

// UpdateMenuReport 更新报菜记录
func (s *MenuReportService) UpdateMenuReport(id uint, storeID uint, req *model.UpdateMenuReportReq) error {
	// 验证记录是否存在
	_, err := s.menuReportModule.GetByID(id, storeID)
	if err != nil {
		return err
	}

	return s.menuReportModule.Update(id, storeID, req)
}

// DeleteMenuReport 删除报菜记录
func (s *MenuReportService) DeleteMenuReport(id uint, storeID uint) error {
	return s.menuReportModule.Delete(id, storeID)
}

// GetStatsByDateRange 获取统计数据
func (s *MenuReportService) GetStatsByDateRange(storeID uint, startDate, endDate time.Time) (*model.MenuReportStats, error) {
	return s.menuReportModule.GetStatsByDateRange(storeID, startDate, endDate)
}

// GetStatsByDateRangeAllStores 获取所有门店统计数据（仅总部）
func (s *MenuReportService) GetStatsByDateRangeAllStores(startDate, endDate time.Time) (*model.MenuReportStats, error) {
	return s.menuReportModule.GetStatsByDateRangeAllStores(startDate, endDate)
}
