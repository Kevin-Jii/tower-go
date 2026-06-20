package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

type InventoryLossService struct {
	lossModule     *module.InventoryLossModule
	productModule  *module.SupplierProductModule
	unitSpecModule *module.ProductUnitSpecModule
	memberModule   *module.MemberModule
	userModule     *module.UserModule
	dictModule     *module.DictModule
}

func NewInventoryLossService(
	lossModule *module.InventoryLossModule,
	productModule *module.SupplierProductModule,
	unitSpecModule *module.ProductUnitSpecModule,
	memberModule *module.MemberModule,
	userModule *module.UserModule,
	dictModule *module.DictModule,
) *InventoryLossService {
	return &InventoryLossService{
		lossModule:     lossModule,
		productModule:  productModule,
		unitSpecModule: unitSpecModule,
		memberModule:   memberModule,
		userModule:     userModule,
		dictModule:     dictModule,
	}
}

func (s *InventoryLossService) CreateOrder(storeID, operatorID uint, req *model.CreateInventoryLossOrderReq, hqUnbound bool) (*model.InventoryLossOrder, error) {
	realStoreID := storeID
	if hqUnbound && req.StoreID > 0 {
		realStoreID = req.StoreID
	}
	if realStoreID == 0 {
		return nil, fmt.Errorf("请选择门店")
	}
	reason, err := s.resolveReason(req.Reason)
	if err != nil {
		return nil, err
	}

	var memberID *uint
	if req.Type == model.InventoryLossTypeGift {
		if req.MemberID == nil || *req.MemberID == 0 {
			return nil, fmt.Errorf("赠送类型必须选择会员")
		}
		member, err := s.memberModule.GetMember(*req.MemberID, realStoreID, true)
		if err != nil || member == nil || member.StoreID != realStoreID {
			return nil, fmt.Errorf("会员不存在或不属于当前门店")
		}
		memberID = req.MemberID
	}

	operatorName := ""
	if user, err := s.userModule.GetByID(operatorID); err == nil && user != nil {
		operatorName = user.Nickname
		if operatorName == "" {
			operatorName = user.Username
		}
	}

	items := make([]model.InventoryLossOrderItem, 0, len(req.Items))
	var totalCost float64
	for _, item := range req.Items {
		unit := strings.TrimSpace(item.Unit)
		if unit == "" {
			return nil, fmt.Errorf("请选择商品规格")
		}

		product, err := s.productModule.GetByID(item.ProductID)
		if err != nil || product == nil {
			return nil, fmt.Errorf("商品不存在")
		}

		baseQuantity, baseUnit := convertToBaseQuantity(s.unitSpecModule, product, item.ProductID, item.Quantity, unit)
		specs, _ := s.unitSpecModule.ListEnabledByProductID(item.ProductID)
		costPrice := resolveUnitCostFromSpecs(unit, specs)
		if costPrice <= 0 {
			costPrice = resolveFallbackCostPrice(unit, product)
		}
		costAmount := costPrice * item.Quantity
		totalCost += costAmount

		items = append(items, model.InventoryLossOrderItem{
			ProductID:    item.ProductID,
			ProductName:  product.Name,
			Unit:         unit,
			Quantity:     item.Quantity,
			BaseQuantity: baseQuantity,
			BaseUnit:     baseUnit,
			CostPrice:    costPrice,
			CostAmount:   costAmount,
			Remark:       item.Remark,
		})
	}

	order := &model.InventoryLossOrder{
		OrderNo:      s.lossModule.GenerateOrderNo(),
		StoreID:      realStoreID,
		Type:         req.Type,
		MemberID:     memberID,
		Reason:       reason,
		TotalCost:    totalCost,
		ItemCount:    len(items),
		OperatorID:   operatorID,
		OperatorName: operatorName,
		Items:        items,
	}

	if err := s.lossModule.CreateWithStockDeduct(order); err != nil {
		return nil, err
	}
	return s.lossModule.GetByIDScoped(order.ID, realStoreID, hqUnbound)
}

func (s *InventoryLossService) UpdateOrder(id, storeID uint, req *model.UpdateInventoryLossOrderReq, hqUnbound bool) (*model.InventoryLossOrder, error) {
	if !hqUnbound && storeID == 0 {
		return nil, fmt.Errorf("current user has no store")
	}
	reason, err := s.resolveReason(req.Reason)
	if err != nil {
		return nil, err
	}
	if err := s.lossModule.UpdateReason(id, storeID, hqUnbound, reason); err != nil {
		return nil, err
	}
	return s.lossModule.GetByIDScoped(id, storeID, hqUnbound)
}

func (s *InventoryLossService) resolveReason(reason string) (string, error) {
	code := strings.TrimSpace(reason)
	if code == "" {
		return "", fmt.Errorf("请选择原因说明")
	}
	if s.dictModule == nil {
		return code, nil
	}
	data, err := s.dictModule.GetDataByTypeAndValue(model.StoreExpenseCategoryDictCode, code)
	if err != nil || data == nil || data.Status != 1 {
		return "", fmt.Errorf("原因说明不存在或已停用")
	}
	return data.Label, nil
}

func resolveFallbackCostPrice(unit string, product *model.SupplierProduct) float64 {
	if product == nil {
		return 0
	}
	if isLargePackUnit(unit) {
		if product.CasePrice > 0 {
			return product.CasePrice
		}
		return product.Price
	}
	if product.BottlePrice > 0 {
		return product.BottlePrice
	}
	return product.Price
}

func (s *InventoryLossService) ListOrders(ctx context.Context, req *model.ListInventoryLossOrderReq) ([]*model.InventoryLossOrder, int64, error) {
	_ = ctx
	return s.lossModule.List(req)
}

func (s *InventoryLossService) GetOrderByIDScoped(id, storeID uint, hqUnbound bool) (*model.InventoryLossOrder, error) {
	if !hqUnbound && storeID == 0 {
		return nil, fmt.Errorf("current user has no store")
	}
	return s.lossModule.GetByIDScoped(id, storeID, hqUnbound)
}

func (s *InventoryLossService) CancelOrder(id, storeID uint, hqUnbound bool) error {
	if !hqUnbound && storeID == 0 {
		return fmt.Errorf("current user has no store")
	}
	return s.lossModule.CancelWithStockRestore(id, storeID, hqUnbound)
}

func (s *InventoryLossService) ListMemberGiftRecords(memberID, storeID uint, hqUnbound bool, req *model.ListMemberGiftRecordsReq) ([]*model.MemberGiftRecord, int64, error) {
	if !hqUnbound && storeID == 0 {
		return nil, 0, fmt.Errorf("current user has no store")
	}
	req.StoreID = storeID
	if hqUnbound {
		member, err := s.memberModule.GetMember(memberID, 0, true)
		if err != nil {
			return nil, 0, err
		}
		if req.StoreID == 0 {
			req.StoreID = member.StoreID
		}
	} else if _, err := s.memberModule.GetMember(memberID, storeID, false); err != nil {
		return nil, 0, err
	}
	return s.lossModule.ListMemberGiftRecords(memberID, req, hqUnbound)
}
