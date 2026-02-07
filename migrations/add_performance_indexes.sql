-- =====================================================
-- 数据库索引优化脚本
-- 生成时间: 2026-02-06
-- 目的: 提升查询性能
-- =====================================================

-- 注意: 执行前请备份数据库！
-- 建议在低峰期执行，避免影响业务

-- =====================================================
-- 1. store_accounts 表索引优化
-- =====================================================

-- 现有索引:
-- - id (主键)
-- - store_id (单列索引)
-- - channel (单列索引)
-- - order_no (单列索引)
-- - tag_code (单列索引)
-- - account_date (单列索引)
-- - deleted_at (软删除索引)

-- 缺失索引 (高优先级):

-- 复合索引: 门店+日期范围查询
-- 场景: 按门店和日期范围统计销售额
CREATE INDEX idx_store_account_date_range ON store_accounts(store_id, account_date);

-- 复合索引: 门店+渠道查询
-- 场景: 按门店和渠道筛选记账
CREATE INDEX idx_store_account_store_channel ON store_accounts(store_id, channel);

-- 复合索引: 日期+渠道查询
-- 场景: 按日期和渠道统计
CREATE INDEX idx_store_account_date_channel ON store_accounts(account_date, channel);

-- 复合索引: 门店+渠道+日期 (最常用查询组合)
-- 场景: 仪表板统计、报表查询
CREATE INDEX idx_store_account_all ON store_accounts(store_id, channel, account_date);

-- =====================================================
-- 2. store_account_items 表索引优化
-- =====================================================

-- 现有索引:
-- - id (主键)
-- - account_id (单列索引)
-- - product_id (单列索引)
-- - deleted_at (软删除索引)

-- 缺失索引 (中优先级):

-- 复合索引: 商品ID+创建时间 (用于商品销售排行)
CREATE INDEX idx_account_items_product_time ON store_account_items(product_id, created_at);

-- =====================================================
-- 3. inventories 表索引优化
-- =====================================================

-- 现有索引:
-- - id (主键)
-- - store_id (单列索引)
-- - product_id (单列索引)
-- - deleted_at (软删除索引)

-- 缺失索引 (高优先级):

-- 复合索引: 门店+商品 (唯一性约束)
-- 场景: 库存查询和更新
CREATE UNIQUE INDEX idx_inventory_unique ON inventories(store_id, product_id);

-- 复合索引: 门店+库存数量 (用于库存预警)
-- 场景: 低库存查询
CREATE INDEX idx_inventory_store_qty ON inventories(store_id, quantity);

-- =====================================================
-- 4. inventory_orders 表索引优化
-- =====================================================

-- 现有索引:
-- - id (主键)
-- - order_no (唯一索引)
-- - store_id (单列索引)
-- - deleted_at (软删除索引)

-- 缺失索引 (高优先级):

-- 复合索引: 门店+类型+日期
-- 场景: 出入库统计
CREATE INDEX idx_inv_order_store_type_date ON inventory_orders(store_id, type, created_at);

-- 复合索引: 类型+日期
-- 场景: 今日入库/出库统计
CREATE INDEX idx_inv_order_type_date ON inventory_orders(type, created_at);

-- =====================================================
-- 5. inventory_order_items 表索引优化
-- =====================================================

-- 现有索引:
-- - id (主键)
-- - order_id (单列索引)
-- - product_id (单列索引)
-- - deleted_at (软删除索引)

-- 缺失索引 (低优先级):

-- 复合索引: 商品ID+数量 (用于库存消耗分析)
CREATE INDEX idx_order_items_product_qty ON inventory_order_items(product_id, quantity);

-- =====================================================
-- 6. supplier_products 表索引优化
-- =====================================================

-- 场景: 商品搜索和筛选
CREATE INDEX idx_supplier_prod_name ON supplier_products(name);
CREATE INDEX idx_supplier_prod_category ON supplier_products(category_id);

-- =====================================================
-- 7. users 表索引优化
-- =====================================================

-- 场景: 用户搜索
CREATE INDEX idx_users_store_name ON users(store_id, username);

-- =====================================================
-- 8. stores 表索引优化
-- =====================================================

-- 场景: 门店搜索
CREATE INDEX idx_stores_name ON stores(name);

-- =====================================================
-- 执行脚本
-- =====================================================

-- 将上述 CREATE INDEX 语句复制到 MySQL 客户端执行
-- 建议分批执行，每批间隔几秒钟

-- 验证索引是否创建成功:
-- SHOW INDEX FROM table_name;
