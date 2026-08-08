package module

import (
	"fmt"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PreOrderModule struct {
	db *gorm.DB
}

func NewPreOrderModule(db *gorm.DB) *PreOrderModule {
	return &PreOrderModule{db: db}
}

func (m *PreOrderModule) GenerateOrderNo(now time.Time) string {
	return fmt.Sprintf("YD%s%04d", now.Format("20060102150405"), now.Nanosecond()%10000)
}

func (m *PreOrderModule) Create(order *model.PreOrder) error {
	return m.db.Create(order).Error
}

func (m *PreOrderModule) GetByID(id uint) (*model.PreOrder, error) {
	var order model.PreOrder
	err := m.db.
		Preload("Store").
		Preload("Customer").
		Preload("Creator").
		Preload("Items").
		Preload("ReminderLogs", func(db *gorm.DB) *gorm.DB { return db.Order("created_at ASC") }).
		First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (m *PreOrderModule) List(req *model.ListPreOrderReq) ([]*model.PreOrder, int64, error) {
	var rows []*model.PreOrder
	var total int64
	q := m.db.Model(&model.PreOrder{})
	if req.StoreID > 0 {
		q = q.Where("pre_orders.store_id = ?", req.StoreID)
	}
	if req.CustomerID > 0 {
		q = q.Where("pre_orders.customer_id = ?", req.CustomerID)
	}
	if req.Status != nil {
		q = q.Where("pre_orders.status = ?", *req.Status)
	}
	if req.StartDate != "" {
		q = q.Where("pre_orders.scheduled_at >= ?", req.StartDate+" 00:00:00")
	}
	if req.EndDate != "" {
		q = q.Where("pre_orders.scheduled_at < DATE_ADD(?, INTERVAL 1 DAY)", req.EndDate)
	}
	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("(pre_orders.order_no LIKE ? OR pre_orders.customer_name LIKE ? OR pre_orders.contact_phone LIKE ? OR pre_orders.delivery_address LIKE ?)", like, like, like, like)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (req.Page - 1) * req.PageSize
	if err := q.
		Preload("Store").
		Preload("Customer").
		Preload("Creator").
		Preload("Items").
		Preload("ReminderLogs", func(db *gorm.DB) *gorm.DB { return db.Order("created_at ASC") }).
		Order("CASE WHEN pre_orders.status IN (1, 2) THEN 0 ELSE 1 END ASC").
		Order("pre_orders.scheduled_at ASC, pre_orders.id DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (m *PreOrderModule) Update(order *model.PreOrder, items []model.PreOrderItem, resetReminders bool) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"customer_id":      order.CustomerID,
			"customer_name":    order.CustomerName,
			"contact_person":   order.ContactPerson,
			"contact_phone":    order.ContactPhone,
			"delivery_address": order.DeliveryAddress,
			"scheduled_at":     order.ScheduledAt,
			"remark":           order.Remark,
		}
		if err := tx.Model(&model.PreOrder{}).Where("id = ?", order.ID).Updates(updates).Error; err != nil {
			return err
		}
		if err := tx.Where("pre_order_id = ?", order.ID).Delete(&model.PreOrderItem{}).Error; err != nil {
			return err
		}
		if resetReminders {
			if err := tx.Where("pre_order_id = ?", order.ID).Delete(&model.PreOrderReminderLog{}).Error; err != nil {
				return err
			}
		}
		if len(items) == 0 {
			return nil
		}
		return tx.Create(&items).Error
	})
}

func (m *PreOrderModule) UpdateStatus(id uint, status int8, at time.Time) error {
	updates := map[string]interface{}{"status": status}
	switch status {
	case model.PreOrderStatusPending:
		updates["prepared_at"] = nil
		updates["delivered_at"] = nil
		updates["cancelled_at"] = nil
	case model.PreOrderStatusPrepared:
		updates["prepared_at"] = at
		updates["delivered_at"] = nil
		updates["cancelled_at"] = nil
	case model.PreOrderStatusDelivered:
		updates["delivered_at"] = at
		updates["cancelled_at"] = nil
	case model.PreOrderStatusCancelled:
		updates["cancelled_at"] = at
	}
	return m.db.Model(&model.PreOrder{}).Where("id = ?", id).Updates(updates).Error
}

func (m *PreOrderModule) Delete(id uint) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("pre_order_id = ?", id).Delete(&model.PreOrderItem{}).Error; err != nil {
			return err
		}
		if err := tx.Where("pre_order_id = ?", id).Delete(&model.PreOrderReminderLog{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.PreOrder{}, id).Error
	})
}

func (m *PreOrderModule) ListForReminder(start, end time.Time) ([]*model.PreOrder, error) {
	var rows []*model.PreOrder
	err := m.db.
		Where("scheduled_at >= ? AND scheduled_at < ?", start, end).
		Where("status IN ?", []int8{model.PreOrderStatusPending, model.PreOrderStatusPrepared}).
		Preload("Items").
		Order("scheduled_at ASC, id ASC").
		Find(&rows).Error
	return rows, err
}

// ClaimReminder uses the unique order/key index to ensure only one process sends a reminder.
func (m *PreOrderModule) ClaimReminder(preOrderID uint, reminderKey string) (bool, error) {
	log := &model.PreOrderReminderLog{
		PreOrderID:  preOrderID,
		ReminderKey: reminderKey,
		Status:      model.PreOrderReminderSending,
	}
	result := m.db.Clauses(clause.OnConflict{DoNothing: true}).Create(log)
	return result.RowsAffected == 1, result.Error
}

func (m *PreOrderModule) CompleteReminder(preOrderID uint, reminderKey string, sendErr error) error {
	updates := map[string]interface{}{}
	if sendErr == nil {
		now := time.Now()
		updates["status"] = model.PreOrderReminderSent
		updates["sent_at"] = now
		updates["error_message"] = ""
	} else {
		updates["status"] = model.PreOrderReminderFailed
		message := sendErr.Error()
		if len(message) > 500 {
			message = message[:500]
		}
		updates["error_message"] = message
	}
	return m.db.Model(&model.PreOrderReminderLog{}).
		Where("pre_order_id = ? AND reminder_key = ?", preOrderID, reminderKey).
		Updates(updates).Error
}
