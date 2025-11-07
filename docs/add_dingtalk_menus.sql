-- =====================================================
-- 钉钉管理菜单和权限初始化脚本
-- 注意: 报菜管理菜单(ID:40-45)已存在,无需重复添加
-- =====================================================

-- 1. 添加钉钉管理父菜单 (ID: 50)
INSERT INTO menus (id, parent_id, name, title, path, component, icon, type, sort, visible, status, permission, created_at, updated_at)
VALUES (
  50,
  0,
  'dingtalk',
  '钉钉管理',
  '/dingtalk',
  NULL,
  'link',  -- TDesign 图标: link (连接) 或 notification (通知)
  1,  -- 1=目录
  50,  -- 排序 (在报菜管理之后)
  1,
  1,
  '',
  NOW(),
  NOW()
);

-- 2. 添加机器人配置子菜单 (ID: 51)
INSERT INTO menus (id, parent_id, name, title, path, component, icon, type, sort, visible, status, permission, created_at, updated_at)
VALUES (
  51,
  50,
  'dingtalk-robot',
  '机器人配置',
  '/dingtalk/robot',
  'dingtalk/robot/index',
  'robot',  -- TDesign 图标: robot
  2,  -- 2=菜单页面
  1,
  1,
  1,
  'dingtalk:robot:list',  -- 页面访问权限
  NOW(),
  NOW()
);

-- 3. 添加机器人配置的操作按钮 (ID: 52-56)
INSERT INTO menus (id, parent_id, name, title, icon, type, sort, visible, status, permission, created_at, updated_at)
VALUES 
  (52, 51, 'dingtalk-robot-add', '新增机器人', '', 3, 1, 1, 1, 'dingtalk:robot:add', NOW(), NOW()),
  (53, 51, 'dingtalk-robot-edit', '编辑机器人', '', 3, 2, 1, 1, 'dingtalk:robot:edit', NOW(), NOW()),
  (54, 51, 'dingtalk-robot-delete', '删除机器人', '', 3, 3, 1, 1, 'dingtalk:robot:delete', NOW(), NOW()),
  (55, 51, 'dingtalk-robot-test', '测试推送', '', 3, 4, 1, 1, 'dingtalk:robot:test', NOW(), NOW()),
  (56, 51, 'dingtalk-robot-status', '启用/禁用', '', 3, 5, 1, 1, 'dingtalk:robot:status', NOW(), NOW());

-- =====================================================
-- 4. 为管理员角色分配钉钉菜单权限
-- =====================================================

-- 方式1: 如果你的 admin 角色 code 是 'admin'
INSERT INTO role_menus (role_id, menu_id, created_at)
SELECT 
  r.id,
  m.id,
  NOW()
FROM roles r
CROSS JOIN menus m
WHERE r.code = 'admin' 
  AND m.id BETWEEN 50 AND 56
  AND NOT EXISTS (
    SELECT 1 FROM role_menus rm 
    WHERE rm.role_id = r.id AND rm.menu_id = m.id
  );

-- 方式2: 如果你的超级管理员 role_id 是 999
INSERT INTO role_menus (role_id, menu_id, created_at)
SELECT 
  999,
  m.id,
  NOW()
FROM menus m
WHERE m.id BETWEEN 50 AND 56
  AND NOT EXISTS (
    SELECT 1 FROM role_menus rm 
    WHERE rm.role_id = 999 AND rm.menu_id = m.id
  );

-- =====================================================
-- 5. 验证查询
-- =====================================================

-- 查看新添加的菜单
SELECT id, parent_id, name, title, type, permission 
FROM menus 
WHERE id BETWEEN 50 AND 56
ORDER BY id;

-- 查看管理员的钉钉权限
SELECT 
  r.name AS role_name,
  m.id AS menu_id,
  m.title AS menu_title,
  m.permission
FROM role_menus rm
JOIN roles r ON rm.role_id = r.id
JOIN menus m ON rm.menu_id = m.id
WHERE m.id BETWEEN 50 AND 56
ORDER BY r.id, m.id;

-- =====================================================
-- 可选: 回滚脚本 (如需删除重新添加)
-- =====================================================

-- DELETE FROM role_menus WHERE menu_id BETWEEN 50 AND 56;
-- DELETE FROM menus WHERE id BETWEEN 50 AND 56;

