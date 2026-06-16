-- 门店记账新增跑腿订单字段（幂等）
SET @db_name = DATABASE();

SET @sql_add_is_errand_order = (
  SELECT IF(
    EXISTS(
      SELECT 1
      FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA = @db_name
        AND TABLE_NAME = 'store_accounts'
        AND COLUMN_NAME = 'is_errand_order'
    ),
    'SELECT ''skip add is_errand_order''',
    'ALTER TABLE store_accounts ADD COLUMN is_errand_order TINYINT NOT NULL DEFAULT 0 COMMENT ''是否跑腿订单 1=是 0=否'' AFTER other_expense_amount'
  )
);
PREPARE stmt_add_is_errand_order FROM @sql_add_is_errand_order;
EXECUTE stmt_add_is_errand_order;
DEALLOCATE PREPARE stmt_add_is_errand_order;

SET @sql_add_errand_fee = (
  SELECT IF(
    EXISTS(
      SELECT 1
      FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA = @db_name
        AND TABLE_NAME = 'store_accounts'
        AND COLUMN_NAME = 'errand_fee'
    ),
    'SELECT ''skip add errand_fee''',
    'ALTER TABLE store_accounts ADD COLUMN errand_fee DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT ''跑腿费用'' AFTER is_errand_order'
  )
);
PREPARE stmt_add_errand_fee FROM @sql_add_errand_fee;
EXECUTE stmt_add_errand_fee;
DEALLOCATE PREPARE stmt_add_errand_fee;

SET @sql_add_idx_store_accounts_is_errand_order = (
  SELECT IF(
    EXISTS(
      SELECT 1
      FROM information_schema.STATISTICS
      WHERE TABLE_SCHEMA = @db_name
        AND TABLE_NAME = 'store_accounts'
        AND INDEX_NAME = 'idx_store_accounts_is_errand_order'
    ),
    'SELECT ''skip add idx_store_accounts_is_errand_order''',
    'CREATE INDEX idx_store_accounts_is_errand_order ON store_accounts(is_errand_order)'
  )
);
PREPARE stmt_add_idx_store_accounts_is_errand_order FROM @sql_add_idx_store_accounts_is_errand_order;
EXECUTE stmt_add_idx_store_accounts_is_errand_order;
DEALLOCATE PREPARE stmt_add_idx_store_accounts_is_errand_order;
