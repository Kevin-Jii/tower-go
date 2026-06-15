package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	updatesPkg "github.com/Kevin-Jii/tower-go/utils/updates"
)

type B2BService struct {
	b2bModule      *module.B2BModule
	storeModule    *module.StoreModule
	productModule  *module.SupplierProductModule
	unitSpecModule *module.ProductUnitSpecModule
	userModule     *module.UserModule
}

func NewB2BService(
	b2bModule *module.B2BModule,
	storeModule *module.StoreModule,
	productModule *module.SupplierProductModule,
	unitSpecModule *module.ProductUnitSpecModule,
	userModule *module.UserModule,
) *B2BService {
	return &B2BService{
		b2bModule:      b2bModule,
		storeModule:    storeModule,
		productModule:  productModule,
		unitSpecModule: unitSpecModule,
		userModule:     userModule,
	}
}

func (s *B2BService) CreateCustomer(storeID uint, req *model.CreateB2BCustomerReq) (*model.B2BCustomer, error) {
	if storeID == 0 {
		return nil, errors.New("store_id is required")
	}
	customer := &model.B2BCustomer{
		StoreID:       storeID,
		Name:          strings.TrimSpace(req.Name),
		CustomerType:  strings.TrimSpace(req.CustomerType),
		ContactPerson: strings.TrimSpace(req.ContactPerson),
		Phone:         strings.TrimSpace(req.Phone),
		Address:       strings.TrimSpace(req.Address),
		Settlement:    normalizeSettlement(req.Settlement),
		PriceLevel:    strings.TrimSpace(req.PriceLevel),
		CreditLimit:   req.CreditLimit,
		Status:        model.B2BCustomerStatusEnabled,
		Remark:        strings.TrimSpace(req.Remark),
	}
	if err := s.b2bModule.CreateCustomer(customer); err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *B2BService) UpdateCustomer(id, storeID uint, isHQ bool, req *model.UpdateB2BCustomerReq) (*model.B2BCustomer, error) {
	customer, err := s.b2bModule.GetCustomer(id)
	if err != nil {
		return nil, errors.New("customer not found")
	}
	if !isHQ && customer.StoreID != storeID {
		return nil, errors.New("permission denied")
	}
	updates := updatesPkg.BuildUpdatesFromReq(req)
	if v, ok := updates["settlement"]; ok {
		updates["settlement"] = normalizeSettlement(fmt.Sprint(v))
	}
	if len(updates) > 0 {
		if err := s.b2bModule.UpdateCustomer(id, updates); err != nil {
			return nil, err
		}
	}
	return s.b2bModule.GetCustomer(id)
}

func (s *B2BService) ListCustomers(_ context.Context, req *model.ListB2BCustomerReq) ([]*model.B2BCustomer, int64, error) {
	return s.b2bModule.ListCustomers(req)
}

func (s *B2BService) UpsertPrice(storeID uint, isHQ bool, req *model.UpsertB2BPriceReq) error {
	if isHQ && req.StoreID > 0 {
		storeID = req.StoreID
	}
	if storeID == 0 {
		return errors.New("store_id is required")
	}
	if req.CustomerID == nil && strings.TrimSpace(req.PriceLevel) == "" {
		return errors.New("customer_id 和 price_level 至少填写一个")
	}
	if req.CustomerID != nil && *req.CustomerID > 0 {
		customer, err := s.b2bModule.GetCustomer(*req.CustomerID)
		if err != nil {
			return errors.New("customer not found")
		}
		if customer.StoreID != storeID {
			return errors.New("customer not in store")
		}
	}
	spec, err := s.unitSpecModule.GetByID(req.UnitSpecID)
	if err != nil || spec == nil {
		return errors.New("unit spec not found")
	}
	if spec.ProductID != req.ProductID {
		return errors.New("unit spec does not match product")
	}
	enabled := true
	if req.IsEnabled != nil {
		enabled = *req.IsEnabled
	}
	minQty := req.MinQuantity
	if minQty <= 0 {
		minQty = 1
	}
	price := &model.B2BCustomerProductPrice{
		StoreID:     storeID,
		CustomerID:  req.CustomerID,
		PriceLevel:  strings.TrimSpace(req.PriceLevel),
		ProductID:   req.ProductID,
		UnitSpecID:  req.UnitSpecID,
		UnitName:    spec.UnitName,
		SupplyPrice: req.SupplyPrice,
		MinQuantity: minQty,
		IsEnabled:   enabled,
		Remark:      strings.TrimSpace(req.Remark),
	}
	return s.b2bModule.UpsertPrice(price)
}

func (s *B2BService) ListPrices(req *model.ListB2BPriceReq) ([]*model.B2BCustomerProductPrice, int64, error) {
	return s.b2bModule.ListPrices(req)
}

func (s *B2BService) DeletePrice(id uint) error {
	return s.b2bModule.DeletePrice(id)
}

func (s *B2BService) CreateSupplyOrder(storeID, operatorID uint, isHQ bool, req *model.CreateB2BSupplyOrderReq) (*model.B2BSupplyOrder, error) {
	if isHQ && req.StoreID > 0 {
		storeID = req.StoreID
	}
	if storeID == 0 {
		return nil, errors.New("store_id is required")
	}
	customer, err := s.b2bModule.GetCustomer(req.CustomerID)
	if err != nil {
		return nil, errors.New("customer not found")
	}
	if customer.StoreID != storeID {
		return nil, errors.New("customer not in store")
	}
	if customer.Status != model.B2BCustomerStatusEnabled {
		return nil, errors.New("customer disabled")
	}

	orderDate := time.Now()
	if strings.TrimSpace(req.OrderDate) != "" {
		if t, err := time.Parse("2006-01-02", req.OrderDate); err == nil {
			orderDate = t
		}
	}

	operatorName := ""
	if user, err := s.userModule.GetByID(operatorID); err == nil && user != nil {
		operatorName = user.Nickname
		if operatorName == "" {
			operatorName = user.Username
		}
	}

	items := make([]model.B2BSupplyOrderItem, 0, len(req.Items))
	var total, costTotal float64
	for _, line := range req.Items {
		product, err := s.productModule.GetByID(line.ProductID)
		if err != nil || product == nil {
			return nil, fmt.Errorf("product %d not found", line.ProductID)
		}
		spec, err := s.unitSpecModule.GetByID(line.UnitSpecID)
		if err != nil || spec == nil {
			return nil, fmt.Errorf("unit spec %d not found", line.UnitSpecID)
		}
		if spec.ProductID != line.ProductID || !spec.IsEnabled {
			return nil, fmt.Errorf("商品【%s】规格不可用", product.Name)
		}
		price := line.SupplyPrice
		if price <= 0 {
			if configured, err := s.b2bModule.ResolvePrice(storeID, customer.ID, product.ID, spec.ID, customer.PriceLevel); err == nil && configured != nil {
				price = configured.SupplyPrice
			}
		}
		if price <= 0 {
			price = spec.SalePrice
		}
		if price <= 0 {
			return nil, fmt.Errorf("商品【%s】未配置供货价", product.Name)
		}

		factor := spec.FactorToBase
		if factor <= 0 {
			factor = 1
		}
		baseQty := line.Quantity * factor
		amount := roundMoney(line.Quantity * price)
		costAmount := roundMoney(line.Quantity * spec.CostPrice)
		item := model.B2BSupplyOrderItem{
			ProductID:    product.ID,
			ProductName:  product.Name,
			UnitSpecID:   spec.ID,
			UnitName:     spec.UnitName,
			FactorToBase: factor,
			Quantity:     line.Quantity,
			BaseQuantity: baseQty,
			SupplyPrice:  price,
			CostPrice:    spec.CostPrice,
			Amount:       amount,
			CostAmount:   costAmount,
			ProfitAmount: roundMoney(amount - costAmount),
			Remark:       strings.TrimSpace(line.Remark),
		}
		items = append(items, item)
		total += amount
		costTotal += costAmount
	}

	total = roundMoney(total)
	costTotal = roundMoney(costTotal)
	paid := roundMoney(req.PaidAmount)
	if paid > total {
		paid = total
	}
	unpaid := roundMoney(total - paid)
	paymentStatus := model.B2BPaymentUnpaid
	if paid >= total {
		paymentStatus = model.B2BPaymentPaid
	} else if paid > 0 {
		paymentStatus = model.B2BPaymentPartial
	}
	deliveryStatus := req.DeliveryStatus
	if deliveryStatus == 0 {
		deliveryStatus = model.B2BDeliveryPending
	}

	order := &model.B2BSupplyOrder{
		OrderNo:        s.b2bModule.GenerateSupplyOrderNo(),
		StoreID:        storeID,
		CustomerID:     customer.ID,
		CustomerName:   customer.Name,
		OrderDate:      orderDate,
		TotalAmount:    total,
		PaidAmount:     paid,
		UnpaidAmount:   unpaid,
		CostAmount:     costTotal,
		ProfitAmount:   roundMoney(total - costTotal),
		PaymentStatus:  paymentStatus,
		DeliveryStatus: deliveryStatus,
		Remark:         strings.TrimSpace(req.Remark),
		OperatorID:     operatorID,
		OperatorName:   operatorName,
		Items:          items,
	}
	if err := s.b2bModule.CreateSupplyOrderWithInventory(order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *B2BService) ListSupplyOrders(req *model.ListB2BSupplyOrderReq) ([]*model.B2BSupplyOrder, int64, error) {
	return s.b2bModule.ListSupplyOrders(req)
}

func (s *B2BService) GetSupplyOrder(id, storeID uint, isHQ bool) (*model.B2BSupplyOrder, error) {
	order, err := s.b2bModule.GetSupplyOrder(id)
	if err != nil {
		return nil, err
	}
	if !isHQ && order.StoreID != storeID {
		return nil, errors.New("permission denied")
	}
	return order, nil
}

func normalizeSettlement(v string) string {
	switch strings.TrimSpace(v) {
	case "week", "month":
		return strings.TrimSpace(v)
	default:
		return "cash"
	}
}

func roundMoney(v float64) float64 {
	return math.Round(v*100) / 100
}
