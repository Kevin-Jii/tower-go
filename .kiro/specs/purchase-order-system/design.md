# Design Document

## Overview

本设计文档描述门店采购报菜单系统的技术架构和实现方案。系统采用Go语言开发，基于Gin框架构建RESTful API，使用GORM作为ORM框架，MySQL作为数据存储。

系统核心功能包括：
- 供应商商品管理（分类、商品CRUD）
- 门店供应商商品绑定（多供应商支持、默认供应商设置）
- 采购单管理（创建、状态流转、按供应商分组查看）

## Architecture

系统采用分层架构设计：

```
┌─────────────────────────────────────────────────────────┐
│                    Controller Layer                      │
│  (HTTP请求处理、参数验证、响应封装)                        │
├─────────────────────────────────────────────────────────┤
│                     Service Layer                        │
│  (业务逻辑处理、事务管理、跨模块协调)                      │
├─────────────────────────────────────────────────────────┤
│                     Module Layer                         │
│  (数据访问、CRUD操作、查询构建)                           │
├─────────────────────────────────────────────────────────┤
│                     Model Layer                          │
│  (数据模型定义、请求/响应结构体)                          │
├─────────────────────────────────────────────────────────┤
│                    Database (MySQL)                      │
└─────────────────────────────────────────────────────────┘
```

### Request Flow

```
HTTP Request → Controller → Service → Module → Database
                  ↓            ↓         ↓
              Validate    Business   SQL Query
              Params      Logic      Execute
```

## Components and Interfaces

### 1. Controller Layer

#### StoreSupplierController
```go
// 门店供应商商品管理
POST   /api/store-suppliers/bind      // 绑定商品
DELETE /api/store-suppliers/unbind    // 解绑商品
PUT    /api/store-suppliers/default   // 设置默认供应商
GET    /api/store-suppliers/:store_id // 获取门店绑定的商品列表
```

#### PurchaseOrderController
```go
// 采购单管理
POST   /api/purchase-orders           // 创建采购单
GET    /api/purchase-orders           // 采购单列表
GET    /api/purchase-orders/:id       // 采购单详情
PUT    /api/purchase-orders/:id       // 更新采购单
DELETE /api/purchase-orders/:id       // 删除采购单
GET    /api/purchase-orders/:id/by-supplier // 按供应商分组查看
```

#### SupplierProductController
```go
// 供应商商品管理
POST   /api/supplier-products         // 创建商品
GET    /api/supplier-products         // 商品列表
GET    /api/supplier-products/:id     // 商品详情
PUT    /api/supplier-products/:id     // 更新商品
DELETE /api/supplier-products/:id     // 删除商品
```

### 2. Service Layer

#### StoreSupplierService
```go
type StoreSupplierService interface {
    BindProducts(storeID uint, productIDs []uint) error
    UnbindProducts(storeID uint, productIDs []uint) error
    SetDefault(storeID, productID uint) error
    ListByStoreID(storeID uint) ([]*model.StoreSupplierProduct, error)
    ValidateStoreProducts(storeID uint, productIDs []uint) error
}
```

#### PurchaseOrderService
```go
type PurchaseOrderService interface {
    CreateOrder(storeID, userID uint, req *model.CreatePurchaseOrderReq) (*model.PurchaseOrder, error)
    GetOrder(id uint) (*model.PurchaseOrder, error)
    ListOrders(req *model.ListPurchaseOrderReq) ([]*model.PurchaseOrder, int64, error)
    UpdateOrder(id uint, req *model.UpdatePurchaseOrderReq) error
    DeleteOrder(id uint) error
    GetOrdersBySupplier(orderID uint) (map[uint][]model.PurchaseOrderItem, error)
}
```

### 3. Module Layer

#### StoreSupplierModule
- `BindProducts`: 批量绑定商品，幂等操作
- `UnbindProducts`: 批量解绑商品
- `SetDefault`: 设置默认供应商，自动取消同名商品其他供应商的默认标记
- `ListByStoreID`: 获取门店所有绑定商品
- `GetStoreProductByName`: 根据商品名获取默认供应商商品

#### PurchaseOrderModule
- `Create`: 创建采购单
- `CreateItems`: 批量创建采购单明细
- `GetByID`: 获取采购单详情（含关联数据）
- `List`: 分页查询采购单列表
- `GetOrdersBySupplier`: 按供应商分组获取明细

## Data Models

### Core Tables

```sql
-- 供应商分类
CREATE TABLE supplier_categories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    supplier_id BIGINT UNSIGNED NOT NULL COMMENT '供应商ID',
    name VARCHAR(100) NOT NULL COMMENT '分类名称',
    sort INT DEFAULT 0 COMMENT '排序',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态 1=启用 0=禁用',
    created_at DATETIME,
    updated_at DATETIME,
    INDEX idx_supplier_id (supplier_id)
);

-- 供应商商品
CREATE TABLE supplier_products (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    supplier_id BIGINT UNSIGNED NOT NULL COMMENT '供应商ID',
    category_id BIGINT UNSIGNED NOT NULL COMMENT '分类ID',
    name VARCHAR(200) NOT NULL COMMENT '商品名称',
    unit VARCHAR(20) NOT NULL DEFAULT '斤' COMMENT '单位',
    price DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '单价',
    spec VARCHAR(100) COMMENT '规格',
    remark VARCHAR(500) COMMENT '备注',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态 1=启用 0=禁用',
    created_at DATETIME,
    updated_at DATETIME,
    INDEX idx_supplier_id (supplier_id),
    INDEX idx_category_id (category_id)
);

-- 门店供应商商品关联
CREATE TABLE store_supplier_products (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    store_id BIGINT UNSIGNED NOT NULL COMMENT '门店ID',
    product_id BIGINT UNSIGNED NOT NULL COMMENT '供应商商品ID',
    is_default BOOLEAN DEFAULT FALSE COMMENT '是否默认供应商',
    created_at DATETIME,
    updated_at DATETIME,
    UNIQUE INDEX idx_store_product (store_id, product_id)
);

-- 采购单
CREATE TABLE purchase_orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_no VARCHAR(32) NOT NULL COMMENT '订单编号',
    store_id BIGINT UNSIGNED NOT NULL COMMENT '门店ID',
    total_amount DECIMAL(12,2) NOT NULL DEFAULT 0 COMMENT '总金额',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态 1=待确认 2=已确认 3=已完成 4=已取消',
    remark VARCHAR(500) COMMENT '备注',
    order_date DATE NOT NULL COMMENT '报菜日期',
    created_by BIGINT UNSIGNED NOT NULL COMMENT '创建人ID',
    created_at DATETIME,
    updated_at DATETIME,
    UNIQUE INDEX idx_order_no (order_no),
    INDEX idx_store_id (store_id)
);

-- 采购单明细
CREATE TABLE purchase_order_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL COMMENT '采购单ID',
    supplier_id BIGINT UNSIGNED NOT NULL COMMENT '供应商ID',
    product_id BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
    quantity DECIMAL(10,2) NOT NULL COMMENT '数量',
    unit_price DECIMAL(10,2) NOT NULL COMMENT '单价',
    amount DECIMAL(12,2) NOT NULL COMMENT '金额',
    remark VARCHAR(200) COMMENT '备注',
    created_at DATETIME,
    updated_at DATETIME,
    INDEX idx_order_id (order_id),
    INDEX idx_supplier_id (supplier_id)
);
```

### Entity Relationships

```
Supplier (1) ──────< SupplierCategory (N)
    │
    └──────────────< SupplierProduct (N) >────── SupplierCategory (1)
                            │
                            │
Store (1) ──────< StoreSupplierProduct (N) >────── SupplierProduct (1)
    │
    └──────────────< PurchaseOrder (N)
                            │
                            └──────< PurchaseOrderItem (N) >─┬── SupplierProduct (1)
                                                             └── Supplier (1)
```



## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

Based on the acceptance criteria analysis, the following correctness properties have been identified:

### Property 1: Binding Idempotency
*For any* store and product combination, binding the same product multiple times SHALL result in exactly one binding record.

**Validates: Requirements 1.2**

### Property 2: Unbind Removes Record
*For any* existing store-product binding, unbinding SHALL result in the binding record being removed and no longer queryable.

**Validates: Requirements 1.3**

### Property 3: Default Supplier Exclusivity
*For any* store with multiple suppliers providing the same-named product, at most one supplier SHALL have `is_default=true` at any time.

**Validates: Requirements 2.2**

### Property 4: Order Number Uniqueness
*For any* set of purchase orders created, all order numbers SHALL be unique.

**Validates: Requirements 3.1**

### Property 5: Supplier Auto-Matching
*For any* purchase order item, the `supplier_id` in the item record SHALL match the `supplier_id` of the referenced product.

**Validates: Requirements 3.2**

### Property 6: Total Amount Calculation Invariant
*For any* purchase order, the `total_amount` SHALL equal the sum of all item amounts, where each item amount equals `quantity * unit_price`.

**Validates: Requirements 3.3**

### Property 7: Unbound Product Rejection
*For any* product ID not in the store's binding list, creating a purchase order with that product SHALL fail with an error.

**Validates: Requirements 3.4**

### Property 8: Supplier Grouping Correctness
*For any* purchase order grouped by supplier, all items within each group SHALL have the same `supplier_id`, and the number of groups SHALL equal the number of distinct suppliers in the order.

**Validates: Requirements 4.1, 4.2**

### Property 9: Initial Status Invariant
*For any* newly created purchase order, the status SHALL be `PurchaseStatusPending` (1).

**Validates: Requirements 5.1**

### Property 10: Delete Constraint
*For any* purchase order with status other than `PurchaseStatusPending` (1) or `PurchaseStatusCancelled` (4), delete operation SHALL fail.

**Validates: Requirements 5.5**

### Property 11: Filter Correctness
*For any* list query with filters (store_id, supplier_id, status, date range), all returned orders SHALL satisfy all specified filter conditions.

**Validates: Requirements 6.1**

### Property 12: Pagination Correctness
*For any* list query with pagination parameters, the number of returned records SHALL not exceed the specified `page_size`.

**Validates: Requirements 6.2**

## Error Handling

### Error Categories

| Error Type | HTTP Status | Description |
|------------|-------------|-------------|
| ValidationError | 400 | 请求参数验证失败 |
| NotFoundError | 404 | 资源不存在 |
| BusinessError | 422 | 业务规则校验失败 |
| InternalError | 500 | 服务器内部错误 |

### Business Error Codes

```go
const (
    ErrProductNotBound     = "PRODUCT_NOT_BOUND"      // 商品未绑定到门店
    ErrInvalidOrderStatus  = "INVALID_ORDER_STATUS"   // 无效的订单状态
    ErrOrderCannotDelete   = "ORDER_CANNOT_DELETE"    // 订单不可删除
    ErrProductNotFound     = "PRODUCT_NOT_FOUND"      // 商品不存在
    ErrSupplierNotFound    = "SUPPLIER_NOT_FOUND"     // 供应商不存在
    ErrInvalidDateFormat   = "INVALID_DATE_FORMAT"    // 日期格式错误
)
```

### Error Response Format

```json
{
    "code": 422,
    "message": "商品未绑定到门店",
    "error_code": "PRODUCT_NOT_BOUND",
    "details": {
        "product_ids": [1, 2, 3]
    }
}
```

## Testing Strategy

### Dual Testing Approach

本系统采用单元测试和属性测试相结合的测试策略：

1. **单元测试**: 验证具体示例、边界条件和错误处理
2. **属性测试**: 验证在所有有效输入下都应成立的通用属性

### Property-Based Testing Framework

使用 [gopter](https://github.com/leanovate/gopter) 作为Go语言的属性测试库。

配置要求：
- 每个属性测试运行至少100次迭代
- 使用注释标记属性测试与设计文档的对应关系

### Test Annotation Format

每个属性测试必须使用以下格式的注释：

```go
// **Feature: purchase-order-system, Property 6: Total Amount Calculation Invariant**
func TestTotalAmountCalculation(t *testing.T) {
    // Property test implementation
}
```

### Unit Test Coverage

单元测试覆盖以下场景：

1. **Controller Layer**
   - 请求参数验证
   - 响应格式正确性
   - 错误处理

2. **Service Layer**
   - 业务逻辑正确性
   - 跨模块协调
   - 事务处理

3. **Module Layer**
   - CRUD操作
   - 查询构建
   - 数据完整性

### Property Test Coverage

属性测试覆盖以下核心属性：

| Property | Test File | Description |
|----------|-----------|-------------|
| Property 1 | store_supplier_test.go | 绑定幂等性 |
| Property 2 | store_supplier_test.go | 解绑删除记录 |
| Property 3 | store_supplier_test.go | 默认供应商互斥性 |
| Property 4 | purchase_order_test.go | 订单编号唯一性 |
| Property 5 | purchase_order_test.go | 供应商自动匹配 |
| Property 6 | purchase_order_test.go | 总金额计算不变量 |
| Property 7 | purchase_order_test.go | 未绑定商品拒绝 |
| Property 8 | purchase_order_test.go | 供应商分组正确性 |
| Property 9 | purchase_order_test.go | 初始状态不变量 |
| Property 10 | purchase_order_test.go | 删除约束 |
| Property 11 | purchase_order_test.go | 筛选正确性 |
| Property 12 | purchase_order_test.go | 分页正确性 |

### Test Data Generation

属性测试需要生成以下类型的测试数据：

```go
// 生成随机门店ID
func GenStoreID() gopter.Gen {
    return gen.UIntRange(1, 1000)
}

// 生成随机商品ID列表
func GenProductIDs() gopter.Gen {
    return gen.SliceOfN(10, gen.UIntRange(1, 100))
}

// 生成随机采购单请求
func GenCreatePurchaseOrderReq() gopter.Gen {
    return gen.Struct(reflect.TypeOf(model.CreatePurchaseOrderReq{}), map[string]gopter.Gen{
        "OrderDate": gen.AnyString().Map(func(s string) string {
            return time.Now().Format("2006-01-02")
        }),
        "Items": gen.SliceOfN(5, GenCreatePurchaseOrderItemReq()),
    })
}
```
