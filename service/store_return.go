package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

type StoreReturnService struct {
	returnModule *module.StoreReturnModule
	userModule   *module.UserModule
}

func NewStoreReturnService(returnModule *module.StoreReturnModule, userModule *module.UserModule) *StoreReturnService {
	return &StoreReturnService{returnModule: returnModule, userModule: userModule}
}

func (s *StoreReturnService) Create(storeID, operatorID uint, req *model.CreateStoreReturnReq, hqUnbound bool) (*model.StoreReturn, error) {
	record, err := s.buildRecord(storeID, operatorID, hqUnbound, req.StoreID, req.ReturnDate, req.LogisticsFee, req.Remark, req.Items)
	if err != nil {
		return nil, err
	}
	record.ClientReqID = strings.TrimSpace(req.ClientReqID)
	if existing, err := s.returnModule.GetByClientReqIDScoped(record.ClientReqID, record.StoreID, true); err == nil && existing != nil {
		record.ID = existing.ID
		record.ReturnNo = existing.ReturnNo
		if err := s.returnModule.Update(record); err != nil {
			return nil, err
		}
		return s.returnModule.GetByIDScoped(record.ID, record.StoreID, true)
	}
	for i := 0; i < 3; i++ {
		record.ReturnNo = s.returnModule.GenerateReturnNo()
		if err := s.returnModule.Create(record); err != nil {
			if module.IsDuplicateKeyError(err) {
				if existing, getErr := s.returnModule.GetByClientReqIDScoped(record.ClientReqID, record.StoreID, true); getErr == nil && existing != nil {
					record.ID = existing.ID
					record.ReturnNo = existing.ReturnNo
					if updateErr := s.returnModule.Update(record); updateErr != nil {
						return nil, updateErr
					}
					return s.returnModule.GetByIDScoped(record.ID, record.StoreID, true)
				}
				continue
			}
			return nil, err
		}
		return s.returnModule.GetByIDScoped(record.ID, record.StoreID, true)
	}
	if err := s.returnModule.Create(record); err != nil {
		return nil, err
	}
	return s.returnModule.GetByIDScoped(record.ID, record.StoreID, true)
}

func (s *StoreReturnService) Update(id, storeID, operatorID uint, req *model.UpdateStoreReturnReq, hqUnbound bool) (*model.StoreReturn, error) {
	existing, err := s.returnModule.GetByIDScoped(id, storeID, hqUnbound)
	if err != nil {
		return nil, err
	}
	if !s.IsReturnEditable(existing) {
		return nil, fmt.Errorf("返厂记录仅允许在录入当天修改")
	}
	record, err := s.buildRecord(storeID, operatorID, hqUnbound, req.StoreID, req.ReturnDate, req.LogisticsFee, req.Remark, req.Items)
	if err != nil {
		return nil, err
	}
	record.ID = existing.ID
	record.ReturnNo = existing.ReturnNo
	if err := s.returnModule.Update(record); err != nil {
		return nil, err
	}
	return s.returnModule.GetByIDScoped(record.ID, record.StoreID, true)
}

func (s *StoreReturnService) buildRecord(
	storeID, operatorID uint,
	hqUnbound bool,
	reqStoreID uint,
	returnDate string,
	logisticsFee float64,
	remark string,
	reqItems []model.CreateStoreReturnItemReq,
) (*model.StoreReturn, error) {
	realStoreID := storeID
	if hqUnbound && reqStoreID > 0 {
		realStoreID = reqStoreID
	}
	if realStoreID == 0 {
		return nil, fmt.Errorf("请选择门店")
	}

	parsedDate, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(returnDate), time.Local)
	if err != nil {
		return nil, fmt.Errorf("返厂日期格式应为 YYYY-MM-DD")
	}

	operatorName := ""
	if s.userModule != nil {
		if user, err := s.userModule.GetByID(operatorID); err == nil && user != nil {
			operatorName = user.Nickname
			if operatorName == "" {
				operatorName = user.Username
			}
		}
	}

	items := make([]model.StoreReturnItem, 0, len(reqItems))
	var totalDeposit float64
	productIDs := make([]uint, 0, len(reqItems))
	for _, item := range reqItems {
		if item.ProductID > 0 {
			productIDs = append(productIDs, item.ProductID)
		}
	}
	productMap, err := s.returnModule.GetProductMap(productIDs, realStoreID, false)
	if err != nil {
		return nil, err
	}

	for _, item := range reqItems {
		name := strings.TrimSpace(item.ProductName)
		deposit := item.Deposit
		if item.ProductID > 0 {
			product := productMap[item.ProductID]
			if product == nil {
				return nil, fmt.Errorf("返厂商品不存在或不属于当前门店")
			}
			if product.Status != 1 {
				return nil, fmt.Errorf("返厂商品【%s】已停用", product.ProductName)
			}
			name = product.ProductName
			deposit = product.Deposit
		}
		if name == "" {
			return nil, fmt.Errorf("商品名称不能为空")
		}
		if deposit < 0 {
			return nil, fmt.Errorf("商品【%s】押金不能小于0", name)
		}
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("商品【%s】数量必须大于0", name)
		}
		items = append(items, model.StoreReturnItem{
			ProductID:   item.ProductID,
			ProductName: name,
			Quantity:    item.Quantity,
			Deposit:     deposit,
			Remark:      strings.TrimSpace(item.Remark),
		})
		totalDeposit += deposit * item.Quantity
	}

	return &model.StoreReturn{
		StoreID:      realStoreID,
		ReturnDate:   parsedDate,
		LogisticsFee: logisticsFee,
		TotalDeposit: totalDeposit,
		ItemCount:    len(items),
		Remark:       strings.TrimSpace(remark),
		OperatorID:   operatorID,
		OperatorName: operatorName,
		Items:        items,
	}, nil
}

func (s *StoreReturnService) List(ctx context.Context, req *model.ListStoreReturnReq) ([]*model.StoreReturn, int64, error) {
	_ = ctx
	return s.returnModule.List(req)
}

func (s *StoreReturnService) Get(id, storeID uint, hqUnbound bool) (*model.StoreReturn, error) {
	if !hqUnbound && storeID == 0 {
		return nil, fmt.Errorf("current user has no store")
	}
	return s.returnModule.GetByIDScoped(id, storeID, hqUnbound)
}

func (s *StoreReturnService) Delete(id, storeID uint, hqUnbound bool) error {
	if !hqUnbound && storeID == 0 {
		return fmt.Errorf("current user has no store")
	}
	existing, err := s.returnModule.GetByIDScoped(id, storeID, hqUnbound)
	if err != nil {
		return err
	}
	if !s.IsReturnEditable(existing) {
		return fmt.Errorf("返厂记录仅允许在录入当天删除")
	}
	return s.returnModule.Delete(id, storeID, hqUnbound)
}

func (s *StoreReturnService) Stats(req *model.ListStoreReturnReq) (*model.StoreReturnStats, error) {
	return s.returnModule.Stats(req)
}

func (s *StoreReturnService) IsReturnEditable(record *model.StoreReturn) bool {
	if record == nil || record.CreatedAt.IsZero() {
		return false
	}
	now := time.Now()
	created := record.CreatedAt.In(now.Location())
	return created.Year() == now.Year() && created.YearDay() == now.YearDay()
}

func (s *StoreReturnService) CreateProduct(storeID uint, req *model.CreateStoreReturnProductReq, hqUnbound bool) (*model.StoreReturnProduct, error) {
	product, err := s.buildProduct(storeID, hqUnbound, req.StoreID, req.ProductName, req.Deposit, req.Remark, req.Status)
	if err != nil {
		return nil, err
	}
	if err := s.returnModule.CreateProduct(product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *StoreReturnService) UpdateProduct(id, storeID uint, req *model.UpdateStoreReturnProductReq, hqUnbound bool) (*model.StoreReturnProduct, error) {
	existing, err := s.returnModule.GetProductByIDScoped(id, storeID, hqUnbound)
	if err != nil {
		return nil, err
	}
	product, err := s.buildProduct(storeID, hqUnbound, req.StoreID, req.ProductName, req.Deposit, req.Remark, req.Status)
	if err != nil {
		return nil, err
	}
	product.ID = existing.ID
	if err := s.returnModule.UpdateProduct(product); err != nil {
		return nil, err
	}
	return s.returnModule.GetProductByIDScoped(product.ID, product.StoreID, true)
}

func (s *StoreReturnService) buildProduct(storeID uint, hqUnbound bool, reqStoreID uint, productName string, deposit float64, remark string, status int) (*model.StoreReturnProduct, error) {
	realStoreID := storeID
	if hqUnbound && reqStoreID > 0 {
		realStoreID = reqStoreID
	}
	if realStoreID == 0 {
		return nil, fmt.Errorf("请选择门店")
	}
	name := strings.TrimSpace(productName)
	if name == "" {
		return nil, fmt.Errorf("商品名称不能为空")
	}
	return &model.StoreReturnProduct{
		StoreID:     realStoreID,
		ProductName: name,
		Deposit:     deposit,
		Remark:      strings.TrimSpace(remark),
		Status:      status,
	}, nil
}

func (s *StoreReturnService) ListProducts(ctx context.Context, req *model.ListStoreReturnProductReq) ([]*model.StoreReturnProduct, int64, error) {
	_ = ctx
	return s.returnModule.ListProducts(req)
}

func (s *StoreReturnService) DeleteProduct(id, storeID uint, hqUnbound bool) error {
	if !hqUnbound && storeID == 0 {
		return fmt.Errorf("current user has no store")
	}
	return s.returnModule.DeleteProduct(id, storeID, hqUnbound)
}
