-- 验证钉钉菜单是否创建成功
-- 在数据库中执行以下查询

-- 1. 查看钉钉菜单 (应该有 7 条记录: ID 50-56)
SELECT id, parent_id, name, title, type, permission, status
FROM menus 
WHERE id BETWEEN 50 AND 56
ORDER BY id;

-- 2. 查看 admin 角色的钉钉权限 (应该有 7 条记录)
SELECT 
  rm.id,
  r.name AS role_name,
  m.id AS menu_id,
  m.title AS menu_title,
  m.permission
FROM role_menus rm
JOIN roles r ON rm.role_id = r.id
JOIN menus m ON rm.menu_id = m.id
WHERE r.code = 'admin'
  AND m.id BETWEEN 50 AND 56
ORDER BY m.id;

-- 3. 查看超级管理员的钉钉权限 (应该有 7 条记录)
SELECT 
  rm.id,
  r.name AS role_name,
  m.id AS menu_id,
  m.title AS menu_title,
  m.permission
FROM role_menus rm
JOIN roles r ON rm.role_id = r.id
JOIN menus m ON rm.menu_id = m.id
WHERE r.id = 999
  AND m.id BETWEEN 50 AND 56
ORDER BY m.id;

-- 4. 查看完整的菜单树结构
SELECT 
  m1.id AS level1_id,
  m1.title AS level1_title,
  m2.id AS level2_id,
  m2.title AS level2_title,
  m3.id AS level3_id,
  m3.title AS level3_title
FROM menus m1
LEFT JOIN menus m2 ON m2.parent_id = m1.id
LEFT JOIN menus m3 ON m3.parent_id = m2.id
WHERE m1.parent_id = 0 AND m1.id = 50
ORDER BY m1.sort, m2.sort, m3.sort;
