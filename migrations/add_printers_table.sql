-- ============================================
-- 打印机表迁移
-- 执行方式: mysql -u用户名 -p密码 数据库名 < migrations/add_printers_table.sql
-- ============================================

-- 创建 printers 表
CREATE TABLE IF NOT EXISTS `printers` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `store_id` INT UNSIGNED NOT NULL COMMENT '关联门店ID',
    `sn` VARCHAR(32) NOT NULL COMMENT '打印机SN号' COLLATE 'utf8mb4_unicode_ci',
    `name` VARCHAR(100) NOT NULL COMMENT '打印机名称' COLLATE 'utf8mb4_unicode_ci',
    `type` TINYINT NOT NULL DEFAULT 1 COMMENT '打印机类型：1=小票，2=标签',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1=正常，2=停用',
    `is_default` TINYINT NOT NULL DEFAULT 0 COMMENT '是否为默认打印机：0=否，1=是',
    `remark` TEXT NOT NULL COMMENT '备注' COLLATE 'utf8mb4_unicode_ci',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `idx_sn` (`sn`),
    INDEX `idx_store_id` (`store_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE='utf8mb4_unicode_ci' COMMENT='打印机表';

-- 添加菜单数据（打印机管理）
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(0, 'printer', '打印机管理', 'printer', '', '', 1, 8, '', 1, 1, NOW(), NOW());
SET @printer_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@printer_id, 'printer-list', '打印机列表', '', '/printer/list', 'printer/list/index', 2, 1, 'printer:list', 1, 1, NOW(), NOW());

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@printer_id, 'printer-bind', '绑定打印机', '', '', '', 3, 1, 'printer:bind', 1, 1, NOW(), NOW()),
(@printer_id, 'printer-unbind', '解绑打印机', '', '', '', 3, 2, 'printer:unbind', 1, 1, NOW(), NOW()),
(@printer_id, 'printer-edit', '编辑打印机', '', '', '', 3, 3, 'printer:edit', 1, 1, NOW(), NOW()),
(@printer_id, 'printer-query', '查询状态', '', '', '', 3, 4, 'printer:query', 1, 1, NOW(), NOW());

-- 为总部管理员角色添加打印机管理权限
INSERT INTO role_menus (role_id, menu_id) VALUES
(1, @printer_id), (1, @printer_id + 1), (1, @printer_id + 2), (1, @printer_id + 3), (1, @printer_id + 4), (1, @printer_id + 5);