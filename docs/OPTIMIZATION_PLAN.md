# Tower-Go 优化方案

本文面向 Tower-Go 当前代码形态：Gin 分层、JWT/RBAC、门店 DataScope、进销存/记账核心业务。目标是在不推翻现有三层架构的前提下，逐步提升模块边界、数据模型表达力、查询性能、算法复杂度和可维护性。

## 1. 当前判断

### 已有基础

- `controller -> service -> module` 分层已经稳定，适合继续渐进演进。
- `AuthContext / DataScope / HQUnboundAdmin / ResolveQueryStoreID` 已经形成门店数据权限主线。
- `internal/datascope` 已开始承载统一 Scope 逻辑，部分列表已迁入。
- `pkg/performance`、`pkg/tenant`、`pkg/repository` 已有雏形，但还没有成为主业务统一入口。
- 核心业务对象已经比较清晰：供应商、门店供应商、采购单、库存、出入库单、门店记账、会员、菜单权限。

### 主要问题

- Controller 中仍有较多权限和门店范围解析逻辑，业务规则散落。
- `module` 既承担数据访问，又承担部分查询组合，Repository 边界不够清晰。
- 多租户仍是“业务层约束”为主，尚未形成全库统一查询入口。
- 部分数据结构通过状态码、位图、字符串字段承载复杂语义，可读性和演进性有限。
- 列表查询、统计查询、详情归属校验还没有完全统一抽象。
- 缓存、索引、分页、批量查询策略没有形成系统化规范。

## 2. 总体优化目标

1. 模块化：让权限、租户范围、业务规则、数据访问各自归位。
2. 数据结构：让核心实体关系更清晰，避免字段语义过载。
3. 数据算法：降低热路径查询复杂度，减少 N+1、重复计算和无索引扫描。
4. 可演进：为后续总部、多门店、组织树、独立权限表、审计、报表打基础。
5. 稳定性：每一步都可小范围提交、测试、回滚，不做一次性大重构。

## 3. 模块化优化

### 3.1 引入 Repository 作为数据访问统一出口

当前 `module` 层直接拼接查询、Scope、Preload、事务。建议新增 `repository` 层，逐步迁移高风险模块。

目标分层：

```text
controller: 参数绑定、HTTP 响应
service: 业务流程、事务编排、领域规则
repository: GORM 查询、DataScope、Preload、分页、锁
module: 逐步弱化，最终作为旧接口兼容层或移除
internal/datascope: 行级权限 Scope 单一实现
```

优先迁移顺序：

1. `purchase_order`
2. `inventory / inventory_order`
3. `store_account`
4. `member`
5. `supplier / store_supplier`
6. `user / menu`

建议接口示例：

```go
type PurchaseOrderRepository interface {
    FindByID(ctx context.Context, id uint) (*model.PurchaseOrder, error)
    List(ctx context.Context, req *model.ListPurchaseOrderReq) ([]*model.PurchaseOrder, int64, error)
    UpdateStatus(ctx context.Context, id uint, status int8, remark string) error
}
```

落地标准：

- Repository 方法必须接收 `context.Context`。
- Repository 内部统一从 `authctx.Context` 或显式 `Scope` 参数应用数据范围。
- Service 不再手写 `store_id` 查询条件。
- Controller 不再判断“是否跨店可见”，只调用 service。

### 3.2 权限模块独立化

当前权限由菜单权限码、角色菜单、门店角色菜单共同承载。建议拆出 `internal/permission` 或 `service/permission` 聚合规则。

优化点：

- `PermissionResolver`：统一解析用户最终权限码。
- `DataScopeResolver`：统一解析行级范围。
- `PermissionCache`：统一管理权限缓存 key、TTL、失效策略。
- `PermissionGuard`：统一提供“是否可进入接口”判断。

中期可引入独立表：

```text
permissions
- id
- code
- name
- resource
- action
- status

role_permissions
- role_id
- permission_id

store_role_permissions
- store_id
- role_id
- permission_id
```

菜单只负责导航展示，权限表负责接口/按钮权限。

### 3.3 领域服务拆分

部分 service 已经变大，建议按领域动作拆分：

```text
service/purchase_order.go
service/purchase_order_state.go
service/purchase_order_notice.go
service/purchase_order_pricing.go

service/store_account.go
service/store_account_payment.go
service/store_account_member.go
service/store_account_notice.go
service/store_account_stats.go
```

拆分原则：

- 创建、状态流转、通知、统计不要混在一个文件里。
- 外部系统对接，如钉钉、打印机、对象存储，放到 adapter 或 integration 包。
- 领域事件统一通过 event bus，不在核心业务里散落 goroutine。

### 3.4 统一错误类型

当前多处直接返回 `errors.New` 字符串，Controller 再统一 500，容易把业务错误误报成服务器错误。

建议新增：

```go
type AppError struct {
    Code int
    Message string
    Cause error
}
```

常用错误：

- `ErrForbidden`
- `ErrNotFound`
- `ErrValidation`
- `ErrConflict`
- `ErrInsufficientStock`
- `ErrInvalidStateTransition`

Controller 只做 `http.WriteError(ctx, err)`。

## 4. 数据结构优化

### 4.1 多租户与组织结构

当前核心隔离维度是 `store_id`。如果后续支持总部、区域、品牌、加盟商，建议引入组织树。

建议表：

```text
orgs
- id
- parent_id
- type          -- hq / region / store_group / store
- code
- name
- path          -- /1/3/9/
- depth
- status

stores
- id
- org_id
- store_code
- name
...
```

优势：

- 支持区域经理看多个门店。
- 支持总部、分区、门店层级报表。
- DataScope 可以从 `store_id` 扩展为 `org_scope`。

建议新增 DataScope：

```text
1 all       全部
2 org       本组织及子组织
3 store     本门店
4 self      本人
```

### 4.2 库存数据结构

当前库存核心是 `inventories` 加出入库单。建议进一步明确“当前库存”和“库存流水”。

建议表：

```text
inventories
- store_id
- product_id
- quantity
- unit
- version

inventory_transactions
- store_id
- product_id
- order_id
- order_no
- direction      -- in / out
- quantity
- before_qty
- after_qty
- operator_id
- occurred_at
```

优势：

- 当前库存用于快速查询。
- 流水用于审计、追溯、回滚、统计。
- 乐观锁或行锁可以集中处理。

### 4.3 记账数据结构

门店记账建议拆出支付与会员绑定关系，减少主表字段膨胀。

建议：

```text
store_accounts
- id
- account_no
- store_id
- account_date
- operator_id
- total_amount
- payment_status
- member_id
- source
- remark

store_account_items
- account_id
- product_id
- custom_name
- quantity
- unit
- unit_price
- amount

store_account_payments
- account_id
- pay_type
- amount
- paid_at
- transaction_no
```

优化点：

- 自定义记账内容使用 `custom_name`，不要强制依赖商品。
- 支付状态与支付明细分离，方便部分支付、组合支付、补录支付。
- 会员绑定保留在主表，但会员消费流水通过事件生成。

### 4.4 权限数据结构

当前 `menus.permission` 与权限位混用，短期可继续用，长期建议菜单与权限解耦。

短期：

- 保留 `menus.permission`。
- 统一权限码命名：`domain:resource:action`，如 `purchase:order:list`。
- 建立权限码常量或生成器，避免字符串散落。

长期：

- `permissions` 存接口/按钮能力。
- `menus` 存导航结构。
- `menu_permissions` 建立菜单和权限的展示关系。

### 4.5 字段与索引规范

建议统一：

- 所有门店业务表保留 `store_id`。
- 所有操作类表保留 `operator_id / created_by`。
- 所有单据表保留 `order_no` 唯一索引。
- 所有列表热字段建立组合索引。

推荐索引：

```text
purchase_orders:       (store_id, order_date, status), (store_id, created_at), unique(order_no)
inventory_orders:      (store_id, type, created_at), unique(order_no)
inventories:           unique(store_id, product_id)
store_accounts:        (store_id, account_date), (store_id, member_id), unique(account_no)
members:               unique(store_id, phone), (store_id, uid)
store_suppliers:       unique(store_id, supplier_id), (supplier_id)
supplier_products:     (supplier_id, category_id, status)
users:                 (store_id, role_id), unique(phone)
wallet_logs:           (member_id, created_at)
```

## 5. 数据算法与查询性能优化

### 5.1 列表查询统一算法

所有列表建议统一为：

```text
1. 解析分页参数
2. 应用 DataScope
3. 应用业务过滤
4. Count
5. 查询 ID 页
6. 按 ID 批量 Preload 或 Join 明细
```

对大表避免：

- 深分页 `OFFSET` 过大。
- Count 和复杂 Join 混在一起。
- Preload 大量无用字段。

高数据量时改为游标分页：

```text
?cursor_id=12345&page_size=50
WHERE id < cursor_id
ORDER BY id DESC
LIMIT 50
```

### 5.2 消除 N+1 查询

重点检查：

- 采购单列表及明细。
- 出入库单列表及商品。
- 门店记账列表及商品/会员/操作人。
- 会员流水和充值单。
- 菜单树与权限码。

优化方式：

- 先查主表 ID。
- 批量查关联表。
- 在 service 层用 map 聚合。
- 大列表避免直接 `Preload("Items.Product.Supplier...")`。

### 5.3 菜单树算法

菜单树当前已有优化思路，应统一使用 O(n) 构建：

```text
1. 一次查询所有菜单
2. map[id]*node
3. 按 parent_id 挂载
4. 对每层按 sort 排序
```

缓存策略：

- 全局菜单树缓存。
- `store_id + role_id` 菜单树缓存。
- `user_id + store_id + role_id` 权限码缓存。
- 菜单/角色/门店角色变更时统一失效。

### 5.4 库存扣减算法

出库必须在事务内保证：

```sql
UPDATE inventories
SET quantity = quantity - ?
WHERE store_id = ? AND product_id = ? AND quantity >= ?
```

并检查 `RowsAffected`。

入库建议：

```sql
INSERT ... ON DUPLICATE KEY UPDATE quantity = quantity + VALUES(quantity)
```

如果数据库为 PostgreSQL，则使用 `ON CONFLICT`。

收益：

- 避免先查后改产生并发窗口。
- 减少锁持有时间。
- 简化库存一致性判断。

### 5.5 统计查询优化

统计类接口建议独立为 read model：

```text
daily_store_stats
- store_id
- stat_date
- purchase_amount
- account_amount
- inventory_in_qty
- inventory_out_qty
- member_recharge_amount
```

生成方式：

- 短期：定时任务每日汇总。
- 中期：业务事件增量更新。
- 长期：异步队列 + 幂等消费。

适用场景：

- 首页 dashboard。
- 门店日报。
- 总部跨店统计。
- 钉钉报表推送。

### 5.6 搜索优化

当前 keyword 搜索多处 `LIKE '%keyword%'`，数据量上来后会慢。

建议：

- 手机号、编码：前缀匹配或精确匹配。
- 商品名/供应商名：维护 normalized 字段。
- 中文搜索：可考虑 MySQL FULLTEXT、PostgreSQL trigram，或接入 Meilisearch/Elasticsearch。
- 搜索条件统一封装为 `search.Filter`，避免各 module 手写。

## 6. 缓存优化

### 6.1 适合缓存的数据

- 菜单树。
- 权限码。
- 门店信息。
- 字典配置。
- 消息模板。
- 商品单位规格。
- 打印机配置。
- 钉钉机器人配置。

### 6.2 不建议直接缓存的数据

- 当前库存数量。
- 会员余额。
- 支付状态。
- 单据状态。

这类数据可以缓存只读视图，但写入必须以数据库事务为准。

### 6.3 缓存 key 规范

```text
tower:{domain}:{entity}:{scope}:{id}

tower:perm:user:{user_id}:{store_id}:{role_id}
tower:menu:store-role:{store_id}:{role_id}
tower:store:{store_id}
tower:dict:{type}
```

## 7. 事务与一致性优化

### 7.1 事务边界

必须放在一个事务内：

- 创建采购单 + 明细。
- 创建出入库单 + 更新库存。
- 记账成功 + 明细 + 会员消费流水。
- 充值支付 + 会员余额 + 钱包流水。
- 删除单据 + 删除明细。

建议由 service 统一开启事务，并将 tx 注入 repository。

### 7.2 幂等设计

外部回调、钉钉消息、打印任务、支付状态变更都应幂等。

建议增加：

```text
idempotency_keys
- key
- business_type
- business_id
- status
- request_hash
- created_at
```

## 8. 接口与 DTO 优化

### 8.1 Request/Response 分离

当前部分接口直接返回 model。建议长期改为 response DTO：

```text
model: 数据库结构
dto/request: 入参结构
dto/response: 出参结构
```

优势：

- 避免数据库字段泄漏。
- 兼容前端字段演进。
- 减少 `gorm:"-"` 临时字段污染 model。

### 8.2 分页响应统一

统一结构：

```json
{
  "list": [],
  "total": 0,
  "page": 1,
  "page_size": 20
}
```

避免有的接口 data 是数组，有的接口 data 是对象加 meta。

## 9. 可观测性优化

建议加入：

- 请求 trace id。
- 慢查询日志。
- 关键业务事件日志。
- 外部系统调用耗时。
- 库存/余额变更审计日志。

核心指标：

```text
接口耗时 P50/P95/P99
数据库查询耗时
Redis 命中率
库存扣减失败率
钉钉/打印发送失败率
登录失败率
权限拒绝次数
```

## 10. 分阶段路线图

### P0：1-3 天，风险收敛

- 扫描并替换 Controller 中剩余不合适的 `IsAdmin` 数据范围判断。
- 给详情/写路径补齐本店归属校验。
- 给供应商、会员、记账、库存、采购单加关键索引。
- 统一错误返回，业务错误不要返回 500。
- 移除或忽略编译产物 `cmd/tower-go`，避免误提交。

### P1：1 周，统一数据权限

- 把 `member/user/supplier/statistics/menu` 列表逐步迁入 `internal/datascope`。
- Service 列表方法统一接收 `context.Context`。
- Controller 只负责 `AttachAuthContextToHTTPRequest`，不再手动注入 DataScope。
- 新增 Repository 接口，从采购单和库存开始试点。

### P2：2-3 周，核心业务结构优化

- 建立库存流水表。
- 记账支付明细表独立。
- 菜单权限和接口权限开始解耦。
- 统计 read model 落地。
- 权限缓存失效策略集中化。

### P3：1-2 月，多租户增强

- 引入 `orgs` 组织树。
- DataScope 从 store 扩展为 org/store/self。
- Repository 层成为全库查询统一入口。
- 独立 permissions 表。
- 审计日志与幂等表完善。

## 11. 优先级建议

最高优先级：

1. 数据权限一致性。
2. 库存/余额/支付事务一致性。
3. 热列表索引与 N+1 查询。
4. 错误类型和接口返回规范。

中优先级：

1. Repository 渐进迁移。
2. DTO 分离。
3. 统计 read model。
4. 菜单与权限解耦。

低优先级：

1. 完整组织树。
2. 搜索引擎。
3. 插件化架构。
4. CQRS/事件溯源。

## 12. 验收标准

### 模块化

- 新增业务不需要在 Controller 手写门店过滤。
- 查询入口统一经过 Repository 或 datascope。
- 权限逻辑不再散落在多个 Controller。

### 数据结构

- 关键业务表有唯一约束和组合索引。
- 库存、余额、支付状态有可追溯流水。
- 菜单和权限职责清晰。

### 数据算法

- 热点列表无明显 N+1 查询。
- 大表分页支持游标或 ID 子查询。
- 库存扣减为单 SQL 条件更新。
- 首页统计不直接扫描大业务表。

### 稳定性

- `go test ./middleware ./controller ./service ./module ./router/api` 通过。
- 核心接口有集成测试。
- 关键事务有并发测试。
- 权限拒绝、跨店访问、库存不足等场景有测试覆盖。
