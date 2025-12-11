# Tower-Go 多门店供应链管理系统

[![Go Version](https://img.shields.io/badge/Go-%E2%89%A51.20-blue)](https://golang.org/)
[![Gin](https://img.shields.io/badge/Framework-Gin-brightgreen)](https://github.com/gin-gonic/gin)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

基于 Go 语言开发的企业级多门店供应链管理系统，支持供应商管理、采购订单、门店管理和完整的 RBAC 权限控制。

## 🚀 核心功能

### 📦 供应链管理
- **多门店管理** - 支持多门店独立管理和数据隔离
- **供应商管理** - 供应商信息、商品分类、商品管理
- **采购管理** - 采购订单创建、审核、跟踪
- **库存管理** - 实时库存监控、库存调整、库存流水记录

### 💰 门店记账
- **电子记账** - 支持一单多商品，灵活的明细管理
- **图片通知** - 自动生成美观的记账通知图片（电子回单风格）
- **渠道管理** - 支持多种销售渠道（微信、堂食、外卖等）
- **标签系统** - 自定义标签分类管理

### 📊 统计分析
- **数据仪表板** - 实时查看库存和销售数据（当日/本周/当月/当季/当年）
- **门店统计** - 多维度门店数据分析
- **趋势分析** - 销售趋势和库存变化追踪

### 🔐 权限与安全
- **RBAC 权限** - 基于角色的访问控制（超级管理员/总部管理员/门店管理员/普通员工）
- **门店级权限** - 支持门店级别的菜单权限定制
- **JWT 认证** - 安全的 Token 认证机制
- **密码加密** - bcrypt 密码加密存储

### 📢 消息通知
- **钉钉集成** - 机器人消息推送、Stream API 实时通信
- **消息模板** - 可配置的钉钉消息模板系统
- **自动通知** - 记账、库存预警等自动消息推送
- **图片消息** - 支持发送图文混合的通知消息

### 📁 文件服务
- **对象存储** - 集成 RustFS/MinIO 对象存储服务
- **图库管理** - 图片上传、分类、管理
- **文件管理** - 支持多种文件类型的存储和管理

### 🛠️ 系统管理
- **字典管理** - 通用的数据字典管理系统
- **用户管理** - 完整的用户生命周期管理
- **角色管理** - 角色创建、权限分配
- **菜单管理** - 动态菜单配置

## 📁 项目结构

```
tower-go/
├── cmd/main.go                 # 程序入口
├── bootstrap/                  # 应用初始化
├── config/                     # 配置管理
├── controller/                 # 控制器层
│   ├── user.go                # 用户管理
│   ├── store.go               # 门店管理
│   ├── supplier.go            # 供应商管理
│   ├── supplier_product.go    # 供应商商品
│   ├── store_supplier.go      # 门店供应商关联
│   ├── purchase_order.go      # 采购订单
│   ├── menu.go                # 菜单权限
│   ├── role.go                # 角色管理
│   ├── dingtalk_bot.go        # 钉钉机器人
│   ├── inventory.go           # 库存管理
│   ├── statistics.go          # 统计分析
│   ├── store_account.go       # 门店记账
│   ├── message_template.go    # 消息模板
│   ├── file.go                # 文件管理
│   └── gallery.go             # 图库管理
├── service/                    # 服务层（业务逻辑）
│   ├── rustfs.go              # RustFS/MinIO 文件服务
│   ├── image_generator.go     # 图片自动生成
│   └── dingtalk.go            # 钉钉集成服务
├── module/                     # 数据访问层（DAO）
├── model/                      # 数据模型
├── middleware/                 # 中间件
│   ├── auth.go                # 认证中间件
│   ├── rbac.go                # 权限中间件
│   └── logger.go              # 日志中间件
├── router/                     # 路由配置
│   ├── api/                   # API 路由模块
│   └── router.go              # 路由注册
├── utils/                      # 工具函数
├── migrations/                 # 数据库迁移脚本
└── docs/                       # Swagger 文档
    ├── docs.go
    ├── swagger.json
    └── swagger.yaml
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
| gg | fogleman/gg | 图片生成库 |
| MinIO/RustFS | - | 对象存储服务 |
| WebSocket | gorilla/websocket | 实时通信 |

## 🎯 快速开始

### 前置条件

- Go 1.20+
- MySQL 8.0+
- Redis

### 安装步骤

1. **克隆项目**
```bash
git clone https://github.com/Kevin-Jii/tower-go.git
cd tower-go
```

2. **安装依赖**
```bash
go mod tidy
```

3. **配置环境变量**
```bash
copy .env.example .env
# 编辑 .env 文件，填写数据库和 Redis 配置
```

4. **创建数据库**
```sql
CREATE DATABASE tower CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

5. **初始化种子数据**
```bash
mysql -u用户名 -p密码 数据库名 < migrations/init_seed_data.sql
```

6. **启动应用**
```bash
go run cmd/main.go
```

7. **访问应用**
- API 地址: `http://localhost:10024`
- Swagger 文档: `http://localhost:10024/api/v1/swagger/index.html`

## 🔐 认证与权限

### JWT Token 认证

```bash
# 登录获取 Token
curl -X POST http://localhost:10024/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone": "130xxxxxxxx", "password": "Admin@123456"}'

# 使用 Token 访问接口
curl -X GET http://localhost:10024/api/v1/users/profile \
  -H "Authorization: Bearer <token>"
```

### 角色权限

| 角色代码 | 角色名称 | 权限说明 |
|----------|----------|----------|
| `super_admin` | 超级管理员 | 系统最高权限，管理所有门店 |
| `admin` | 总部管理员 | 管理所有门店、查看汇总数据 |
| `store_admin` | 门店管理员 | 管理本门店数据 |
| `staff` | 普通员工 | 查看和操作本门店数据 |

### 权限位说明
```
1  = 0001 = 仅删除
2  = 0010 = 仅修改
4  = 0100 = 仅新增
8  = 1000 = 仅查看
15 = 1111 = 全部权限
```

## 📊 数据模型

### 核心数据表

#### 用户权限
- **users** - 用户表（6位工号自动生成）
- **roles** - 角色表（超级管理员/总部管理员/门店管理员/普通员工）
- **menus** - 菜单表
- **role_menus** - 角色菜单关联（权限控制）
- **store_role_menus** - 门店角色菜单（门店级权限定制）

#### 供应链与门店
- **stores** - 门店表（JW+4位编码）
- **suppliers** - 供应商表
- **supplier_categories** - 供应商商品分类
- **supplier_products** - 供应商商品
- **store_supplier_products** - 门店供应商商品关联
- **purchase_orders** - 采购订单
- **purchase_order_items** - 采购订单明细

#### 记账与库存
- **store_accounts** - 门店记账主表（支持一单多商品）
- **store_account_items** - 记账明细表
- **inventories** - 库存表
- **inventory_transactions** - 库存流水记录

#### 字典与配置
- **dict_types** - 字典类型表
- **dict_datas** - 字典数据表
- **message_templates** - 消息模板表
- **dingtalk_bots** - 钉钉机器人配置

#### 文件管理
- **galleries** - 图库管理表
- **files** - 文件管理表

## ⚙️ 配置说明

### 环境变量 (.env)

```env
# 应用配置
APP_NAME=tower-go
APP_PORT=10024

# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=your_password
DB_NAME=tower

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_ENABLED=true

# JWT配置
JWT_SECRET=your_jwt_secret_at_least_32_characters

# 钉钉配置
DINGTALK_CLIENT_ID=your_client_id
DINGTALK_CLIENT_SECRET=your_client_secret
DINGTALK_ROBOT_CODE=your_robot_code
DINGTALK_ROBOT_WEBHOOK=https://oapi.dingtalk.com/robot/send

# RustFS/MinIO 配置（可选）
RUSTFS_ENABLED=true
RUSTFS_ENDPOINT=your.rustfs.server:9000
RUSTFS_ACCESS_KEY=your_access_key
RUSTFS_SECRET_KEY=your_secret_key
RUSTFS_BUCKET=tower
RUSTFS_NOTIFY_BUCKET=tower-notify
RUSTFS_USE_SSL=false
```

## 🐳 Docker 部署

```bash
# 构建镜像
docker build -t tower-go:latest .

# 启动服务
docker-compose up -d
```

## 📝 默认账号

系统初始化后的默认账号：

| 账号 | 密码 | 角色 |
|------|------|------|
| 138 xxxx xxxx | Admin@123456 | 超级管理员 |

> ⚠️ 请在首次登录后立即修改默认密码！

## 🤝 贡献指南

提交规范：
- `fix:` Bug 修复
- `feat:` 新功能
- `docs:` 文档更新
- `refactor:` 代码重构

## 📝 许可证

[MIT License](LICENSE)

---

## 🎨 记账通知图片功能

系统支持自动生成美观的记账通知图片，以电子回单的形式发送给门店管理人员。

### 功能特点

- **自动生成** - 记账完成后自动生成图片，无需人工干预
- **现代设计** - 采用电子回单风格设计，包含：
  - 门店信息和记账编号
  - 商品明细列表（支持多商品）
  - 合计金额和笔数统计
  - 渠道、操作人和日期信息
  - "已入账"电子印章
  - 生成时间戳
- **对象存储** - 生成的图片自动上传到 RustFS/MinIO
- **钉钉推送** - 优先发送图片+文字的通知消息

### 技术实现

- 使用 [gg](https://github.com/fogleman/gg) 图形库进行图片绘制
- 支持中文字体渲染（Windows/Unix 系统自动适配）
- 高质量 PNG 输出（2倍分辨率）
- 自动适配商品数量，动态计算图片高度

---

**最后更新**: 2025-12-12
