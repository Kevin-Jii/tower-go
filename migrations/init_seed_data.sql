-- ============================================
-- Tower Go 业务种子数据（仅 INSERT / 幂等数据写入）
-- 表结构、列补丁、CREATE TABLE 等见 migrations/init.sql
-- 说明：
-- - 可在 dev/prod 重复执行以增量补齐缺失数据
-- - 尽量使用 ON DUPLICATE KEY UPDATE / NOT EXISTS 保证幂等
-- - 文末含「演示模拟数据」供应商/商品/会员等，按需导入
-- - 应用启动若设置 SKIP_SEED_DATA=1，则不会自动执行本文件（请 mysql 手工 source）
-- ============================================

-- 角色数据（使用主键 id 幂等）
INSERT INTO roles (id, name, code, data_scope, description, created_at, updated_at) VALUES
(1, '总部管理员', 'admin', 1, '系统最高权限角色', NOW(), NOW()),
(2, '门店管理员', 'store_admin', 3, '门店维度管理权限角色', NOW(), NOW()),
(3, '普通员工', 'staff', 4, '基础操作权限角色', NOW(), NOW()),
(999, '超级管理员', 'super_admin', 1, '系统最高权限，不可删除', NOW(), NOW())
ON DUPLICATE KEY UPDATE
  name=VALUES(name),
  code=VALUES(code),
  data_scope=VALUES(data_scope),
  description=VALUES(description),
  updated_at=NOW();

-- ============================================
-- 菜单数据（menus 表无唯一约束，采用 NOT EXISTS + SELECT id 的方式幂等创建）
-- ============================================

-- 系统管理（目录）
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT 0, 'system', '系统管理', 'setting', '', '', 1, 1, '', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=0 AND name='system' AND type=1);
SET @system_id = (SELECT id FROM menus WHERE parent_id=0 AND name='system' AND type=1 ORDER BY id LIMIT 1);

-- 用户管理
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @system_id, 'user', '用户管理', 'user', '/system/user', 'system/user/index', 2, 1, 'system:user:list', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@system_id AND name='user' AND type=2);
SET @user_id = (SELECT id FROM menus WHERE parent_id=@system_id AND name='user' AND type=2 ORDER BY id LIMIT 1);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @user_id, 'user-add', '新增用户', '', '', '', 3, 1, 'system:user:add', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@user_id AND name='user-add' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @user_id, 'user-edit', '编辑用户', '', '', '', 3, 2, 'system:user:edit', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@user_id AND name='user-edit' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @user_id, 'user-delete', '删除用户', '', '', '', 3, 3, 'system:user:delete', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@user_id AND name='user-delete' AND type=3);

-- 角色管理
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @system_id, 'role', '角色管理', 'usergroup', '/system/role', 'system/role/index', 2, 2, 'system:role:list', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@system_id AND name='role' AND type=2);
SET @role_id = (SELECT id FROM menus WHERE parent_id=@system_id AND name='role' AND type=2 ORDER BY id LIMIT 1);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @role_id, 'role-add', '新增角色', '', '', '', 3, 1, 'system:role:add', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@role_id AND name='role-add' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @role_id, 'role-edit', '编辑角色', '', '', '', 3, 2, 'system:role:edit', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@role_id AND name='role-edit' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @role_id, 'role-delete', '删除角色', '', '', '', 3, 3, 'system:role:delete', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@role_id AND name='role-delete' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @role_id, 'role-menu', '分配菜单', '', '', '', 3, 4, 'system:role:menu', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@role_id AND name='role-menu' AND type=3);

-- 菜单管理
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @system_id, 'menu', '菜单管理', 'menu-fold', '/system/menu', 'system/menu/index', 2, 3, 'system:menu:list', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@system_id AND name='menu' AND type=2);
SET @menu_id = (SELECT id FROM menus WHERE parent_id=@system_id AND name='menu' AND type=2 ORDER BY id LIMIT 1);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @menu_id, 'menu-add', '新增菜单', '', '', '', 3, 1, 'system:menu:add', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@menu_id AND name='menu-add' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @menu_id, 'menu-edit', '编辑菜单', '', '', '', 3, 2, 'system:menu:edit', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@menu_id AND name='menu-edit' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @menu_id, 'menu-delete', '删除菜单', '', '', '', 3, 3, 'system:menu:delete', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@menu_id AND name='menu-delete' AND type=3);

-- 数据字典
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @system_id, 'dict', '数据字典', 'read', '/system/dict', 'system/dict/index', 2, 4, 'system:dict:list', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@system_id AND name='dict' AND type=2);
SET @dict_id = (SELECT id FROM menus WHERE parent_id=@system_id AND name='dict' AND type=2 ORDER BY id LIMIT 1);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @dict_id, 'dict-type-add', '新增字典类型', '', '', '', 3, 1, 'system:dict:type:add', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@dict_id AND name='dict-type-add' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @dict_id, 'dict-type-edit', '编辑字典类型', '', '', '', 3, 2, 'system:dict:type:edit', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@dict_id AND name='dict-type-edit' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @dict_id, 'dict-type-delete', '删除字典类型', '', '', '', 3, 3, 'system:dict:type:delete', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@dict_id AND name='dict-type-delete' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @dict_id, 'dict-data-add', '新增字典数据', '', '', '', 3, 4, 'system:dict:data:add', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@dict_id AND name='dict-data-add' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @dict_id, 'dict-data-edit', '编辑字典数据', '', '', '', 3, 5, 'system:dict:data:edit', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@dict_id AND name='dict-data-edit' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @dict_id, 'dict-data-delete', '删除字典数据', '', '', '', 3, 6, 'system:dict:data:delete', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@dict_id AND name='dict-data-delete' AND type=3);

-- 图库管理
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @system_id, 'gallery', '图库管理', 'picture', '/system/gallery', 'system/gallery/index', 2, 5, 'system:gallery:list', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@system_id AND name='gallery' AND type=2);
SET @gallery_id = (SELECT id FROM menus WHERE parent_id=@system_id AND name='gallery' AND type=2 ORDER BY id LIMIT 1);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @gallery_id, 'gallery-upload', '上传图片', '', '', '', 3, 1, 'system:gallery:upload', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@gallery_id AND name='gallery-upload' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @gallery_id, 'gallery-delete', '删除图片', '', '', '', 3, 2, 'system:gallery:delete', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@gallery_id AND name='gallery-delete' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @gallery_id, 'gallery-edit', '编辑图片', '', '', '', 3, 3, 'system:gallery:edit', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@gallery_id AND name='gallery-edit' AND type=3);

-- 消息模板（系统管理下）
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @system_id, 'message-template', '消息模板', 'Message', '/system/message-template', 'system/message-template/index', 2, 6, 'message:template:list', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@system_id AND name='message-template' AND type=2);
SET @msg_tpl_id = (SELECT id FROM menus WHERE parent_id=@system_id AND name='message-template' AND type=2 ORDER BY id LIMIT 1);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @msg_tpl_id, 'message-template-add', '新增模板', '', '', '', 3, 1, 'message:template:add', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@msg_tpl_id AND name='message-template-add' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @msg_tpl_id, 'message-template-edit', '编辑模板', '', '', '', 3, 2, 'message:template:edit', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@msg_tpl_id AND name='message-template-edit' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @msg_tpl_id, 'message-template-delete', '删除模板', '', '', '', 3, 3, 'message:template:delete', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@msg_tpl_id AND name='message-template-delete' AND type=3);

-- 操作日志（系统管理下）
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @system_id, 'audit-log', '操作日志', 'history', '/system/audit-log', 'system/audit-log/index', 2, 7, 'system:audit-log:list', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@system_id AND name='audit-log' AND type=2);
SET @audit_log_id = (SELECT id FROM menus WHERE parent_id=@system_id AND name='audit-log' AND type=2 ORDER BY id LIMIT 1);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @audit_log_id, 'audit-log-detail', '查看日志详情', '', '', '', 3, 1, 'system:audit-log:detail', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@audit_log_id AND name='audit-log-detail' AND type=3);

-- 门店管理（目录）
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT 0, 'store', '门店管理', 'shop', '', '', 1, 2, '', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=0 AND name='store' AND type=1);
SET @store_id = (SELECT id FROM menus WHERE parent_id=0 AND name='store' AND type=1 ORDER BY id LIMIT 1);

-- 门店列表
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_id, 'store-list', '门店列表', 'view-list', '/store/list', 'store/list/index', 2, 1, 'store:list', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_id AND name='store-list' AND type=2);
SET @store_list_id = (SELECT id FROM menus WHERE parent_id=@store_id AND name='store-list' AND type=2 ORDER BY id LIMIT 1);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_list_id, 'store-add', '新增门店', '', '', '', 3, 1, 'store:add', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_list_id AND name='store-add' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_list_id, 'store-edit', '编辑门店', '', '', '', 3, 2, 'store:edit', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_list_id AND name='store-edit' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_list_id, 'store-delete', '删除门店', '', '', '', 3, 3, 'store:delete', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_list_id AND name='store-delete' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_list_id, 'store-menu', '配置权限', '', '', '', 3, 4, 'store:menu', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_list_id AND name='store-menu' AND type=3);

-- 供应商管理
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_id, 'supplier', '供应商管理', 'Shop', '/store/supplier', 'store/supplier/index', 2, 2, 'supplier:list', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_id AND name='supplier' AND type=2);
SET @supplier_id = (SELECT id FROM menus WHERE parent_id=@store_id AND name='supplier' AND type=2 ORDER BY id LIMIT 1);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @supplier_id, 'supplier-add', '新增供应商', '', '', '', 3, 1, 'supplier:add', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@supplier_id AND name='supplier-add' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @supplier_id, 'supplier-edit', '编辑供应商', '', '', '', 3, 2, 'supplier:edit', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@supplier_id AND name='supplier-edit' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @supplier_id, 'supplier-delete', '删除供应商', '', '', '', 3, 3, 'supplier:delete', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@supplier_id AND name='supplier-delete' AND type=3);

-- 采购管理
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_id, 'purchase', '采购管理', 'Document', '/store/purchase', 'store/purchase/index', 2, 4, 'purchase:list', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_id AND name='purchase' AND type=2);
SET @purchase_id = (SELECT id FROM menus WHERE parent_id=@store_id AND name='purchase' AND type=2 ORDER BY id LIMIT 1);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @purchase_id, 'purchase-add', '新增采购单', '', '', '', 3, 1, 'purchase:add', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@purchase_id AND name='purchase-add' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @purchase_id, 'purchase-edit', '编辑采购单', '', '', '', 3, 2, 'purchase:edit', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@purchase_id AND name='purchase-edit' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @purchase_id, 'purchase-delete', '删除采购单', '', '', '', 3, 3, 'purchase:delete', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@purchase_id AND name='purchase-delete' AND type=3);

-- 库存管理
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_id, 'inventory', '库存管理', 'Box', '/store/inventory', 'store/inventory/index', 2, 5, 'inventory:list', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_id AND name='inventory' AND type=2);
SET @inventory_id = (SELECT id FROM menus WHERE parent_id=@store_id AND name='inventory' AND type=2 ORDER BY id LIMIT 1);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @inventory_id, 'inventory-in', '入库', '', '', '', 3, 1, 'inventory:in', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@inventory_id AND name='inventory-in' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @inventory_id, 'inventory-out', '出库', '', '', '', 3, 2, 'inventory:out', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@inventory_id AND name='inventory-out' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @inventory_id, 'inventory-record', '出入库记录', '', '', '', 3, 3, 'inventory:record', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@inventory_id AND name='inventory-record' AND type=3);

-- 门店记账
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_id, 'store-account', '门店记账', 'Wallet', '/store/account', 'store/account/index', 2, 6, 'store:account:list', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_id AND name='store-account' AND type=2);
SET @account_id = (SELECT id FROM menus WHERE parent_id=@store_id AND name='store-account' AND type=2 ORDER BY id LIMIT 1);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @account_id, 'store-account-add', '新增记账', '', '', '', 3, 1, 'store:account:add', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@account_id AND name='store-account-add' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @account_id, 'store-account-edit', '编辑记账', '', '', '', 3, 2, 'store:account:edit', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@account_id AND name='store-account-edit' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @account_id, 'store-account-delete', '删除记账', '', '', '', 3, 3, 'store:account:delete', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@account_id AND name='store-account-delete' AND type=3);

-- 门店返厂管理
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_id, 'store-return', '门店返厂管理', 'Van', '/store/return', 'store/return/index', 2, 7, 'store:return:list', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_id AND name='store-return' AND type=2);
SET @store_return_id = (SELECT id FROM menus WHERE parent_id=@store_id AND name='store-return' AND type=2 ORDER BY id LIMIT 1);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_return_id, 'store-return-add', '新增返厂', '', '', '', 3, 1, 'store:return:add', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_return_id AND name='store-return-add' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_return_id, 'store-return-edit', '编辑返厂', '', '', '', 3, 2, 'store:return:edit', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_return_id AND name='store-return-edit' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_return_id, 'store-return-delete', '删除返厂', '', '', '', 3, 3, 'store:return:delete', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_return_id AND name='store-return-delete' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_return_id, 'store-return-product', '维护返厂商品', '', '', '', 3, 4, 'store:return:product', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_return_id AND name='store-return-product' AND type=3);

-- 会员管理
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_id, 'store-member', '会员管理', 'User', '/store/member', 'store/member/index', 2, 7, 'store:member:list', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_id AND name='store-member' AND type=2);
SET @store_member_id = (SELECT id FROM menus WHERE parent_id=@store_id AND name='store-member' AND type=2 ORDER BY id LIMIT 1);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_member_id, 'store-member-add', '新增会员', '', '', '', 3, 1, 'store:member:add', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_member_id AND name='store-member-add' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_member_id, 'store-member-edit', '编辑会员', '', '', '', 3, 2, 'store:member:edit', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_member_id AND name='store-member-edit' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_member_id, 'store-member-delete', '删除会员', '', '', '', 3, 3, 'store:member:delete', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_member_id AND name='store-member-delete' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_member_id, 'store-member-balance', '调整余额', '', '', '', 3, 4, 'store:member:balance', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_member_id AND name='store-member-balance' AND type=3);

-- B2B 供货业务
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_id, 'b2b-supply', 'B2B供货', 'apps', '', '', 1, 8, '', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_id AND name='b2b-supply' AND type=1);
SET @b2b_id = (SELECT id FROM menus WHERE parent_id=@store_id AND name='b2b-supply' AND type=1 ORDER BY id LIMIT 1);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @b2b_id, 'b2b-dashboard', '供货管理', '', '/store/b2b', 'store/b2b/index', 2, 1, 'b2b:order:list', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@b2b_id AND name='b2b-dashboard' AND type=2);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @b2b_id, 'b2b-customer-list', '客户列表', '', '', '', 3, 1, 'b2b:customer:list', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@b2b_id AND name='b2b-customer-list' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @b2b_id, 'b2b-customer-add', '新增客户', '', '', '', 3, 2, 'b2b:customer:add', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@b2b_id AND name='b2b-customer-add' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @b2b_id, 'b2b-customer-edit', '编辑客户', '', '', '', 3, 3, 'b2b:customer:edit', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@b2b_id AND name='b2b-customer-edit' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @b2b_id, 'b2b-price-list', '供货价格', '', '', '', 3, 4, 'b2b:price:list', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@b2b_id AND name='b2b-price-list' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @b2b_id, 'b2b-price-edit', '编辑供货价', '', '', '', 3, 5, 'b2b:price:edit', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@b2b_id AND name='b2b-price-edit' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @b2b_id, 'b2b-price-delete', '删除供货价', '', '', '', 3, 6, 'b2b:price:delete', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@b2b_id AND name='b2b-price-delete' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @b2b_id, 'b2b-order-add', '新增供货单', '', '', '', 3, 7, 'b2b:order:add', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@b2b_id AND name='b2b-order-add' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @b2b_id, 'b2b-order-edit', '修改供货单状态', '', '', '', 3, 8, 'b2b:order:edit', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@b2b_id AND name='b2b-order-edit' AND type=3);

-- 打印机管理（目录）
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_id, 'printer', '打印机管理', 'printer', '', '', 1, 9, '', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_id AND name='printer' AND type=1);
SET @printer_id = (SELECT id FROM menus WHERE parent_id=@store_id AND name='printer' AND type=1 ORDER BY id LIMIT 1);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @printer_id, 'printer-list', '打印机列表', '', '/printer/list', 'printer/list/index', 2, 1, 'printer:list', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@printer_id AND name='printer-list' AND type=2);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @printer_id, 'printer-bind', '绑定打印机', '', '', '', 3, 1, 'printer:bind', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@printer_id AND name='printer-bind' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @printer_id, 'printer-unbind', '解绑打印机', '', '', '', 3, 2, 'printer:unbind', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@printer_id AND name='printer-unbind' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @printer_id, 'printer-edit', '编辑打印机', '', '', '', 3, 3, 'printer:edit', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@printer_id AND name='printer-edit' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @printer_id, 'printer-query', '查询状态', '', '', '', 3, 4, 'printer:query', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@printer_id AND name='printer-query' AND type=3);

-- 数据统计（门店下）
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_id, 'statistics-dash', '数据统计', 'DataBoard', '/store/statistics', 'store/statistics/index', 2, 9, 'statistics:dashboard', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_id AND name='statistics-dash' AND type=2);

-- 价目单（门店下）
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @store_id, 'price-list', '价目单', 'Tickets', '/store/price-list', 'store/price-list/index', 2, 10, 'price:list', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@store_id AND name='price-list' AND type=2);
SET @price_list_id = (SELECT id FROM menus WHERE parent_id=@store_id AND name='price-list' AND type=2 ORDER BY id LIMIT 1);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @price_list_id, 'price-add', '新增价目单', '', '', '', 3, 1, 'price:add', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@price_list_id AND name='price-add' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @price_list_id, 'price-edit', '编辑价目单', '', '', '', 3, 2, 'price:edit', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@price_list_id AND name='price-edit' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @price_list_id, 'price-delete', '删除价目单', '', '', '', 3, 3, 'price:delete', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@price_list_id AND name='price-delete' AND type=3);

-- 钉钉管理（目录）
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT 0, 'dingtalk', '钉钉管理', 'link', '', '', 1, 50, '', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=0 AND name='dingtalk' AND type=1);
SET @dingtalk_id = (SELECT id FROM menus WHERE parent_id=0 AND name='dingtalk' AND type=1 ORDER BY id LIMIT 1);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @dingtalk_id, 'dingtalk-robot', '机器人配置', 'robot', '/dingtalk/robot', 'dingtalk/robot/index', 2, 1, 'dingtalk:robot:list', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@dingtalk_id AND name='dingtalk-robot' AND type=2);
SET @robot_id = (SELECT id FROM menus WHERE parent_id=@dingtalk_id AND name='dingtalk-robot' AND type=2 ORDER BY id LIMIT 1);

INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @robot_id, 'dingtalk-robot-add', '新增机器人', '', '', '', 3, 1, 'dingtalk:robot:add', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@robot_id AND name='dingtalk-robot-add' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @robot_id, 'dingtalk-robot-edit', '编辑机器人', '', '', '', 3, 2, 'dingtalk:robot:edit', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@robot_id AND name='dingtalk-robot-edit' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @robot_id, 'dingtalk-robot-delete', '删除机器人', '', '', '', 3, 3, 'dingtalk:robot:delete', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@robot_id AND name='dingtalk-robot-delete' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @robot_id, 'dingtalk-robot-test', '测试推送', '', '', '', 3, 4, 'dingtalk:robot:test', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@robot_id AND name='dingtalk-robot-test' AND type=3);
INSERT INTO menus (parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at)
SELECT @robot_id, 'dingtalk-robot-status', '启用/禁用', '', '', '', 3, 5, 'dingtalk:robot:status', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE parent_id=@robot_id AND name='dingtalk-robot-status' AND type=3);

-- 门店数据（使用主键 id 幂等）
INSERT INTO stores (id, store_code, name, address, phone, business_hours, status, contact_person, remark, created_at, updated_at) VALUES
(999, 'JW9999', '总部', '系统默认总部地址', '13082848180', '全天', 1, '超级管理员', '系统默认总部门店', NOW(), NOW()),
(1, 'JW0001', '示例门店1', '杭州市西湖区文三路100号', '13800000001', '09:00-22:00', 1, '张三', '示例门店', NOW(), NOW()),
(2, 'JW0002', '示例门店2', '杭州市余杭区勾庄路200号', '13800000002', '08:00-21:00', 1, '李四', '示例门店', NOW(), NOW())
ON DUPLICATE KEY UPDATE
  store_code=VALUES(store_code),
  name=VALUES(name),
  address=VALUES(address),
  phone=VALUES(phone),
  business_hours=VALUES(business_hours),
  status=VALUES(status),
  contact_person=VALUES(contact_person),
  remark=VALUES(remark),
  updated_at=NOW();

-- 用户数据 (密码: Admin@123456)（使用主键 id 幂等；同时 phone/employee_no 也有唯一索引）
INSERT INTO users (id, employee_no, username, phone, password, nickname, email, store_id, role_id, status, gender, created_at, updated_at) VALUES
(999, '999999', 'admin', '13082848180', '$2a$10$6xWaEeNOICc0wmCcTS8Ac.5Iam7.zR4W.vWoUVbjHIsobDTB6L02W', '超级管理员', 'admin@tower.com', 0, 999, 1, 1, NOW(), NOW()),
(1, '000001', 'store1_admin', '13800000001', '$2a$10$6xWaEeNOICc0wmCcTS8Ac.5Iam7.zR4W.vWoUVbjHIsobDTB6L02W', '门店1管理员', 'store1@tower.com', 1, 2, 1, 1, NOW(), NOW()),
(2, '000002', 'store2_admin', '13800000002', '$2a$10$6xWaEeNOICc0wmCcTS8Ac.5Iam7.zR4W.vWoUVbjHIsobDTB6L02W', '门店2管理员', 'store2@tower.com', 2, 2, 1, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE
  employee_no=VALUES(employee_no),
  username=VALUES(username),
  phone=VALUES(phone),
  password=VALUES(password),
  nickname=VALUES(nickname),
  email=VALUES(email),
  store_id=VALUES(store_id),
  role_id=VALUES(role_id),
  status=VALUES(status),
  gender=VALUES(gender),
  updated_at=NOW();

-- 字典数据 - 销售渠道
INSERT INTO dict_types (code, name, remark, status, created_at, updated_at) VALUES
('sales_channel', '销售渠道', '门店记账-销售渠道', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE name=VALUES(name), remark=VALUES(remark), status=VALUES(status), updated_at=NOW();
SET @channel_type_id = (SELECT id FROM dict_types WHERE code = 'sales_channel' ORDER BY id LIMIT 1);

INSERT INTO dict_data (type_id, type_code, label, value, sort, status, created_at, updated_at) VALUES
(@channel_type_id, 'sales_channel', '线下门店', 'offline', 1, 1, NOW(), NOW()),
(@channel_type_id, 'sales_channel', '美团外卖', 'meituan', 2, 1, NOW(), NOW()),
(@channel_type_id, 'sales_channel', '饿了么', 'eleme', 3, 1, NOW(), NOW()),
(@channel_type_id, 'sales_channel', '抖音', 'douyin', 4, 1, NOW(), NOW()),
(@channel_type_id, 'sales_channel', '小红书', 'xiaohongshu', 5, 1, NOW(), NOW()),
(@channel_type_id, 'sales_channel', '微信小程序', 'wechat_mini', 6, 1, NOW(), NOW()),
(@channel_type_id, 'sales_channel', '其他', 'other', 99, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE
  label=VALUES(label),
  sort=VALUES(sort),
  status=VALUES(status),
  updated_at=NOW();

-- 字典数据 - 商品单位（用于商品规格换算）
INSERT INTO dict_types (code, name, remark, status, created_at, updated_at) VALUES
('product_unit', '商品单位', '商品规格单位（瓶/箱/桶/L等）', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE name=VALUES(name), remark=VALUES(remark), status=VALUES(status), updated_at=NOW();
SET @product_unit_type_id = (SELECT id FROM dict_types WHERE code = 'product_unit' ORDER BY id LIMIT 1);

INSERT INTO dict_data (type_id, type_code, label, value, sort, status, created_at, updated_at) VALUES
(@product_unit_type_id, 'product_unit', '瓶', 'bottle', 1, 1, NOW(), NOW()),
(@product_unit_type_id, 'product_unit', '箱', 'case', 2, 1, NOW(), NOW()),
(@product_unit_type_id, 'product_unit', '桶', 'barrel', 3, 1, NOW(), NOW()),
(@product_unit_type_id, 'product_unit', '升', 'L', 4, 1, NOW(), NOW()),
(@product_unit_type_id, 'product_unit', '毫升', 'ml', 5, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE
  label=VALUES(label),
  sort=VALUES(sort),
  status=VALUES(status),
  updated_at=NOW();

-- 字典数据 - 订单来源
INSERT INTO dict_types (code, name, remark, status, created_at, updated_at) VALUES
('order_source', '订单来源', '门店记账-订单来源', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE name=VALUES(name), remark=VALUES(remark), status=VALUES(status), updated_at=NOW();
SET @source_type_id = (SELECT id FROM dict_types WHERE code = 'order_source' ORDER BY id LIMIT 1);

INSERT INTO dict_data (type_id, type_code, label, value, sort, status, created_at, updated_at) VALUES
(@source_type_id, 'order_source', '堂食', 'dine_in', 1, 1, NOW(), NOW()),
(@source_type_id, 'order_source', '外卖', 'takeout', 2, 1, NOW(), NOW()),
(@source_type_id, 'order_source', '自提', 'pickup', 3, 1, NOW(), NOW()),
(@source_type_id, 'order_source', '团购', 'group_buy', 4, 1, NOW(), NOW()),
(@source_type_id, 'order_source', '预订', 'reservation', 5, 1, NOW(), NOW()),
(@source_type_id, 'order_source', '其他', 'other', 99, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE
  label=VALUES(label),
  sort=VALUES(sort),
  status=VALUES(status),
  updated_at=NOW();

-- 字典数据 - 出入库原因（与 bootstrap/dict_seed.go 保持一致，便于纯 SQL 初始化库）
INSERT INTO dict_types (code, name, remark, status, created_at, updated_at) VALUES
('inventory_reason', '出入库原因', '库存管理-出入库原因', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE name=VALUES(name), remark=VALUES(remark), status=VALUES(status), updated_at=NOW();
SET @inventory_reason_type_id = (SELECT id FROM dict_types WHERE code = 'inventory_reason' ORDER BY id LIMIT 1);

INSERT INTO dict_data (type_id, type_code, label, value, sort, status, created_at, updated_at) VALUES
(@inventory_reason_type_id, 'inventory_reason', '采购入库', 'purchase_in', 1, 1, NOW(), NOW()),
(@inventory_reason_type_id, 'inventory_reason', '退货入库', 'return_in', 2, 1, NOW(), NOW()),
(@inventory_reason_type_id, 'inventory_reason', '调拨入库', 'transfer_in', 3, 1, NOW(), NOW()),
(@inventory_reason_type_id, 'inventory_reason', '盘盈入库', 'inventory_in', 4, 1, NOW(), NOW()),
(@inventory_reason_type_id, 'inventory_reason', '销售出库', 'sale_out', 10, 1, NOW(), NOW()),
(@inventory_reason_type_id, 'inventory_reason', '报损出库', 'loss_out', 11, 1, NOW(), NOW()),
(@inventory_reason_type_id, 'inventory_reason', '调拨出库', 'transfer_out', 12, 1, NOW(), NOW()),
(@inventory_reason_type_id, 'inventory_reason', '盘亏出库', 'inventory_out', 13, 1, NOW(), NOW()),
(@inventory_reason_type_id, 'inventory_reason', '其他', 'other', 99, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE
  label=VALUES(label),
  sort=VALUES(sort),
  status=VALUES(status),
  updated_at=NOW();

-- 消息模板数据（code 有唯一索引）
INSERT INTO message_templates (code, name, title, content, description, variables, is_enabled, created_at, updated_at) VALUES
('store_account_created', '记账通知', '📝 新记账通知 - {{.StoreName}}',
'## 📝 新记账通知 - {{.StoreName}}

**记账编号：** {{.AccountNo}}
**渠道来源：** {{.ChannelName}}
**记账日期：** {{.AccountDate}}
**操作人：** {{.OperatorName}}

**记账明细：**
{{range .Items}}
- {{.ProductName}} × {{.Quantity}} = {{.Amount}}
{{end}}

**合计：** {{.TotalAmount}}

> 本消息由系统自动发送',
'门店记账创建通知模板',
'["StoreName","AccountNo","ChannelName","AccountDate","OperatorName","Items","TotalAmount"]',
1, NOW(), NOW()),
('bot_help', '帮助菜单', '📚 帮助中心',
'## 📚 帮助中心

发送以下指令可快速使用功能：

- **菜单**：查看可用菜单
- **记账**：创建记账
- **库存**：库存相关操作
- **供应商**：供应商管理

如需更多帮助请联系管理员',
'钉钉机器人帮助菜单模板',
'[]',
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
  is_enabled=VALUES(is_enabled),
  updated_at=NOW();

-- 角色菜单权限（role_menus 一般会有联合唯一键；即使没有，重复执行也不会影响数据正确性）
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
SELECT 2, id, 15 FROM menus
WHERE name LIKE 'store%'
   OR name LIKE 'supplier%'
   OR name LIKE 'purchase%'
   OR name LIKE 'inventory%'
   OR name LIKE 'printer%'
   OR name LIKE 'b2b-%'
ON DUPLICATE KEY UPDATE permissions=15;

-- 门店管理员赋予会员管理权限
INSERT INTO role_menus (role_id, menu_id, permissions)
SELECT 2, id, 15 FROM menus WHERE name IN ('store-member', 'store-member-add', 'store-member-edit', 'store-member-delete', 'store-member-balance')
ON DUPLICATE KEY UPDATE permissions=15;

-- 门店管理员：数据统计、价目单
INSERT INTO role_menus (role_id, menu_id, permissions)
SELECT 2, id, 15 FROM menus WHERE permission IN ('statistics:dashboard', 'price:list', 'price:add', 'price:edit', 'price:delete')
ON DUPLICATE KEY UPDATE permissions=15;

-- 门店管理员：系统管理目录 + 图库（与路由 /galleries 的 system:gallery:* 鉴权一致；未配置 store_role_menus 时走角色默认菜单）
INSERT INTO role_menus (role_id, menu_id, permissions)
SELECT 2, id, 15 FROM menus
WHERE name IN ('system', 'gallery', 'gallery-upload', 'gallery-delete', 'gallery-edit')
ON DUPLICATE KEY UPDATE permissions=VALUES(permissions);

-- 门店管理员：本店用户管理（仅受后端 store_id 隔离影响，不可跨店）
INSERT INTO role_menus (role_id, menu_id, permissions)
SELECT 2, id, 15 FROM menus
WHERE name IN ('user', 'user-add', 'user-edit', 'user-delete', 'user-reset-password')
   OR permission IN ('system:user:list', 'system:user:add', 'system:user:edit', 'system:user:delete')
ON DUPLICATE KEY UPDATE permissions=VALUES(permissions);

-- 普通员工：默认菜单/权限与门店管理员一致（未配置 store_role_menus 时生效；各门店可在「配置权限」中再裁剪）
INSERT INTO role_menus (role_id, menu_id, permissions)
SELECT 3, menu_id, permissions FROM role_menus WHERE role_id = 2
ON DUPLICATE KEY UPDATE permissions = VALUES(permissions);

-- ============================================
-- 演示模拟数据（供应商 / 分类 / 商品 / 门店绑定 / 库存 / 会员）
-- 依赖：stores 已有 id=1、2；与业务种子独立，按编码/手机号幂等
-- ============================================

INSERT INTO suppliers (supplier_code, supplier_name, contact_person, contact_phone, supplier_address, remark, status, created_at, updated_at)
SELECT 'SEED_DEMO_A', '演示供应商·清泉配送', '张敏', '13800138001', '杭州市余杭区物流园A区', 'init_seed_data 模拟', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM suppliers WHERE supplier_code = 'SEED_DEMO_A');
INSERT INTO suppliers (supplier_code, supplier_name, contact_person, contact_phone, supplier_address, remark, status, created_at, updated_at)
SELECT 'SEED_DEMO_B', '演示供应商·鲜达农贸', '李强', '13800138002', '上海市嘉定区江桥市场', 'init_seed_data 模拟', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM suppliers WHERE supplier_code = 'SEED_DEMO_B');

SET @seed_sup_a = (SELECT id FROM suppliers WHERE supplier_code = 'SEED_DEMO_A' ORDER BY id LIMIT 1);
SET @seed_sup_b = (SELECT id FROM suppliers WHERE supplier_code = 'SEED_DEMO_B' ORDER BY id LIMIT 1);

INSERT INTO supplier_categories (supplier_id, name, sort, status, created_at, updated_at)
SELECT @seed_sup_a, '桶装水', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM supplier_categories WHERE supplier_id = @seed_sup_a AND name = '桶装水');
INSERT INTO supplier_categories (supplier_id, name, sort, status, created_at, updated_at)
SELECT @seed_sup_a, '饮料', 2, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM supplier_categories WHERE supplier_id = @seed_sup_a AND name = '饮料');
INSERT INTO supplier_categories (supplier_id, name, sort, status, created_at, updated_at)
SELECT @seed_sup_b, '蔬菜', 1, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM supplier_categories WHERE supplier_id = @seed_sup_b AND name = '蔬菜');
INSERT INTO supplier_categories (supplier_id, name, sort, status, created_at, updated_at)
SELECT @seed_sup_b, '豆制品', 2, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM supplier_categories WHERE supplier_id = @seed_sup_b AND name = '豆制品');

SET @seed_cat_a_water = (SELECT id FROM supplier_categories WHERE supplier_id = @seed_sup_a AND name = '桶装水' ORDER BY id LIMIT 1);
SET @seed_cat_a_drink = (SELECT id FROM supplier_categories WHERE supplier_id = @seed_sup_a AND name = '饮料' ORDER BY id LIMIT 1);
SET @seed_cat_b_veg = (SELECT id FROM supplier_categories WHERE supplier_id = @seed_sup_b AND name = '蔬菜' ORDER BY id LIMIT 1);
SET @seed_cat_b_tofu = (SELECT id FROM supplier_categories WHERE supplier_id = @seed_sup_b AND name = '豆制品' ORDER BY id LIMIT 1);

INSERT INTO supplier_products (supplier_id, category_id, name, unit, price, bottle_price, case_price, bottles_per_case, spec, remark, status, created_at, updated_at)
SELECT @seed_sup_a, @seed_cat_a_water, '农夫山泉 19L 桶装水', '桶', 25.00, 25.00, 140.00, 6, '19L', '模拟商品', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM supplier_products WHERE supplier_id = @seed_sup_a AND name = '农夫山泉 19L 桶装水');
INSERT INTO supplier_products (supplier_id, category_id, name, unit, price, bottle_price, case_price, bottles_per_case, spec, remark, status, created_at, updated_at)
SELECT @seed_sup_a, @seed_cat_a_water, '怡宝 18.9L 桶装水', '桶', 22.00, 22.00, 120.00, 6, '18.9L', '模拟商品', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM supplier_products WHERE supplier_id = @seed_sup_a AND name = '怡宝 18.9L 桶装水');
INSERT INTO supplier_products (supplier_id, category_id, name, unit, price, bottle_price, case_price, bottles_per_case, spec, remark, status, created_at, updated_at)
SELECT @seed_sup_a, @seed_cat_a_drink, '可乐 500ml', '瓶', 3.00, 3.00, 36.00, 12, '500ml', '模拟商品', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM supplier_products WHERE supplier_id = @seed_sup_a AND name = '可乐 500ml');
INSERT INTO supplier_products (supplier_id, category_id, name, unit, price, bottle_price, case_price, bottles_per_case, spec, remark, status, created_at, updated_at)
SELECT @seed_sup_b, @seed_cat_b_veg, '生菜', '斤', 4.50, 4.50, 0, 1, '新鲜', '模拟商品', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM supplier_products WHERE supplier_id = @seed_sup_b AND name = '生菜');
INSERT INTO supplier_products (supplier_id, category_id, name, unit, price, bottle_price, case_price, bottles_per_case, spec, remark, status, created_at, updated_at)
SELECT @seed_sup_b, @seed_cat_b_tofu, '嫩豆腐', '盒', 3.20, 3.20, 0, 1, '350g', '模拟商品', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM supplier_products WHERE supplier_id = @seed_sup_b AND name = '嫩豆腐');

SET @seed_prod_nfs = (SELECT id FROM supplier_products WHERE supplier_id = @seed_sup_a AND name = '农夫山泉 19L 桶装水' ORDER BY id LIMIT 1);

INSERT INTO store_suppliers (store_id, supplier_id, status, created_at, updated_at)
SELECT 1, @seed_sup_a, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM store_suppliers WHERE store_id = 1 AND supplier_id = @seed_sup_a);
INSERT INTO store_suppliers (store_id, supplier_id, status, created_at, updated_at)
SELECT 2, @seed_sup_a, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM store_suppliers WHERE store_id = 2 AND supplier_id = @seed_sup_a);
INSERT INTO store_suppliers (store_id, supplier_id, status, created_at, updated_at)
SELECT 1, @seed_sup_b, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM store_suppliers WHERE store_id = 1 AND supplier_id = @seed_sup_b);
INSERT INTO store_suppliers (store_id, supplier_id, status, created_at, updated_at)
SELECT 2, @seed_sup_b, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM store_suppliers WHERE store_id = 2 AND supplier_id = @seed_sup_b);

INSERT INTO inventories (store_id, product_id, quantity, unit, created_at, updated_at)
SELECT 1, @seed_prod_nfs, 120.00, '桶', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM inventories WHERE store_id = 1 AND product_id = @seed_prod_nfs);
INSERT INTO inventories (store_id, product_id, quantity, unit, created_at, updated_at)
SELECT 2, @seed_prod_nfs, 48.00, '桶', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM inventories WHERE store_id = 2 AND product_id = @seed_prod_nfs);

INSERT INTO t_member (uid, name, phone, balance, points, level, version, created_at, updated_at)
SELECT 'SEED-M-09001', '演示会员·阿林', '13900009001', 188.00, 50, 2, 0, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM t_member WHERE phone = '13900009001');
INSERT INTO t_member (uid, name, phone, balance, points, level, version, created_at, updated_at)
SELECT 'SEED-M-09002', '演示会员·周姐', '13900009002', 56.50, 10, 1, 0, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM t_member WHERE phone = '13900009002');
