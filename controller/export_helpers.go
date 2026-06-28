package controller

import (
	"fmt"
	"strings"

	"github.com/Kevin-Jii/tower-go/model"
)

const exportPageSize = 100000

func exportDateQuery(ctxDate, fallbackStart, fallbackEnd string) (string, string, string) {
	date := strings.TrimSpace(ctxDate)
	if date != "" {
		return date, date, date
	}
	return strings.TrimSpace(fallbackStart), strings.TrimSpace(fallbackEnd), ""
}

func storeName(store *model.Store) string {
	if store == nil || strings.TrimSpace(store.Name) == "" {
		return "-"
	}
	return store.Name
}

func userName(user *model.User, fallback string) string {
	if user != nil {
		if strings.TrimSpace(user.Nickname) != "" {
			return user.Nickname
		}
		if strings.TrimSpace(user.Username) != "" {
			return user.Username
		}
		if strings.TrimSpace(user.Phone) != "" {
			return user.Phone
		}
	}
	if strings.TrimSpace(fallback) != "" {
		return fallback
	}
	return "-"
}

func memberName(member *model.Member) string {
	if member == nil {
		return "-"
	}
	if strings.TrimSpace(member.Name) != "" && strings.TrimSpace(member.Phone) != "" {
		return fmt.Sprintf("%s(%s)", member.Name, member.Phone)
	}
	if strings.TrimSpace(member.Name) != "" {
		return member.Name
	}
	if strings.TrimSpace(member.Phone) != "" {
		return member.Phone
	}
	return fmt.Sprintf("会员%d", member.ID)
}

func storeAccountPaymentLabel(status int) string {
	if status == model.StoreAccountPaymentPaid {
		return "已支付"
	}
	if status == model.StoreAccountPaymentUnpaid {
		return "未支付"
	}
	return fmt.Sprintf("未知(%d)", status)
}

func b2bPaymentLabel(status int) string {
	switch status {
	case model.B2BPaymentUnpaid:
		return "未收"
	case model.B2BPaymentPartial:
		return "部分收款"
	case model.B2BPaymentPaid:
		return "已收"
	default:
		return fmt.Sprintf("未知(%d)", status)
	}
}

func b2bDeliveryLabel(status int) string {
	switch status {
	case model.B2BDeliveryPending:
		return "待配送"
	case model.B2BDeliveryDone:
		return "已配送"
	case model.B2BDeliveryCancel:
		return "已取消"
	default:
		return fmt.Sprintf("未知(%d)", status)
	}
}

func inventoryTypeLabel(t int8) string {
	if t == model.InventoryTypeIn {
		return "入库"
	}
	if t == model.InventoryTypeOut {
		return "出库"
	}
	return fmt.Sprintf("未知(%d)", t)
}

func inventoryLossTypeLabel(t string) string {
	switch t {
	case model.InventoryLossTypeLoss:
		return "报损"
	case model.InventoryLossTypeSelfUse:
		return "自用"
	case model.InventoryLossTypeGift:
		return "赠送"
	default:
		if strings.TrimSpace(t) == "" {
			return "-"
		}
		return t
	}
}

func yesNo(v bool) string {
	if v {
		return "是"
	}
	return "否"
}

func accountItemsText(items []model.StoreAccountItem) string {
	parts := make([]string, 0, len(items))
	for _, item := range items {
		name := strings.TrimSpace(item.ProductName)
		if name == "" {
			name = fmt.Sprintf("商品%d", item.ProductID)
		}
		unit := strings.TrimSpace(item.Unit)
		if unit == "" {
			unit = strings.TrimSpace(item.Spec)
		}
		parts = append(parts, fmt.Sprintf("%s x%g%s %s", name, item.Quantity, unit, formatAmount(item.Amount)))
	}
	return strings.Join(parts, "；")
}

func accountConsumablesText(items []model.StoreAccountConsumable) string {
	parts := make([]string, 0, len(items))
	for _, item := range items {
		name := strings.TrimSpace(item.ProductName)
		if name == "" {
			name = fmt.Sprintf("消耗品%d", item.ProductID)
		}
		parts = append(parts, fmt.Sprintf("%s x%g %s", name, item.Quantity, formatAmount(item.Amount)))
	}
	return strings.Join(parts, "；")
}

func b2bItemsText(items []model.B2BSupplyOrderItem) string {
	parts := make([]string, 0, len(items))
	for _, item := range items {
		parts = append(parts, fmt.Sprintf("%s/%s x%g %s", item.ProductName, item.UnitName, item.Quantity, formatAmount(item.Amount)))
	}
	return strings.Join(parts, "；")
}

func formatAmount(v float64) string {
	return fmt.Sprintf("%.2f", v)
}
