# 数据范围（DataScope）盘点与约定 — D1 / D4

本文档对应重构路线 **D1（盘点）** 与 **D4（全库 / 超管单一真源）**，与代码路径一致时可随实现更新。

---

## D4：行级范围的单一真源（必读）

| 来源 | 职责 |
|------|------|
| `middleware.GetDataScope(c)` | **唯一**根据 Gin 上下文计算 `model.DataScope*` 数值 |
| `middleware.HQUnboundAdmin(c)` | **唯一**判断「总部未绑店可跨店选店」会话（`admin|super_admin` 且 `token.store_id==0`） |
| `middleware.ResolveQueryStoreID(c, key)` | 总部未绑店时解析 query 中的门店；绑店用户恒为 token 门店 |
| `authctx.Context.EffectiveDataScope` | 与 `GetDataScope` **同步**（在 `AuthMiddleware` 内赋值） |
| `service/list_authctx.go` | HTTP 列表将 `EffectiveDataScope` 写入各 `List*Req` 的 `DataScope` 字段 |

**禁止**：在 Controller / Service 用 `GetRoleCode() == super_admin` 自行推断「全库」行级范围；全库仅当 `GetDataScope == DataScopeAll`（由 `HQUnboundAdmin` 推导）。

---

## D1：统一列表 Scope（已收敛）

| 入口 | 实现 | 说明 |
|------|------|------|
| `pkg/datascope.ApplyPurchaseOrdersList` | `internal/datascope.PurchaseOrderListScope` | `purchase_orders` |
| `pkg/datascope.ApplyStoreAccountsList` | `internal/datascope.StoreAccountListScope` | `store_accounts`，本人=`operator_id` |
| `pkg/datascope.ApplyInventoriesList` | `internal/datascope.InventoryListScope` | 别名 `i`，SELF≈门店 |
| `pkg/datascope.ApplyInventoryOrdersList` | `internal/datascope.InventoryOrderListScope` | `inventory_orders` |

公共逻辑：`internal/datascope/list_scope_common.go` 的 `listDataScopeScope`。  
列名注册：`internal/datascope/policies.go` 的 `Policy*`（**D2**）。

---

## D1：仍含 `store_id` / 数据隔离逻辑的热点（待逐步收敛）

以下路径**尚未**统一走 `internal/datascope` 列表 Scope，后续可按业务拆 `Policy*` + `listDataScopeScope` 或 `DataPermission`：

| 区域 | 文件（示例） | 说明 |
|------|----------------|------|
| 统计 | `module/statistics.go` | 多处 `Where("store_id = ?", …)` |
| 会员 | `module/member.go` | 门店过滤 |
| 用户 | `module/user.go` | 门店用户列表/计数 |
| 供应商 | `module/supplier.go` | `store_suppliers` 关联 |
| 菜单 | `module/menu.go` | `store_role_menus.store_id` |
| 门店记账其它 | `module/store_account.go` | `GetStatsByDateRange` 等 |
| 库存写路径 | `module/inventory.go` | 事务内 `store_id` 条件（业务正确性，非列表 Scope） |
| Controller | `controller/statistics.go`、`user.go`、`store_supplier.go` 等 | `ResolveQueryStoreID` + 传参 |

钉钉机器人等无 HTTP `AuthContext` 的调用链：仍由调用方构造 `List*Req` 的 `StoreID` / `DataScope`（见 `service/dingtalk_command_handler.go`）。

---

## D3：模块与 Repository 约定

- **当前**：`module` 通过 `pkg/datascope.Apply*` 挂 Scope（已委托 `internal/datascope`）。
- **目标**：新增 `repository` 层后，由 Repository **唯一**组装 `db.Scopes(...)`，module 仅调 Repository；见 `internal/datascope/doc.go`。

---

## 索引建议（与 D1 热点表相关）

大表列表/统计在 `store_id`（及 `created_at`/`status`）上保持组合索引，避免「全表扫描 + 应用层筛店」。具体见 `cmd/apply_indexes` 或迁移脚本。
