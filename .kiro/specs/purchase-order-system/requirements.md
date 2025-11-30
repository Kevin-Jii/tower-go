# Requirements Document

## Introduction

本系统是一个门店采购报菜单管理系统，用于管理门店与供应商之间的商品采购关系。系统支持门店从多个供应商采购商品，同一商品可以有多个供应商提供，门店可以设置默认供应商。系统自动根据门店配置的供应商关系，将报菜单明细归类到对应供应商。

## Glossary

- **Store（门店）**: 进行采购的业务单位
- **Supplier（供应商）**: 提供商品的供货商
- **SupplierCategory（供应商分类）**: 供应商自定义的商品分类，如"鱼类"、"蔬菜"
- **SupplierProduct（供应商商品）**: 供应商提供的具体商品，包含价格、单位等信息
- **StoreSupplierProduct（门店供应商商品关联）**: 门店与供应商商品的绑定关系，标记门店从哪个供应商进哪些货
- **PurchaseOrder（采购单/报菜单）**: 门店提交的采购订单
- **PurchaseOrderItem（采购单明细）**: 采购单中的具体商品项
- **DefaultSupplier（默认供应商）**: 同一商品有多个供应商时，门店标记的首选供应商

## Data Model Relationships

```
┌─────────────┐     ┌──────────────────┐     ┌─────────────┐
│   stores    │     │ supplier_categories │   │  suppliers  │
│   门店表    │     │   供应商分类表      │     │  供应商表   │
└─────────────┘     └──────────────────┘     └─────────────┘
                            │                       │
                            └───────────┬───────────┘
                                        ▼
                           ┌──────────────────┐
                           │ supplier_products │
                           │   供应商商品表    │
                           │ (供应商+分类+商品) │
                           └──────────────────┘
                                        │
                                        ▼
┌─────────────┐     ┌──────────────────────────┐
│   stores    │────▶│ store_supplier_products  │
│   门店表    │     │    门店供应商商品关联     │
└─────────────┘     │   (门店从哪家进哪些货)    │
                    │   is_default: 默认供应商  │
                    └──────────────────────────┘
                                        │
                                        ▼
                    ┌──────────────────────────┐
                    │     purchase_orders      │
                    │       采购报菜单         │
                    └──────────────────────────┘
                                        │
                                        ▼
                    ┌──────────────────────────┐
                    │   purchase_order_items   │
                    │       采购单明细         │
                    │  (supplier_id+product_id) │
                    └──────────────────────────┘
```

### Business Scenario Example

```
门店A ──┬── 供应商1(鱼类) ── 黑鱼、草鱼
        └── 供应商2(蔬菜) ── 白菜、萝卜

门店B ──┬── 供应商1(鱼类) ── 黑鱼
        └── 供应商3(鱼类) ── 草鱼  ← 同品类不同供应商
```

### Key Design Points

1. **同商品多供应商**: `store_supplier_products.is_default` 标记默认供应商
2. **报菜时自动匹配**: 根据门店配置的供应商关系，自动归类到对应供应商
3. **价格可能不同**: 同一商品不同供应商价格可能不同，存在 `supplier_products` 中

## Requirements

### Requirement 1

**User Story:** As a 门店管理员, I want to 绑定供应商商品到门店, so that 门店可以从指定供应商采购商品。

#### Acceptance Criteria

1. WHEN 门店管理员选择供应商商品并提交绑定请求 THEN the StoreSupplierProduct SHALL 创建门店与商品的关联记录
2. WHEN 门店管理员绑定已存在的商品关联 THEN the StoreSupplierProduct SHALL 保持现有记录不变并返回成功
3. WHEN 门店管理员解绑供应商商品 THEN the StoreSupplierProduct SHALL 删除对应的关联记录
4. WHEN 门店绑定同名商品的多个供应商 THEN the StoreSupplierProduct SHALL 允许创建多条关联记录

### Requirement 2

**User Story:** As a 门店管理员, I want to 设置默认供应商, so that 报菜时系统自动选择首选供应商。

#### Acceptance Criteria

1. WHEN 门店管理员设置某商品的默认供应商 THEN the StoreSupplierProduct SHALL 将该商品的is_default标记为true
2. WHEN 门店管理员设置新的默认供应商 THEN the StoreSupplierProduct SHALL 将同名商品的其他供应商is_default标记为false
3. WHEN 门店首次绑定某商品 THEN the StoreSupplierProduct SHALL 自动将该商品设为默认供应商

### Requirement 3

**User Story:** As a 门店员工, I want to 创建采购单, so that 我可以向供应商报菜。

#### Acceptance Criteria

1. WHEN 门店员工提交采购单请求 THEN the PurchaseOrder SHALL 生成唯一订单编号并创建采购单记录
2. WHEN 采购单包含商品明细 THEN the PurchaseOrderItem SHALL 根据商品ID自动关联对应的供应商ID
3. WHEN 采购单创建成功 THEN the PurchaseOrder SHALL 计算并存储所有明细的总金额
4. WHEN 采购单明细中的商品不在门店绑定列表中 THEN the PurchaseOrder SHALL 拒绝创建并返回错误信息
5. WHEN 采购单创建时 THEN the PurchaseOrder SHALL 记录创建人ID和报菜日期

### Requirement 4

**User Story:** As a 门店管理员, I want to 按供应商分组查看采购单明细, so that 我可以分别向各供应商发送订单。

#### Acceptance Criteria

1. WHEN 查询采购单详情 THEN the PurchaseOrder SHALL 返回按供应商分组的商品明细列表
2. WHEN 采购单包含多个供应商的商品 THEN the PurchaseOrderItem SHALL 按供应商ID分组展示
3. WHEN 展示分组明细 THEN the PurchaseOrderItem SHALL 包含供应商名称、商品名称、数量、单价、金额

### Requirement 5

**User Story:** As a 门店管理员, I want to 管理采购单状态, so that 我可以跟踪采购流程。

#### Acceptance Criteria

1. WHEN 采购单创建后 THEN the PurchaseOrder SHALL 初始状态为"待确认"
2. WHEN 门店管理员确认采购单 THEN the PurchaseOrder SHALL 状态变更为"已确认"
3. WHEN 采购完成 THEN the PurchaseOrder SHALL 状态变更为"已完成"
4. WHEN 门店管理员取消采购单 THEN the PurchaseOrder SHALL 状态变更为"已取消"
5. WHEN 删除采购单 THEN the PurchaseOrder SHALL 仅允许删除"待确认"或"已取消"状态的订单

### Requirement 6

**User Story:** As a 门店管理员, I want to 查询采购单列表, so that 我可以查看历史采购记录。

#### Acceptance Criteria

1. WHEN 查询采购单列表 THEN the PurchaseOrder SHALL 支持按门店ID、供应商ID、状态、日期范围筛选
2. WHEN 返回列表结果 THEN the PurchaseOrder SHALL 支持分页查询
3. WHEN 展示列表 THEN the PurchaseOrder SHALL 包含订单编号、门店名称、总金额、状态、报菜日期

### Requirement 7

**User Story:** As a 供应商管理员, I want to 管理商品分类, so that 我可以组织商品目录。

#### Acceptance Criteria

1. WHEN 供应商创建分类 THEN the SupplierCategory SHALL 创建属于该供应商的分类记录
2. WHEN 供应商更新分类 THEN the SupplierCategory SHALL 支持修改名称和排序
3. WHEN 供应商禁用分类 THEN the SupplierCategory SHALL 将状态设为禁用但保留记录

### Requirement 8

**User Story:** As a 供应商管理员, I want to 管理商品信息, so that 门店可以选择采购。

#### Acceptance Criteria

1. WHEN 供应商创建商品 THEN the SupplierProduct SHALL 关联供应商ID和分类ID
2. WHEN 供应商更新商品价格 THEN the SupplierProduct SHALL 更新价格字段
3. WHEN 供应商禁用商品 THEN the SupplierProduct SHALL 将状态设为禁用
4. WHEN 查询商品列表 THEN the SupplierProduct SHALL 支持按供应商、分类、关键词筛选
