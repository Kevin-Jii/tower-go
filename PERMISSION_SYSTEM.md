# 菜单权限系统使用说明

## 概述

菜单权限系统包含两个层面的权限控制：

1. **菜单访问权限**：控制用户可以访问哪些菜单（由 `user-menus` 和 `user-permissions` 接口提供）
   - `user-menus`: 返回用户可访问的菜单树（用于渲染侧边栏）
   - `user-permissions`: 返回按钮级权限标识列表（如 `user:create`, `dish:update` 等，用于控制页面内的按钮显示）

2. **菜单配置权限**（本次新增）：使用二进制位控制对菜单配置的增删改查权限
   - 用于权限管理页面，控制管理员是否可以对某个菜单进行配置操作
   - 使用 4 位二进制精确控制：查看、新增、修改、删除

## 权限位说明

使用 4 位二进制表示权限（从左到右）：

```
位置:  第1位  第2位  第3位  第4位
权限:  查看   新增   修改   删除
示例:   1     1     1     0
```

### 权限常量

```go
PermView   = 8  // 1000 - 查看权限
PermCreate = 4  // 0100 - 新增权限
PermUpdate = 2  // 0010 - 修改权限
PermDelete = 1  // 0001 - 删除权限
PermAll    = 15 // 1111 - 所有权限
```

### 常见权限组合

| 二进制 | 十进制 | 说明 |
|--------|--------|------|
| 0000 | 0 | 无权限 |
| 1000 | 8 | 只读（仅查看） |
| 1100 | 12 | 查看 + 新增 |
| 1010 | 10 | 查看 + 修改 |
| 1110 | 14 | 查看 + 新增 + 修改 |
| 1111 | 15 | 所有权限 |

## API 使用示例

### 1. 为角色分配菜单权限

```bash
POST /api/v1/menus/assign-role
Content-Type: application/json

{
  "role_id": 2,
  "menu_ids": [1, 2, 3, 4],
  "perms": {
    "1": 15,  // 菜单1：所有权限 (1111)
    "2": 8,   // 菜单2：只读 (1000)
    "3": 14,  // 菜单3：查看+新增+修改 (1110)
    "4": 12   // 菜单4：查看+新增 (1100)
  }
}
```

### 2. 为门店角色分配菜单权限

```bash
POST /api/v1/menus/assign-store-role
Content-Type: application/json

{
  "store_id": 1,
  "role_id": 2,
  "menu_ids": [1, 2, 3],
  "perms": {
    "1": 15,  // 所有权限
    "2": 8,   // 只读
    "3": 14   // 查看+新增+修改
  }
}
```

### 3. 获取角色的菜单权限映射

```bash
GET /api/v1/menus/role-permissions?role_id=2

响应:
{
  "code": 200,
  "data": {
    "1": 15,
    "2": 8,
    "3": 14,
    "4": 12
  }
}
```

### 4. 获取门店角色的菜单权限映射

```bash
GET /api/v1/menus/store-role-permissions?store_id=1&role_id=2

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

### 5. 获取用户菜单（包含权限位）

```bash
GET /api/v1/menus/user-menus

响应:
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "name": "user-management",
      "title": "用户管理",
      "permissions": 15,  // 该用户对此菜单配置的权限位（用于权限管理页面）
      "children": [...]
    }
  ]
}
```

### 6. 获取用户按钮级权限（不受影响）

```bash
GET /api/v1/menus/user-permissions

响应:
{
  "code": 200,
  "data": [
    "user:view",
    "user:create",
    "user:update",
    "dish:view",
    "dish:create"
  ]
}

说明: 此接口返回的是按钮级权限标识，用于控制页面内的按钮显示，与菜单配置权限位无关。
```

## 使用场景说明

### 场景1: 按钮级权限控制（原有功能，不受影响）

用于控制页面内的操作按钮是否显示，使用 `user-permissions` 接口返回的权限标识：

```javascript
// 获取用户权限标识列表
const permissions = ['user:view', 'user:create', 'user:update'];

// 检查是否有新增用户权限
if (permissions.includes('user:create')) {
  // 显示"新增用户"按钮
}
```

### 场景2: 菜单配置权限控制（本次新增）

用于权限管理页面，控制管理员是否可以对某个菜单进行配置操作：

```javascript
// 获取菜单列表（包含权限位）
const menus = await getMenus(); // 每个菜单包含 permissions 字段

// 在权限管理页面中
menus.forEach(menu => {
  // 根据权限位控制操作按钮
  if (hasViewPermission(menu.permissions)) {
    // 可以查看此菜单配置
  }
  if (hasUpdatePermission(menu.permissions)) {
    // 可以修改此菜单的权限配置
  }
});
```

## 前端使用示例

### 1. 菜单配置权限检查工具函数

```javascript
// 权限常量
const PERM_VIEW = 8;    // 1000
const PERM_CREATE = 4;  // 0100
const PERM_UPDATE = 2;  // 0010
const PERM_DELETE = 1;  // 0001
const PERM_ALL = 15;    // 1111

// 检查是否有指定权限
function hasPermission(perms, perm) {
  return (perms & perm) === perm;
}

// 检查各项权限
function hasViewPermission(perms) {
  return hasPermission(perms, PERM_VIEW);
}

function hasCreatePermission(perms) {
  return hasPermission(perms, PERM_CREATE);
}

function hasUpdatePermission(perms) {
  return hasPermission(perms, PERM_UPDATE);
}

function hasDeletePermission(perms) {
  return hasPermission(perms, PERM_DELETE);
}
```

### 2. Vue 权限管理页面中使用

```vue
<template>
  <div class="permission-management">
    <!-- 菜单列表 -->
    <el-table :data="menus">
      <el-table-column prop="title" label="菜单名称" />
      
      <!-- 操作列：根据菜单配置权限位控制按钮显示 -->
      <el-table-column label="操作">
        <template #default="{ row }">
          <!-- 查看菜单配置 -->
          <el-button 
            v-if="hasViewPermission(row.permissions)" 
            @click="viewMenuConfig(row)">
            查看
          </el-button>
          
          <!-- 修改菜单权限配置 -->
          <el-button 
            v-if="hasUpdatePermission(row.permissions)" 
            @click="editMenuPermission(row)">
            配置权限
          </el-button>
          
          <!-- 删除菜单 -->
          <el-button 
            v-if="hasDeletePermission(row.permissions)" 
            type="danger"
            @click="deleteMenu(row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
export default {
  data() {
    return {
      menus: [
        {
          id: 1,
          title: '用户管理',
          permissions: 14 // 1110: 可查看、可配置权限、不可删除
        },
        {
          id: 2,
          title: '报表查看',
          permissions: 8  // 1000: 只能查看，不能修改或删除
        }
      ]
    }
  },
  methods: {
    hasViewPermission(perms) {
      return (perms & 8) === 8;
    },
    hasUpdatePermission(perms) {
      return (perms & 2) === 2;
    },
    hasDeletePermission(perms) {
      return (perms & 1) === 1;
    },
    viewMenuConfig(menu) {
      // 查看菜单配置
    },
    editMenuPermission(menu) {
      // 编辑菜单权限配置
    },
    deleteMenu(menu) {
      // 删除菜单
    }
  }
}
</script>
```

### 3. 权限编辑组件

```vue
<template>
  <div class="permission-editor">
    <el-checkbox-group v-model="selectedPerms">
      <el-checkbox :label="8">查看</el-checkbox>
      <el-checkbox :label="4">新增</el-checkbox>
      <el-checkbox :label="2">修改</el-checkbox>
      <el-checkbox :label="1">删除</el-checkbox>
    </el-checkbox-group>
    
    <div>权限值: {{ permissionValue }} ({{ permissionBinary }})</div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      selectedPerms: [8, 4, 2] // 默认选中查看、新增、修改
    }
  },
  computed: {
    permissionValue() {
      return this.selectedPerms.reduce((sum, val) => sum + val, 0);
    },
    permissionBinary() {
      const val = this.permissionValue;
      return [
        (val & 8) ? '1' : '0',
        (val & 4) ? '1' : '0',
        (val & 2) ? '1' : '0',
        (val & 1) ? '1' : '0'
      ].join('');
    }
  },
  methods: {
    // 从权限值初始化选中状态
    initFromPermValue(perm) {
      this.selectedPerms = [];
      if (perm & 8) this.selectedPerms.push(8);
      if (perm & 4) this.selectedPerms.push(4);
      if (perm & 2) this.selectedPerms.push(2);
      if (perm & 1) this.selectedPerms.push(1);
    }
  }
}
</script>
```

## 后端代码使用示例

### 1. 检查权限

```go
import (
    "github.com/Kevin-Jii/tower-go/model"
    "github.com/Kevin-Jii/tower-go/utils/permission"
)

// 检查用户是否有查看权限
if permission.HasViewPermission(menu.Permissions) {
    // 允许查看
}

// 检查用户是否有新增权限
if permission.HasCreatePermission(menu.Permissions) {
    // 允许新增
}

// 检查用户是否有修改权限
if permission.HasUpdatePermission(menu.Permissions) {
    // 允许修改
}

// 检查用户是否有删除权限
if permission.HasDeletePermission(menu.Permissions) {
    // 允许删除
}

// 通用权限检查
if model.HasPermission(menu.Permissions, model.PermCreate) {
    // 有新增权限
}
```

### 2. 设置权限

```go
// 设置所有权限
perms := model.PermAll // 15

// 设置只读权限
perms := model.PermView // 8

// 组合权限：查看 + 新增
perms := model.PermView | model.PermCreate // 12

// 组合权限：查看 + 新增 + 修改
perms := model.PermView | model.PermCreate | model.PermUpdate // 14
```

## 数据库迁移

执行以下 SQL 添加权限字段：

```sql
-- 为 role_menus 表添加 permissions 字段
ALTER TABLE `role_menus` 
ADD COLUMN `permissions` TINYINT UNSIGNED NOT NULL DEFAULT 15 
COMMENT '权限位：bit0=查看,bit1=新增,bit2=修改,bit3=删除' 
AFTER `menu_id`;

-- 为 store_role_menus 表添加 permissions 字段
ALTER TABLE `store_role_menus` 
ADD COLUMN `permissions` TINYINT UNSIGNED NOT NULL DEFAULT 15 
COMMENT '权限位：bit0=查看,bit1=新增,bit2=修改,bit3=删除' 
AFTER `menu_id`;
```

## 注意事项

1. **默认权限**: 新分配的菜单如果没有指定权限位，默认给予所有权限（15）
2. **权限继承**: 子菜单的权限不会自动继承父菜单，需要单独设置
3. **缓存更新**: 修改权限后会自动清除相关缓存
4. **位运算**: 使用位运算进行权限检查，性能高效
5. **兼容性**: 现有代码不受影响，未指定权限位时默认为全部权限

## 权限位计算器

可以使用以下在线工具或代码快速计算权限值：

```go
// 从二进制字符串转换为权限值
func ParsePermissionString(permStr string) uint8 {
    var perm uint8 = 0
    if len(permStr) >= 4 {
        if permStr[0] == '1' { perm |= 8 } // 查看
        if permStr[1] == '1' { perm |= 4 } // 新增
        if permStr[2] == '1' { perm |= 2 } // 修改
        if permStr[3] == '1' { perm |= 1 } // 删除
    }
    return perm
}

// 示例
perm := ParsePermissionString("1110") // 返回 14
```
