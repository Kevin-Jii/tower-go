# Tower-Go 多门店供应链管理系统

Tower-Go 是一个面向连锁门店经营的供应链与门店运营管理系统，后端使用 Go + Gin + GORM，前端管理端使用 Vue 3 + Vite + Arco Design。项目覆盖门店、供应商、采购、库存、记账、退货、会员、B2B 供货、第三方订单、钉钉通知、云打印、图库文件等业务。

## 项目概览

- 后端 API：`cmd/main.go` 启动，默认端口 `10024`
- 管理端：`web-admin`，默认 Vite 端口 `5173`
- 数据库：MySQL，启动时可执行 AutoMigrate、种子数据和默认字典初始化
- 缓存：Redis，可通过环境变量开关
- 文件服务：RustFS/MinIO S3 兼容对象存储
- API 文档：Swagger + Scalar
- 权限模型：JWT + RBAC + 门店业务隔离

## 核心功能

### 门店与权限

- 门店资料、门店编码、门店关联第三方账号
- 用户、角色、菜单、角色菜单权限
- 门店级业务权限与数据隔离
- 动态菜单、动态前端路由

### 供应链

- 供应商、供应商分类、供应商商品
- 商品规格单位、门店供应商商品关联
- 采购订单、采购明细、采购流程
- 价格清单、门店价格管理

### 库存与门店业务

- 库存台账、库存订单、库存流水
- 报损、自用、赠送等库存损耗单
- 门店记账、记账明细、通知图片生成
- 门店退货单
- 经营统计与仪表盘

### 会员与 B2B

- 会员资料、钱包流水、充值订单
- B2B 客户、客户价格、供货订单
- 适用于门店对客户供货、批发或企业客户管理

### 外部集成

- 钉钉机器人、钉钉 Stream 客户端、消息模板
- 美团 AI 建议能力
- 第三方账号池、第三方订单、物流路线导入与历史查询
- 芯烨云打印机、打印机状态同步定时任务
- RustFS/MinIO 文件上传、图库管理、通知图片存储

### 工程能力

- Swagger 自动生成
- 健康检查：`/health`、`/ready`、`/live`
- WebSocket：`/ws`
- 慢查询、缓存、本地热路径等性能配置
- 数据权限策略、租户/门店上下文、请求日志、恢复中间件

## 技术栈

| 层级      | 技术                                           |
| --------- | ---------------------------------------------- |
| 后端语言  | Go `1.25.3`                                    |
| Web 框架  | Gin `v1.11.0`                                  |
| ORM       | GORM `v1.31.0`                                 |
| 数据库    | MySQL                                          |
| 缓存      | Redis                                          |
| 认证      | JWT                                            |
| 日志      | Zap + Lumberjack                               |
| API 文档  | swaggo/gin-swagger + Scalar                    |
| 对象存储  | MinIO SDK / RustFS                             |
| 实时通信  | gorilla/websocket、DingTalk Stream SDK         |
| 前端      | Vue 3、Vite、TypeScript、Pinia、Vue Router     |
| UI 与图表 | Arco Design Vue、UnoCSS、ECharts、Handsontable |

## 目录结构

```text
tower-go/
├── cmd/                         # 命令入口
│   ├── main.go                  # API 服务入口
│   ├── apply_indexes/           # 应用索引工具
│   ├── verify_indexes/          # 索引校验工具
│   └── init_dingtalk_menu/      # 钉钉菜单初始化工具
├── bootstrap/                   # 应用启动、配置、数据库、迁移、路由装配
├── config/                      # 环境变量配置与性能配置
├── controller/                  # HTTP 控制器
├── service/                     # 业务服务层
├── module/                      # 数据访问与仓储封装
├── model/                       # GORM 模型
├── router/                      # 路由注册
│   └── api/                     # 各业务模块 API 路由
├── middleware/                  # 认证、权限、门店业务守卫、日志、恢复
├── internal/                    # 内部领域能力，如数据权限上下文
├── pkg/                         # 可复用基础包
│   ├── auth/                    # JWT、密码工具
│   ├── performance/             # 性能优化组件与说明
│   ├── tenant/                  # 租户/门店作用域
│   ├── search/                  # 查询构建与优化
│   └── xpyun/                   # 芯烨云打印 SDK 封装
├── utils/                       # 数据库、缓存、日志、通用工具
├── migrations/                  # 初始化 SQL 与种子数据
├── docs/                        # Swagger 文档与架构资料
├── cron/                        # 定时任务
└── web-admin/                   # Vue 3 管理端
```

## 快速开始

### 1. 准备环境

- Go `1.25.3` 或与 `go.mod` 兼容的版本
- Node.js `20+`，npm
- MySQL `8.0+`
- Redis，可选但默认开启
- RustFS/MinIO，可选，启用文件/图库/通知图片时需要
- `swag`，用于生成 Swagger 文档

安装 `swag`：

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 2. 配置后端

```bash
cp .env.example .env
```

编辑 `.env`，至少确认这些配置：

```env
APP_PORT=10024

DB_DRIVER=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=your_password
DB_NAME=tower

REDIS_ENABLED=true
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

JWT_SECRET=replace_with_a_long_random_secret
```

创建数据库：

```sql
CREATE DATABASE tower CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

导入初始化 SQL：

```bash
mysql -u root -p tower < migrations/init.sql
mysql -u root -p tower < migrations/init_seed_data.sql
```

启动后端：

```bash
go mod tidy
go run cmd/main.go
```

也可以使用 Makefile：

```bash
make run
```

### 3. 配置前端

```bash
cd web-admin
npm install
npm run dev
```

默认访问：

- 管理端：`http://localhost:5173`
- 后端 API：`http://localhost:10024`
- Swagger：`http://localhost:10024/swagger/index.html`
- Scalar 文档：`http://localhost:10024/docs`

`web-admin/vite.config.ts` 默认把 `/api` 代理到线上地址。需要本地联调时，将代理目标改为：

```ts
target: "http://localhost:10024";
```

## 常用命令

```bash
# 后端开发启动，先生成 Swagger
make run

# 只生成 Swagger
make docs

# 构建后端二进制
make build

# 直接启动后端
go run cmd/main.go

# 后端测试
go test ./...

# 前端开发
cd web-admin && npm run dev

# 前端构建
cd web-admin && npm run build
```

## 环境变量说明

| 变量                               | 说明                                 | 默认/示例                |
| ---------------------------------- | ------------------------------------ | ------------------------ |
| `APP_NAME`                         | 应用名称                             | `tower-go`               |
| `APP_PORT`                         | API 端口                             | `10024`                  |
| `DB_DRIVER`                        | 数据库驱动                           | `mysql`                  |
| `DB_HOST`                          | 数据库地址                           | `127.0.0.1`              |
| `DB_PORT`                          | 数据库端口                           | `3306`                   |
| `DB_USERNAME`                      | 数据库用户                           | `root`                   |
| `DB_PASSWORD`                      | 数据库密码                           | 自行配置                 |
| `DB_NAME`                          | 数据库名                             | `tower`                  |
| `REDIS_ENABLED`                    | 是否启用 Redis                       | `true`                   |
| `REDIS_HOST`                       | Redis 地址                           | `127.0.0.1`              |
| `REDIS_PORT`                       | Redis 端口                           | `6379`                   |
| `REDIS_DB`                         | Redis DB 编号                        | `0`                      |
| `JWT_SECRET`                       | JWT 签名密钥                         | 建议 32 位以上随机字符串 |
| `SKIP_AUTO_MIGRATE`                | 跳过启动 AutoMigrate                 | `1` 表示跳过             |
| `SKIP_SEED_DATA`                   | 跳过启动种子数据、默认字典、默认模板 | `1` 表示跳过             |
| `SWAG_AUTO`                        | 控制启动时自动生成 Swagger           | `0` 表示禁用             |
| `RUSTFS_ENABLED`                   | 是否启用 RustFS/MinIO                | `true`/`false`           |
| `RUSTFS_ENDPOINT`                  | S3 兼容服务地址                      | `127.0.0.1:9000`         |
| `RUSTFS_PUBLIC_BASE_URL`           | 对外访问根地址                       | `https://example.com`    |
| `RUSTFS_BUCKET`                    | 默认文件 bucket                      | `images`                 |
| `RUSTFS_NOTIFY_BUCKET`             | 通知图片 bucket                      | `notify`                 |
| `DINGTALK_CLIENT_ID`               | 钉钉 Stream Client ID                | 可选                     |
| `DINGTALK_CLIENT_SECRET`           | 钉钉 Stream Client Secret            | 可选                     |
| `DINGTALK_MENU_REPORT_WEBHOOK_URL` | 钉钉报菜通知 Webhook                 | 可选                     |
| `XPYUN_USER`                       | 芯烨云账号                           | 可选                     |
| `XPYUN_USER_KEY`                   | 芯烨云 UserKey                       | 可选                     |
| `DEEPSEEK_API_KEY`                 | 美团 AI 建议使用的模型密钥           | 可选                     |

更多性能相关变量见 `.env.example` 和 `config/performance.go`。

## API 模块

后端统一前缀为 `/api/v1`，主要模块如下：

| 模块                 | 路由前缀                                                            |
| -------------------- | ------------------------------------------------------------------- |
| 认证                 | `/auth`                                                             |
| 用户                 | `/users`                                                            |
| 角色                 | `/roles`                                                            |
| 权限                 | `/permission`                                                       |
| 门店                 | `/stores`                                                           |
| 菜单                 | `/menus`                                                            |
| 字典类型/字典数据    | `/dict-types`、`/dict-data`                                         |
| 供应商               | `/suppliers`                                                        |
| 供应商商品/分类/规格 | `/supplier-products`、`/supplier-categories`、`/product-unit-specs` |
| 门店供应商           | `/store-suppliers`                                                  |
| 采购订单             | `/purchase-orders`                                                  |
| 库存                 | `/inventories`、`/inventory-orders`                                 |
| 库存损耗             | `/inventory-loss-orders`                                            |
| 门店记账             | `/store-accounts`                                                   |
| 门店退货             | `/store-returns`                                                    |
| 会员                 | `/members`、`/wallet-logs`、`/recharge-orders`                      |
| B2B                  | `/b2b`                                                              |
| 价格清单             | `/price-lists`                                                      |
| 统计                 | `/statistics`                                                       |
| 美团 AI              | `/meituan-ai`                                                       |
| 钉钉                 | `/dingtalk`                                                         |
| 消息模板             | `/message-templates`                                                |
| 打印机               | `/printers`                                                         |
| 第三方账号/路线      | `/third-party-accounts`、`/third-party-routes`                      |
| 文件/图库            | `/files`、`/galleries`                                              |
| 公开供应商档案       | `/public/suppliers`                                                 |

## 前端模块

`web-admin/src/views` 当前包含：

- 登录、个人资料、经营数据大屏
- 门店列表、供应商、采购、库存、报损、自用、赠送
- 门店记账、门店退货、会员、B2B、美团 AI
- 系统用户、角色、菜单、字典、图库、消息模板
- 钉钉机器人、打印机
- 第三方账号池、物流路线、导入订单、历史物流单
- 公开供应商档案页

## 认证与权限

登录接口：

```bash
curl -X POST http://localhost:10024/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13000000000","password":"your_password"}'
```

携带 Token：

```bash
curl http://localhost:10024/api/v1/users/profile \
  -H "Authorization: Bearer <token>"
```

权限由 `middleware.AuthMiddleware()`、`middleware.Permission()` 和 `middleware.StoreBusinessGuard()` 共同控制：

- `AuthMiddleware` 校验 JWT 并写入用户上下文
- `Permission` 校验菜单/按钮权限编码
- `StoreBusinessGuard` 保护门店业务数据边界

权限位：

```text
1  = 删除
2  = 修改
4  = 新增
8  = 查看
15 = 全部权限
```

## 数据初始化与迁移

项目启动时会执行以下初始化逻辑：

- 加载 `.env`
- 初始化数据库连接
- 初始化 Redis 缓存
- 根据环境变量决定是否执行 AutoMigrate
- 执行 SQL 种子、默认字典、默认消息模板
- 修正超级管理员 `store_id=0` 相关历史数据
- 初始化事件订阅、会话管理、定时任务、钉钉 Stream 客户端

生产环境建议：

- 使用正式迁移流程管理表结构
- 设置 `SKIP_AUTO_MIGRATE=1`
- 设置 `SKIP_SEED_DATA=1`
- 避免每次启动重复执行初始化数据

## 文件与图片服务

文件服务依赖 RustFS/MinIO。启用后可使用：

- 文件上传：`/api/v1/files`
- 图库管理：`/api/v1/galleries`
- 记账通知图片生成与上传
- 通知图片独立 bucket

未启用 `RUSTFS_ENABLED` 时，相关控制器可能不会被初始化，文件/图库接口不可用。

## 性能配置

性能相关能力集中在：

- `config/performance.go`
- `pkg/performance/`
- `pkg/search/`
- `utils/cache/`

支持的方向包括：

- 数据库连接池参数
- 慢查询阈值与查询日志
- 本地缓存、Redis 缓存、缓存锁
- worker pool 与并发限制
- 性能指标采集
- pprof 开关

详细说明可参考：

- `pkg/performance/README.md`
- `pkg/performance/CACHE_MANAGER_README.md`
- `pkg/performance/PERFORMANCE_REPORT.md`

## 开发约定

- 新业务接口建议按 `model -> module -> service -> controller -> router/api` 的层次补齐
- 需要权限控制的接口在路由层显式挂 `middleware.Permission("xxx")`
- 门店业务接口优先加 `middleware.StoreBusinessGuard()`
- 公共能力放入 `pkg/`，应用内通用工具放入 `utils/`
- 新增前端页面时同步菜单、权限编码和动态路由配置
- 修改 API 后运行 `make docs` 更新 Swagger
- 涉及表结构时同步维护 `migrations/init.sql` 与必要种子数据

## 部署提示

后端构建：

```bash
make build
```

前端构建：

```bash
cd web-admin
npm run build
```

部署时需要确认：

- `.env` 使用生产配置，不提交到版本控制
- MySQL、Redis、RustFS/MinIO 网络可达
- `JWT_SECRET`、数据库密码、对象存储密钥、钉钉密钥已替换
- 前端代理或 Nginx 反向代理正确指向后端 `/api`
- 生产环境关闭不必要的调试、自动迁移和重复种子写入
