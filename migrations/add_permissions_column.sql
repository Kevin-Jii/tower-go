-- 为 role_menus 表添加 permissions 字段
ALTER TABLE `role_menus` ADD COLUMN `permissions` TINYINT UNSIGNED NOT NULL DEFAULT 15 COMMENT '权限位：bit0=查看,bit1=新增,bit2=修改,bit3=删除' AFTER `menu_id`;

-- 为 store_role_menus 表添加 permissions 字段
ALTER TABLE `store_role_menus` ADD COLUMN `permissions` TINYINT UNSIGNED NOT NULL DEFAULT 15 COMMENT '权限位：bit0=查看,bit1=新增,bit2=修改,bit3=删除' AFTER `menu_id`;

-- 更新现有数据，默认给予所有权限（1111 = 15）
UPDATE `role_menus` SET `permissions` = 15 WHERE `permissions` = 0;
UPDATE `store_role_menus` SET `permissions` = 15 WHERE `permissions` = 0;
