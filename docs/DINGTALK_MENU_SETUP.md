# 钉钉管理菜单自动初始化说明

## ✅ 已完成

### 1. 自动初始化功能
已集成到启动流程中,**无需手动执行 SQL**:

```go
// bootstrap/migrate.go
func AutoMigrateAndSeeds() {
    // ... 其他初始化
    
    // 自动初始化钉钉管理菜单
    if err := utils.InitDingTalkMenuSeeds(utils.DB); err != nil {
        utils.LogError("钉钉管理菜单初始化失败", zap.Error(err))
    }
}
```

### 2. 初始化内容
**菜单 ID: 50-56**

| ID | 类型 | 名称 | 路径 | 权限标识 |
|----|------|------|------|---------|
| 50 | 目录 | 钉钉管理 | `/dingtalk` | - |
| 51 | 菜单 | 机器人配置 | `/dingtalk/robot` | `dingtalk:robot:list` |
| 52 | 按钮 | 新增机器人 | - | `dingtalk:robot:add` |
| 53 | 按钮 | 编辑机器人 | - | `dingtalk:robot:edit` |
| 54 | 按钮 | 删除机器人 | - | `dingtalk:robot:delete` |
| 55 | 按钮 | 测试推送 | - | `dingtalk:robot:test` |
| 56 | 按钮 | 启用/禁用 | - | `dingtalk:robot:status` |

### 3. 自动分配权限
- ✅ `admin` 角色 → 自动获得钉钉管理权限
- ✅ `super_admin` (ID:999) → 自动获得钉钉管理权限

## 🚀 使用方式

### 方式1: 重启服务器 (推荐)
```bash
# 停止当前服务器 (Ctrl+C)

# 重新启动
./tower-go.exe

# 或者使用 go run
go run ./cmd/main.go
```

**启动时会自动检测并初始化钉钉菜单**:
```
✅ 钉钉管理菜单创建成功 | count=7
✅ admin 角色钉钉权限分配成功 | role_id=1
✅ 超级管理员钉钉权限分配成功 | role_id=999
```

### 方式2: 手动执行 SQL (备用)
如果自动初始化失败,可以使用备用 SQL:

```bash
# 连接数据库
mysql -u root -p tower_db

# 执行 SQL 文件
source docs/add_dingtalk_menus.sql
```

## 🔍 验证

### 1. 查看菜单是否创建成功
```bash
curl http://localhost:10024/api/v1/menus \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

查找 ID 50-56 的菜单:
```json
{
  "id": 50,
  "name": "dingtalk",
  "title": "钉钉管理",
  "icon": "link",
  "type": 1
}
```

### 2. 查看用户菜单树
```bash
curl http://localhost:10024/api/v1/menus/user-menus \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

应该能看到"钉钉管理"目录和"机器人配置"子菜单。

### 3. 查看日志
```bash
tail -f logs/app.log | grep -i dingtalk
```

成功日志:
```
✅ 钉钉管理菜单创建成功 | count=7
✅ admin 角色钉钉权限分配成功
```

已存在跳过:
```
钉钉管理菜单已存在，跳过初始化 | menu_id=50
```

## 🎯 前端集成

### 1. 路由配置
前端需要创建对应的页面组件:

```typescript
// src/views/dingtalk/robot/index.vue
<template>
  <div class="dingtalk-robot-page">
    <t-card title="机器人配置">
      <!-- 机器人列表表格 -->
      <t-table
        :data="botList"
        :columns="columns"
        :loading="loading"
      >
        <!-- 操作列 -->
        <template #operation="{ row }">
          <t-button 
            v-if="hasPermission('dingtalk:robot:edit')"
            @click="handleEdit(row)"
          >编辑</t-button>
          
          <t-button 
            v-if="hasPermission('dingtalk:robot:test')"
            @click="handleTest(row)"
          >测试</t-button>
          
          <t-button 
            v-if="hasPermission('dingtalk:robot:delete')"
            theme="danger"
            @click="handleDelete(row)"
          >删除</t-button>
        </template>
      </t-table>
      
      <!-- 新增按钮 -->
      <t-button 
        v-if="hasPermission('dingtalk:robot:add')"
        @click="handleAdd"
      >
        <template #icon><add-icon /></template>
        新增机器人
      </t-button>
    </t-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { usePermission } from '@/hooks/usePermission';

const { hasPermission } = usePermission();
const botList = ref([]);
const loading = ref(false);

// 获取机器人列表
const fetchBotList = async () => {
  loading.value = true;
  try {
    const res = await fetch('/api/v1/dingtalk-bots', {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });
    const data = await res.json();
    botList.value = data.data.list;
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchBotList();
});
</script>
```

### 2. 图标说明
菜单使用的 TDesign 图标:
- `link`: 连接/链接图标 (钉钉管理目录)
- `robot`: 机器人图标 (机器人配置页面)

如果需要更换图标,修改数据库或 `seed_dingtalk_menus.go`:
```go
Icon: "notification",  // 通知图标
Icon: "chat",          // 聊天图标
Icon: "service",       // 服务图标
```

## 📋 注意事项

### 1. 幂等性
- ✅ 重复启动不会创建重复菜单
- ✅ 重复分配不会创建重复权限
- ✅ 检测到已存在自动跳过

### 2. ID 冲突
如果数据库已有 ID 50-56 的菜单:
- 自动初始化会失败
- 需要手动修改 `seed_dingtalk_menus.go` 中的 ID
- 或者删除冲突的菜单后重启

### 3. 角色权限
如果需要为其他角色分配钉钉权限:

**方式1: 通过后台管理界面**
```
系统管理 → 角色管理 → 选择角色 → 分配菜单 → 勾选"钉钉管理"
```

**方式2: 直接操作数据库**
```sql
-- 为门店管理员(store_admin)分配钉钉权限
INSERT INTO role_menus (role_id, menu_id, created_at)
SELECT 
  (SELECT id FROM roles WHERE code = 'store_admin'),
  id,
  NOW()
FROM menus 
WHERE id BETWEEN 50 AND 56
  AND NOT EXISTS (
    SELECT 1 FROM role_menus rm 
    WHERE rm.role_id = (SELECT id FROM roles WHERE code = 'store_admin')
      AND rm.menu_id = menus.id
  );
```

## 🐛 故障排查

### 问题1: 菜单没有出现
**排查步骤:**
1. 查看日志: `tail -f logs/app.log | grep "钉钉"`
2. 检查数据库: `SELECT * FROM menus WHERE id BETWEEN 50 AND 56;`
3. 检查权限: `SELECT * FROM role_menus WHERE menu_id BETWEEN 50 AND 56;`

### 问题2: 前端菜单报404
**原因:** 前端路由未配置

**解决:**
1. 确认前端存在 `src/views/dingtalk/robot/index.vue`
2. 检查路由配置是否包含 `/dingtalk/robot`

### 问题3: 按钮权限不生效
**原因:** 前端未检查权限

**解决:**
```vue
<t-button 
  v-if="hasPermission('dingtalk:robot:add')"
  @click="handleAdd"
>新增</t-button>
```

## 📚 相关文档

- [钉钉 Stream 模式技术文档](./DINGTALK_STREAM_MODE.md)
- [钉钉快速入门指南](./DINGTALK_QUICK_START.md)
- [实现总结](./DINGTALK_IMPLEMENTATION_SUMMARY.md)
- [备用 SQL 脚本](./add_dingtalk_menus.sql)

## ✅ 总结

✅ **自动化**: 重启服务器自动创建菜单  
✅ **幂等性**: 重复执行不会出错  
✅ **权限自动分配**: admin 角色自动获得权限  
✅ **备用方案**: 提供 SQL 脚本手动执行  

**现在只需重启服务器,钉钉管理菜单就会自动出现在后台!** 🎉
