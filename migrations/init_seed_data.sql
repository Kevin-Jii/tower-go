-- ============================================
-- Tower Go 初始化种子数据
-- 执行方式: mysql -u用户名 -p密码 数据库名 < migrations/init_seed_data.sql
-- ============================================

-- 1. 角色数据
INSERT INTO roles (id, name, code, description, created_at, updated_at) VALUES
(1, '总部管理员', 'admin', '系统最高权限角色', NOW(), NOW()),
(2, '门店管理员', 'store_admin', '门店维度管理权限角色', NOW(), NOW()),
(3, '普通员工', 'staff', '基础操作权限角色', NOW(), NOW()),
(999, '超级管理员', 'super_admin', '系统最高权限，不可删除', NOW(), NOW())
ON DUPLICATE KEY UPDATE name=VALUES(name), description=VALUES(description), updated_at=NOW();

-- 2. 菜单数据
INSERT INTO menus (id, parent_id, name, title, icon, path, component, type, sort, permission, visible, status, created_at, updated_at) VALUES
-- 系统管理
(1, 0, 'system', '系统管理', 'setting', '', '', 1, 1, '', 1, 1, NOW(), NOW()),
(2, 1, 'user', '用户管理', 'user', '/system/user', 'system/user/index', 2, 1, 'system:user:list', 1, 1, NOW(), NOW()),
(3, 2, 'user-add', '新增用户', '', '', '', 3, 1, 'system:user:add', 1, 1, NOW(), NOW()),
(4, 2, 'user-edit', '编辑用户', '', '', '', 3, 2, 'system:user:edit', 1, 1, NOW(), NOW()),
(5, 2, 'user-delete', '删除用户', '', '', '', 3, 3, 'system:user:delete', 1, 1, NOW(), NOW()),
(6, 1, 'role', '角色管理', 'usergroup', '/system/role', 'system/role/index', 2, 2, 'system:role:list', 1, 1, NOW(), NOW()),
(7, 6, 'role-add', '新增角色', '', '', '', 3, 1, 'system:role:add', 1, 1, NOW(), NOW()),
(8, 6, 'role-edit', '编辑角色', '', '', '', 3, 2, 'system:role:edit', 1, 1, NOW(), NOW()),
(9, 6, 'role-delete', '删除角色', '', '', '', 3, 3, 'system:role:delete', 1, 1, NOW(), NOW()),
(10, 6, 'role-menu', '分配菜单', '', '', '', 3, 4, 'system:role:menu', 1, 1, NOW(), NOW()),
(11, 1, 'menu', '菜单管理', 'menu-fold', '/system/menu', 'system/menu/index', 2, 3, 'system:menu:list', 1, 1, NOW(), NOW()),
(12, 11, 'menu-add', '新增菜单', '', '', '', 3, 1, 'system:menu:add', 1, 1, NOW(), NOW()),
(13, 11, 'menu-edit', '编辑菜单', '', '', '', 3, 2, 'system:menu:edit', 1, 1, NOW(), NOW()),
(14, 11, 'menu-delete', '删除菜单', '', '', '', 3, 3, 'system:menu:delete', 1, 1, NOW(), NOW()),
-- 门店管理
(20, 0, 'store', '门店管理', 'shop', '', '', 1, 2, '', 1, 1, NOW(), NOW()),
(21, 20, 'store-list', '门店列表', 'view-list', '/store/list', 'store/list/index', 2, 1, 'store:list', 1, 1, NOW(), NOW()),
(22, 21, 'store-add', '新增门店', '', '', '', 3, 1, 'store:add', 1, 1, NOW(), NOW()),
(23, 21, 'store-edit', '编辑门店', '', '', '', 3, 2, 'store:edit', 1, 1, NOW(), NOW()),
(24, 21, 'store-delete', '删除门店', '', '', '', 3, 3, 'store:delete', 1, 1, NOW(), NOW()),
(25, 21, 'store-menu', '配置权限', '', '', '', 3, 4, 'store:menu', 1, 1, NOW(), NOW()),
-- 菜品管理
(30, 0, 'dish', '菜品管理', 'food', '', '', 1, 3, '', 1, 1, NOW(), NOW()),
(31, 30, 'dish-list', '菜品列表', 'view-list', '/dish/list', 'dish/list/index', 2, 1, 'dish:list', 1, 1, NOW(), NOW()),
(32, 31, 'dish-add', '新增菜品', '', '', '', 3, 1, 'dish:add', 1, 1, NOW(), NOW()),
(33, 31, 'dish-edit', '编辑菜品', '', '', '', 3, 2, 'dish:edit', 1, 1, NOW(), NOW()),
(34, 31, 'dish-delete', '删除菜品', '', '', '', 3, 3, 'dish:delete', 1, 1, NOW(), NOW()),
(35, 31, 'dish-status', '上下架', '', '', '', 3, 4, 'dish:status', 1, 1, NOW(), NOW()),
-- 报菜管理
(40, 0, 'report', '报菜管理', 'file-paste', '', '', 1, 4, '', 1, 1, NOW(), NOW()),
(41, 40, 'report-list', '报菜记录', 'view-list', '/report/list', 'report/list/index', 2, 1, 'report:list', 1, 1, NOW(), NOW()),
(42, 41, 'report-add', '创建报菜', '', '', '', 3, 1, 'report:add', 1, 1, NOW(), NOW()),
(43, 41, 'report-edit', '编辑报菜', '', '', '', 3, 2, 'report:edit', 1, 1, NOW(), NOW()),
(44, 41, 'report-delete', '删除报菜', '', '', '', 3, 3, 'report:delete', 1, 1, NOW(), NOW()),
(45, 40, 'report-stats', '数据统计', 'chart-bar', '/report/statistics', 'report/statistics/index', 2, 2, 'report:statistics', 1, 1, NOW(), NOW()),
-- 钉钉管理
(50, 0, 'dingtalk', '钉钉管理', 'link', '', '', 1, 50, '', 1, 1, NOW(), NOW()),
(51, 50, 'dingtalk-robot', '机器人配置', 'robot', '/dingtalk/robot', 'dingtalk/robot/index', 2, 1, 'dingtalk:robot:list', 1, 1, NOW(), NOW()),
(52, 51, 'dingtalk-robot-add', '新增机器人', '', '', '', 3, 1, 'dingtalk:robot:add', 1, 1, NOW(), NOW()),
(53, 51, 'dingtalk-robot-edit', '编辑机器人', '', '', '', 3, 2, 'dingtalk:robot:edit', 1, 1, NOW(), NOW()),
(54, 51, 'dingtalk-robot-delete', '删除机器人', '', '', '', 3, 3, 'dingtalk:robot:delete', 1, 1, NOW(), NOW()),
(55, 51, 'dingtalk-robot-test', '测试推送', '', '', '', 3, 4, 'dingtalk:robot:test', 1, 1, NOW(), NOW()),
(56, 51, 'dingtalk-robot-status', '启用/禁用', '', '', '', 3, 5, 'dingtalk:robot:status', 1, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE title=VALUES(title), icon=VALUES(icon), path=VALUES(path), component=VALUES(component), updated_at=NOW();

-- 3. 总部门店
INSERT INTO stores (id, store_code, name, address, phone, business_hours, status, contact_person, remark, created_at, updated_at) VALUES
(999, 'JW9999', '总部', '系统默认总部', '13082848180', '全天', 1, '超级管理员', '系统默认总部门店', NOW(), NOW())
ON DUPLICATE KEY UPDATE name=VALUES(name), updated_at=NOW();

-- 4. 超级管理员用户 (密码: Admin@123456)
-- 密码哈希值需要用 bcrypt 生成，这里使用预生成的哈希
INSERT INTO users (id, employee_no, username, phone, password, nickname, email, store_id, role_id, status, gender, created_at, updated_at) VALUES
(999, '999999', 'admin', '13082848180', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQb9tTmMYgDKwANaIJ/Ld9Ld9Ld9', '超级管理员', 'admin@tower.com', 999, 999, 1, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE username=VALUES(username), updated_at=NOW();

-- 5. 添加权限字段（如果不存在）
-- permissions: bit0=查看, bit1=新增, bit2=修改, bit3=删除, 15=全部权限
ALTER TABLE `role_menus` ADD COLUMN IF NOT EXISTS `permissions` TINYINT UNSIGNED NOT NULL DEFAULT 15 COMMENT '权限位：bit0=查看,bit1=新增,bit2=修改,bit3=删除' AFTER `menu_id`;
ALTER TABLE `store_role_menus` ADD COLUMN IF NOT EXISTS `permissions` TINYINT UNSIGNED NOT NULL DEFAULT 15 COMMENT '权限位：bit0=查看,bit1=新增,bit2=修改,bit3=删除' AFTER `menu_id`;

-- 6. 角色菜单权限（带权限位，15=全部权限）
-- 总部管理员(ID:1): 所有权限
INSERT INTO role_menus (role_id, menu_id, permissions) 
SELECT 1, id, 15 FROM menus WHERE id <= 56
ON DUPLICATE KEY UPDATE permissions=15;

-- 超级管理员(ID:999): 所有权限
INSERT INTO role_menus (role_id, menu_id, permissions) 
SELECT 999, id, 15 FROM menus WHERE id <= 56
ON DUPLICATE KEY UPDATE permissions=15;

-- 门店管理员(ID:2): 门店、菜品、报菜权限
INSERT INTO role_menus (role_id, menu_id, permissions) 
SELECT 2, id, 15 FROM menus WHERE id >= 20 AND id <= 56
ON DUPLICATE KEY UPDATE permissions=15;

-- 普通员工(ID:3): 菜品和报菜权限（不含删除，permissions=7: 查看+新增+修改）
INSERT INTO role_menus (role_id, menu_id, permissions) VALUES
(3, 30, 15), (3, 31, 7), (3, 32, 7), (3, 33, 7), (3, 35, 7),
(3, 40, 15), (3, 41, 7), (3, 42, 7), (3, 43, 7), (3, 45, 7)
ON DUPLICATE KEY UPDATE permissions=VALUES(permissions);

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
