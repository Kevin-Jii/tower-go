# 菜单权限系统改造总结

## 改造目标

在原有的菜单访问权限基础上，新增**菜单配置权限**控制。

### 权限体系说明

1. **菜单访问权限**（原有功能，不受影响）
   - 控制用户可以访问哪些菜单
   - `user-menus` 接口：返回用户可访问的菜单树
   - `user-permissions` 接口：返回按钮级权限标识（如 `user:create`, `dish:update`）

2. **菜单配置权限**（本次新增）
   - 使用二进制位（4位）控制对菜单配置的增删改查权限
   - 用于权限管理页面，控制管理员是否可以对某个菜单进行配置操作
   - 第1位：查看权限 (1000 = 8)
   - 第2位：新增权限 (0100 = 4)
   - 第3位：修改权限 (0010 = 2)
   - 第4位：删除权限 (0001 = 1)

## 修改的文件

### 1. 模型层 (Model)

#### `model/menu.go`
- 添加 `Permissions uint8` 字段到 `Menu` 结构体（运行时填充）

#### `model/role_menu.go`
- 添加 `Permissions uint8` 字段到 `RoleMenu` 结构体
- 添加权限常量：`PermView`, `PermCreate`, `PermUpdate`, `PermDelete`, `PermAll`
- 添加 `HasPermission()` 函数用于权限检查
- 修改 `AssignMenusToRoleReq` 添加 `Perms map[uint]uint8` 字段

#### `model/store_role_menu.go`
- 添加 `Permissions uint8` 字段到 `StoreRoleMenu` 结构体
- 修改 `AssignStoreMenusReq` 添加 `Perms map[uint]uint8` 字段

### 2. 模块层 (Module)

#### `module/role_menu.go`
- 修改 `AssignMenusToRole()` 支持权限位参数
- 添加 `GetMenuPermissionsByRoleID()` 获取权限映射

#### `module/store_role_menu.go`
- 修改 `AssignMenusToStoreRole()` 支持权限位参数
- 添加 `GetMenuPermissionsByStoreAndRole()` 获取权限映射
- 修改 `CopyStoreMenus()` 复制权限位

#### `module/menu.go`
- 修改 `GetMenusByRoleID()` 查询时包含权限位
- 修改 `GetMenusByStoreAndRole()` 查询时包含权限位

### 3. 服务层 (Service)

#### `service/menu.go`
- 修改 `AssignMenusToRole()` 传递权限位参数
- 修改 `AssignMenusToStoreRole()` 传递权限位参数
- 添加 `GetRoleMenuPermissions()` 获取角色权限映射
- 添加 `GetStoreRoleMenuPermissions()` 获取门店角色权限映射

### 4. 控制器层 (Controller)

#### `controller/menu.go`
- 添加 `GetRoleMenuPermissions()` 接口
- 添加 `GetStoreRoleMenuPermissions()` 接口

### 5. 路由层 (Routes)

#### `bootstrap/routes.go`
- 添加 `GET /menus/role-permissions` 路由
- 添加 `GET /menus/store-role-permissions` 路由

### 6. 工具层 (Utils)

#### `utils/permission/permission.go` (新建)
- `CheckMenuPermission()` - 检查菜单权限
- `HasViewPermission()` - 检查查看权限
- `HasCreatePermission()` - 检查新增权限
- `HasUpdatePermission()` - 检查修改权限
- `HasDeletePermission()` - 检查删除权限
- `ParsePermissionString()` - 解析二进制字符串
- `FormatPermissionBits()` - 格式化权限位
- `GetPermissionDescription()` - 获取权限描述

### 7. 数据库迁移

#### `migrations/add_permissions_column.sql` (新建)
- 为 `role_menus` 表添加 `permissions` 字段
- 为 `store_role_menus` 表添加 `permissions` 字段
- 更新现有数据默认权限为 15（所有权限）

### 8. 文档

#### `PERMISSION_SYSTEM.md` (新建)
- 完整的权限系统使用说明
- API 使用示例
- 前端集成示例
- 后端代码示例

#### `PERMISSION_QUICK_REFERENCE.md` (新建)
- 权限位对照表
- 常用权限值
- API 端点列表
- 代码示例速查

#### `examples/permission_example.go` (新建)
- 权限系统使用示例代码
- 角色菜单分配示例

## 新增 API 接口

1. `GET /api/v1/menus/role-permissions?role_id=X`
   - 获取角色的菜单权限映射
   - 返回：`map[uint]uint8` (菜单ID -> 权限位)

2. `GET /api/v1/menus/store-role-permissions?store_id=X&role_id=Y`
   - 获取门店角色的菜单权限映射
   - 返回：`map[uint]uint8` (菜单ID -> 权限位)

## 兼容性说明

1. **向后兼容**：现有 API 接口保持不变
2. **默认权限**：未指定权限位时默认为 15（所有权限）
3. **可选参数**：`perms` 字段为可选，不传则使用默认值
4. **数据迁移**：现有数据会被设置为默认权限 15

## 使用示例

### 分配权限（带权限位）

```json
POST /api/v1/menus/assign-role
{
  "role_id": 2,
  "menu_ids": [1, 2, 3],
  "perms": {
    "1": 15,  // 所有权限
    "2": 8,   // 只读
    "3": 14   // 查看+新增+修改
  }
}
```

### 查询权限映射

```json
GET /api/v1/menus/role-permissions?role_id=2

响应:
{
  "code": 200,
  "data": {
    "1": 15,
    "2": 8,
    "3": 14
  }
}
```

### 前端权限检查

```javascript
// 检查是否有新增权限
if ((menu.permissions & 4) === 4) {
  // 显示新增按钮
}
```

## 数据库变更

执行 `migrations/add_permissions_column.sql` 文件中的 SQL 语句：

```sql
ALTER TABLE `role_menus` 
ADD COLUMN `permissions` TINYINT UNSIGNED NOT NULL DEFAULT 15;

ALTER TABLE `store_role_menus` 
ADD COLUMN `permissions` TINYINT UNSIGNED NOT NULL DEFAULT 15;
```

## 测试建议

1. 测试权限分配接口
2. 测试权限查询接口
3. 测试用户菜单获取（验证权限位正确返回）
4. 测试前端权限检查逻辑
5. 测试数据库迁移脚本
6. 测试权限复制功能

## 注意事项

1. 权限位使用 `uint8` 类型，范围 0-255，当前只使用低 4 位
2. 使用位运算进行权限检查，性能高效
3. 前端需要实现对应的权限检查逻辑
4. 建议在前端使用常量定义权限值，避免魔法数字
5. 权限修改后会自动清除相关缓存

## 下一步工作

1. 执行数据库迁移脚本
2. 更新前端代码实现权限检查
3. 更新 API 文档（Swagger）
4. 编写单元测试
5. 进行集成测试
