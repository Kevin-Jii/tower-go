# 权限系统快速参考

## 权限说明

### 两种权限类型

1. **菜单访问权限**（原有功能）
   - 控制用户可以访问哪些菜单
   - `user-permissions` 接口返回按钮级权限标识（如 `user:create`）

2. **菜单配置权限**（本次新增）
   - 控制对菜单配置的增删改查权限
   - 用于权限管理页面

## 菜单配置权限位对照表

| 二进制 | 十进制 | 查看 | 新增 | 修改 | 删除 | 说明 |
|--------|--------|------|------|------|------|------|
| 0000 | 0 | ❌ | ❌ | ❌ | ❌ | 无权限 |
| 1000 | 8 | ✅ | ❌ | ❌ | ❌ | 只读 |
| 1100 | 12 | ✅ | ✅ | ❌ | ❌ | 查看+新增 |
| 1010 | 10 | ✅ | ❌ | ✅ | ❌ | 查看+修改 |
| 1001 | 9 | ✅ | ❌ | ❌ | ✅ | 查看+删除 |
| 1110 | 14 | ✅ | ✅ | ✅ | ❌ | 查看+新增+修改 |
| 1101 | 13 | ✅ | ✅ | ❌ | ✅ | 查看+新增+删除 |
| 1011 | 11 | ✅ | ❌ | ✅ | ✅ | 查看+修改+删除 |
| 1111 | 15 | ✅ | ✅ | ✅ | ✅ | 所有权限 |

## 常用权限值

```go
model.PermView   = 8   // 1000 - 查看
model.PermCreate = 4   // 0100 - 新增
model.PermUpdate = 2   // 0010 - 修改
model.PermDelete = 1   // 0001 - 删除
model.PermAll    = 15  // 1111 - 所有权限
```

## API 端点

### 角色权限管理
- `POST /api/v1/menus/assign-role` - 分配角色菜单权限
- `GET /api/v1/menus/role-permissions?role_id=X` - 获取角色权限映射

### 门店角色权限管理
- `POST /api/v1/menus/assign-store-role` - 分配门店角色菜单权限
- `GET /api/v1/menus/store-role-permissions?store_id=X&role_id=Y` - 获取门店角色权限映射

### 用户权限查询
- `GET /api/v1/menus/user-menus` - 获取当前用户菜单（含菜单配置权限位）
- `GET /api/v1/menus/user-permissions` - 获取当前用户按钮级权限标识列表（不受影响）

## 代码示例

### Go 后端

```go
// 检查权限
if permission.HasViewPermission(menu.Permissions) {
    // 有查看权限
}

// 组合权限
perms := model.PermView | model.PermCreate | model.PermUpdate // 14

// 权限检查
hasCreate := model.HasPermission(perms, model.PermCreate)
```

### JavaScript 前端

```javascript
// 权限常量
const PERM_VIEW = 8, PERM_CREATE = 4, PERM_UPDATE = 2, PERM_DELETE = 1;

// 检查权限
function hasPermission(perms, perm) {
  return (perms & perm) === perm;
}

// 使用
if (hasPermission(menu.permissions, PERM_CREATE)) {
  // 显示新增按钮
}
```

## 请求示例

```bash
# 分配权限
curl -X POST http://localhost:8080/api/v1/menus/assign-role \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "role_id": 2,
    "menu_ids": [1, 2, 3],
    "perms": {
      "1": 15,
      "2": 8,
      "3": 14
    }
  }'

# 查询权限
curl -X GET "http://localhost:8080/api/v1/menus/role-permissions?role_id=2" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 数据库字段

```sql
-- role_menus 表
permissions TINYINT UNSIGNED NOT NULL DEFAULT 15

-- store_role_menus 表
permissions TINYINT UNSIGNED NOT NULL DEFAULT 15
```

## 位运算速查

| 操作 | 符号 | 示例 | 结果 |
|------|------|------|------|
| 按位或 | \| | 8 \| 4 | 12 |
| 按位与 | & | 14 & 4 | 4 |
| 检查权限 | & | (14 & 4) == 4 | true |
| 添加权限 | \|= | perms \|= 4 | 添加新增权限 |
| 移除权限 | &= ~ | perms &= ~4 | 移除新增权限 |
