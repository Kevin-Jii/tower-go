# Tower-Go 多门店供应链管理系统

[![Go Version](https://img.shields.io/badge/Go-%E2%89%A51.20-blue)](https://golang.org/)
[![Gin](https://img.shields.io/badge/Framework-Gin-brightgreen)](https://github.com/gin-gonic/gin)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

基于 Go 语言开发的企业级多门店供应链管理系统，支持供应商管理、采购订单、门店管理和完整的 RBAC 权限控制。

## 🚀 核心功能

- **多门店管理** - 支持多门店独立管理和数据隔离
- **供应商管理** - 供应商信息、商品分类、商品管理
- **采购管理** - 采购订单创建、审核、跟踪
- **用户权限** - 基于角色的访问控制（RBAC），支持门店级权限定制
- **钉钉集成** - 机器人消息推送、Stream API 支持
- **安全认证** - JWT Token 认证，bcrypt 密码加密

## 📁 项目结构

```
tower-go/
├── cmd/main.go              # 程序入口
├── bootstrap/               # 应用初始化
├── config/                  # 配置管理
├── controller/              # 控制器层
│   ├── user.go             # 用户管理
│   ├── store.go            # 门店管理
│   ├── supplier.go         # 供应商管理
│   ├── supplier_product.go # 供应商商品
│   ├── store_supplier.go   # 门店供应商关联
│   ├── purchase_order.go   # 采购订单
│   ├── menu.go             # 菜单权限
│   ├── role.go             # 角色管理
│   └── dingtalk_bot.go     # 钉钉机器人
├── service/                 # 服务层（业务逻辑）
├── module/                  # 数据访问层（DAO）
├── model/                   # 数据模型
├── middleware/              # 中间件
├── router/                  # 路由配置
│   └── api/                # API 路由模块
├── utils/                   # 工具函数
└── migrations/              # 数据库迁移脚本
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

- **users** - 用户表（6位工号自动生成）
- **stores** - 门店表（JW+4位编码）
- **roles** - 角色表
- **menus** - 菜单表
- **role_menus** - 角色菜单关联
- **store_role_menus** - 门店角色菜单（门店级权限定制）
- **suppliers** - 供应商表
- **supplier_categories** - 供应商商品分类
- **supplier_products** - 供应商商品
- **store_supplier_products** - 门店供应商商品关联
- **purchase_orders** - 采购订单
- **purchase_order_items** - 采购订单明细
- **dingtalk_bots** - 钉钉机器人配置

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

**最后更新**: 2025-12-02
