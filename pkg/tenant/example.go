package tenant

/*
门店数据隔离使用示例

1. 在 module 层使用 Scopes：

	// 原来的写法
	func (m *PurchaseOrderModule) List(req *model.ListPurchaseOrderReq) ([]*model.PurchaseOrder, int64, error) {
		query := m.db.Model(&model.PurchaseOrder{})
		if req.StoreID > 0 {
			query = query.Where("store_id = ?", req.StoreID)
		}
		// ...
	}

	// 使用租户隔离后
	func (m *PurchaseOrderModule) List(tenant *TenantContext, req *model.ListPurchaseOrderReq) ([]*model.PurchaseOrder, int64, error) {
		query := m.db.Model(&model.PurchaseOrder{}).Scopes(TenantScopes(tenant, NewStoreIsolationStrategy()))
		// 不需要手动判断 store_id，自动隔离
		// ...
	}

2. 在 service 层使用：

	func (s *PurchaseOrderService) ListOrders(ctx context.Context, req *model.ListPurchaseOrderReq) ([]*model.PurchaseOrder, int64, error) {
		tenant := GetTenant(ctx)
		return s.orderModule.List(tenant, req)
	}

3. 在 controller 层使用：

	func (c *PurchaseOrderController) ListOrders(ctx *gin.Context) {
		tenant := GetTenantFromGin(ctx)
		// 或者从 request context 获取
		// tenant := GetTenant(ctx.Request.Context())

		orders, total, err := c.orderService.ListOrders(tenant, &req)
		// ...
	}

4. 使用泛型仓储：

	// 创建带隔离的仓储
	orderRepo := NewTenantRepository[model.PurchaseOrder](db, NewStoreIsolationStrategy())

	// 查询时自动隔离
	orders, err := orderRepo.FindAll(tenant)

	// 更新时验证权限
	err := orderRepo.Update(tenant, orderID, map[string]interface{}{"status": 2})
	// 如果订单不属于该门店，返回 ErrAccessDenied

5. 便捷的 Scope 函数：

	// 门店隔离
	db.Scopes(StoreScope(storeID)).Find(&orders)

	// 供应商隔离
	db.Scopes(SupplierScope(supplierID)).Find(&products)

	// 管理员（无隔离）
	db.Scopes(AdminScope()).Find(&allOrders)

6. 组合多种隔离策略：

	// 同时按门店和供应商隔离
	composite := NewCompositeIsolationStrategy(
		NewStoreIsolationStrategy(),
		NewSupplierIsolationStrategy(),
	)
	db.Scopes(TenantScopes(tenant, composite)).Find(&items)
*/

// ExampleUsage 示例用法（仅供参考，不实际运行）
func ExampleUsage() {
	// 这个函数仅作为文档示例
}
