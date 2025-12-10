package service

import (
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

type StoreAccountService struct {
	storeAccountModule *module.StoreAccountModule
	productModule      *module.SupplierProductModule
}

func NewStoreAccountService(storeAccountModule *module.StoreAccountModule, productModule *module.SupplierProductModule) *StoreAccountService {
	return &StoreAccountService{
		storeAccountModule: storeAccountModule,
		productModule:      productModule,
	}
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

	// 构建明细
	var items []model.StoreAccountItem
	var totalAmount float64

	for _, item := range req.Items {
		// 获取商品名称
		productName := ""
		unit := item.Unit
		if s.productModule != nil {
			if product, err := s.productModule.GetByID(item.ProductID); err == nil && product != nil {
				productName = product.Name
				if unit == "" {
					unit = product.Unit
				}
			}
		}

		// 计算金额
		amount := item.Amount
		if amount == 0 && item.Price > 0 && item.Quantity > 0 {
			amount = item.Price * item.Quantity
		}

		items = append(items, model.StoreAccountItem{
			ProductID:   item.ProductID,
			ProductName: productName,
			Spec:        item.Spec,
			Quantity:    item.Quantity,
			Unit:        unit,
			Price:       item.Price,
			Amount:      amount,
			Remark:      item.Remark,
		})

		totalAmount += amount
	}

	account := &model.StoreAccount{
		AccountNo:   accountNo,
		StoreID:     storeID,
		Channel:     req.Channel,
		OrderNo:     req.OrderNo,
		TotalAmount: totalAmount,
		ItemCount:   len(items),
		TagCode:     req.TagCode,
		TagName:     req.TagName,
		Remark:      req.Remark,
		OperatorID:  operatorID,
		AccountDate: accountDate,
		Items:       items,
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

	if req.Channel != "" {
		updates["channel"] = req.Channel
	}
	if req.OrderNo != "" {
		updates["order_no"] = req.OrderNo
	}
	if req.TagCode != "" {
		updates["tag_code"] = req.TagCode
	}
	if req.TagName != "" {
		updates["tag_name"] = req.TagName
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
