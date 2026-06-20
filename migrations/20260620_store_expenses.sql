-- 门店支出记录：用于记录已支付的外卖推广充值、平台费用、维修维护等支出。
CREATE TABLE IF NOT EXISTS `store_expenses` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `expense_no` VARCHAR(50) NOT NULL COMMENT '支出单号',
  `store_id` BIGINT UNSIGNED NOT NULL COMMENT '门店ID',
  `expense_date` DATE DEFAULT NULL COMMENT '支出日期',
  `category_code` VARCHAR(100) NOT NULL COMMENT '支出分类编码(字典:EXPENDITURECLASS)',
  `category_name` VARCHAR(100) NOT NULL COMMENT '支出分类名称',
  `amount` DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '支出金额',
  `remark` VARCHAR(500) DEFAULT NULL COMMENT '备注说明',
  `operator_id` BIGINT UNSIGNED NOT NULL COMMENT '操作人ID',
  `operator_name` VARCHAR(100) DEFAULT NULL COMMENT '操作人名称',
  `created_at` DATETIME(3) DEFAULT NULL,
  `updated_at` DATETIME(3) DEFAULT NULL,
  `deleted_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_store_expenses_expense_no` (`expense_no`),
  KEY `idx_store_expenses_store_id` (`store_id`),
  KEY `idx_store_expenses_expense_date` (`expense_date`),
  KEY `idx_store_expenses_category_code` (`category_code`),
  KEY `idx_store_expenses_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='门店支出记录';

INSERT INTO `dict_types` (`code`, `name`, `remark`, `status`, `created_at`, `updated_at`)
SELECT 'STOREEXPENSES', '门店支出', '门店支出模块字典', 1, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `dict_types` WHERE `code` = 'STOREEXPENSES');

INSERT INTO `dict_types` (`code`, `name`, `remark`, `status`, `created_at`, `updated_at`)
SELECT 'EXPENDITURECLASS', '门店支出分类', '门店支出分类字典', 1, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM `dict_types` WHERE `code` = 'EXPENDITURECLASS');

SET @expense_type_id = (SELECT id FROM dict_types WHERE code = 'EXPENDITURECLASS' LIMIT 1);

INSERT INTO dict_data (`type_id`, `type_code`, `label`, `value`, `sort`, `css_class`, `list_class`, `is_default`, `remark`, `status`, `created_at`, `updated_at`)
SELECT @expense_type_id, 'EXPENDITURECLASS', '外卖平台推广充值', 'takeout_promotion', 1, '', 'warning', true, '用于外卖推广费用和 ROI 统计', 1, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM dict_data WHERE type_code = 'EXPENDITURECLASS' AND value = 'takeout_promotion');

INSERT INTO dict_data (`type_id`, `type_code`, `label`, `value`, `sort`, `css_class`, `list_class`, `is_default`, `remark`, `status`, `created_at`, `updated_at`)
SELECT @expense_type_id, 'EXPENDITURECLASS', '平台服务费', 'platform_service_fee', 2, '', 'info', false, '', 1, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM dict_data WHERE type_code = 'EXPENDITURECLASS' AND value = 'platform_service_fee');

INSERT INTO dict_data (`type_id`, `type_code`, `label`, `value`, `sort`, `css_class`, `list_class`, `is_default`, `remark`, `status`, `created_at`, `updated_at`)
SELECT @expense_type_id, 'EXPENDITURECLASS', '物料采购', 'materials', 3, '', 'success', false, '', 1, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM dict_data WHERE type_code = 'EXPENDITURECLASS' AND value = 'materials');

INSERT INTO dict_data (`type_id`, `type_code`, `label`, `value`, `sort`, `css_class`, `list_class`, `is_default`, `remark`, `status`, `created_at`, `updated_at`)
SELECT @expense_type_id, 'EXPENDITURECLASS', '维修维护', 'maintenance', 4, '', 'warning', false, '', 1, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM dict_data WHERE type_code = 'EXPENDITURECLASS' AND value = 'maintenance');

INSERT INTO dict_data (`type_id`, `type_code`, `label`, `value`, `sort`, `css_class`, `list_class`, `is_default`, `remark`, `status`, `created_at`, `updated_at`)
SELECT @expense_type_id, 'EXPENDITURECLASS', '其他', 'other', 99, '', 'info', false, '', 1, NOW(3), NOW(3)
WHERE NOT EXISTS (SELECT 1 FROM dict_data WHERE type_code = 'EXPENDITURECLASS' AND value = 'other');

SET @store_id = (SELECT id FROM menus WHERE name = 'store' AND type = 1 ORDER BY id LIMIT 1);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_id, 'store-expenses', '门店支出', 'WalletCards', '/store/expenses', 'store/expenses/index', 2, 9, 'store:expenses:list', 1, 1, NOW(), NOW()
WHERE @store_id IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM menus WHERE parent_id = @store_id AND name = 'store-expenses' AND type = 2);

SET @store_expenses_id = (SELECT id FROM menus WHERE parent_id = @store_id AND name = 'store-expenses' AND type = 2 ORDER BY id LIMIT 1);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_expenses_id, 'store-expenses-add', '新增支出', '', '', '', 3, 1, 'store:expenses:add', 1, 1, NOW(), NOW()
WHERE @store_expenses_id IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM menus WHERE parent_id = @store_expenses_id AND name = 'store-expenses-add' AND type = 3);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_expenses_id, 'store-expenses-edit', '编辑支出', '', '', '', 3, 2, 'store:expenses:edit', 1, 1, NOW(), NOW()
WHERE @store_expenses_id IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM menus WHERE parent_id = @store_expenses_id AND name = 'store-expenses-edit' AND type = 3);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_expenses_id, 'store-expenses-delete', '删除支出', '', '', '', 3, 3, 'store:expenses:delete', 1, 1, NOW(), NOW()
WHERE @store_expenses_id IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM menus WHERE parent_id = @store_expenses_id AND name = 'store-expenses-delete' AND type = 3);

INSERT INTO role_menus (role_id, menu_id, permissions)
SELECT 1, id, 15 FROM menus WHERE name LIKE 'store-expenses%'
ON DUPLICATE KEY UPDATE permissions = 15;

INSERT INTO role_menus (role_id, menu_id, permissions)
SELECT 999, id, 15 FROM menus WHERE name LIKE 'store-expenses%'
ON DUPLICATE KEY UPDATE permissions = 15;

INSERT INTO role_menus (role_id, menu_id, permissions)
SELECT 2, id, 15 FROM menus WHERE name LIKE 'store-expenses%'
ON DUPLICATE KEY UPDATE permissions = 15;

INSERT INTO role_menus (role_id, menu_id, permissions)
SELECT 3, id, 15 FROM menus WHERE name LIKE 'store-expenses%'
ON DUPLICATE KEY UPDATE permissions = 15;
