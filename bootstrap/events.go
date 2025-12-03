package bootstrap

import (
	"fmt"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils"
	"github.com/Kevin-Jii/tower-go/utils/logging"
)

// InitEventSubscribers 初始化事件订阅
func InitEventSubscribers() {
	// 订阅订单创建事件
	utils.GlobalEventBus.Subscribe(utils.EventOrderCreated, func(data interface{}) {
		if order, ok := data.(*model.PurchaseOrder); ok {
			logging.LogInfo(fmt.Sprintf("📦 新采购单创建: %s, 门店ID: %d", order.OrderNo, order.StoreID))
			// 可以在这里添加更多逻辑，如发送通知等
		}
	})

	// 订阅订单确认事件
	utils.GlobalEventBus.Subscribe(utils.EventOrderConfirmed, func(data interface{}) {
		logging.LogInfo("✅ 采购单已确认")
		// 可以在这里添加通知供应商的逻辑
	})

	// 订阅订单完成事件
	utils.GlobalEventBus.Subscribe(utils.EventOrderCompleted, func(data interface{}) {
		logging.LogInfo("🎉 采购单已完成")
		// 可以在这里添加统计、报表等逻辑
	})

	// 订阅订单取消事件
	utils.GlobalEventBus.Subscribe(utils.EventOrderCancelled, func(data interface{}) {
		logging.LogInfo("❌ 采购单已取消")
		// 可以在这里添加库存回滚等逻辑
	})

	// 订阅供应商绑定事件
	utils.GlobalEventBus.Subscribe(utils.EventSupplierBound, func(data interface{}) {
		if info, ok := data.(map[string]interface{}); ok {
			logging.LogInfo(fmt.Sprintf("🔗 门店 %v 绑定了供应商: %v", info["store_id"], info["supplier_ids"]))
		}
	})

	fmt.Println("📡 事件订阅初始化完成")
}
