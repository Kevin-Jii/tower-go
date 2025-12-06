package service

import (
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

type StoreAccountService struct {
	storeAccountModule *module.StoreAccountModule
}

func NewStoreAccountService(storeAccountModule *module.StoreAccountModule) *StoreAccountService {
	return &StoreAccountService{storeAccountModule: storeAccountModule}
}

// Create 创建记账
func (s *StoreAccountService) Create(storeID, operatorID uint, req *model.CreateStoreAccountReq) (*model.StoreAccount, error) {
	accountNo := s.storeAccountModule.GenerateAccountNo()

	// 解析记账日期
	accountDate := time.Now()
	if req.AccountDate != "" {
		if t, err := time.Parse("2006-01-02", req.AccountDate); err == nil {
			accountDate = t
		}
	}

	// 计算金额
	amount := req.Amount
	if amount == 0 && req.Price > 0 && req.Quantity > 0 {
		amount = req.Price * req.Quantity
	}

	account := &model.StoreAccount{
		AccountNo:   accountNo,
		StoreID:     storeID,
		ProductID:   req.ProductID,
		Spec:        req.Spec,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
		Price:       req.Price,
		Amount:      amount,
		Channel:     req.Channel,
		OrderSource: req.OrderSource,
		OrderNo:     req.OrderNo,
		Remark:      req.Remark,
		OperatorID:  operatorID,
		AccountDate: accountDate,
	}

	if err := s.storeAccountModule.Create(account); err != nil {
		return nil, err
	}

	return account, nil
}

// Get 获取记账详情
func (s *StoreAccountService) Get(id uint) (*model.StoreAccount, error) {
	return s.storeAccountModule.GetByID(id)
}

// List 记账列表
func (s *StoreAccountService) List(req *model.ListStoreAccountReq) ([]*model.StoreAccount, int64, error) {
	return s.storeAccountModule.List(req)
}

// Update 更新记账
func (s *StoreAccountService) Update(id uint, req *model.UpdateStoreAccountReq) error {
	updates := make(map[string]interface{})

	if req.ProductID != nil {
		updates["product_id"] = *req.ProductID
	}
	if req.Spec != "" {
		updates["spec"] = req.Spec
	}
	if req.Quantity != nil {
		updates["quantity"] = *req.Quantity
	}
	if req.Unit != "" {
		updates["unit"] = req.Unit
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}
	if req.Channel != "" {
		updates["channel"] = req.Channel
	}
	if req.OrderSource != "" {
		updates["order_source"] = req.OrderSource
	}
	if req.OrderNo != "" {
		updates["order_no"] = req.OrderNo
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	if req.AccountDate != "" {
		if t, err := time.Parse("2006-01-02", req.AccountDate); err == nil {
			updates["account_date"] = t
		}
	}

	if len(updates) == 0 {
		return nil
	}

	return s.storeAccountModule.Update(id, updates)
}

// Delete 删除记账
func (s *StoreAccountService) Delete(id uint) error {
	return s.storeAccountModule.Delete(id)
}

// GetStats 获取统计
func (s *StoreAccountService) GetStats(storeID uint, startDate, endDate string) (map[string]interface{}, error) {
	totalAmount, count, err := s.storeAccountModule.GetStatsByDateRange(storeID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_amount": totalAmount,
		"count":        count,
	}, nil
}
