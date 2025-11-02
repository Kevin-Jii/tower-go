# Tower-Go 多门店报菜管理系统 API 指南

## 项目概述

这是一个基于 Go + Gin + GORM 的多门店报菜管理系统，支持：
- ✅ 多门店数据隔离
- ✅ 角色权限管理（总部管理员、门店管理员、员工）
- ✅ JWT 认证
- ✅ 门店管理
- ✅ 菜品管理
- ✅ 报菜记录管理
- ✅ 统计分析功能
- ✅ Swagger API 文档

## 快速启动

### 1. 启动服务

```bash
# 编译
go build ./cmd/main.go

# 运行
./main.exe
```

服务默认运行在 `http://localhost:8080`

### 2. 访问 Swagger 文档

访问地址：http://localhost:8080/swagger/index.html

在 Swagger UI 中可以查看所有 API 接口文档和进行接口测试。

## API 架构

### 认证机制

系统使用 JWT Token 进行认证，Token 包含以下信息：
- `userID`: 用户 ID
- `username`: 用户名
- `storeID`: 门店 ID（用于数据隔离）
- `roleCode`: 角色代码（admin/store_admin/staff）

### 角色说明

| 角色代码 | 角色名称 | 权限说明 |
|---------|---------|---------|
| `admin` | 总部管理员 | 可以管理所有门店、查看汇总数据 |
| `store_admin` | 门店管理员 | 只能管理自己门店的数据 |
| `staff` | 普通员工 | 只能查看和操作自己门店的数据 |

## API 端点列表

### 1. 认证相关 (`/api/v1/auth`)

| 方法 | 路径 | 说明 | 需要认证 |
|-----|------|------|---------|
| POST | `/auth/register` | 用户注册 | ❌ |
| POST | `/auth/login` | 用户登录 | ❌ |

### 2. 用户管理 (`/api/v1/users`)

| 方法 | 路径 | 说明 | 需要认证 | 数据隔离 |
|-----|------|------|---------|---------|
| GET | `/users/profile` | 获取个人信息 | ✅ | - |
| PUT | `/users/profile` | 更新个人信息 | ✅ | - |
| POST | `/users` | 创建用户 | ✅ | ✅ |
| GET | `/users` | 用户列表 | ✅ | ✅ |
| GET | `/users/:id` | 用户详情 | ✅ | ✅ |
| PUT | `/users/:id` | 更新用户 | ✅ | ✅ |
| DELETE | `/users/:id` | 删除用户 | ✅ | ✅ |

### 3. 门店管理 (`/api/v1/stores`)

| 方法 | 路径 | 说明 | 需要认证 | 权限要求 |
|-----|------|------|---------|---------|
| POST | `/stores` | 创建门店 | ✅ | 仅管理员 |
| GET | `/stores` | 门店列表 | ✅ | 所有用户 |
| GET | `/stores/:id` | 门店详情 | ✅ | 所有用户 |
| PUT | `/stores/:id` | 更新门店 | ✅ | 仅管理员 |
| DELETE | `/stores/:id` | 删除门店 | ✅ | 仅管理员 |

### 4. 菜品管理 (`/api/v1/dishes`)

| 方法 | 路径 | 说明 | 需要认证 | 数据隔离 |
|-----|------|------|---------|---------|
| POST | `/dishes` | 创建菜品 | ✅ | ✅ |
| GET | `/dishes` | 菜品列表 | ✅ | ✅ |
| GET | `/dishes?category=xxx` | 按分类查询 | ✅ | ✅ |
| GET | `/dishes/:id` | 菜品详情 | ✅ | ✅ |
| PUT | `/dishes/:id` | 更新菜品 | ✅ | ✅ |
| DELETE | `/dishes/:id` | 删除菜品 | ✅ | ✅ |

### 5. 报菜记录管理 (`/api/v1/menu-reports`)

| 方法 | 路径 | 说明 | 需要认证 | 数据隔离 |
|-----|------|------|---------|---------|
| POST | `/menu-reports` | 创建报菜记录 | ✅ | ✅ |
| GET | `/menu-reports` | 报菜记录列表 | ✅ | ✅ |
| GET | `/menu-reports?start_date=xxx&end_date=xxx` | 按日期范围查询 | ✅ | ✅ |
| GET | `/menu-reports/statistics` | 获取统计数据 | ✅ | ✅ (管理员可跨店) |
| GET | `/menu-reports/:id` | 报菜记录详情 | ✅ | ✅ |
| PUT | `/menu-reports/:id` | 更新报菜记录 | ✅ | ✅ |
| DELETE | `/menu-reports/:id` | 删除报菜记录 | ✅ | ✅ |

## 使用示例

### 1. 登录获取 Token

```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "phone": "13800138000",
  "password": "123456"
}
```

响应：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "张三",
      "store_id": 1,
      "role_code": "store_admin"
    }
  }
}
```

### 2. 使用 Token 访问受保护接口

在请求头中添加：
```
Authorization: Bearer <your_token_here>
```

### 3. 创建门店（仅管理员）

```bash
POST /api/v1/stores
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "name": "北京朝阳店",
  "address": "北京市朝阳区xxx路xxx号",
  "phone": "010-12345678",
  "manager": "李经理"
}
```

### 4. 创建菜品

```bash
POST /api/v1/dishes
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "宫保鸡丁",
  "category": "川菜",
  "price": 38.00,
  "description": "经典川菜",
  "is_available": true
}
```

### 5. 创建报菜记录

```bash
POST /api/v1/menu-reports
Authorization: Bearer <token>
Content-Type: application/json

{
  "dish_id": 1,
  "quantity": 50,
  "remark": "今日备货"
}
```

### 6. 查询统计数据

```bash
GET /api/v1/menu-reports/statistics?start_date=2025-11-01&end_date=2025-11-02
Authorization: Bearer <token>
```

## 数据隔离说明

### 门店级数据隔离

系统通过 JWT Token 中的 `storeID` 实现自动数据隔离：

1. **用户**: 只能看到和操作本门店的用户
2. **菜品**: 只能看到和操作本门店的菜品
3. **报菜记录**: 只能看到和操作本门店的报菜记录

### 总部管理员特权

角色代码为 `admin` 的用户可以：
- 创建、修改、删除门店
- 查看所有门店的统计数据（通过 `/menu-reports/statistics` 接口）
- 访问跨门店数据（未来扩展）

## 数据库表结构

### 核心表

1. **roles** - 角色表
2. **stores** - 门店表
3. **users** - 用户表
4. **dishes** - 菜品表
5. **menu_reports** - 报菜记录表

### 外键关系

- `users.store_id` → `stores.id`
- `users.role_id` → `roles.id`
- `dishes.store_id` → `stores.id`
- `menu_reports.store_id` → `stores.id`
- `menu_reports.dish_id` → `dishes.id`
- `menu_reports.user_id` → `users.id`

## 开发说明

### 项目结构

```
tower-go/
├── cmd/
│   └── main.go           # 程序入口
├── config/               # 配置文件
├── controller/           # 控制器层（HTTP 处理）
├── service/              # 服务层（业务逻辑）
├── module/               # 数据访问层（DAO）
├── model/                # 数据模型
├── middleware/           # 中间件
├── utils/                # 工具函数
└── docs/                 # Swagger 文档

```

### 开发规范

1. **三层架构**: Controller → Service → Module
2. **数据隔离**: 所有门店相关数据查询都需要带 `store_id` 过滤
3. **错误处理**: 使用统一的 `utils.Response` 返回格式
4. **认证**: 受保护路由必须使用 `middleware.AuthMiddleware()` 或 `middleware.StoreAuthMiddleware()`

### 重新生成 Swagger 文档

```bash
swag init -g cmd/main.go --output docs
```

## 常见问题

### Q1: 数据库迁移外键约束错误

**问题**: AutoMigrate 时出现外键约束错误

**原因**: 已有用户数据但没有对应的门店数据

**解决方案**:
1. 清空旧数据重新开始
2. 或手动修复数据：先创建门店，再更新用户的 `store_id`

### Q2: Token 认证失败

**检查清单**:
1. Token 格式是否正确：`Bearer <token>`
2. Token 是否过期
3. Token 解析是否成功

### Q3: 访问数据为空

**检查清单**:
1. 确认 Token 中的 `storeID` 是否正确
2. 确认数据是否确实属于该门店
3. 如果是管理员，确认是否需要特殊处理

## 后续扩展

- [ ] 菜品销售排行榜
- [ ] 多维度统计报表
- [ ] 数据导出功能
- [ ] 图表可视化
- [ ] 移动端适配
- [ ] 通知推送功能

## 技术栈

- **后端框架**: Gin v1.10.0
- **ORM**: GORM v1.31.0
- **数据库**: MySQL 8.0+
- **认证**: JWT (golang-jwt/jwt v5)
- **文档**: Swagger (swaggo)
- **Go版本**: 1.25+

---

**开发团队**: Tower-Go Team  
**最后更新**: 2025-11-02
