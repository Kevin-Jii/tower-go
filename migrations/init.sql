-- ============================================
-- Tower Go 初始化脚本
-- 执行方式: mysql -u用户名 -p密码 数据库名 < migrations/init.sql
-- ============================================

-- ============================================
-- 1. 性能索引优化
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
-- 2. 扩展表结构
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
  UNIQUE KEY `idx_phone` (`phone`)
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
-- 3. 种子数据
-- ============================================

-- 角色数据
INSERT INTO roles (id, name, code, description, created_at, updated_at) VALUES
(1, '总部管理员', 'admin', '系统最高权限角色', NOW(), NOW()),
(2, '门店管理员', 'store_admin', '门店维度管理权限角色', NOW(), NOW()),
(3, '普通员工', 'staff', '基础操作权限角色', NOW(), NOW()),
(999, '超级管理员', 'super_admin', '系统最高权限，不可删除', NOW(), NOW())
ON DUPLICATE KEY UPDATE name=VALUES(name), description=VALUES(description), updated_at=NOW();

-- 菜单数据
-- 系统管理
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(0, 'system', '系统管理', 'setting', '', '', 1, 1, '', 1, 1, NOW(), NOW());
SET @system_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@system_id, 'user', '用户管理', 'user', '/system/user', 'system/user/index', 2, 1, 'system:user:list', 1, 1, NOW(), NOW());
SET @user_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@user_id, 'user-add', '新增用户', '', '', '', 3, 1, 'system:user:add', 1, 1, NOW(), NOW()),
(@user_id, 'user-edit', '编辑用户', '', '', '', 3, 2, 'system:user:edit', 1, 1, NOW(), NOW()),
(@user_id, 'user-delete', '删除用户', '', '', '', 3, 3, 'system:user:delete', 1, 1, NOW(), NOW());

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@system_id, 'role', '角色管理', 'usergroup', '/system/role', 'system/role/index', 2, 2, 'system:role:list', 1, 1, NOW(), NOW());
SET @role_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@role_id, 'role-add', '新增角色', '', '', '', 3, 1, 'system:role:add', 1, 1, NOW(), NOW()),
(@role_id, 'role-edit', '编辑角色', '', '', '', 3, 2, 'system:role:edit', 1, 1, NOW(), NOW()),
(@role_id, 'role-delete', '删除角色', '', '', '', 3, 3, 'system:role:delete', 1, 1, NOW(), NOW()),
(@role_id, 'role-menu', '分配菜单', '', '', '', 3, 4, 'system:role:menu', 1, 1, NOW(), NOW());

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@system_id, 'menu', '菜单管理', 'menu-fold', '/system/menu', 'system/menu/index', 2, 3, 'system:menu:list', 1, 1, NOW(), NOW());
SET @menu_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@menu_id, 'menu-add', '新增菜单', '', '', '', 3, 1, 'system:menu:add', 1, 1, NOW(), NOW()),
(@menu_id, 'menu-edit', '编辑菜单', '', '', '', 3, 2, 'system:menu:edit', 1, 1, NOW(), NOW()),
(@menu_id, 'menu-delete', '删除菜单', '', '', '', 3, 3, 'system:menu:delete', 1, 1, NOW(), NOW());

-- 数据字典
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@system_id, 'dict', '数据字典', 'read', '/system/dict', 'system/dict/index', 2, 4, 'system:dict:list', 1, 1, NOW(), NOW());
SET @dict_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@dict_id, 'dict-type-add', '新增字典类型', '', '', '', 3, 1, 'system:dict:type:add', 1, 1, NOW(), NOW()),
(@dict_id, 'dict-type-edit', '编辑字典类型', '', '', '', 3, 2, 'system:dict:type:edit', 1, 1, NOW(), NOW()),
(@dict_id, 'dict-type-delete', '删除字典类型', '', '', '', 3, 3, 'system:dict:type:delete', 1, 1, NOW(), NOW()),
(@dict_id, 'dict-data-add', '新增字典数据', '', '', '', 3, 4, 'system:dict:data:add', 1, 1, NOW(), NOW()),
(@dict_id, 'dict-data-edit', '编辑字典数据', '', '', '', 3, 5, 'system:dict:data:edit', 1, 1, NOW(), NOW()),
(@dict_id, 'dict-data-delete', '删除字典数据', '', '', '', 3, 6, 'system:dict:data:delete', 1, 1, NOW(), NOW());

-- 图库管理
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@system_id, 'gallery', '图库管理', 'picture', '/system/gallery', 'system/gallery/index', 2, 5, 'system:gallery:list', 1, 1, NOW(), NOW());
SET @gallery_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@gallery_id, 'gallery-upload', '上传图片', '', '', '', 3, 1, 'system:gallery:upload', 1, 1, NOW(), NOW()),
(@gallery_id, 'gallery-delete', '删除图片', '', '', '', 3, 2, 'system:gallery:delete', 1, 1, NOW(), NOW()),
(@gallery_id, 'gallery-edit', '编辑图片', '', '', '', 3, 3, 'system:gallery:edit', 1, 1, NOW(), NOW());

-- 门店管理
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(0, 'store', '门店管理', 'shop', '', '', 1, 2, '', 1, 1, NOW(), NOW());
SET @store_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@store_id, 'store-list', '门店列表', 'view-list', '/store/list', 'store/list/index', 2, 1, 'store:list', 1, 1, NOW(), NOW());
SET @store_list_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@store_list_id, 'store-add', '新增门店', '', '', '', 3, 1, 'store:add', 1, 1, NOW(), NOW()),
(@store_list_id, 'store-edit', '编辑门店', '', '', '', 3, 2, 'store:edit', 1, 1, NOW(), NOW()),
(@store_list_id, 'store-delete', '删除门店', '', '', '', 3, 3, 'store:delete', 1, 1, NOW(), NOW()),
(@store_list_id, 'store-menu', '配置权限', '', '', '', 3, 4, 'store:menu', 1, 1, NOW(), NOW());

-- 供应商管理
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@store_id, 'supplier', '供应商管理', 'Shop', '/store/supplier', 'store/supplier/index', 2, 2, 'supplier:list', 1, 1, NOW(), NOW());
SET @supplier_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@supplier_id, 'supplier-add', '新增供应商', '', '', '', 3, 1, 'supplier:add', 1, 1, NOW(), NOW()),
(@supplier_id, 'supplier-edit', '编辑供应商', '', '', '', 3, 2, 'supplier:edit', 1, 1, NOW(), NOW()),
(@supplier_id, 'supplier-delete', '删除供应商', '', '', '', 3, 3, 'supplier:delete', 1, 1, NOW(), NOW());

-- 采购管理
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@store_id, 'purchase', '采购管理', 'Document', '/store/purchase', 'store/purchase/index', 2, 4, 'purchase:list', 1, 1, NOW(), NOW());
SET @purchase_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@purchase_id, 'purchase-add', '新增采购单', '', '', '', 3, 1, 'purchase:add', 1, 1, NOW(), NOW()),
(@purchase_id, 'purchase-edit', '编辑采购单', '', '', '', 3, 2, 'purchase:edit', 1, 1, NOW(), NOW()),
(@purchase_id, 'purchase-delete', '删除采购单', '', '', '', 3, 3, 'purchase:delete', 1, 1, NOW(), NOW());

-- 库存管理
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@store_id, 'inventory', '库存管理', 'Box', '/store/inventory', 'store/inventory/index', 2, 5, 'inventory:list', 1, 1, NOW(), NOW());
SET @inventory_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@inventory_id, 'inventory-in', '入库', '', '', '', 3, 1, 'inventory:in', 1, 1, NOW(), NOW()),
(@inventory_id, 'inventory-out', '出库', '', '', '', 3, 2, 'inventory:out', 1, 1, NOW(), NOW()),
(@inventory_id, 'inventory-record', '出入库记录', '', '', '', 3, 3, 'inventory:record', 1, 1, NOW(), NOW());

-- 门店记账
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@store_id, 'store-account', '门店记账', 'Wallet', '/store/account', 'store/account/index', 2, 6, 'store:account:list', 1, 1, NOW(), NOW());
SET @account_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@account_id, 'store-account-add', '新增记账', '', '', '', 3, 1, 'store:account:add', 1, 1, NOW(), NOW()),
(@account_id, 'store-account-edit', '编辑记账', '', '', '', 3, 2, 'store:account:edit', 1, 1, NOW(), NOW()),
(@account_id, 'store-account-delete', '删除记账', '', '', '', 3, 3, 'store:account:delete', 1, 1, NOW(), NOW());

-- 会员管理
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@store_id, 'store-member', '会员管理', 'User', '/store/member', 'store/member/index', 2, 7, 'store:member:list', 1, 1, NOW(), NOW());
SET @store_member_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@store_member_id, 'store-member-add', '新增会员', '', '', '', 3, 1, 'store:member:add', 1, 1, NOW(), NOW()),
(@store_member_id, 'store-member-edit', '编辑会员', '', '', '', 3, 2, 'store:member:edit', 1, 1, NOW(), NOW()),
(@store_member_id, 'store-member-delete', '删除会员', '', '', '', 3, 3, 'store:member:delete', 1, 1, NOW(), NOW()),
(@store_member_id, 'store-member-balance', '调整余额', '', '', '', 3, 4, 'store:member:balance', 1, 1, NOW(), NOW());

-- 打印机管理
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@store_id, 'printer', '打印机管理', 'printer', '', '', 1, 8, '', 1, 1, NOW(), NOW());
SET @printer_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@printer_id, 'printer-list', '打印机列表', '', '/printer/list', 'printer/list/index', 2, 1, 'printer:list', 1, 1, NOW(), NOW()),
(@printer_id, 'printer-bind', '绑定打印机', '', '', '', 3, 1, 'printer:bind', 1, 1, NOW(), NOW()),
(@printer_id, 'printer-unbind', '解绑打印机', '', '', '', 3, 2, 'printer:unbind', 1, 1, NOW(), NOW()),
(@printer_id, 'printer-edit', '编辑打印机', '', '', '', 3, 3, 'printer:edit', 1, 1, NOW(), NOW()),
(@printer_id, 'printer-query', '查询状态', '', '', '', 3, 4, 'printer:query', 1, 1, NOW(), NOW());

-- 钉钉管理
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(0, 'dingtalk', '钉钉管理', 'link', '', '', 1, 50, '', 1, 1, NOW(), NOW());
SET @dingtalk_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@dingtalk_id, 'dingtalk-robot', '机器人配置', 'robot', '/dingtalk/robot', 'dingtalk/robot/index', 2, 1, 'dingtalk:robot:list', 1, 1, NOW(), NOW());
SET @robot_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@robot_id, 'dingtalk-robot-add', '新增机器人', '', '', '', 3, 1, 'dingtalk:robot:add', 1, 1, NOW(), NOW()),
(@robot_id, 'dingtalk-robot-edit', '编辑机器人', '', '', '', 3, 2, 'dingtalk:robot:edit', 1, 1, NOW(), NOW()),
(@robot_id, 'dingtalk-robot-delete', '删除机器人', '', '', '', 3, 3, 'dingtalk:robot:delete', 1, 1, NOW(), NOW()),
(@robot_id, 'dingtalk-robot-test', '测试推送', '', '', '', 3, 4, 'dingtalk:robot:test', 1, 1, NOW(), NOW()),
(@robot_id, 'dingtalk-robot-status', '启用/禁用', '', '', '', 3, 5, 'dingtalk:robot:status', 1, 1, NOW(), NOW());

-- 门店数据
INSERT INTO stores (id, store_code, name, address, phone, business_hours, status, contact_person, remark, created_at, updated_at) VALUES
(999, 'JW9999', '总部', '系统默认总部地址', '13082848180', '全天', 1, '超级管理员', '系统默认总部门店', NOW(), NOW()),
(1, 'JW0001', '示例门店1', '杭州市西湖区文三路100号', '13800000001', '09:00-22:00', 1, '张三', '示例门店', NOW(), NOW()),
(2, 'JW0002', '示例门店2', '杭州市余杭区勾庄路200号', '13800000002', '08:00-21:00', 1, '李四', '示例门店', NOW(), NOW())
ON DUPLICATE KEY UPDATE name=name;

-- 用户数据 (密码: Admin@123456)
INSERT INTO users (id, employee_no, username, phone, password, nickname, email, store_id, role_id, status, gender, created_at, updated_at) VALUES
(999, '999999', 'admin', '13082848180', '$2a$10$6xWaEeNOICc0wmCcTS8Ac.5Iam7.zR4W.vWoUVbjHIsobDTB6L02W', '超级管理员', 'admin@tower.com', 999, 999, 1, 1, NOW(), NOW()),
(1, '000001', 'store1_admin', '13800000001', '$2a$10$6xWaEeNOICc0wmCcTS8Ac.5Iam7.zR4W.vWoUVbjHIsobDTB6L02W', '门店1管理员', 'store1@tower.com', 1, 2, 1, 1, NOW(), NOW()),
(2, '000002', 'store2_admin', '13800000002', '$2a$10$6xWaEeNOICc0wmCcTS8Ac.5Iam7.zR4W.vWoUVbjHIsobDTB6L02W', '门店2管理员', 'store2@tower.com', 2, 2, 1, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE username=username;

-- 字典数据 - 销售渠道
INSERT INTO dict_types (code, name, remark, status, created_at, updated_at) VALUES
('sales_channel', '销售渠道', '门店记账-销售渠道', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE name=VALUES(name);
SET @channel_type_id = (SELECT id FROM dict_types WHERE code = 'sales_channel');

INSERT INTO dict_data (type_id, type_code, label, value, sort, status, created_at, updated_at) VALUES
(@channel_type_id, 'sales_channel', '线下门店', 'offline', 1, 1, NOW(), NOW()),
(@channel_type_id, 'sales_channel', '美团外卖', 'meituan', 2, 1, NOW(), NOW()),
(@channel_type_id, 'sales_channel', '饿了么', 'eleme', 3, 1, NOW(), NOW()),
(@channel_type_id, 'sales_channel', '抖音', 'douyin', 4, 1, NOW(), NOW()),
(@channel_type_id, 'sales_channel', '小红书', 'xiaohongshu', 5, 1, NOW(), NOW()),
(@channel_type_id, 'sales_channel', '微信小程序', 'wechat_mini', 6, 1, NOW(), NOW()),
(@channel_type_id, 'sales_channel', '其他', 'other', 99, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE label=VALUES(label);

-- 字典数据 - 订单来源
INSERT INTO dict_types (code, name, remark, status, created_at, updated_at) VALUES
('order_source', '订单来源', '门店记账-订单来源', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE name=VALUES(name);
SET @source_type_id = (SELECT id FROM dict_types WHERE code = 'order_source');

INSERT INTO dict_data (type_id, type_code, label, value, sort, status, created_at, updated_at) VALUES
(@source_type_id, 'order_source', '堂食', 'dine_in', 1, 1, NOW(), NOW()),
(@source_type_id, 'order_source', '外卖', 'takeout', 2, 1, NOW(), NOW()),
(@source_type_id, 'order_source', '自提', 'pickup', 3, 1, NOW(), NOW()),
(@source_type_id, 'order_source', '团购', 'group_buy', 4, 1, NOW(), NOW()),
(@source_type_id, 'order_source', '预订', 'reservation', 5, 1, NOW(), NOW()),
(@source_type_id, 'order_source', '其他', 'other', 99, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE label=VALUES(label);

-- 消息模板数据
INSERT INTO message_templates (code, name, title, content, description, variables, is_enabled, created_at, updated_at) VALUES
('store_account_created', '记账通知', '📝 新记账通知 - {{.StoreName}}',
'## 📝 新记账通知 - {{.StoreName}}

**记账编号：** {{.AccountNo}}
**渠道来源：** {{.ChannelName}}
**记账日期：** {{.AccountDate}}
**操作人：** {{.OperatorName}}

### 商品明细
{{.ItemList}}

**合计金额：** ¥{{.TotalAmount}}
**商品数量：** {{.ItemCount}} 项

---
{{.CreateTime}}',
'新增记账时发送给门店负责人的通知',
'["StoreName","AccountNo","ChannelName","AccountDate","OperatorName","ItemList","TotalAmount","ItemCount","CreateTime"]',
1, NOW(), NOW()),

('inventory_created', '入库通知', '📦 新入库通知 - {{.StoreName}}',
'## 📦 新入库通知 - {{.StoreName}}

**入库单号：** {{.OrderNo}}
**入库类型：** {{.OrderType}}
**入库日期：** {{.OrderDate}}
**操作人：** {{.OperatorName}}

### 入库明细
{{.ItemList}}

**合计金额：** ¥{{.TotalAmount}}
**商品数量：** {{.ItemCount}} 项

---
{{.CreateTime}}',
'新增入库时发送给门店负责人的通知',
'["StoreName","OrderNo","OrderType","OrderDate","OperatorName","ItemList","TotalAmount","ItemCount","CreateTime"]',
1, NOW(), NOW()),

('bot_help', '机器人帮助菜单', '📋 功能菜单',
'## 📋 功能菜单

您可以发送以下命令：

**库存相关**
- 库存查询 - 查看当前库存
- 查询库存 商品名 - 搜索指定商品

**记账相关**
- 今日记账 - 查看今日记账汇总

**入库相关**
- 今日入库 - 查看今日入库记录

---
发送 **帮助** 可再次查看此菜单',
'钉钉机器人帮助菜单',
'[]',
1, NOW(), NOW()),

('bot_inventory_query', '库存查询回复', '📦 库存查询',
'## 📦 库存查询

**门店库存（共{{.Total}}项）**

{{.ItemList}}

---
{{.CreateTime}}',
'钉钉机器人库存查询回复',
'["Total","ItemList","CreateTime"]',
1, NOW(), NOW()),

('bot_today_account', '今日记账回复', '📝 今日记账',
'## 📝 今日记账汇总

**日期：** {{.Date}}
**记账笔数：** {{.Count}} 笔
**总金额：** ¥{{.TotalAmount}}

---
{{.CreateTime}}',
'钉钉机器人今日记账查询回复',
'["Date","Count","TotalAmount","CreateTime"]',
1, NOW(), NOW()),

('bot_today_inventory', '今日入库回复', '📦 今日入库',
'## 📦 今日入库汇总

**日期：** {{.Date}}
**入库单数：** {{.Count}} 单
**总入库数量：** {{.TotalQuantity}}

**入库明细：**
{{.ItemList}}

---
{{.CreateTime}}',
'钉钉机器人今日入库查询回复',
'["Date","Count","TotalQuantity","ItemList","CreateTime"]',
1, NOW(), NOW()),

('bot_search_result', '搜索结果回复', '🔍 库存搜索',
'## 🔍 库存搜索

**关键词：** {{.Keyword}}
**搜索结果（共{{.Total}}项）**

{{.ItemList}}

---
{{.CreateTime}}',
'钉钉机器人库存搜索回复',
'["Keyword","Total","ItemList","CreateTime"]',
1, NOW(), NOW()),

('bot_unknown', '未知命令回复', '🤖 智能助手',
'## 🤖 智能助手

您发送的内容：{{.Content}}

抱歉，我暂时无法理解您的意思。

发送 **帮助** 或 **菜单** 查看可用功能',
'钉钉机器人未知命令回复',
'["Content"]',
1, NOW(), NOW())
ON DUPLICATE KEY UPDATE
    name=VALUES(name),
    title=VALUES(title),
    content=VALUES(content),
    description=VALUES(description),
    variables=VALUES(variables),
    updated_at=NOW();

-- 角色菜单权限
-- 总部管理员(ID:1): 所有权限
INSERT INTO role_menus (role_id, menu_id, permissions)
SELECT 1, id, 15 FROM menus
ON DUPLICATE KEY UPDATE permissions=15;

-- 超级管理员(ID:999): 所有权限
INSERT INTO role_menus (role_id, menu_id, permissions)
SELECT 999, id, 15 FROM menus
ON DUPLICATE KEY UPDATE permissions=15;

-- 门店管理员(ID:2): 门店相关权限
INSERT INTO role_menus (role_id, menu_id, permissions)
SELECT 2, id, 15 FROM menus WHERE name LIKE 'store%' OR name LIKE 'supplier%' OR name LIKE 'purchase%' OR name LIKE 'inventory%'
ON DUPLICATE KEY UPDATE permissions=15;

-- 门店管理员赋予会员管理权限
INSERT INTO role_menus (role_id, menu_id, permissions)
SELECT 2, id, 15 FROM menus WHERE name IN ('store-member', 'store-member-add', 'store-member-edit', 'store-member-delete', 'store-member-balance')
ON DUPLICATE KEY UPDATE permissions=15;

-- ============================================
-- 4. Mock测试数据
-- ============================================

-- 供应商数据
INSERT INTO suppliers (id, supplier_code, supplier_name, contact_person, contact_phone, contact_email, supplier_address, remark, status, created_at, updated_at) VALUES
(1, 'GYS001', '杭州鲜蔬农产品有限公司', '张经理', '13800001001', 'zhang@xiansu.com', '杭州市余杭区农贸市场A区', '主营蔬菜类', 1, NOW(), NOW()),
(2, 'GYS002', '浙江海鲜水产批发中心', '李总', '13800001002', 'li@haixian.com', '杭州市萧山区水产市场', '主营海鲜水产', 1, NOW(), NOW()),
(3, 'GYS003', '金华火腿食品厂', '王师傅', '13800001003', 'wang@jinhua.com', '金华市婺城区工业园', '主营肉类制品', 1, NOW(), NOW()),
(4, 'GYS004', '温州调味品贸易公司', '陈老板', '13800001004', 'chen@tiaoweipin.com', '温州市鹿城区调味品市场', '主营调味品', 1, NOW(), NOW()),
(5, 'GYS005', '宁波粮油批发商行', '刘经理', '13800001005', 'liu@liangyou.com', '宁波市海曙区粮油市场', '主营粮油干货', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE supplier_name=VALUES(supplier_name);

-- 供应商分类数据
INSERT INTO supplier_categories (id, supplier_id, name, sort, status, created_at, updated_at) VALUES
(1, 1, '叶菜类', 1, 1, NOW(), NOW()),
(2, 1, '根茎类', 2, 1, NOW(), NOW()),
(3, 1, '瓜果类', 3, 1, NOW(), NOW()),
(4, 1, '菌菇类', 4, 1, NOW(), NOW()),
(5, 2, '鱼类', 1, 1, NOW(), NOW()),
(6, 2, '虾蟹类', 2, 1, NOW(), NOW()),
(7, 2, '贝类', 3, 1, NOW(), NOW()),
(8, 3, '猪肉', 1, 1, NOW(), NOW()),
(9, 3, '牛肉', 2, 1, NOW(), NOW()),
(10, 3, '禽类', 3, 1, NOW(), NOW()),
(11, 4, '酱油醋', 1, 1, NOW(), NOW()),
(12, 4, '香料', 2, 1, NOW(), NOW()),
(13, 4, '酱料', 3, 1, NOW(), NOW()),
(14, 5, '大米', 1, 1, NOW(), NOW()),
(15, 5, '食用油', 2, 1, NOW(), NOW()),
(16, 5, '干货', 3, 1, NOW(), NOW());

-- 供应商商品数据
INSERT INTO supplier_products (id, supplier_id, category_id, name, unit, price, spec, remark, status, created_at, updated_at) VALUES
-- 叶菜类
(1, 1, 1, '大白菜', '斤', 2.50, '新鲜', '当季蔬菜', 1, NOW(), NOW()),
(2, 1, 1, '小白菜', '斤', 3.00, '嫩叶', '', 1, NOW(), NOW()),
(3, 1, 1, '菠菜', '斤', 4.50, '新鲜', '', 1, NOW(), NOW()),
(4, 1, 1, '生菜', '斤', 5.00, '罗马生菜', '', 1, NOW(), NOW()),
(5, 1, 1, '油麦菜', '斤', 4.00, '新鲜', '', 1, NOW(), NOW()),
(6, 1, 1, '空心菜', '斤', 3.50, '新鲜', '', 1, NOW(), NOW()),
-- 根茎类
(7, 1, 2, '土豆', '斤', 2.00, '黄心', '', 1, NOW(), NOW()),
(8, 1, 2, '红薯', '斤', 2.50, '红心', '', 1, NOW(), NOW()),
(9, 1, 2, '胡萝卜', '斤', 3.00, '新鲜', '', 1, NOW(), NOW()),
(10, 1, 2, '白萝卜', '斤', 1.50, '新鲜', '', 1, NOW(), NOW()),
(11, 1, 2, '莲藕', '斤', 6.00, '九孔', '', 1, NOW(), NOW()),
(12, 1, 2, '山药', '斤', 8.00, '铁棍山药', '', 1, NOW(), NOW()),
-- 瓜果类
(13, 1, 3, '黄瓜', '斤', 3.00, '新鲜', '', 1, NOW(), NOW()),
(14, 1, 3, '西红柿', '斤', 4.00, '普罗旺斯', '', 1, NOW(), NOW()),
(15, 1, 3, '茄子', '斤', 3.50, '长茄', '', 1, NOW(), NOW()),
(16, 1, 3, '青椒', '斤', 4.00, '新鲜', '', 1, NOW(), NOW()),
(17, 1, 3, '冬瓜', '斤', 1.50, '新鲜', '', 1, NOW(), NOW()),
(18, 1, 3, '南瓜', '斤', 2.00, '贝贝南瓜', '', 1, NOW(), NOW()),
-- 菌菇类
(19, 1, 4, '香菇', '斤', 12.00, '干香菇', '', 1, NOW(), NOW()),
(20, 1, 4, '平菇', '斤', 6.00, '新鲜', '', 1, NOW(), NOW()),
(21, 1, 4, '金针菇', '包', 4.00, '150g/包', '', 1, NOW(), NOW()),
(22, 1, 4, '杏鲍菇', '斤', 8.00, '新鲜', '', 1, NOW(), NOW()),
-- 鱼类
(23, 2, 5, '草鱼', '斤', 12.00, '活鱼', '', 1, NOW(), NOW()),
(24, 2, 5, '鲈鱼', '斤', 25.00, '活鱼', '', 1, NOW(), NOW()),
(25, 2, 5, '黄鱼', '斤', 35.00, '野生', '', 1, NOW(), NOW()),
(26, 2, 5, '带鱼', '斤', 18.00, '新鲜', '', 1, NOW(), NOW()),
(27, 2, 5, '三文鱼', '斤', 45.00, '挪威进口', '', 1, NOW(), NOW()),
-- 虾蟹类
(28, 2, 6, '基围虾', '斤', 38.00, '活虾', '', 1, NOW(), NOW()),
(29, 2, 6, '大闸蟹', '只', 35.00, '3两母蟹', '', 1, NOW(), NOW()),
(30, 2, 6, '皮皮虾', '斤', 45.00, '活虾', '', 1, NOW(), NOW()),
(31, 2, 6, '小龙虾', '斤', 28.00, '4-6钱', '', 1, NOW(), NOW()),
-- 贝类
(32, 2, 7, '花蛤', '斤', 8.00, '新鲜', '', 1, NOW(), NOW()),
(33, 2, 7, '蛏子', '斤', 15.00, '新鲜', '', 1, NOW(), NOW()),
(34, 2, 7, '扇贝', '斤', 12.00, '新鲜', '', 1, NOW(), NOW()),
(35, 2, 7, '生蚝', '个', 5.00, '大号', '', 1, NOW(), NOW()),
-- 猪肉
(36, 3, 8, '五花肉', '斤', 18.00, '带皮', '', 1, NOW(), NOW()),
(37, 3, 8, '里脊肉', '斤', 22.00, '精选', '', 1, NOW(), NOW()),
(38, 3, 8, '排骨', '斤', 28.00, '肋排', '', 1, NOW(), NOW()),
(39, 3, 8, '猪蹄', '斤', 20.00, '前蹄', '', 1, NOW(), NOW()),
-- 牛肉
(40, 3, 9, '牛腩', '斤', 45.00, '新鲜', '', 1, NOW(), NOW()),
(41, 3, 9, '牛腱子', '斤', 55.00, '新鲜', '', 1, NOW(), NOW()),
(42, 3, 9, '肥牛卷', '斤', 48.00, '火锅用', '', 1, NOW(), NOW()),
-- 禽类
(43, 3, 10, '三黄鸡', '只', 35.00, '约2斤', '', 1, NOW(), NOW()),
(44, 3, 10, '鸡翅', '斤', 18.00, '新鲜', '', 1, NOW(), NOW()),
(45, 3, 10, '鸭腿', '斤', 15.00, '新鲜', '', 1, NOW(), NOW()),
-- 调味品
(46, 4, 11, '生抽', '瓶', 12.00, '500ml', '海天', 1, NOW(), NOW()),
(47, 4, 11, '老抽', '瓶', 10.00, '500ml', '海天', 1, NOW(), NOW()),
(48, 4, 11, '陈醋', '瓶', 8.00, '500ml', '山西老陈醋', 1, NOW(), NOW()),
(49, 4, 12, '八角', '斤', 35.00, '干货', '', 1, NOW(), NOW()),
(50, 4, 12, '花椒', '斤', 45.00, '四川汉源', '', 1, NOW(), NOW()),
(51, 4, 12, '桂皮', '斤', 25.00, '干货', '', 1, NOW(), NOW()),
(52, 4, 13, '豆瓣酱', '瓶', 15.00, '500g', '郫县豆瓣', 1, NOW(), NOW()),
(53, 4, 13, '甜面酱', '瓶', 10.00, '300g', '', 1, NOW(), NOW()),
-- 粮油
(54, 5, 14, '东北大米', '袋', 65.00, '10kg', '五常大米', 1, NOW(), NOW()),
(55, 5, 14, '泰国香米', '袋', 85.00, '10kg', '进口', 1, NOW(), NOW()),
(56, 5, 15, '花生油', '桶', 120.00, '5L', '鲁花', 1, NOW(), NOW()),
(57, 5, 15, '菜籽油', '桶', 80.00, '5L', '金龙鱼', 1, NOW(), NOW()),
(58, 5, 16, '干木耳', '斤', 45.00, '东北', '', 1, NOW(), NOW()),
(59, 5, 16, '干腐竹', '斤', 18.00, '优质', '', 1, NOW(), NOW()),
(60, 5, 16, '粉丝', '包', 8.00, '500g', '龙口粉丝', 1, NOW(), NOW());

-- 会员示例数据
INSERT INTO `t_member` (id, uid, phone, balance, points, level, version, created_at, updated_at) VALUES
(1, 'U001', '13800000001', 1000.00, 100, 1, 0, NOW(), NOW()),
(2, 'U002', '13800000002', 500.50, 50, 1, 0, NOW(), NOW()),
(3, 'U003', '13800000003', 2000.00, 200, 2, 0, NOW(), NOW())
ON DUPLICATE KEY UPDATE balance=VALUES(balance);

-- 流水示例数据
INSERT INTO `t_member_wallet_log` (member_id, change_type, change_amount, balance_after, related_order_no, remark, created_at) VALUES
(1, 1, 1000.00, 1000.00, '', '初始充值', NOW()),
(1, 2, -100.00, 900.00, 'SO202401090001', '消费', NOW()),
(2, 1, 500.50, 500.50, '', '初始充值', NOW()),
(3, 1, 2000.00, 2000.00, '', '初始充值', NOW());

-- 充值单示例数据
INSERT INTO `t_recharge_order` (order_no, member_id, pay_amount, gift_amount, total_amount, pay_status, pay_time, created_at) VALUES
('RE202401090001', 1, 100.00, 10.00, 110.00, 1, NOW(), NOW()),
('RE202401090002', 2, 200.00, 20.00, 220.00, 0, NULL, NOW()),
('RE202401090003', 3, 500.00, 50.00, 550.00, 1, NOW(), NOW());

-- ============================================
-- 初始化完成
-- 默认账号: 13082848180
-- 默认密码: Admin@123456 (请立即修改!)
-- ============================================
