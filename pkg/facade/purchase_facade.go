package facade

import (
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/pkg/statemachine"
	"github.com/Kevin-Jii/tower-go/utils"
)

// PurchaseFacade 采购业务外观
// 封装采购相关的复杂业务逻辑，提供简单统一的接口
type PurchaseFacade struct {
	orderService    OrderService
	supplierService SupplierService
	notifyService   NotifyService
	stateMachine    *statemachine.StateMachine
}

// OrderService 订单服务接口
type OrderService interface {
	CreateOrder(storeID, userID uint, req *model.CreatePurchaseOrderReq) (*model.PurchaseOrder, error)
	GetOrder(id uint) (*model.PurchaseOrder, error)
	UpdateOrder(id uint, req *model.UpdatePurchaseOrderReq) error
}

// SupplierService 供应商服务接口
type SupplierService interface {
	ListSuppliersByStoreID(storeID uint) ([]*model.StoreSupplier, error)
}

// NotifyService 通知服务接口
type NotifyService interface {
	NotifySupplier(supplierID uint, message string) error
	NotifyStore(storeID uint, message string) error
}

// NewPurchaseFacade 创建采购外观
func NewPurchaseFacade(
	orderService OrderService,
	supplierService SupplierService,
	notifyService NotifyService,
) *PurchaseFacade {
	sm := statemachine.NewOrderStateMachine()

	// 注册状态变更钩子
	sm.OnAction(statemachine.ActionConfirm, func(from, to statemachine.State) error {
		// 确认时发送事件
		utils.GlobalEventBus.Publish(utils.EventOrderConfirmed, nil)
		return nil
	})

	sm.OnAction(statemachine.ActionComplete, func(from, to statemachine.State) error {
		utils.GlobalEventBus.Publish(utils.EventOrderCompleted, nil)
		return nil
	})

	return &PurchaseFacade{
		orderService:    orderService,
		supplierService: supplierService,
		notifyService:   notifyService,
		stateMachine:    sm,
	}
}

// CreateAndNotify 创建采购单并通知供应商
func (f *PurchaseFacade) CreateAndNotify(storeID, userID uint, req *model.CreatePurchaseOrderReq) (*model.PurchaseOrder, error) {
	// 1. 创建订单
	order, err := f.orderService.CreateOrder(storeID, userID, req)
	if err != nil {
		return nil, err
	}

	// 2. 发布事件
	utils.GlobalEventBus.Publish(utils.EventOrderCreated, order)

	// 3. 通知供应商（异步）
	go func() {
		if f.notifyService != nil {
			// 获取订单涉及的供应商并通知
			// f.notifyService.NotifySupplier(...)
		}
	}()

	return order, nil
}

// ConfirmOrder 确认采购单
func (f *PurchaseFacade) ConfirmOrder(orderID uint) error {
	order, err := f.orderService.GetOrder(orderID)
	if err != nil {
		return err
	}

	// 使用状态机验证转换
	newStatus, err := f.stateMachine.Execute(
		statemachine.State(order.Status),
		statemachine.ActionConfirm,
	)
	if err != nil {
		return err
	}

	// 更新状态
	status := int8(newStatus)
	return f.orderService.UpdateOrder(orderID, &model.UpdatePurchaseOrderReq{
		Status: &status,
	})
}

// CancelOrder 取消采购单
func (f *PurchaseFacade) CancelOrder(orderID uint, reason string) error {
	order, err := f.orderService.GetOrder(orderID)
	if err != nil {
		return err
	}

	newStatus, err := f.stateMachine.Execute(
		statemachine.State(order.Status),
		statemachine.ActionCancel,
	)
	if err != nil {
		return err
	}

	status := int8(newStatus)
	return f.orderService.UpdateOrder(orderID, &model.UpdatePurchaseOrderReq{
		Status: &status,
		Remark: &reason,
	})
}

// GetAvailableActions 获取订单可用操作
func (f *PurchaseFacade) GetAvailableActions(orderID uint) ([]string, error) {
	order, err := f.orderService.GetOrder(orderID)
	if err != nil {
		return nil, err
	}

	actions := f.stateMachine.GetAvailableActions(statemachine.State(order.Status))
	result := make([]string, len(actions))
	for i, a := range actions {
		result[i] = string(a)
	}
	return result, nil
}
