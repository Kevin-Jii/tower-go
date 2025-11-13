# Tower-Go 多门店报菜管理系统

[![Go Version](https://img.shields.io/badge/Go-%E2%89%A51.20-blue)](https://golang.org/)
[![Gin](https://img.shields.io/badge/Framework-Gin-brightgreen)](https://github.com/gin-gonic/gin)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

一个基于 Go 语言开发的企业级多门店报菜管理系统，支持完整的菜品管理、报菜记录、数据统计和权限控制功能。

## 🚀 项目特性

### 核心功能
- ✨ **多门店管理** - 支持多门店的独立管理和数据隔离
- 🍽️ **菜品管理** - 菜品分类、新品添加、编辑和删除
- 📊 **报菜记录** - 每日报菜、历史记录查询、数据统计
- 👥 **用户权限** - 基于角色的访问控制（RBAC）
- 🔐 **安全认证** - JWT Token 认证，数据加密存储
- 📱 **钉钉集成** - 报菜后自动发送钉钉通知

### 技术特性
- 🏗️ **三层架构** - Controller-Service-Module 分层设计
- 📚 **完整 API 文档** - 自动生成 Swagger 文档
- 📝 **结构化日志** - 基于 Zap 的高性能日志系统
- 🗄️ **数据库 ORM** - GORM 框架，支持复杂查询
- ⚡ **高性能缓存** - Redis 集成
- 🎯 **事件驱动** - 内置事件总线，支持异步任务

## 📁 项目结构

```
tower-go/
├── cmd/                      # 程序入口
│   └── main.go              # 主程序入口
├── bootstrap/               # 应用初始化
├── config/                  # 配置管理
│   └── config.go
├── controller/              # 控制器层
│   ├── user.go             # 用户管理
│   ├── store.go            # 门店管理
│   ├── dish.go             # 菜品管理
│   ├── menu_report.go      # 报菜记录
│   ├── category.go         # 菜品分类
│   ├── menu.go             # 菜单权限
│   ├── dingtalk_bot.go     # 钉钉机器人
│   └── report_bot.go       # 报菜机器人
├── service/                 # 服务层（业务逻辑）
├── module/                  # 数据访问层（DAO）
├── model/                   # 数据模型
├── middleware/              # 中间件
│   └── auth.go             # 认证中间件
├── pkg/                     # 可复用包
│   ├── search/            # 查询优化器
│   └── utils/             # 通用工具
└── utils/                   # 工具函数
    ├── events/            # 事件总线
    ├── http/              # HTTP 工具
    ├── logging/           # 日志系统
    └── cache/             # 缓存封装
```

## 🛠️ 技术栈

| 技术 | 版本 | 说明 |
|------|------|------|
| Go | 1.20+ | 编程语言 |
| Gin | v1.11.0 | Web 框架 |
| GORM | v1.31.0 | ORM 框架 |
| MySQL | 8.0+ | 数据库 |
| Redis | - | 缓存 |
| JWT | golang-jwt/jwt/v5 | 认证 |
| Zap | - | 日志库 |
| Swagger | swaggo | API 文档 |

完整依赖请查看 [go.mod](go.mod)

## 🎯 快速开始

### 前置条件

- Go 1.20 或更高版本
- MySQL 8.0+
- Redis
- Git

### 安装步骤

1. **克隆项目**

```bash
git clone https://github.com/your-org/tower-go.git
cd tower-go
```

2. **安装依赖**

```bash
go mod tidy
```

3. **配置环境变量**

复制配置文件模板：

```bash
copy .env.example .env
```

编辑 `.env` 文件，填写数据库和 Redis 配置：

```env
# 应用配置
APP_NAME=tower-go
SERVER_PORT=10024

# MySQL 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=tower_go

# Redis 配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
```

4. **创建数据库**

```sql
CREATE DATABASE tower_go CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

5. **数据库迁移**

系统首次启动时会自动执行数据库迁移

6. **生成 Swagger 文档**

```bash
swag init -g cmd/main.go
```

7. **启动应用**

```bash
# 使用 make
go run cmd/main.go

# 或使用 make
make run
```

8. **访问应用**

- API 地址: `http://localhost:10024`
- Swagger 文档: `http://localhost:10024/api/v1/swagger/index.html`

## 📱 API 概览

### 认证接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/auth/register` | 用户注册 |
| POST | `/api/v1/auth/login` | 用户登录 |

### 用户管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/users/profile` | 获取个人信息 |
| PUT | `/api/v1/users/profile` | 更新个人信息 |
| GET | `/api/v1/users` | 用户列表 |
| POST | `/api/v1/users` | 创建用户 |
| GET | `/api/v1/users/:id` | 用户详情 |
| PUT | `/api/v1/users/:id` | 更新用户 |
| DELETE | `/api/v1/users/:id` | 删除用户 |

### 门店管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/stores` | 门店列表 |
| POST | `/api/v1/stores` | 创建门店 |
| GET | `/api/v1/stores/:id` | 门店详情 |
| PUT | `/api/v1/stores/:id` | 更新门店 |
| DELETE | `/api/v1/stores/:id` | 删除门店 |

### 菜品管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/dishes` | 菜品列表 |
| POST | `/api/v1/dishes` | 创建菜品 |
| GET | `/api/v1/dishes/:id` | 菜品详情 |
| PUT | `/api/v1/dishes/:id` | 更新菜品 |
| DELETE | `/api/v1/dishes/:id` | 删除菜品 |

### 菜品分类

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/dish-categories` | 分类列表 |
| POST | `/api/v1/dish-categories` | 创建分类 |
| PUT | `/api/v1/dish-categories/:id` | 更新分类 |
| DELETE | `/api/v1/dish-categories/:id` | 删除分类 |

### 报菜记录

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/menu-reports` | 报菜列表 |
| POST | `/api/v1/menu-reports` | 创建报菜 |
| GET | `/api/v1/menu-reports/statistics` | 统计数据 |
| GET | `/api/v1/menu-reports/:id` | 记录详情 |
| PUT | `/api/v1/menu-reports/:id` | 更新记录 |
| DELETE | `/api/v1/menu-reports/:id` | 删除记录 |

### 菜单权限

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/menus` | 菜单列表 |
| POST | `/api/v1/menus` | 创建菜单 |
| GET | `/api/v1/menus/tree` | 菜单树 |
| GET | `/api/v1/menus/user-menus` | 用户菜单 |
| POST | `/api/v1/menus/assign-role` | 分配角色菜单 |
| POST | `/api/v1/menus/assign-store-role` | 分配门店菜单 |

### 钉钉集成

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/dingtalk/robots` | 机器人列表 |
| POST | `/api/v1/dingtalk/robots` | 创建机器人 |
| POST | `/api/v1/dingtalk/robots/:id/test` | 测试机器人 |

完整 API 文档请访问 Swagger: `http://localhost:10024/api/v1/swagger/index.html`

## 🔐 认证机制

### JWT Token 认证

1. **登录获取 Token**

```bash
curl -X POST http://localhost:10024/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "13800138000",
    "password": "your_password"
  }'
```

响应示例：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

2. **使用 Token 访问受保护接口**

```bash
curl -X GET http://localhost:10024/api/v1/users/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Token 包含的信息

Token Payload 包含：
- `userID`: 用户 ID
- `username`: 用户名
- `storeID`: 门店 ID（用于数据隔离）
- `roleCode`: 角色代码
- `roleID`: 角色 ID

### 角色和权限

| 角色代码 | 角色名称 | 权限说明 |
|----------|----------|----------|
| `admin` | 总部管理员 | 管理所有门店、查看汇总数据 |
| `store_admin` | 门店管理员 | 管理本门店数据 |
| `staff` | 普通员工 | 查看和操作本门店数据 |

### 中间件

```go
// 基础认证中间件
r.Use(middleware.AuthMiddleware())

// 门店级认证（别名）
r.Use(middleware.StoreAuthMiddleware())
```

### 获取当前用户信息

```go
// 获取用户 ID
userID := middleware.GetUserID(ctx)

// 获取门店 ID
storeID := middleware.GetStoreID(ctx)

// 获取用户名
username := middleware.GetUsername(ctx)

// 判断是否为管理员
isAdmin := middleware.IsAdmin(ctx)
```

## 🏗️ 架构设计

### 三层架构

本项目采用经典的三层架构设计：

```
┌─────────────────────────────────────┐
│        Controller 层                │
│  (HTTP 请求处理、参数解析)            │
└─────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────┐
│         Service 层                  │
│  (业务逻辑、数据校验、事务管理)      │
└─────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────┐
│         Module 层 (DAO)             │
│  (数据库交互、查询构建)              │
└─────────────────────────────────────┘
```

### 分层职责

#### Controller 层
- HTTP 请求处理
- 参数解析和验证
- 调用 Service 层
- 响应封装

#### Service 层
- 业务逻辑实现
- 数据校验
- 事务管理
- 调用 Module 层

#### Module 层
- 数据库交互
- GORM 查询构建
- 数据模型映射

### 数据流向

```
HTTP 请求
    ↓
Router (Gin)
    ↓
Controller（解析参数）
    ↓
Service（业务逻辑 + 事务）
    ↓
Module（数据库操作）
    ↓
MySQL/Redis
```

### 数据隔离机制

```go
// Service 层自动获取当前用户的 storeID
storeID := middleware.GetStoreID(ctx)

// Module 层自动添加 store_id 过滤条件
func (m *DishModule) GetDishesByStoreID(storeID uint, params map[string]interface{}) ([]*model.Dish, error) {
    query := m.db.Where("store_id = ?", storeID)
    // ... 其他条件
}
```

## 📊 数据模型

### 核心数据表

#### 1. 用户表 (users)

```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    employee_no VARCHAR(6) UNIQUE NOT NULL COMMENT '6位工号',
    username VARCHAR(50) NOT NULL COMMENT '姓名',
    phone VARCHAR(11) UNIQUE NOT NULL COMMENT '手机号',
    password VARCHAR(255) NOT NULL COMMENT '加密密码',
    store_id BIGINT NOT NULL COMMENT '所属门店ID',
    role_id BIGINT NOT NULL COMMENT '角色ID',
    status TINYINT DEFAULT 1 COMMENT '状态：1正常 0禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_phone (phone),
    INDEX idx_store_id (store_id)
);
```

#### 2. 门店表 (stores)

```sql
CREATE TABLE stores (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    code VARCHAR(50) UNIQUE NOT NULL COMMENT '门店编码',
    name VARCHAR(100) NOT NULL COMMENT '门店名称',
    address VARCHAR(255) COMMENT '门店地址',
    status TINYINT DEFAULT 1 COMMENT '状态：1正常 0停业',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (status)
);
```

#### 3. 菜品表 (dishes)

```sql
CREATE TABLE dishes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT NOT NULL COMMENT '门店ID',
    category_id BIGINT NOT NULL COMMENT '分类ID',
    name VARCHAR(100) NOT NULL COMMENT '菜品名称',
    unit VARCHAR(10) COMMENT '单位',
    is_active TINYINT DEFAULT 1 COMMENT '是否启用：1是 0否',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_store_category_name (store_id, category_id, name),
    INDEX idx_store_id (store_id),
    INDEX idx_category_id (category_id)
);
```

#### 4. 菜品分类表 (dish_categories)

```sql
CREATE TABLE dish_categories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT NOT NULL COMMENT '门店ID',
    name VARCHAR(50) NOT NULL COMMENT '分类名称',
    sort_order INT DEFAULT 0 COMMENT '排序',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_store_id (store_id)
);
```

#### 5. 报菜记录表 (menu_reports)

```sql
CREATE TABLE menu_reports (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT NOT NULL COMMENT '门店ID',
    dish_id BIGINT NOT NULL COMMENT '菜品ID',
    report_date DATE NOT NULL COMMENT '报菜日期',
    quantity DECIMAL(10,2) NOT NULL COMMENT '备餐数量',
    reporter_id BIGINT NOT NULL COMMENT '报菜人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_store_dish_date (store_id, dish_id, report_date),
    INDEX idx_store_id (store_id),
    INDEX idx_report_date (report_date),
    INDEX idx_reporter_id (reporter_id)
);
```

#### 6. 角色表 (roles)

```sql
CREATE TABLE roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL COMMENT '角色名称',
    code VARCHAR(20) UNIQUE NOT NULL COMMENT '角色编码',
    description VARCHAR(255) COMMENT '角色描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### 数据关系

```
users -- 所属 --> stores
users -- 拥有 --> roles

stores -- 拥有 --> dishes
dish_categories -- 包含 --> dishes

stores -- 产生 --> menu_reports
dishes -- 被报 --> menu_reports
users -- 报 --> menu_reports

stores -- 拥有 --> menus
roles -- 可访问 --> menus
```

## 🌐 事件驱动架构

### 事件总线

**文件**: `/utils/events/event_bus.go`

基于发布/订阅模式的事件总线：

```go
// 发布事件
events.Publish("menu_report.created", map[string]interface{}{
    "store_id": storeID,
    "dish_name": dishName,
    "quantity": quantity,
    "reporter_name": reporterName,
})

// 订阅事件
events.Subscribe("menu_report.created", func(data map[string]interface{}) {
    // 处理事件
})
```

### 内置事件

#### 1. menu_report.created

**触发时机**: 创建报菜记录后

**事件数据**:
```json
{
  "store_id": 1,
  "store_name": "南山店",
  "dish_id": 10,
  "dish_name": "宫保鸡丁",
  "quantity": 50,
  "reporter_id": 5,
  "reporter_name": "张三",
  "report_date": "2025-11-11"
}
```

**应用场景**: 自动发送钉钉通知

## 📞 钉钉集成

### 功能特性

**文件**: `/service/dingtalk.go`

- 发送文本消息
- 发送 Markdown 格式消息
- 支持 Stream API
- 机器人管理

### 配置

在 `.env` 文件中配置钉钉凭证：

```env
# 钉钉 Stream API 配置
DINGTALK_APP_ID=your_app_id
DINGTALK_APP_SECRET=your_app_secret
```

### 使用示例

```go
// 发送文本消息
dingtalkService.SendTextMessage(ctx, "robot_id", "报菜成功：宫保鸡丁 50份")

// 发送 Markdown 消息
dingtalkService.SendMarkdownMessage(ctx, "robot_id", "报菜通知", markdownContent)
```

### 报菜通知示例

当门店完成报菜后，自动发送钉钉通知：

```
📌 报菜通知 - 南山店
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📅 报菜日期: 2025-11-11
👤 报菜人: 张三
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🍽️ 今日备餐清单:
  • 宫保鸡丁 - 50 份
  • 麻婆豆腐 - 30 份
  • 回锅肉 - 40 份
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ 报菜已完成，请厨师长安排备餐
```

## 📝 日志系统

### 日志配置

**文件**: `/utils/logging/logger.go`

基于 Zap 的高性能日志系统，支持日志轮转：

```go
// 日志等级
debug < info < warn < error < fatal

// 日志输出格式
{
  "level": "info",
  "ts": "2025-11-11T10:00:00.000+0800",
  "caller": "controller/user.go:50",
  "msg": "用户登录成功",
  "user_id": 1,
  "username": "admin"
}
```

### 日志文件

- **日志文件**: `logs/app.log`
- **日志轮转**: 按天分割，保留 7 天
- **日志级别**: 生产环境默认 info，开发环境可调整为 debug

### 使用方法

```go
// 基础日志
log.Info("用户登录成功")
log.Error("数据库连接失败", zap.Error(err))

// 结构化日志
log.Info("创建菜品",
    zap.Uint("dish_id", dish.ID),
    zap.String("dish_name", dish.Name),
    zap.Uint("store_id", dish.StoreID),
)

// HTTP 请求日志（自动记录）
log.HTTPRequest(ctx, "POST /api/v1/dishes", zap.Int("status_code", 200))
```

### 日志分类

- **业务日志**: 记录业务操作
- **HTTP 日志**: 记录 HTTP 请求
- **数据库日志**: 记录 SQL 查询
- **认证日志**: 记录登录、权限验证
- **错误日志**: 记录异常和错误

## 💻 开发指南

### 项目初始化

```bash
# 克隆项目
git clone https://github.com/your-org/tower-go.git
cd tower-go

# 安装依赖
go mod tidy

# 安装开发工具
go install github.com/swaggo/swag/cmd/swag@latest

# 生成 Swagger 文档
swag init -g cmd/main.go
```

### 目录规范

```
controller/  # 仅处理 HTTP，调用 Service
  ├── user.go
  └── ...

service/     # 业务逻辑，调用 Module
  ├── user.go
  └── ...

module/      # 数据库操作，返回 Model
  ├── user_module.go
  └── ...

model/       # 仅定义数据结构
  └── user.go
```

### 添加新功能

以添加"库存管理"功能为例：

**1. 创建数据模型** (`model/inventory.go`)

```go
type Inventory struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    StoreID   uint      `json:"store_id" gorm:"index"`
    DishID    uint      `json:"dish_id"`
    Quantity  float64   `json:"quantity"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

**2. 创建 Module** (`module/inventory_module.go`)

```go
type InventoryModule struct {
    db *gorm.DB
}

func NewInventoryModule(db *gorm.DB) *InventoryModule {
    return &InventoryModule{db: db}
}

func (m *InventoryModule) Create(inventory *model.Inventory) error {
    return m.db.Create(inventory).Error
}

func (m *InventoryModule) GetByStoreID(storeID uint) ([]*model.Inventory, error) {
    var inventories []*model.Inventory
    err := m.db.Where("store_id = ?", storeID).Find(&inventories).Error
    return inventories, err
}
```

**3. 创建 Service** (`service/inventory.go`)

```go
type InventoryService struct {
    inventoryModule *module.InventoryModule
    dishModule      *module.DishModule
}

func NewInventoryService(
    inventoryModule *module.InventoryModule,
    dishModule *module.DishModule,
) *InventoryService {
    return &InventoryService{
        inventoryModule: inventoryModule,
        dishModule:      dishModule,
    }
}

func (s *InventoryService) CreateInventory(ctx *gin.Context, params map[string]interface{}) error {
    storeID := middleware.GetStoreID(ctx)

    inventory := &model.Inventory{
        StoreID:  storeID,
        DishID:   params["dish_id"].(uint),
        Quantity: params["quantity"].(float64),
    }

    return s.inventoryModule.Create(inventory)
}
```

**4. 创建 Controller** (`controller/inventory.go`)

```go
type InventoryController struct {
    inventoryService *service.InventoryService
}

func NewInventoryController(inventoryService *service.InventoryService) *InventoryController {
    return &InventoryController{
        inventoryService: inventoryService,
    }
}

// CreateInventory godoc
// @Summary 创建库存记录
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]interface{} true "库存信息"
// @Success 200 {object} utils.Response
// @Router /api/v1/inventories [post]
func (c *InventoryController) Create(ctx *gin.Context) {
    var params map[string]interface{}
    if err := ctx.ShouldBindJSON(&params); err != nil {
        utils.Error(ctx, "参数错误")
        return
    }

    if err := c.inventoryService.CreateInventory(ctx, params); err != nil {
        log.Error("创建库存失败", zap.Error(err))
        utils.Error(ctx, "创建失败")
        return
    }

    utils.Success(ctx, "创建成功", nil)
}
```

**5. 注册路由** (`bootstrap/router.go`)

```go
// 创建实例
inventoryModule := module.NewInventoryModule(db)
inventoryService := service.NewInventoryService(inventoryModule, dishModule)
inventoryController := controller.NewInventoryController(inventoryService)

// 注册路由
group := r.Group("/api/v1")
group.Use(middleware.AuthMiddleware())
{
    group.POST("/inventories", inventoryController.Create)
    group.GET("/inventories", inventoryController.List)
}
```

**6. 生成 Swagger 文档**

```bash
swag init -g cmd/main.go
```

### 代码规范

#### 1. 错误处理

```go
// ✅ 推荐：记录错误日志，返回友好提示
if err != nil {
    log.Error("数据库查询失败", zap.Error(err))
    return utils.Error(ctx, "获取数据失败")
}

// ❌ 不推荐：直接返回原始错误
if err != nil {
    return utils.Error(ctx, err.Error())
}
```

#### 2. 日志记录

```go
// ✅ 推荐：结构化日志
log.Info("用户登录成功",
    zap.Uint("user_id", user.ID),
    zap.String("username", user.Username),
)

// ❌ 不推荐：拼接字符串
log.Info(fmt.Sprintf("用户 %d 登录成功", user.ID))
```

#### 3. 参数验证

```go
// ✅ 推荐：在 Service 层验证
func (s *UserService) CreateUser(ctx *gin.Context, params map[string]interface{}) error {
    phone := params["phone"].(string)
    if len(phone) != 11 {
        return errors.New("手机号格式错误")
    }
    // ...
}

// ❌ 不推荐：在 Controller 层验证业务逻辑
```

#### 4. 数据隔离

```go
// ✅ 推荐：在 Service 层统一获取 storeID
func (s *DishService) GetDishes(ctx *gin.Context) ([]*model.Dish, error) {
    storeID := middleware.GetStoreID(ctx)
    return s.dishModule.GetByStoreID(storeID)
}

// ❌ 不推荐：在 Controller 层处理
```

### 数据库迁移

本项目使用 GORM 的 AutoMigrate 功能，启动时自动迁移数据库结构：

```go
// bootstrap/db.go
db.AutoMigrate(
    &model.User{},
    &model.Store{},
    &model.Dish{},
    &model.DishCategory{},
    &model.MenuReport{},
    &model.Role{},
    &model.Menu{},
    &model.RoleMenu{},
    // ...
)
```

如果需要手动管理迁移，可以使用 GORM 的 Migrator：

```go
// 创建表
migrator := db.Migrator()
migrator.CreateTable(&model.User{})

// 添加字段
db.Migrator().AddColumn(&model.User{}, "Avatar")

// 创建索引
db.Migrator().CreateIndex(&model.User{}, "idx_phone")
```

### 测试

#### 单元测试

```go
// service/user_test.go
func TestUserService_Login(t *testing.T) {
    // 准备测试数据
    mockUserModule := &mockUserModule{}
    userService := service.NewUserService(mockUserModule, nil, nil)

    // 执行测试
    token, err := userService.Login(ctx, "13800138000", "password")

    // 验证结果
    assert.NoError(t, err)
    assert.NotEmpty(t, token)
}
```

#### API 测试

使用工具：
- [Postman](https://www.postman.com/)
- [HTTPie](https://httpie.io/)
- [curl](https://curl.se/)

测试示例：

```bash
# 登录
token=$(curl -s -X POST http://localhost:10024/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800138000","password":"password"}' |
  jq -r '.data.token')

# 使用 token 访问受保护接口
curl -X GET http://localhost:10024/api/v1/users/profile \
  -H "Authorization: Bearer $token"
```

## 🐳 Docker 部署

### Dockerfile

```dockerfile
# 构建阶段
FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o tower-go cmd/main.go

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/tower-go .
COPY --from=builder /app/.env .

EXPOSE 10024
CMD ["./tower-go"]
```

### docker-compose.yml

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "10024:10024"
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=root
      - DB_PASSWORD=example
      - DB_NAME=tower_go
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    depends_on:
      - mysql
      - redis
    volumes:
      - ./logs:/app/logs

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: example
      MYSQL_DATABASE: tower_go
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  mysql_data:
  redis_data:
```

### 部署步骤

1. **构建镜像**

```bash
docker build -t tower-go:latest .
```

2. **启动服务**

```bash
docker-compose up -d
```

3. **查看日志**

```bash
docker-compose logs -f app
```

4. **停止服务**

```bash
docker-compose down
```

## 🔧 配置说明

### 配置文件

项目使用 `.env` 文件管理配置，支持以下配置项：

#### 应用配置

```env
# 应用名称
APP_NAME=tower-go

# 服务器端口
SERVER_PORT=10024

# 运行模式: debug, release
gin_mode=release
```

#### MySQL 数据库

```env
# 数据库主机
DB_HOST=localhost

# 数据库端口
DB_PORT=3306

# 数据库用户名
DB_USER=root

# 数据库密码
DB_PASSWORD=your_password

# 数据库名称
DB_NAME=tower_go

# 连接池最大连接数
DB_MAX_OPEN_CONNS=100

# 连接池最大空闲连接数
DB_MAX_IDLE_CONNS=10

# 连接最大存活时间（秒）
DB_CONN_MAX_LIFETIME=3600
```

#### Redis 缓存

```env
# Redis 主机
REDIS_HOST=localhost

# Redis 端口
REDIS_PORT=6379

# Redis 密码
REDIS_PASSWORD=

# Redis 数据库
REDIS_DB=0

# 连接池大小
REDIS_POOL_SIZE=10
```

#### JWT 配置

```env
# JWT 签名密钥（生产环境必须修改！）
JWT_SECRET=your_jwt_secret_key_here_change_in_production

# Token 过期时间（小时）
JWT_EXPIRE_HOURS=24
```

#### 钉钉配置

```env
# 钉钉 Stream API 配置
DINGTALK_APP_ID=your_app_id
DINGTALK_APP_SECRET=your_app_secret

# 钉钉机器人 Webhook（可选）
DINGTALK_WEBHOOK=https://oapi.dingtalk.com/robot/send?access_token=xxx
```

#### 日志配置

```env
# 日志文件路径
LOG_FILE=logs/app.log

# 日志等级: debug, info, warn, error
LOG_LEVEL=info

# 日志保留天数
LOG_MAX_AGE=7

# 是否启用日志轮转
LOG_ROTATION=true
```

#### 安全配置

```env
# 是否启用 HTTPS
ENABLE_TLS=false

# TLS 证书路径
TLS_CERT_FILE=
TLS_KEY_FILE=

# CORS 配置
CORS_ALLOW_ORIGINS=http://localhost:3000,http://localhost:8080
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOW_HEADERS=*
```

### 环境变量优先级

配置加载优先级（从高到低）：

1. **环境变量** (优先级最高)
2. **`.env` 文件** (项目根目录)
3. **`config/.env` 文件**
4. **默认值** (优先级最低)

### 生产环境配置建议

```env
# 生产环境必须修改的配置

# 1. JWT 密钥（必须设置复杂密钥）
JWT_SECRET=生成一个32位以上的随机字符串

# 2. 数据库使用内网地址
DB_HOST=10.0.1.10
DB_PASSWORD=强密码

# 3. 启用 HTTPS
ENABLE_TLS=true
TLS_CERT_FILE=/path/to/cert.pem
TLS_KEY_FILE=/path/to/key.pem

# 4. 日志级别调整为 warn 或 error
LOG_LEVEL=warn

# 5. Gin 运行模式
gin_mode=release

# 6. 限制 CORS 域名
CORS_ALLOW_ORIGINS=https://your-domain.com
```

## 📊 性能优化

### 1. 数据库优化

#### 索引优化

```sql
-- users 表索引
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_store_id ON users(store_id);
CREATE INDEX idx_users_role_id ON users(role_id);

-- dishes 表索引
CREATE INDEX idx_dishes_store_id ON dishes(store_id);
CREATE INDEX idx_dishes_category_id ON dishes(category_id);
CREATE INDEX idx_dishes_name ON dishes(name);

-- menu_reports 表索引
CREATE INDEX idx_menu_reports_store_id ON menu_reports(store_id);
CREATE INDEX idx_menu_reports_report_date ON menu_reports(report_date);
CREATE INDEX idx_menu_reports_dish_id ON menu_reports(dish_id);

-- 联合索引
CREATE INDEX idx_menu_reports_store_date ON menu_reports(store_id, report_date);
```

#### 查询优化

```go
// ✅ 推荐：使用预加载减少 N+1 查询
dishes, err := s.dishModule.GetDishesWithCategory(storeID)

// ❌ 不推荐：循环查询
for _, dish := range dishes {
    category := s.categoryModule.GetByID(dish.CategoryID)
    dish.Category = category
}
```

### 2. Redis 缓存

#### 缓存策略

```go
// 1. 菜品列表缓存（1小时）
func (s *DishService) GetCacheDishes(storeID uint) ([]*model.Dish, error) {
    cacheKey := fmt.Sprintf("dishes:store:%d", storeID)

    // 先查缓存
    if data, err := cache.Get(cacheKey); err == nil {
        dishes := []*model.Dish{}
        json.Unmarshal(data, &dishes)
        return dishes, nil
    }

    // 缓存未命中，查询数据库
    dishes, err := s.dishModule.GetByStoreID(storeID)
    if err != nil {
        return nil, err
    }

    // 写入缓存（1小时过期）
    if data, err := json.Marshal(dishes); err == nil {
        cache.Set(cacheKey, data, time.Hour)
    }

    return dishes, nil
}
```

#### 缓存更新策略

```go
// 创建菜品后清除缓存
func (s *DishService) CreateDish(ctx *gin.Context, dish *model.Dish) error {
    if err := s.dishModule.Create(dish); err != nil {
        return err
    }

    // 清除缓存
    cacheKey := fmt.Sprintf("dishes:store:%d", dish.StoreID)
    cache.Delete(cacheKey)

    return nil
}
```

### 3. 连接池配置

#### MySQL 连接池

```go
// bootstrap/db.go
sqlDB, err := db.DB()
if err != nil {
    log.Fatal("获取数据库实例失败", zap.Error(err))
}

// 设置连接池参数
sqlDB.SetMaxOpenConns(100)        // 最大连接数
sqlDB.SetMaxIdleConns(10)          // 最大空闲连接数
sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期
```

#### Redis 连接池

```go
// bootstrap/redis.go
redisClient := redis.NewClient(&redis.Options{
    Addr:        "localhost:6379",
    Password:    "",
    DB:          0,
    PoolSize:    10,               // 连接池大小
    MinIdleConns: 5,               // 最小空闲连接
})
```

### 4. 接口优化

#### 分页查询

```go
// ✅ 推荐：强制分页
type PageParams struct {
    Page     int `form:"page" binding:"required,min=1"`
    PageSize int `form:"page_size" binding:"required,min=1,max=100"` // 限制最大100条
}

func (s *DishService) GetDishes(storeID uint, page, pageSize int) ([]*model.Dish, int64, error) {
    var total int64

    // 先查询总数
    s.db.Model(&model.Dish{}).Where("store_id = ?", storeID).Count(&total)

    // 再查询分页数据
    var dishes []*model.Dish
    offset := (page - 1) * pageSize
    err := s.db.Where("store_id = ?", storeID).
        Limit(pageSize).
        Offset(offset).
        Find(&dishes).Error

    return dishes, total, err
}
```

#### 批量操作

```go
// ✅ 推荐：批量创建
func (m *MenuReportModule) BatchCreate(reports []*model.MenuReport) error {
    return m.db.CreateInBatches(reports, 100).Error
}

// ❌ 不推荐：循环创建
for _, report := range reports {
    m.db.Create(report)
}
```

## 🔍 监控与调试

### 1. 日志监控

```bash
# 实时查看日志
tail -f logs/app.log

# 过滤错误日志
grep "ERROR" logs/app.log

# 统计某个接口的访问次数
grep "POST /api/v1/menu-reports" logs/app.log | wc -l
```

### 2. 健康检查

```bash
# HTTP 健康检查
curl -f http://localhost:10024/health || echo "服务异常"
```

### 3. 性能分析

```go
// 启用 pprof（开发环境）
import _ "net/http/pprof"

go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

访问：
- `http://localhost:6060/debug/pprof/` - pprof 首页
- `http://localhost:6060/debug/pprof/goroutine` - goroutine 分析
- `http://localhost:6060/debug/pprof/heap` - 内存分析

### 4. 慢查询日志

```sql
-- MySQL 开启慢查询日志
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 1; -- 超过 1 秒记录
SET GLOBAL slow_query_log_file = '/var/log/mysql/slow.log';
```

## 🔐 安全加固

### 1. 密码安全

```go
// ✅ 推荐：使用 bcrypt 加密（已在项目中实现）
import "golang.org/x/crypto/bcrypt"

passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// ✅ 验证密码
err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
```

### 2. SQL 注入防护

```go
// ✅ 推荐：使用 GORM 预编译
db.Where("phone = ?", phone).First(&user)

// ❌ 不推荐：字符串拼接
db.Raw(fmt.Sprintf("SELECT * FROM users WHERE phone = '%s'", phone))
```

### 3. XSS 防护

```go
// ✅ 推荐：使用 html.EscapeString
escaped := html.EscapeString(userInput)

// ✅ 推荐：输出 JSON 时自动转义
ctx.JSON(200, gin.H{"data": escaped})
```

### 4. CORS 配置

```go
// ✅ 生产环境严格限制域名
config := cors.Config{
    AllowOrigins:     []string{"https://your-domain.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}
```

### 5. 限流

```go
// 使用 gin-limit 中间件
import "github.com/ulule/limiter/v3"

rateLimiter := mgin.NewMiddleware(limiter.New(
    limiter.Rate{
        Period: 1 * time.Minute,
        Limit:  100, // 每分钟最多 100 次请求
    },
))

r.Use(rateLimiter)
```

### 6. SSL/TLS

```go
// ✅ 生产环境启用 HTTPS
r.RunTLS(":10024", "cert.pem", "key.pem")
```

## 📚 完整 API 文档

参见 [API_GUIDE.md](./API_GUIDE.md)

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

### 提交规范

- **Bug 修复**: `fix: 修复 XX 问题`
- **新功能**: `feat: 添加 XX 功能`
- **文档更新**: `docs: 更新 XX 文档`
- **代码重构**: `refactor: 重构 XX 模块`

### 开发流程

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📝 许可证

本项目基于 [MIT 许可证](LICENSE) 开源

## 📞 联系方式

- 项目地址: [https://github.com/your-org/tower-go](https://github.com/your-org/tower-go)
- 问题反馈: [Issues](https://github.com/your-org/tower-go/issues)
- Email: your-email@example.com

## 🙏 致谢

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [Zap](https://github.com/uber-go/zap)
- [JWT-Go](https://github.com/golang-jwt/jwt)
- [Swagger](https://swagger.io/)

---

**文档版本**: v1.0.0
**最后更新**: 2025-11-11
**项目版本**: 基于 commit ea1443f
