package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/pkg/apicode"
	"github.com/Kevin-Jii/tower-go/utils/logging"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	PreOrderReminderSlot0930 = "0930"
	PreOrderReminderSlot1600 = "1600"
)

var preOrderLocation = func() *time.Location {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.FixedZone("Asia/Shanghai", 8*60*60)
	}
	return location
}()

type PreOrderReminderTarget struct {
	Start        time.Time
	End          time.Time
	ReminderKey  string
	RelativeDate string
}

type PreOrderService struct {
	preOrderModule  *module.PreOrderModule
	b2bModule       *module.B2BModule
	productModule   *module.SupplierProductModule
	unitSpecModule  *module.ProductUnitSpecModule
	storeModule     *module.StoreModule
	botModule       *module.DingTalkBotModule
	dingTalkService *DingTalkService
}

func NewPreOrderService(
	preOrderModule *module.PreOrderModule,
	b2bModule *module.B2BModule,
	productModule *module.SupplierProductModule,
	unitSpecModule *module.ProductUnitSpecModule,
	storeModule *module.StoreModule,
	botModule *module.DingTalkBotModule,
	dingTalkService *DingTalkService,
) *PreOrderService {
	return &PreOrderService{
		preOrderModule:  preOrderModule,
		b2bModule:       b2bModule,
		productModule:   productModule,
		unitSpecModule:  unitSpecModule,
		storeModule:     storeModule,
		botModule:       botModule,
		dingTalkService: dingTalkService,
	}
}

func (s *PreOrderService) Create(storeID, userID uint, req *model.CreatePreOrderReq) (*model.PreOrder, error) {
	if storeID == 0 {
		return nil, apicode.New(apicode.StoreRequired)
	}
	scheduledAt, err := parsePreOrderTime(req.ScheduledAt)
	if err != nil {
		return nil, err
	}
	order, items, err := s.buildOrder(storeID, req.CustomerID, scheduledAt, req.ContactPerson, req.ContactPhone, req.DeliveryAddress, req.Remark, req.Items)
	if err != nil {
		return nil, err
	}
	order.OrderNo = s.preOrderModule.GenerateOrderNo(time.Now().In(preOrderLocation))
	order.CreatedBy = userID
	order.Items = items
	if err := s.preOrderModule.Create(order); err != nil {
		return nil, err
	}
	return s.preOrderModule.GetByID(order.ID)
}

func (s *PreOrderService) Update(id, storeID uint, hqUnbound bool, req *model.UpdatePreOrderReq) (*model.PreOrder, error) {
	existing, err := s.getScoped(id, storeID, hqUnbound)
	if err != nil {
		return nil, err
	}
	if existing.Status == model.PreOrderStatusDelivered || existing.Status == model.PreOrderStatusCancelled {
		return nil, apicode.New(apicode.OrderStateConflict)
	}
	scheduledAt, err := parsePreOrderTime(req.ScheduledAt)
	if err != nil {
		return nil, err
	}
	order, items, err := s.buildOrder(existing.StoreID, req.CustomerID, scheduledAt, req.ContactPerson, req.ContactPhone, req.DeliveryAddress, req.Remark, req.Items)
	if err != nil {
		return nil, err
	}
	order.ID = existing.ID
	for i := range items {
		items[i].PreOrderID = existing.ID
	}
	resetReminders := !samePreOrderDate(existing.ScheduledAt, scheduledAt)
	if err := s.preOrderModule.Update(order, items, resetReminders); err != nil {
		return nil, err
	}
	return s.preOrderModule.GetByID(id)
}

func (s *PreOrderService) buildOrder(
	storeID, customerID uint,
	scheduledAt time.Time,
	contactPerson, contactPhone, address, remark string,
	reqItems []model.CreatePreOrderItemReq,
) (*model.PreOrder, []model.PreOrderItem, error) {
	customer, err := s.b2bModule.GetCustomer(customerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, apicode.New(apicode.CustomerNotFound)
		}
		return nil, nil, err
	}
	if customer.StoreID != storeID {
		return nil, nil, apicode.New(apicode.CustomerNotFound)
	}
	if customer.Status != model.B2BCustomerStatusEnabled {
		return nil, nil, apicode.New(apicode.CustomerDisabled)
	}

	if strings.TrimSpace(contactPerson) == "" {
		contactPerson = customer.ContactPerson
	}
	if strings.TrimSpace(contactPhone) == "" {
		contactPhone = customer.Phone
	}
	if strings.TrimSpace(address) == "" {
		address = customer.Address
	}
	items := make([]model.PreOrderItem, 0, len(reqItems))
	seen := make(map[string]struct{}, len(reqItems))
	for _, reqItem := range reqItems {
		key := fmt.Sprintf("%d:%d", reqItem.ProductID, reqItem.UnitSpecID)
		if _, exists := seen[key]; exists {
			return nil, nil, apicode.Newf(apicode.ValidationFailed, "同一商品规格不能重复添加")
		}
		seen[key] = struct{}{}

		product, err := s.productModule.GetByID(reqItem.ProductID)
		if err != nil || product.Status != 1 {
			return nil, nil, apicode.New(apicode.ProductNotFound)
		}
		spec, err := s.unitSpecModule.GetByID(reqItem.UnitSpecID)
		if err != nil {
			return nil, nil, apicode.New(apicode.UnitSpecNotFound)
		}
		if spec.ProductID != product.ID {
			return nil, nil, apicode.New(apicode.UnitSpecMismatch)
		}
		if !spec.IsEnabled {
			return nil, nil, apicode.Newf(apicode.ValidationFailed, "商品【%s】规格【%s】已停用", product.Name, spec.UnitName)
		}
		configuredPrice, err := s.b2bModule.GetConfiguredPrice(storeID, customer.ID, product.ID, spec.ID, customer.PriceLevel)
		if err != nil {
			return nil, nil, err
		}
		if configuredPrice == nil || !configuredPrice.IsEnabled {
			return nil, nil, apicode.Newf(apicode.ValidationFailed, "商品【%s】规格【%s】未给该客户启用", product.Name, spec.UnitName)
		}
		items = append(items, model.PreOrderItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			UnitSpecID:  spec.ID,
			UnitName:    spec.UnitName,
			Quantity:    reqItem.Quantity,
			Remark:      strings.TrimSpace(reqItem.Remark),
		})
	}
	if len(items) == 0 {
		return nil, nil, apicode.Newf(apicode.ValidationFailed, "至少添加一项商品")
	}
	return &model.PreOrder{
		StoreID:         storeID,
		CustomerID:      customer.ID,
		CustomerName:    customer.Name,
		ContactPerson:   strings.TrimSpace(contactPerson),
		ContactPhone:    strings.TrimSpace(contactPhone),
		DeliveryAddress: strings.TrimSpace(address),
		ScheduledAt:     scheduledAt,
		Status:          model.PreOrderStatusPending,
		Remark:          strings.TrimSpace(remark),
	}, items, nil
}

func (s *PreOrderService) Get(id, storeID uint, hqUnbound bool) (*model.PreOrder, error) {
	return s.getScoped(id, storeID, hqUnbound)
}

func (s *PreOrderService) List(req *model.ListPreOrderReq) ([]*model.PreOrder, int64, error) {
	return s.preOrderModule.List(req)
}

func (s *PreOrderService) UpdateStatus(id, storeID uint, hqUnbound bool, status int8) (*model.PreOrder, error) {
	order, err := s.getScoped(id, storeID, hqUnbound)
	if err != nil {
		return nil, err
	}
	if !preOrderStatusTransitionAllowed(order.Status, status) {
		return nil, apicode.New(apicode.OrderStateConflict)
	}
	if err := s.preOrderModule.UpdateStatus(id, status, time.Now().In(preOrderLocation)); err != nil {
		return nil, err
	}
	return s.preOrderModule.GetByID(id)
}

func (s *PreOrderService) Delete(id, storeID uint, hqUnbound bool) error {
	order, err := s.getScoped(id, storeID, hqUnbound)
	if err != nil {
		return err
	}
	if order.Status != model.PreOrderStatusPending && order.Status != model.PreOrderStatusCancelled {
		return apicode.New(apicode.OrderDeletionDenied)
	}
	return s.preOrderModule.Delete(id)
}

func (s *PreOrderService) getScoped(id, storeID uint, hqUnbound bool) (*model.PreOrder, error) {
	order, err := s.preOrderModule.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apicode.New(apicode.OrderNotFound)
		}
		return nil, err
	}
	if !hqUnbound && (storeID == 0 || order.StoreID != storeID) {
		return nil, apicode.New(apicode.OrderNotFound)
	}
	return order, nil
}

func preOrderStatusTransitionAllowed(from, to int8) bool {
	if from == to {
		return true
	}
	switch from {
	case model.PreOrderStatusPending:
		return to == model.PreOrderStatusPrepared || to == model.PreOrderStatusCancelled
	case model.PreOrderStatusPrepared:
		return to == model.PreOrderStatusPending || to == model.PreOrderStatusDelivered || to == model.PreOrderStatusCancelled
	default:
		return false
	}
}

func parsePreOrderTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	for _, layout := range []string{"2006-01-02 15:04:05", "2006-01-02 15:04"} {
		if parsed, err := time.ParseInLocation(layout, value, preOrderLocation); err == nil {
			return parsed, nil
		}
	}
	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return parsed.In(preOrderLocation), nil
	}
	return time.Time{}, apicode.Newf(apicode.InvalidDate, "计划配送时间格式无效")
}

func samePreOrderDate(left, right time.Time) bool {
	left = left.In(preOrderLocation)
	right = right.In(preOrderLocation)
	return left.Year() == right.Year() && left.Month() == right.Month() && left.Day() == right.Day()
}

func BuildPreOrderReminderTargets(now time.Time, slot string) []PreOrderReminderTarget {
	now = now.In(preOrderLocation)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, preOrderLocation)
	tomorrow := today.AddDate(0, 0, 1)
	dueKey := model.PreOrderReminderDueDay0930
	previousKey := model.PreOrderReminderPreviousDay0930
	if slot == PreOrderReminderSlot1600 {
		dueKey = model.PreOrderReminderDueDay1600
		previousKey = model.PreOrderReminderPreviousDay1600
	}
	return []PreOrderReminderTarget{
		{Start: tomorrow, End: tomorrow.AddDate(0, 0, 1), ReminderKey: previousKey, RelativeDate: "明天"},
		{Start: today, End: tomorrow, ReminderKey: dueKey, RelativeDate: "今天"},
	}
}

func PreOrderStatusAllowsReminder(status int8) bool {
	return status == model.PreOrderStatusPending || status == model.PreOrderStatusPrepared
}

func (s *PreOrderService) ProcessReminderSlot(now time.Time, slot string) error {
	var firstErr error
	for _, target := range BuildPreOrderReminderTargets(now, slot) {
		orders, err := s.preOrderModule.ListForReminder(target.Start, target.End)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		for _, order := range orders {
			if !PreOrderStatusAllowsReminder(order.Status) {
				continue
			}
			claimed, err := s.preOrderModule.ClaimReminder(order.ID, target.ReminderKey)
			if err != nil {
				if firstErr == nil {
					firstErr = err
				}
				continue
			}
			if !claimed {
				continue
			}
			sendErr := s.sendReminder(order, target.RelativeDate)
			if err := s.preOrderModule.CompleteReminder(order.ID, target.ReminderKey, sendErr); err != nil && firstErr == nil {
				firstErr = err
			}
			if sendErr != nil {
				logging.LogWarn("预订单钉钉提醒发送失败", zap.Uint("pre_order_id", order.ID), zap.String("reminder_key", target.ReminderKey), zap.Error(sendErr))
				if firstErr == nil {
					firstErr = sendErr
				}
			}
		}
	}
	return firstErr
}

func (s *PreOrderService) sendReminder(order *model.PreOrder, relativeDate string) error {
	bot, err := s.botModule.GetByStoreID(order.StoreID)
	if err != nil {
		return fmt.Errorf("get DingTalk bot: %w", err)
	}
	store, err := s.storeModule.GetByID(order.StoreID)
	if err != nil {
		return fmt.Errorf("get store: %w", err)
	}
	title := fmt.Sprintf("预订单提醒｜%s配送", relativeDate)
	text := buildPreOrderReminderMarkdown(order, store.Name, relativeDate)
	if strings.EqualFold(bot.BotType, "stream") {
		if strings.TrimSpace(store.Phone) == "" {
			return fmt.Errorf("store phone is required for stream reminder")
		}
		return s.dingTalkService.SendStreamMarkdownToMobile(bot, title, text, store.Phone)
	}
	return s.dingTalkService.SendMarkdownToBot(bot, title, text)
}

func buildPreOrderReminderMarkdown(order *model.PreOrder, storeName, relativeDate string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "### 预订单%s提醒\n\n", relativeDate)
	fmt.Fprintf(&b, "- **门店：** %s\n", storeName)
	fmt.Fprintf(&b, "- **客户：** %s\n", order.CustomerName)
	fmt.Fprintf(&b, "- **配送时间：** %s\n", order.ScheduledAt.In(preOrderLocation).Format("2006-01-02 15:04"))
	contact := strings.TrimSpace(strings.Join([]string{order.ContactPerson, order.ContactPhone}, " "))
	if contact != "" {
		fmt.Fprintf(&b, "- **联系人：** %s\n", contact)
	}
	if order.DeliveryAddress != "" {
		fmt.Fprintf(&b, "- **配送地址：** %s\n", order.DeliveryAddress)
	}
	b.WriteString("\n**备货明细**\n\n")
	for _, item := range order.Items {
		fmt.Fprintf(&b, "- %s / %s × %g", item.ProductName, item.UnitName, item.Quantity)
		if item.Remark != "" {
			fmt.Fprintf(&b, "（%s）", item.Remark)
		}
		b.WriteString("\n")
	}
	if order.Remark != "" {
		fmt.Fprintf(&b, "\n- **备注：** %s\n", order.Remark)
	}
	fmt.Fprintf(&b, "\n预订单号：%s", order.OrderNo)
	return b.String()
}
