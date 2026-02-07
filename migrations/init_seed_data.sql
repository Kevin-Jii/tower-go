-- ============================================
-- Tower Go 初始化种子数据
-- 执行方式: mysql -u用户名 -p密码 数据库名 < migrations/init_seed_data.sql
-- 注意: 此脚本只应执行一次，重复执行前请先清空相关表
-- ============================================

-- 检查是否已初始化（如果菜单表有数据则跳过）
-- 如需重新初始化，请先执行: TRUNCATE TABLE menus; TRUNCATE TABLE role_menus;

-- 1. 角色数据
INSERT INTO roles (id, name, code, description, created_at, updated_at) VALUES
(1, '总部管理员', 'admin', '系统最高权限角色', NOW(), NOW()),
(2, '门店管理员', 'store_admin', '门店维度管理权限角色', NOW(), NOW()),
(3, '普通员工', 'staff', '基础操作权限角色', NOW(), NOW()),
(999, '超级管理员', 'super_admin', '系统最高权限，不可删除', NOW(), NOW())
ON DUPLICATE KEY UPDATE name=VALUES(name), description=VALUES(description), updated_at=NOW();

-- 2. 菜单数据（使用自增ID）
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

-- 供应商管理（门店管理子菜单）
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@store_id, 'supplier', '供应商管理', 'Shop', '/store/supplier', 'store/supplier/index', 2, 2, 'supplier:list', 1, 1, NOW(), NOW());
SET @supplier_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@supplier_id, 'supplier-add', '新增供应商', '', '', '', 3, 1, 'supplier:add', 1, 1, NOW(), NOW()),
(@supplier_id, 'supplier-edit', '编辑供应商', '', '', '', 3, 2, 'supplier:edit', 1, 1, NOW(), NOW()),
(@supplier_id, 'supplier-delete', '删除供应商', '', '', '', 3, 3, 'supplier:delete', 1, 1, NOW(), NOW());

-- 采购管理（门店管理子菜单）
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@store_id, 'purchase', '采购管理', 'Document', '/store/purchase', 'store/purchase/index', 2, 4, 'purchase:list', 1, 1, NOW(), NOW());
SET @purchase_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@purchase_id, 'purchase-add', '新增采购单', '', '', '', 3, 1, 'purchase:add', 1, 1, NOW(), NOW()),
(@purchase_id, 'purchase-edit', '编辑采购单', '', '', '', 3, 2, 'purchase:edit', 1, 1, NOW(), NOW()),
(@purchase_id, 'purchase-delete', '删除采购单', '', '', '', 3, 3, 'purchase:delete', 1, 1, NOW(), NOW());

-- 库存管理（门店管理子菜单）
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@store_id, 'inventory', '库存管理', 'Box', '/store/inventory', 'store/inventory/index', 2, 5, 'inventory:list', 1, 1, NOW(), NOW());
SET @inventory_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@inventory_id, 'inventory-in', '入库', '', '', '', 3, 1, 'inventory:in', 1, 1, NOW(), NOW()),
(@inventory_id, 'inventory-out', '出库', '', '', '', 3, 2, 'inventory:out', 1, 1, NOW(), NOW()),
(@inventory_id, 'inventory-record', '出入库记录', '', '', '', 3, 3, 'inventory:record', 1, 1, NOW(), NOW());

-- 门店记账（门店管理子菜单）
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@store_id, 'store-account', '门店记账', 'Wallet', '/store/account', 'store/account/index', 2, 6, 'store:account:list', 1, 1, NOW(), NOW());
SET @account_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@account_id, 'store-account-add', '新增记账', '', '', '', 3, 1, 'store:account:add', 1, 1, NOW(), NOW()),
(@account_id, 'store-account-edit', '编辑记账', '', '', '', 3, 2, 'store:account:edit', 1, 1, NOW(), NOW()),
(@account_id, 'store-account-delete', '删除记账', '', '', '', 3, 3, 'store:account:delete', 1, 1, NOW(), NOW());

-- 会员管理（门店管理子菜单）
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@store_id, 'store-member', '会员管理', 'User', '/store/member', 'store/member/index', 2, 7, 'store:member:list', 1, 1, NOW(), NOW());
SET @store_member_id = LAST_INSERT_ID();

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
(@store_member_id, 'store-member-add', '新增会员', '', '', '', 3, 1, 'store:member:add', 1, 1, NOW(), NOW()),
(@store_member_id, 'store-member-edit', '编辑会员', '', '', '', 3, 2, 'store:member:edit', 1, 1, NOW(), NOW()),
(@store_member_id, 'store-member-delete', '删除会员', '', '', '', 3, 3, 'store:member:delete', 1, 1, NOW(), NOW()),
(@store_member_id, 'store-member-balance', '调整余额', '', '', '', 3, 4, 'store:member:balance', 1, 1, NOW(), NOW());

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

-- 3. 门店数据（总部 + 示例门店）
INSERT INTO stores (id, store_code, name, address, phone, business_hours, status, contact_person, remark, created_at, updated_at) VALUES
(999, 'JW9999', '总部', '系统默认总部地址', '13082848180', '全天', 1, '超级管理员', '系统默认总部门店', NOW(), NOW()),
(1, 'JW0001', '示例门店1', '杭州市西湖区文三路100号', '13800000001', '09:00-22:00', 1, '张三', '示例门店', NOW(), NOW()),
(2, 'JW0002', '示例门店2', '杭州市余杭区勾庄路200号', '13800000002', '08:00-21:00', 1, '李四', '示例门店', NOW(), NOW())
ON DUPLICATE KEY UPDATE name=name;

-- 4. 超级管理员用户 (密码: Admin@123456)
INSERT INTO users (id, employee_no, username, phone, password, nickname, email, store_id, role_id, status, gender, created_at, updated_at) VALUES
(999, '999999', 'admin', '13082848180', '$2a$10$6xWaEeNOICc0wmCcTS8Ac.5Iam7.zR4W.vWoUVbjHIsobDTB6L02W', '超级管理员', 'admin@tower.com', 999, 999, 1, 1, NOW(), NOW()),
(1, '000001', 'store1_admin', '13800000001', '$2a$10$6xWaEeNOICc0wmCcTS8Ac.5Iam7.zR4W.vWoUVbjHIsobDTB6L02W', '门店1管理员', 'store1@tower.com', 1, 2, 1, 1, NOW(), NOW()),
(2, '000002', 'store2_admin', '13800000002', '$2a$10$6xWaEeNOICc0wmCcTS8Ac.5Iam7.zR4W.vWoUVbjHIsobDTB6L02W', '门店2管理员', 'store2@tower.com', 2, 2, 1, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE username=username;

-- 5. 字典数据
-- 销售渠道
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

-- 订单来源
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

-- 6. 角色菜单权限（所有菜
-- 总部管理员(ID:1): 所有权限
INSERT INTO role_menus (role_id, menu_id, permissions) 
SELECT 1, id, 15 FROM menus
ON DUPLICATE KEY UPDATE permissions=15;

-- 超级管理员(ID:999): 所有权限
INSERT INTO role_menus (role_id, menu_id, permissions) 
SELECT 999, id, 15 FROM menus
ON DUPLICATE KEY UPDATE permissions=15;

-- 门店管理员(ID:2): 门店相关权限（通过name匹配）
INSERT INTO role_menus (role_id, menu_id, permissions) 
SELECT 2, id, 15 FROM menus WHERE name LIKE 'store%' OR name LIKE 'supplier%' OR name LIKE 'purchase%' OR name LIKE 'inventory%'
ON DUPLICATE KEY UPDATE permissions=15;

-- ============================================
-- 初始化完成
-- 默认超级管理员账号: 13082848180
-- 默认密码: Admin@123456 (请立即修改!)
-- 
-- 权限位说明:
-- 1  = 0001 = 仅查看
-- 3  = 0011 = 查看+新增
-- 7  = 0111 = 查看+新增+修改
-- 15 = 1111 = 全部权限
-- ============================================

-- 7. 消息模板数据
INSERT INTO message_templates (code, name, title, content, description, variables, is_enabled, created_at, updated_at) VALUES
-- 业务通知模板
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

-- 钉钉机器人命令回复模板
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

-- ============================================
-- 会员管理模块（表结构和示例数据）
-- ============================================

-- 会员管理表结构
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

-- 流水表
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



-- 13. 门店管理员赋予会员管理权限
INSERT INTO role_menus (role_id, menu_id, permissions)
SELECT 2, id, 15 FROM menus WHERE name IN ('store-member', 'store-member-add', 'store-member-edit', 'store-member-delete', 'store-member-balance')
ON DUPLICATE KEY UPDATE permissions=15;

-- ============================================
-- 会员管理模块初始化完成
-- ============================================
