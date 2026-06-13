-- ============================================
-- Tower Go 数据库结构初始化（CREATE / ALTER / 索引，不含业务 INSERT）
-- 建议顺序：
--   1) mysql ... < migrations/init.sql（空库可直接执行；含核心表 CREATE IF NOT EXISTS）
--   2) mysql ... < migrations/init_seed_data.sql（仅种子数据）
-- 若已由应用 GORM AutoMigrate 建表，仍可执行本脚本做索引与列补丁。
-- ============================================

-- ============================================
-- 0. 核心业务表（CREATE TABLE IF NOT EXISTS，与 model 包字段对齐，供 §1 索引与 init_seed_data 依赖）
-- ============================================

CREATE TABLE IF NOT EXISTS `roles` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(50) NOT NULL,
  `code` VARCHAR(50) NOT NULL,
  `data_scope` TINYINT NOT NULL DEFAULT 3 COMMENT '数据范围 1=全部 2=租户 3=门店 4=仅本人',
  `status` TINYINT(1) NOT NULL DEFAULT 1,
  `description` VARCHAR(255) DEFAULT NULL,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_roles_name` (`name`),
  UNIQUE KEY `idx_roles_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

CREATE TABLE IF NOT EXISTS `stores` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `store_code` VARCHAR(6) DEFAULT NULL,
  `name` VARCHAR(100) NOT NULL,
  `address` VARCHAR(255) DEFAULT NULL,
  `administrative_unit` VARCHAR(100) DEFAULT NULL COMMENT '归属区',
  `phone` VARCHAR(20) DEFAULT NULL,
  `business_hours` VARCHAR(100) DEFAULT NULL,
  `status` INT NOT NULL DEFAULT 1,
  `contact_person` VARCHAR(50) DEFAULT NULL,
  `remark` TEXT,
  `third_party_account_id` BIGINT UNSIGNED DEFAULT NULL,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_stores_store_code` (`store_code`),
  KEY `idx_stores_third_party_account_id` (`third_party_account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='门店表';

CREATE TABLE IF NOT EXISTS `menus` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `name` VARCHAR(50) NOT NULL,
  `title` VARCHAR(50) NOT NULL,
  `icon` VARCHAR(100) DEFAULT NULL,
  `path` VARCHAR(200) DEFAULT NULL,
  `component` VARCHAR(200) DEFAULT NULL,
  `type` INT NOT NULL DEFAULT 1,
  `sort` INT NOT NULL DEFAULT 0,
  `permission` VARCHAR(200) DEFAULT NULL,
  `visible` INT NOT NULL DEFAULT 1,
  `status` INT NOT NULL DEFAULT 1,
  `remark` VARCHAR(500) DEFAULT NULL,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_menus_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='菜单表';

CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `employee_no` VARCHAR(6) NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `phone` VARCHAR(20) NOT NULL,
  `username` VARCHAR(191) NOT NULL,
  `nickname` VARCHAR(100) DEFAULT NULL,
  `email` VARCHAR(255) DEFAULT NULL,
  `store_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `role_id` BIGINT UNSIGNED NOT NULL DEFAULT 3,
  `status` INT NOT NULL DEFAULT 1,
  `gender` INT NOT NULL DEFAULT 1,
  `last_login_at` DATETIME(3) DEFAULT NULL,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_employee_no` (`employee_no`),
  UNIQUE KEY `idx_users_phone` (`phone`),
  UNIQUE KEY `idx_store_username` (`store_id`, `username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

CREATE TABLE IF NOT EXISTS `dict_types` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(100) NOT NULL,
  `name` VARCHAR(100) NOT NULL,
  `remark` VARCHAR(500) DEFAULT NULL,
  `status` TINYINT NOT NULL DEFAULT 1,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_dict_types_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='字典类型';

CREATE TABLE IF NOT EXISTS `dict_data` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `type_id` BIGINT UNSIGNED NOT NULL,
  `type_code` VARCHAR(100) NOT NULL,
  `label` VARCHAR(100) NOT NULL,
  `value` VARCHAR(100) NOT NULL,
  `sort` INT NOT NULL DEFAULT 0,
  `css_class` VARCHAR(100) DEFAULT NULL,
  `list_class` VARCHAR(100) DEFAULT NULL,
  `is_default` TINYINT(1) NOT NULL DEFAULT 0,
  `remark` VARCHAR(500) DEFAULT NULL,
  `status` TINYINT NOT NULL DEFAULT 1,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_dict_data_type_id` (`type_id`),
  KEY `idx_dict_data_type_code` (`type_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='字典数据';

CREATE TABLE IF NOT EXISTS `message_templates` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(50) NOT NULL,
  `name` VARCHAR(100) NOT NULL,
  `title` VARCHAR(200) DEFAULT NULL,
  `content` TEXT NOT NULL,
  `description` VARCHAR(500) DEFAULT NULL,
  `variables` TEXT,
  `is_enabled` TINYINT(1) NOT NULL DEFAULT 1,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_message_templates_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息模板';

CREATE TABLE IF NOT EXISTS `role_menus` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `role_id` BIGINT UNSIGNED NOT NULL,
  `menu_id` BIGINT UNSIGNED NOT NULL,
  `permissions` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '权限位',
  `created_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_menus_role_menu` (`role_id`, `menu_id`),
  KEY `idx_role_menus_role_id` (`role_id`),
  KEY `idx_role_menus_menu_id` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色菜单关联';

CREATE TABLE IF NOT EXISTS `suppliers` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `supplier_code` VARCHAR(50) NOT NULL,
  `supplier_name` VARCHAR(200) NOT NULL,
  `contact_person` VARCHAR(100) DEFAULT NULL,
  `contact_phone` VARCHAR(20) DEFAULT NULL,
  `contact_email` VARCHAR(100) DEFAULT NULL,
  `supplier_address` VARCHAR(500) DEFAULT NULL,
  `remark` TEXT,
  `status` TINYINT NOT NULL DEFAULT 1,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  `deleted_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_suppliers_supplier_code` (`supplier_code`),
  KEY `idx_suppliers_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='供应商';

CREATE TABLE IF NOT EXISTS `supplier_categories` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `supplier_id` BIGINT UNSIGNED NOT NULL,
  `name` VARCHAR(100) NOT NULL,
  `sort` INT NOT NULL DEFAULT 0,
  `status` TINYINT NOT NULL DEFAULT 1,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_supplier_categories_supplier_id` (`supplier_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='供应商分类';

CREATE TABLE IF NOT EXISTS `supplier_products` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `supplier_id` BIGINT UNSIGNED NOT NULL,
  `category_id` BIGINT UNSIGNED NOT NULL,
  `name` VARCHAR(200) NOT NULL,
  `unit` VARCHAR(20) NOT NULL DEFAULT '斤',
  `price` DECIMAL(10,2) NOT NULL DEFAULT 0,
  `bottle_price` DECIMAL(10,2) NOT NULL DEFAULT 0,
  `case_price` DECIMAL(10,2) NOT NULL DEFAULT 0,
  `bottles_per_case` INT NOT NULL DEFAULT 1,
  `spec` VARCHAR(100) DEFAULT NULL,
  `remark` VARCHAR(500) DEFAULT NULL,
  `status` TINYINT NOT NULL DEFAULT 1,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_supplier_products_supplier_id` (`supplier_id`),
  KEY `idx_supplier_products_category_id` (`category_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='供应商商品';

CREATE TABLE IF NOT EXISTS `store_suppliers` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `store_id` BIGINT UNSIGNED NOT NULL,
  `supplier_id` BIGINT UNSIGNED NOT NULL,
  `status` TINYINT NOT NULL DEFAULT 1,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_store_supplier` (`store_id`, `supplier_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='门店供应商关联';

CREATE TABLE IF NOT EXISTS `store_accounts` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `account_no` VARCHAR(50) NOT NULL,
  `store_id` BIGINT UNSIGNED NOT NULL,
  `member_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '关联会员ID',
  `payment_status` TINYINT NOT NULL DEFAULT 1 COMMENT '支付状态 1=已支付 2=未支付',
  `channel` VARCHAR(50) DEFAULT NULL,
  `order_no` VARCHAR(100) DEFAULT NULL,
  `total_amount` DECIMAL(10,2) DEFAULT NULL,
  `other_expense_amount` DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '其他支出金额',
  `net_income_amount` DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '净收入金额',
  `item_count` INT DEFAULT NULL,
  `tag_code` VARCHAR(50) DEFAULT NULL,
  `tag_name` VARCHAR(100) DEFAULT NULL,
  `remark` TEXT,
  `operator_id` BIGINT UNSIGNED NOT NULL,
  `account_date` DATE DEFAULT NULL,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  `deleted_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_store_accounts_account_no` (`account_no`),
  KEY `idx_store_accounts_store_id` (`store_id`),
  KEY `idx_store_accounts_member_id` (`member_id`),
  KEY `idx_store_accounts_payment_status` (`payment_status`),
  KEY `idx_store_accounts_channel` (`channel`),
  KEY `idx_store_accounts_order_no` (`order_no`),
  KEY `idx_store_accounts_tag_code` (`tag_code`),
  KEY `idx_store_accounts_account_date` (`account_date`),
  KEY `idx_store_accounts_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='门店记账';

CREATE TABLE IF NOT EXISTS `store_account_items` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `account_id` BIGINT UNSIGNED NOT NULL,
  `product_id` BIGINT UNSIGNED NOT NULL,
  `product_name` VARCHAR(200) DEFAULT NULL,
  `spec` VARCHAR(100) DEFAULT NULL,
  `quantity` DECIMAL(10,2) NOT NULL DEFAULT 1,
  `unit` VARCHAR(20) DEFAULT NULL,
  `price` DECIMAL(10,2) DEFAULT NULL,
  `amount` DECIMAL(10,2) DEFAULT NULL,
  `remark` VARCHAR(500) DEFAULT NULL,
  `created_at` DATETIME(3) DEFAULT NULL,
  `deleted_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_store_account_items_account_id` (`account_id`),
  KEY `idx_store_account_items_product_id` (`product_id`),
  KEY `idx_store_account_items_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='门店记账明细';

CREATE TABLE IF NOT EXISTS `store_account_consumables` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `account_id` BIGINT UNSIGNED NOT NULL,
  `product_id` BIGINT UNSIGNED NOT NULL,
  `product_name` VARCHAR(200) DEFAULT NULL,
  `quantity` DECIMAL(10,2) NOT NULL DEFAULT 1,
  `unit` VARCHAR(20) DEFAULT NULL,
  `price` DECIMAL(10,2) DEFAULT NULL,
  `amount` DECIMAL(10,2) DEFAULT NULL,
  `remark` VARCHAR(500) DEFAULT NULL,
  `created_at` DATETIME(3) DEFAULT NULL,
  `deleted_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_store_account_consumables_account_id` (`account_id`),
  KEY `idx_store_account_consumables_product_id` (`product_id`),
  KEY `idx_store_account_consumables_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='门店记账消耗品明细';

CREATE TABLE IF NOT EXISTS `third_party_accounts` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `platform_name` VARCHAR(50) NOT NULL DEFAULT 'tsbeer',
  `name` VARCHAR(100) NOT NULL,
  `login_name` VARCHAR(100) NOT NULL,
  `phone` VARCHAR(30) DEFAULT NULL,
  `password` VARCHAR(255) NOT NULL,
  `application_key` VARCHAR(128) NOT NULL,
  `login_type` VARCHAR(10) NOT NULL DEFAULT '2',
  `channel` VARCHAR(20) NOT NULL DEFAULT 'WEB',
  `shop_id` VARCHAR(64) DEFAULT NULL,
  `customer_id` VARCHAR(64) DEFAULT NULL,
  `is_enabled` TINYINT(1) NOT NULL DEFAULT 1,
  `last_test_ok` TINYINT(1) NOT NULL DEFAULT 0,
  `last_test_msg` VARCHAR(500) DEFAULT NULL,
  `last_token` TEXT,
  `token_valid_time` BIGINT DEFAULT NULL,
  `last_test_at` DATETIME(3) DEFAULT NULL,
  `last_sync_at` DATETIME(3) DEFAULT NULL,
  `last_sync_msg` VARCHAR(500) DEFAULT NULL,
  `last_sync_count` INT NOT NULL DEFAULT 0,
  `remark` VARCHAR(500) DEFAULT NULL,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  `deleted_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_third_party_accounts_login_name` (`login_name`),
  KEY `idx_third_party_accounts_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='第三方账号池';

CREATE TABLE IF NOT EXISTS `third_party_orders` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `account_id` BIGINT UNSIGNED NOT NULL,
  `platform_name` VARCHAR(50) NOT NULL,
  `order_no` VARCHAR(100) NOT NULL,
  `place_time` DATETIME(3) DEFAULT NULL,
  `place_date` VARCHAR(10) DEFAULT NULL,
  `order_trade_status` VARCHAR(64) DEFAULT NULL,
  `status_name` VARCHAR(100) DEFAULT NULL,
  `pay_amount` DECIMAL(12,2) DEFAULT NULL,
  `total_amount` DECIMAL(12,2) DEFAULT NULL,
  `total_item_num` DECIMAL(12,2) DEFAULT NULL,
  `raw_json` LONGTEXT,
  `synced_at` DATETIME(3) DEFAULT NULL,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  `deleted_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tp_order_no` (`order_no`),
  KEY `idx_third_party_orders_account_id` (`account_id`),
  KEY `idx_third_party_orders_place_date` (`place_date`),
  KEY `idx_third_party_orders_synced_at` (`synced_at`),
  KEY `idx_third_party_orders_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='第三方同步订单';

CREATE TABLE IF NOT EXISTS `meituan_ai_accounts` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `store_id` BIGINT UNSIGNED NOT NULL,
  `shop_name` VARCHAR(100) NOT NULL,
  `shop_id` VARCHAR(100) DEFAULT NULL,
  `login_name` VARCHAR(100) DEFAULT NULL,
  `auth_status` VARCHAR(20) NOT NULL DEFAULT 'manual',
  `is_enabled` TINYINT(1) NOT NULL DEFAULT 1,
  `last_imported_at` DATETIME(3) DEFAULT NULL,
  `remark` VARCHAR(500) DEFAULT NULL,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  `deleted_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_meituan_ai_accounts_store_id` (`store_id`),
  KEY `idx_meituan_ai_accounts_shop_id` (`shop_id`),
  KEY `idx_meituan_ai_accounts_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='美团AI运营账号';

CREATE TABLE IF NOT EXISTS `meituan_ai_orders` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `store_id` BIGINT UNSIGNED NOT NULL,
  `account_id` BIGINT UNSIGNED NOT NULL,
  `order_no` VARCHAR(100) NOT NULL,
  `order_time` DATETIME(3) NOT NULL,
  `customer_name` VARCHAR(100) DEFAULT NULL,
  `product_summary` VARCHAR(500) DEFAULT NULL,
  `original_amount` DECIMAL(10,2) NOT NULL DEFAULT 0,
  `actual_amount` DECIMAL(10,2) NOT NULL DEFAULT 0,
  `delivery_fee` DECIMAL(10,2) NOT NULL DEFAULT 0,
  `platform_fee` DECIMAL(10,2) NOT NULL DEFAULT 0,
  `refund_amount` DECIMAL(10,2) NOT NULL DEFAULT 0,
  `status` VARCHAR(50) DEFAULT NULL,
  `store_account_id` BIGINT UNSIGNED DEFAULT NULL,
  `imported_at` DATETIME(3) DEFAULT NULL,
  `raw_json` LONGTEXT,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  `deleted_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_meituan_ai_orders_account_order` (`account_id`, `order_no`),
  KEY `idx_meituan_ai_orders_store_id` (`store_id`),
  KEY `idx_meituan_ai_orders_account_id` (`account_id`),
  KEY `idx_meituan_ai_orders_order_no` (`order_no`),
  KEY `idx_meituan_ai_orders_order_time` (`order_time`),
  KEY `idx_meituan_ai_orders_store_account_id` (`store_account_id`),
  KEY `idx_meituan_ai_orders_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='美团AI运营订单';

CREATE TABLE IF NOT EXISTS `meituan_ai_reviews` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `store_id` BIGINT UNSIGNED NOT NULL,
  `account_id` BIGINT UNSIGNED NOT NULL,
  `review_id` VARCHAR(100) DEFAULT NULL,
  `order_no` VARCHAR(100) DEFAULT NULL,
  `rating` INT NOT NULL DEFAULT 0,
  `content` TEXT,
  `sentiment` VARCHAR(20) DEFAULT NULL,
  `tags` VARCHAR(255) DEFAULT NULL,
  `suggested_reply` TEXT,
  `review_time` DATETIME(3) NOT NULL,
  `reply_status` VARCHAR(20) DEFAULT 'pending',
  `imported_at` DATETIME(3) DEFAULT NULL,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  `deleted_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_meituan_ai_reviews_account_review` (`account_id`, `review_id`),
  KEY `idx_meituan_ai_reviews_store_id` (`store_id`),
  KEY `idx_meituan_ai_reviews_account_id` (`account_id`),
  KEY `idx_meituan_ai_reviews_review_id` (`review_id`),
  KEY `idx_meituan_ai_reviews_order_no` (`order_no`),
  KEY `idx_meituan_ai_reviews_review_time` (`review_time`),
  KEY `idx_meituan_ai_reviews_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='美团AI运营评价';

CREATE TABLE IF NOT EXISTS `meituan_ai_suggestions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `store_id` BIGINT UNSIGNED NOT NULL,
  `account_id` BIGINT UNSIGNED NOT NULL,
  `type` VARCHAR(30) NOT NULL,
  `title` VARCHAR(120) NOT NULL,
  `content` TEXT,
  `reason` TEXT,
  `impact_score` INT NOT NULL DEFAULT 0,
  `status` VARCHAR(20) NOT NULL DEFAULT 'pending',
  `action_payload` TEXT,
  `generated_at` DATETIME(3) DEFAULT NULL,
  `approved_at` DATETIME(3) DEFAULT NULL,
  `done_at` DATETIME(3) DEFAULT NULL,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  `deleted_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_meituan_ai_suggestions_store_id` (`store_id`),
  KEY `idx_meituan_ai_suggestions_account_id` (`account_id`),
  KEY `idx_meituan_ai_suggestions_type` (`type`),
  KEY `idx_meituan_ai_suggestions_status` (`status`),
  KEY `idx_meituan_ai_suggestions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='美团AI运营建议';

CREATE TABLE IF NOT EXISTS `third_party_routes` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  `remark` VARCHAR(500) DEFAULT NULL,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='第三方物流路线';

CREATE TABLE IF NOT EXISTS `third_party_route_stores` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `route_id` BIGINT UNSIGNED NOT NULL,
  `store_id` BIGINT UNSIGNED NOT NULL,
  `sort` INT NOT NULL DEFAULT 0,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_route_store` (`route_id`, `store_id`),
  KEY `idx_route_stores_route_id` (`route_id`),
  KEY `idx_route_stores_store_id` (`store_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='第三方物流路线-门店';

CREATE TABLE IF NOT EXISTS `third_party_logistics_sheets` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `route_id` BIGINT UNSIGNED NOT NULL,
  `sheet_date` VARCHAR(10) NOT NULL,
  `start_date` VARCHAR(10) NOT NULL,
  `end_date` VARCHAR(10) NOT NULL,
  `headers_json` LONGTEXT,
  `rows_json` LONGTEXT,
  `products_json` LONGTEXT,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_route_logistics_period` (`route_id`, `start_date`, `end_date`),
  KEY `idx_logistics_sheets_route_id` (`route_id`),
  KEY `idx_logistics_sheets_sheet_date` (`sheet_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='第三方物流单';

CREATE TABLE IF NOT EXISTS `inventories` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `store_id` BIGINT UNSIGNED NOT NULL,
  `product_id` BIGINT UNSIGNED NOT NULL,
  `quantity` DECIMAL(10,2) NOT NULL DEFAULT 0,
  `unit` VARCHAR(20) DEFAULT NULL,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  `deleted_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_inventories_store_id` (`store_id`),
  KEY `idx_inventories_product_id` (`product_id`),
  KEY `idx_inventories_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='门店库存';

CREATE TABLE IF NOT EXISTS `inventory_orders` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `order_no` VARCHAR(50) NOT NULL,
  `type` TINYINT NOT NULL,
  `store_id` BIGINT UNSIGNED NOT NULL,
  `store_name` VARCHAR(100) DEFAULT NULL,
  `reason` VARCHAR(100) DEFAULT NULL,
  `remark` TEXT,
  `total_quantity` DECIMAL(10,2) DEFAULT NULL,
  `item_count` INT DEFAULT NULL,
  `operator_id` BIGINT UNSIGNED NOT NULL,
  `operator_name` VARCHAR(50) DEFAULT NULL,
  `operator_phone` VARCHAR(20) DEFAULT NULL,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  `deleted_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_inventory_orders_order_no` (`order_no`),
  KEY `idx_inventory_orders_store_id` (`store_id`),
  KEY `idx_inventory_orders_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='出入库单';

CREATE TABLE IF NOT EXISTS `inventory_order_items` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `order_id` BIGINT UNSIGNED NOT NULL,
  `product_id` BIGINT UNSIGNED NOT NULL,
  `product_name` VARCHAR(200) DEFAULT NULL,
  `quantity` DECIMAL(10,2) NOT NULL,
  `unit` VARCHAR(20) DEFAULT NULL,
  `production_date` DATE DEFAULT NULL,
  `expiry_date` DATE DEFAULT NULL,
  `remark` VARCHAR(500) DEFAULT NULL,
  `created_at` DATETIME(3) DEFAULT NULL,
  `deleted_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_inventory_order_items_order_id` (`order_id`),
  KEY `idx_inventory_order_items_product_id` (`product_id`),
  KEY `idx_inventory_order_items_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='出入库单明细';

CREATE TABLE IF NOT EXISTS `inventory_loss_orders` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `order_no` VARCHAR(50) NOT NULL,
  `store_id` BIGINT UNSIGNED NOT NULL,
  `type` VARCHAR(20) NOT NULL,
  `member_id` BIGINT UNSIGNED DEFAULT NULL,
  `reason` VARCHAR(200) DEFAULT NULL,
  `total_cost` DECIMAL(10,2) NOT NULL DEFAULT 0,
  `item_count` INT NOT NULL DEFAULT 0,
  `operator_id` BIGINT UNSIGNED NOT NULL,
  `operator_name` VARCHAR(50) DEFAULT NULL,
  `is_canceled` TINYINT(1) NOT NULL DEFAULT 0,
  `canceled_at` DATETIME(3) DEFAULT NULL,
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  `deleted_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_inventory_loss_orders_order_no` (`order_no`),
  KEY `idx_inventory_loss_orders_store_id` (`store_id`),
  KEY `idx_inventory_loss_orders_type` (`type`),
  KEY `idx_inventory_loss_orders_member_id` (`member_id`),
  KEY `idx_inventory_loss_orders_is_canceled` (`is_canceled`),
  KEY `idx_inventory_loss_orders_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存报损/自用/赠送单';

CREATE TABLE IF NOT EXISTS `inventory_loss_order_items` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `order_id` BIGINT UNSIGNED NOT NULL,
  `product_id` BIGINT UNSIGNED NOT NULL,
  `product_name` VARCHAR(200) DEFAULT NULL,
  `unit` VARCHAR(50) DEFAULT NULL,
  `quantity` DECIMAL(10,2) NOT NULL,
  `base_quantity` DECIMAL(10,2) NOT NULL,
  `base_unit` VARCHAR(20) DEFAULT NULL,
  `cost_price` DECIMAL(10,2) NOT NULL DEFAULT 0,
  `cost_amount` DECIMAL(10,2) NOT NULL DEFAULT 0,
  `remark` VARCHAR(500) DEFAULT NULL,
  `created_at` DATETIME(3) DEFAULT NULL,
  `deleted_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_inventory_loss_order_items_order_id` (`order_id`),
  KEY `idx_inventory_loss_order_items_product_id` (`product_id`),
  KEY `idx_inventory_loss_order_items_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存报损/自用/赠送明细';

-- ============================================
-- 1. 性能索引优化（依赖业务表已存在）
-- ============================================

-- store_accounts 复合索引
CREATE INDEX idx_store_account_date_range ON store_accounts(store_id, account_date);
CREATE INDEX idx_store_account_store_channel ON store_accounts(store_id, channel);
CREATE INDEX idx_store_account_date_channel ON store_accounts(account_date, channel);
CREATE INDEX idx_store_account_all ON store_accounts(store_id, channel, account_date);

-- store_account_items 索引
CREATE INDEX idx_account_items_product_time ON store_account_items(product_id, created_at);

-- inventories 索引
CREATE UNIQUE INDEX idx_inventory_unique ON inventories(store_id, product_id);
CREATE INDEX idx_inventory_store_qty ON inventories(store_id, quantity);

-- inventory_orders 索引
CREATE INDEX idx_inv_order_store_type_date ON inventory_orders(store_id, type, created_at);
CREATE INDEX idx_inv_order_type_date ON inventory_orders(type, created_at);

-- inventory_order_items 索引
CREATE INDEX idx_order_items_product_qty ON inventory_order_items(product_id, quantity);

-- supplier_products 索引
CREATE INDEX idx_supplier_prod_name ON supplier_products(name);
CREATE INDEX idx_supplier_prod_category ON supplier_products(category_id);

-- users 索引
CREATE INDEX idx_users_store_name ON users(store_id, username);

-- stores 索引
CREATE INDEX idx_stores_name ON stores(name);

-- ============================================
-- 2. 扩展表结构（CREATE TABLE IF NOT EXISTS）
-- ============================================

-- 打印机表
CREATE TABLE IF NOT EXISTS `printers` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `store_id` INT UNSIGNED NOT NULL COMMENT '关联门店ID',
    `sn` VARCHAR(32) NOT NULL COMMENT '打印机SN号',
    `name` VARCHAR(100) NOT NULL COMMENT '打印机名称',
    `type` TINYINT NOT NULL DEFAULT 1 COMMENT '打印机类型：1=小票，2=标签',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1=正常，2=停用',
    `is_default` TINYINT NOT NULL DEFAULT 0 COMMENT '是否为默认打印机：0=否，1=是',
    `online` TINYINT NOT NULL DEFAULT 0 COMMENT '在线状态：0=离线，1=在线，2=异常',
    `last_heartbeat` DATETIME DEFAULT NULL COMMENT '最后心跳时间',
    `remark` TEXT NOT NULL COMMENT '备注',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_sn` (`sn`),
    INDEX `idx_store_id` (`store_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='打印机表';

-- 会员表
CREATE TABLE IF NOT EXISTS `t_member` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `store_id` int unsigned NOT NULL DEFAULT '0' COMMENT '所属门店ID',
  `uid` varchar(32) NOT NULL DEFAULT '' COMMENT '用户唯一标识',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '会员姓名',
  `phone` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号',
  `balance` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '余额',
  `points` int NOT NULL DEFAULT '0' COMMENT '积分',
  `level` int NOT NULL DEFAULT '1' COMMENT '等级',
  `version` int NOT NULL DEFAULT '0' COMMENT '乐观锁版本号',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uid` (`uid`),
  UNIQUE KEY `idx_store_phone` (`store_id`, `phone`),
  KEY `idx_store_id` (`store_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员表';

-- 会员流水表
CREATE TABLE IF NOT EXISTS `t_member_wallet_log` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `member_id` int unsigned NOT NULL DEFAULT '0' COMMENT '会员ID',
  `change_type` tinyint NOT NULL DEFAULT '0' COMMENT '变动类型: 1=充值 2=消费 3=退款 4=调增 5=调减',
  `change_amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '变动金额',
  `balance_after` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '变动后余额',
  `related_order_no` varchar(64) NOT NULL DEFAULT '' COMMENT '关联单号',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_member_id` (`member_id`),
  KEY `idx_related_order_no` (`related_order_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员流水表';

-- 充值单表
CREATE TABLE IF NOT EXISTS `t_recharge_order` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `order_no` varchar(32) NOT NULL DEFAULT '' COMMENT '单号',
  `member_id` int unsigned NOT NULL DEFAULT '0' COMMENT '会员ID',
  `pay_amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '实付金额',
  `gift_amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '赠送金额',
  `total_amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '总金额',
  `pay_status` tinyint NOT NULL DEFAULT '0' COMMENT '支付状态: 0=待支付 1=已支付 2=已取消 3=已退款',
  `pay_time` datetime DEFAULT NULL COMMENT '支付时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_order_no` (`order_no`),
  KEY `idx_member_id` (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='充值单表';

-- ============================================
-- 2b. 列补丁与补充表（幂等，兼容历史库升级）
-- ============================================

SET @db_name = DATABASE();

SET @sql_add_other_expense_amount = (
  SELECT IF(
    EXISTS(
      SELECT 1
      FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA = @db_name
        AND TABLE_NAME = 'store_accounts'
        AND COLUMN_NAME = 'other_expense_amount'
    ),
    'SELECT ''skip add other_expense_amount''',
    'ALTER TABLE store_accounts ADD COLUMN other_expense_amount DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT ''其他支出金额'' AFTER total_amount'
  )
);
PREPARE stmt_add_other_expense_amount FROM @sql_add_other_expense_amount;
EXECUTE stmt_add_other_expense_amount;
DEALLOCATE PREPARE stmt_add_other_expense_amount;

SET @sql_add_net_income_amount = (
  SELECT IF(
    EXISTS(
      SELECT 1
      FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA = @db_name
        AND TABLE_NAME = 'store_accounts'
        AND COLUMN_NAME = 'net_income_amount'
    ),
    'SELECT ''skip add net_income_amount''',
    'ALTER TABLE store_accounts ADD COLUMN net_income_amount DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT ''净收入金额'' AFTER other_expense_amount'
  )
);
PREPARE stmt_add_net_income_amount FROM @sql_add_net_income_amount;
EXECUTE stmt_add_net_income_amount;
DEALLOCATE PREPARE stmt_add_net_income_amount;

SET @sql_add_bottle_price = (
  SELECT IF(
    EXISTS(
      SELECT 1
      FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA = @db_name
        AND TABLE_NAME = 'supplier_products'
        AND COLUMN_NAME = 'bottle_price'
    ),
    'SELECT ''skip add bottle_price''',
    'ALTER TABLE supplier_products ADD COLUMN bottle_price DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT ''单瓶价格'' AFTER price'
  )
);
PREPARE stmt_add_bottle_price FROM @sql_add_bottle_price;
EXECUTE stmt_add_bottle_price;
DEALLOCATE PREPARE stmt_add_bottle_price;

SET @sql_add_case_price = (
  SELECT IF(
    EXISTS(
      SELECT 1
      FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA = @db_name
        AND TABLE_NAME = 'supplier_products'
        AND COLUMN_NAME = 'case_price'
    ),
    'SELECT ''skip add case_price''',
    'ALTER TABLE supplier_products ADD COLUMN case_price DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT ''整箱价格'' AFTER bottle_price'
  )
);
PREPARE stmt_add_case_price FROM @sql_add_case_price;
EXECUTE stmt_add_case_price;
DEALLOCATE PREPARE stmt_add_case_price;

SET @sql_add_bottles_per_case = (
  SELECT IF(
    EXISTS(
      SELECT 1
      FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA = @db_name
        AND TABLE_NAME = 'supplier_products'
        AND COLUMN_NAME = 'bottles_per_case'
    ),
    'SELECT ''skip add bottles_per_case''',
    'ALTER TABLE supplier_products ADD COLUMN bottles_per_case INT NOT NULL DEFAULT 1 COMMENT ''每箱瓶数'' AFTER case_price'
  )
);
PREPARE stmt_add_bottles_per_case FROM @sql_add_bottles_per_case;
EXECUTE stmt_add_bottles_per_case;
DEALLOCATE PREPARE stmt_add_bottles_per_case;

-- 商品多单位换算与价格表
CREATE TABLE IF NOT EXISTS product_unit_specs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  product_id BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
  unit_code VARCHAR(50) NOT NULL COMMENT '单位编码，如 bottle/case/barrel/liter',
  unit_name VARCHAR(50) NOT NULL COMMENT '单位名称，如 瓶/箱/桶/L',
  factor_to_base DECIMAL(12,6) NOT NULL DEFAULT 1 COMMENT '换算到基础单位L的系数',
  `precision` INT NOT NULL DEFAULT 0 COMMENT '数量精度(小数位)',
  cost_price DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '单位成本价',
  sale_price DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '单位售价',
  is_enabled TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (id),
  KEY idx_product_unit_specs_product_id (product_id),
  UNIQUE KEY uk_product_unit_specs_product_unit_name (product_id, unit_code, unit_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品多单位换算与价格表';

SET @sql_add_product_unit_precision = (
  SELECT IF(
    EXISTS(
      SELECT 1
      FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA = DATABASE()
        AND TABLE_NAME = 'product_unit_specs'
        AND COLUMN_NAME = 'precision'
    ),
    'SELECT ''skip add product_unit_specs.precision''',
    'ALTER TABLE product_unit_specs ADD COLUMN `precision` INT NOT NULL DEFAULT 0 COMMENT ''数量精度(小数位)'' AFTER factor_to_base'
  )
);
PREPARE stmt_add_product_unit_precision FROM @sql_add_product_unit_precision;
EXECUTE stmt_add_product_unit_precision;
DEALLOCATE PREPARE stmt_add_product_unit_precision;

SET @sql_add_roles_data_scope = (
  SELECT IF(
    EXISTS(
      SELECT 1
      FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA = @db_name
        AND TABLE_NAME = 'roles'
        AND COLUMN_NAME = 'data_scope'
    ),
    'SELECT ''skip add roles.data_scope''',
    'ALTER TABLE roles ADD COLUMN data_scope TINYINT NOT NULL DEFAULT 3 COMMENT ''数据范围 1=全部 2=租户 3=门店 4=仅本人'' AFTER code'
  )
);
PREPARE stmt_add_roles_data_scope FROM @sql_add_roles_data_scope;
EXECUTE stmt_add_roles_data_scope;
DEALLOCATE PREPARE stmt_add_roles_data_scope;

SET @sql_add_stores_third_party_account_id = (
  SELECT IF(
    EXISTS(
      SELECT 1
      FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA = @db_name
        AND TABLE_NAME = 'stores'
        AND COLUMN_NAME = 'third_party_account_id'
    ),
    'SELECT ''skip add stores.third_party_account_id''',
    'ALTER TABLE stores ADD COLUMN third_party_account_id BIGINT UNSIGNED DEFAULT NULL COMMENT ''绑定第三方账号池ID'' AFTER remark'
  )
);
PREPARE stmt_add_stores_third_party_account_id FROM @sql_add_stores_third_party_account_id;
EXECUTE stmt_add_stores_third_party_account_id;
DEALLOCATE PREPARE stmt_add_stores_third_party_account_id;

SET @sql_add_stores_third_party_account_idx = (
  SELECT IF(
    EXISTS(
      SELECT 1
      FROM information_schema.STATISTICS
      WHERE TABLE_SCHEMA = @db_name
        AND TABLE_NAME = 'stores'
        AND INDEX_NAME = 'idx_stores_third_party_account_id'
    ),
    'SELECT ''skip add idx_stores_third_party_account_id''',
    'ALTER TABLE stores ADD INDEX idx_stores_third_party_account_id (third_party_account_id)'
  )
);
PREPARE stmt_add_stores_third_party_account_idx FROM @sql_add_stores_third_party_account_idx;
EXECUTE stmt_add_stores_third_party_account_idx;
DEALLOCATE PREPARE stmt_add_stores_third_party_account_idx;

-- 超级管理员 store_id=0 不引用 stores，移除 GORM 自动创建的外键（若存在）
SET @db_name = DATABASE();
SET @sql_drop_users_store_fk = (
  SELECT IF(
    EXISTS(
      SELECT 1
      FROM information_schema.TABLE_CONSTRAINTS
      WHERE TABLE_SCHEMA = @db_name
        AND TABLE_NAME = 'users'
        AND CONSTRAINT_TYPE = 'FOREIGN KEY'
        AND CONSTRAINT_NAME = 'fk_users_store'
    ),
    'ALTER TABLE `users` DROP FOREIGN KEY `fk_users_store`',
    'SELECT ''skip drop fk_users_store'''
  )
);
PREPARE stmt_drop_users_store_fk FROM @sql_drop_users_store_fk;
EXECUTE stmt_drop_users_store_fk;
DEALLOCATE PREPARE stmt_drop_users_store_fk;

-- 历史超级管理员账号归一化为未绑店
UPDATE users u
INNER JOIN roles r ON u.role_id = r.id
SET u.store_id = 0
WHERE r.code = 'super_admin' AND u.store_id <> 0;

-- 列补齐后的历史行回填（非 INSERT，随结构脚本执行一次即可）
UPDATE supplier_products
SET bottle_price = price
WHERE (bottle_price IS NULL OR bottle_price = 0) AND price > 0;

UPDATE supplier_products
SET bottles_per_case = 1
WHERE bottles_per_case IS NULL OR bottles_per_case <= 0;

-- ============================================
-- 3. 业务种子数据（INSERT）
-- ============================================
-- 已全部迁移至 migrations/init_seed_data.sql

-- ============================================
-- 4. 价目单模块表结构
-- ============================================

CREATE TABLE IF NOT EXISTS `price_lists` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `store_id` bigint unsigned NOT NULL COMMENT '门店ID',
  `name` varchar(200) NOT NULL COMMENT '价目单名称',
  `logo` varchar(500) DEFAULT NULL COMMENT '品牌LOGO URL',
  `description` varchar(500) DEFAULT NULL COMMENT '描述',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1=启用 0=禁用',
  `is_default` tinyint NOT NULL DEFAULT '0' COMMENT '是否默认 1=是 0=否',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_store_id` (`store_id`),
  KEY `idx_is_default` (`is_default`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='价目单表';

CREATE TABLE IF NOT EXISTS `price_list_categories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `price_list_id` bigint unsigned NOT NULL COMMENT '价目单ID',
  `main_title` varchar(100) NOT NULL COMMENT '主标题',
  `sub_title` varchar(100) DEFAULT NULL COMMENT '副标题',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_price_list_id` (`price_list_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='价目单分类表';

CREATE TABLE IF NOT EXISTS `price_list_items` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `category_id` bigint unsigned NOT NULL COMMENT '价目单分类ID',
  `product_id` bigint unsigned NOT NULL COMMENT '供应商商品ID',
  `display_name` varchar(200) DEFAULT NULL COMMENT '显示名称（可覆盖商品名称）',
  `price` decimal(10,2) NOT NULL COMMENT '价格',
  `unit` varchar(20) DEFAULT NULL COMMENT '单位',
  `spec` varchar(100) DEFAULT NULL COMMENT '规格说明',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1=显示 0=隐藏',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='价目单商品表';

INSERT INTO `dict_types` (`code`, `name`, `remark`, `status`, `created_at`, `updated_at`)
SELECT 'ADMINISTRATIVEUNIT', '归属区', '门店归属区字典', 1, NOW(3), NOW(3)
WHERE NOT EXISTS (
  SELECT 1 FROM `dict_types` WHERE `code` = 'ADMINISTRATIVEUNIT'
);
