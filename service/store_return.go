package service

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/pkg/apicode"
)

type StoreReturnService struct {
	returnModule *module.StoreReturnModule
	userModule   *module.UserModule
}

func NewStoreReturnService(returnModule *module.StoreReturnModule, userModule *module.UserModule) *StoreReturnService {
	return &StoreReturnService{returnModule: returnModule, userModule: userModule}
}

func (s *StoreReturnService) Create(storeID, operatorID uint, req *model.CreateStoreReturnReq, hqUnbound bool) (*model.StoreReturn, error) {
	record, err := s.buildRecord(storeID, operatorID, hqUnbound, req.StoreID, req.ReturnDate, req.LogisticsFee, req.Photos, req.Remark, req.Items)
	if err != nil {
		return nil, err
	}
	clientReqID := strings.TrimSpace(req.ClientReqID)
	record.ClientReqID = optionalStoreReturnClientReqID(clientReqID)
	if clientReqID != "" {
		if existing, err := s.returnModule.GetByClientReqIDScoped(clientReqID, record.StoreID, true); err == nil && existing != nil {
			record.ID = existing.ID
			record.ReturnNo = existing.ReturnNo
			if err := s.returnModule.Update(record); err != nil {
				return nil, err
			}
			return s.returnModule.GetByIDScoped(record.ID, record.StoreID, true)
		}
	}
	for i := 0; i < 3; i++ {
		record.ReturnNo = s.returnModule.GenerateReturnNo()
		if err := s.returnModule.Create(record); err != nil {
			if module.IsDuplicateKeyError(err) {
				if clientReqID != "" {
					if existing, getErr := s.returnModule.GetByClientReqIDScoped(clientReqID, record.StoreID, true); getErr == nil && existing != nil {
						record.ID = existing.ID
						record.ReturnNo = existing.ReturnNo
						if updateErr := s.returnModule.Update(record); updateErr != nil {
							return nil, updateErr
						}
						return s.returnModule.GetByIDScoped(record.ID, record.StoreID, true)
					}
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

func optionalStoreReturnClientReqID(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func (s *StoreReturnService) Update(id, storeID, operatorID uint, req *model.UpdateStoreReturnReq, hqUnbound bool) (*model.StoreReturn, error) {
	existing, err := s.returnModule.GetByIDScoped(id, storeID, hqUnbound)
	if err != nil {
		return nil, err
	}
	if !s.IsReturnEditable(existing) {
		return nil, apicode.Newf(apicode.OrderStateConflict, "返厂记录仅允许在录入当天修改")
	}
	record, err := s.buildRecord(storeID, operatorID, hqUnbound, req.StoreID, req.ReturnDate, req.LogisticsFee, req.Photos, req.Remark, req.Items)
	if err != nil {
		return nil, err
	}
	record.ID = existing.ID
	record.ReturnNo = existing.ReturnNo
	record.ClientReqID = existing.ClientReqID
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
	photos []string,
	remark string,
	reqItems []model.CreateStoreReturnItemReq,
) (*model.StoreReturn, error) {
	realStoreID := storeID
	if hqUnbound && reqStoreID > 0 {
		realStoreID = reqStoreID
	}
	if realStoreID == 0 {
		return nil, apicode.New(apicode.StoreRequired)
	}

	parsedDate, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(returnDate), time.Local)
	if err != nil {
		return nil, apicode.New(apicode.ReturnDateInvalid)
	}
	photoURLs, err := normalizeStoreReturnPhotos(photos)
	if err != nil {
		return nil, err
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
				return nil, apicode.Newf(apicode.ProductNotFound, "返厂商品不存在或不属于当前门店")
			}
			if product.Status != 1 {
				return nil, apicode.Newf(apicode.ValidationFailed, "返厂商品【%s】已停用", product.ProductName)
			}
			name = product.ProductName
			deposit = product.Deposit
		}
		if name == "" {
			return nil, apicode.Newf(apicode.ValidationFailed, "商品名称不能为空")
		}
		if deposit < 0 {
			return nil, apicode.Newf(apicode.ValidationFailed, "商品【%s】押金不能小于0", name)
		}
		if item.Quantity <= 0 {
			return nil, apicode.Newf(apicode.ValidationFailed, "商品【%s】数量必须大于0", name)
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
		Photos:       photoURLs,
		Remark:       strings.TrimSpace(remark),
		OperatorID:   operatorID,
		OperatorName: operatorName,
		Items:        items,
	}, nil
}

func normalizeStoreReturnPhotos(photos []string) (model.StringList, error) {
	if len(photos) > 3 {
		return nil, apicode.Newf(apicode.ValidationFailed, "返厂照片最多上传3张")
	}

	result := make(model.StringList, 0, len(photos))
	seen := make(map[string]struct{}, len(photos))
	for _, raw := range photos {
		photoURL := strings.TrimSpace(raw)
		if photoURL == "" {
			continue
		}
		if len(photoURL) > 500 {
			return nil, apicode.Newf(apicode.ValidationFailed, "返厂照片地址不能超过500个字符")
		}
		parsed, err := url.Parse(photoURL)
		if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
			return nil, apicode.Newf(apicode.ValidationFailed, "返厂照片地址无效")
		}
		if _, ok := seen[photoURL]; ok {
			continue
		}
		seen[photoURL] = struct{}{}
		result = append(result, photoURL)
	}
	return result, nil
}

func (s *StoreReturnService) List(ctx context.Context, req *model.ListStoreReturnReq) ([]*model.StoreReturn, int64, error) {
	_ = ctx
	return s.returnModule.List(req)
}

func (s *StoreReturnService) Get(id, storeID uint, hqUnbound bool) (*model.StoreReturn, error) {
	if !hqUnbound && storeID == 0 {
		return nil, apicode.New(apicode.StoreRequired)
	}
	return s.returnModule.GetByIDScoped(id, storeID, hqUnbound)
}

func (s *StoreReturnService) Delete(id, storeID uint, hqUnbound bool) error {
	if !hqUnbound && storeID == 0 {
		return apicode.New(apicode.StoreRequired)
	}
	existing, err := s.returnModule.GetByIDScoped(id, storeID, hqUnbound)
	if err != nil {
		return err
	}
	if !s.IsReturnEditable(existing) {
		return apicode.Newf(apicode.OrderStateConflict, "返厂记录仅允许在录入当天删除")
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
		return nil, apicode.New(apicode.StoreRequired)
	}
	name := strings.TrimSpace(productName)
	if name == "" {
		return nil, apicode.Newf(apicode.ValidationFailed, "商品名称不能为空")
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
		return apicode.New(apicode.StoreRequired)
	}
	return s.returnModule.DeleteProduct(id, storeID, hqUnbound)
}
