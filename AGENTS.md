# Tower-Go 项目功能与开发约束

> 本文件是仓库级 AI 编码约束，适用于本目录及全部子目录。后续 AI 在分析、设计、修改、审查代码时必须优先阅读并遵守本文件。若具体子目录将来存在更细粒度的 `AGENTS.md`，则同时遵守，冲突时以更具体的规则为准。

## 1. 项目定位与核心功能

Tower-Go 是面向连锁门店的供应链与经营管理系统：

- 后端：Go + Gin + GORM + MySQL，入口为 `cmd/main.go`，依赖装配在 `bootstrap/` 和 `router/api/container.go`。
- 管理端：`web-admin/`，使用 Vue 3 + TypeScript + Vite + Arco Design Vue + Pinia + TanStack Vue Query。
- 基础设施：Redis、Zap/Lumberjack、JWT、RustFS/MinIO、OpenAPI、WebSocket、钉钉、云打印等。
- 主要业务：门店、用户角色菜单、供应商、采购、库存及损耗、门店记账、退货、会员钱包、B2B、价格、统计、短信、第三方订单、钉钉、打印、文件图库、AI 助手。

系统不是普通单店 CRUD 后台。任何改动都必须同时考虑：

1. 身份认证与接口权限；
2. 总部、门店与行级数据范围；
3. 门店营业状态；
4. 自然日与营业日差异；
5. 库存、余额、状态迁移的事务与并发；
6. 审计、幂等和敏感信息保护；
7. 后端、前端、菜单权限与 OpenAPI 契约的一致性。

## 2. 最高优先级规则

### 2.1 修改前必须先检索复用能力

**禁止先写代码再找公共模块。** 实现功能前必须按顺序执行：

1. 用 `rg` 搜索同类业务、类型名、错误码、权限码和接口路径；
2. 查看同业务的 `model -> module -> service -> controller -> router/api` 完整链路；
3. 检查本文件第 4 节列出的公共能力；
4. 选取最接近的既有模块作为实现范例；
5. 仅在确认现有能力无法满足时新增抽象，并说明为何不能复用。

同类功能已有两套实现时，不得创建第三套，也不得仅凭目录名选择“看起来更新”的实现；应沿用目标业务当前使用的主链路。

### 2.2 中文强制要求

- **本次新写或语义上修改的人工代码注释必须使用简体中文。** 不要求为无关变更翻译存量英文注释。
- **所有新增或修改的运行时打印、日志 message、告警和面向用户的提示必须使用简体中文。** 这是新代码的收敛要求；仓库中的存量英文日志或调试打印不代表推荐范例。
- Go 导出符号注释可保留标识符开头，例如：`// DataPermission 根据鉴权上下文追加数据范围条件。`
- 日志结构化字段名、JSON 字段、协议字段、数据库列名保持既有英文命名；完整描述句必须为中文。
- JWT、HTTP、SQL、Redis、OpenAPI、ID 等技术名词可以保留英文，但不能写成整句英文提示。
- Swagger 指令、第三方协议固定字面量和自动生成文件不要求翻译；不要手工修改生成文件。
- 禁止新增 `fmt.Print*`、`log.Print*`、`console.log` 作为业务日志；后端统一使用 `utils/logging`，前端用户反馈统一使用 `src/feedback`。
- 禁止将 validator、数据库、第三方 SDK 的原始英文错误直接展示给用户。

正确示例：

```go
logging.LogError("库存扣减失败",
    zap.Uint("store_id", storeID),
    zap.Uint("product_id", productID),
    zap.Error(err),
)
```

错误示例：

```go
log.Printf("update inventory failed: %v", err)
logging.LogError(fmt.Sprintf("库存扣减失败: %v", err))
```

### 2.3 保持变更边界

- 不修改与当前需求无关的文件，不顺手大规模格式化或重构。
- 不覆盖、回退或整理用户已有的未提交修改。
- 修改前后都要检查 `git status --short` 和目标文件 diff。
- 发现既有缺陷但不属于当前需求时，记录风险，不擅自扩大范围。
- 不提交 `.env`、密钥、Token、密码、生产地址或真实业务数据。

## 3. 后端架构与分层

标准请求链路：

```text
router/api -> controller -> service -> module -> GORM/MySQL
```

新增完整业务通常按以下文件对齐：

```text
model/<domain>.go
module/<domain>.go
service/<domain>.go
controller/<domain>.go
router/api/<domain>.go
```

### 3.1 `model/`

- 放置 GORM 实体、请求 DTO、列表参数和响应结构。
- 请求结构使用 `CreateXxxReq`、`UpdateXxxReq`、`ListXxxReq` 等明确命名。
- 输入 DTO 与持久化实体必须区分，禁止把用户输入直接用于不受控的全字段更新。
- 金额字段先沿用目标业务已有类型和数据库精度：钱包、余额等已使用 `decimal.Decimal` 的链路必须继续使用 decimal；存量 `float64 + decimal(...)` 链路不得在普通需求中顺手迁移。金额类型迁移必须作为独立的全链路改造处理。
- 索引、唯一约束、软删除和字段精度应在 GORM tag 与迁移中保持一致。

### 3.2 `module/`

本项目的 `module` 实际承担 DAO/Repository 职责：

- 负责 GORM 查询、CRUD、Preload、分页、Scope、事务、锁和原子更新。
- 完整业务原子写操作应封装在此层，不应由 Controller 拼接数据库流程。
- 事务闭包内必须始终使用 `tx`，不得误用 `m.db`。
- 单条门店资源优先提供或复用 `GetByIDScoped`、`Update*Scoped`、`Delete*Scoped` 一类方法。
- 不要新增与 `module` 平行且语义相同的 DAO 层；需要可测试性时可由 Service 定义窄接口，由 module 实现。

### 3.3 `service/`

- 负责编排业务流程、状态机、跨 module 协作、业务校验、外部服务调用和 DTO 构造。
- 不直接处理 Gin 响应。
- 不向上层泄露数据库或第三方原始错误。
- 优先依赖窄接口，以便使用 fake 编写单元测试。
- 外部网络调用不要长时间放在数据库事务中；使用提交后发送、可重试状态或 Outbox 类方案。

### 3.4 `controller/`

- 只做 HTTP 适配：绑定和校验参数、读取鉴权上下文、调用 Service、输出统一响应。
- 禁止写 SQL、开启业务事务或堆积复杂业务规则。
- JSON 优先使用 `utils/http.BindJSON`，ID 优先使用 `utils/http.ParseUintParam`，Query 使用带 binding tag 的 DTO。
- Service 返回错误后使用 `http.ErrorFrom(ctx, err)`。
- 禁止直接向客户端返回 `err.Error()`。

### 3.5 `router/api/`

- 只负责路径、HTTP 方法、中间件、权限码和 Handler 注册。
- 新 Controller/Service/Module 在 `router/api/container.go` 中集中装配，禁止新增随意的全局单例。
- 私有 `/api/v1` 路由必须明确认证、权限和门店营业状态策略。
- 新权限码须与菜单/种子数据、前端权限保持一致。

## 4. 后端公共能力复用清单

开始编码前必须检查本节。表中“当前主链路”表示普通业务新代码应优先沿用。

| 领域 | 当前主链路/位置 | 约束 |
| --- | --- | --- |
| 业务错误码 | `pkg/apicode` | 新错误集中定义；Service/Module 使用 `New`、`Newf`、`Wrap` |
| HTTP 响应 | `utils/http` | 使用 `Success`、`SuccessWithPagination`、`ErrorFrom`；禁止新增响应结构 |
| JWT | `utils/auth/jwt.go`、`middleware/auth.go` | 不自行解析 Token，不信任客户端身份 Header |
| 接口权限 | `middleware/permission.go` | 路由使用 `Permission`、`PermissionAny`、`PermissionOrRoles` |
| 总部身份 | `middleware.HQUnboundAdmin` | 不能把普通 `admin` 直接等同于跨店管理员 |
| 鉴权上下文 | `internal/authctx` | AuthMiddleware 注入 Gin Context；标准 Context 需显式桥接 |
| 数据范围 | `internal/datascope`、`pkg/datascope` | 列表使用集中 Scope；缺少上下文时必须 fail-closed |
| 营业状态 | `middleware/store_business_guard.go` | 门店写业务复用 `StoreBusinessGuard` |
| 日志 | `utils/logging` | 中文 message + Zap 结构化字段 |
| 数据库事务 | GORM `db.Transaction`、`utils/database` | 原子业务统一事务，事务内只用 `tx` |
| 营业日 | `utils/businessdate` | 经营口径使用 `Date`/`DateString`，不能自行格式化当前日期 |
| 分页响应 | `utils/http.SuccessWithPagination` | 当前字段为 `list/total/page/page_size/page_num` |
| 参数校验 | Gin binding、`utils/http` | 结构校验在 DTO，业务校验在 Service，返回中文稳定错误 |
| 审计 | `middleware/audit_log*.go` | 写接口默认纳入审计，敏感信息必须脱敏 |
| 请求防重 | `middleware/duplicate_request.go` | 只是短期防抖，不能替代业务幂等 |
| 批处理 | `utils/batch` | 批量操作先检查已有批处理与事务封装 |
| 缓存 | `utils/cache`、`pkg/performance` | 复用既有缓存、锁和性能配置，不能自建不受管控缓存 |
| 事件 | `utils/events/event_bus.go`、`utils/event_bus.go` | 目前是双轨预留能力，未形成稳定主链路；新增前先确认调用方，禁止再造第三套 |
| Excel/XML | `utils/excelxml` | 导入导出前先复用现有解析能力和 `controller/export_helpers.go` |
| 安全/加密 | `utils/encryption`、`utils/auth/password.go` | 复用已有 AES/密码能力；禁止自行实现密码、签名或加密算法 |
| 云打印 | `pkg/xpyun`、既有 printer service | 禁止绕过 SDK 封装直接散落第三方请求 |
| 阿里云短信 | `pkg/aliyunsms`、既有短信 service | 模板、发送、状态与幂等走现有链路 |
| 客户端来源 | `pkg/clientsource` | 保持与前端 `X-Client-Source` 约定一致 |

### 4.1 双轨能力的特别警告

仓库存在尚未统一的平行实现：

- `utils/http` 与 `pkg/http/request`、`pkg/http/response`；
- `middleware/auth.go` / `utils/auth` 与 `pkg/middleware/auth` / `pkg/auth`；
- 主流分页写法与 `utils/pagination`；
- 主流 Gin binding 与较少使用的 `utils/validation`；
- 当前 module 模式与 `pkg/repository` 的预留接口。

在没有全局迁移设计和完整回归验证时：

- 普通业务必须沿用当前模块已经使用的主链路；
- 不得在单个新接口中切换响应协议、JWT Claims、Context key 或 HTTP 状态语义；
- 不得为了“统一”而一次性改动全仓库；
- 不得创建第三套 helper。

当前 `utils/http` 即使业务失败也返回 HTTP 200，真实业务状态位于 JSON `code`。这是现有前后端契约；除非执行全局迁移，不得局部改成另一套 HTTP status 语义。

## 5. 错误码、响应与日志

### 5.1 错误码

`pkg/apicode` 是唯一业务错误码真源，格式为 HTTP 语义前三位加两位业务序号：

- `400xx` 参数错误；
- `401xx` 认证错误；
- `403xx` 权限或业务禁止；
- `404xx` 资源不存在；
- `409xx` 重复、库存不足、状态冲突；
- `410xx` 资源已过期；
- `413xx` 请求或上传内容过大；
- `422xx` 业务数据校验失败；
- `500xx` 服务端内部错误；
- `502xx` 外部服务调用失败；
- `503xx` 服务不可用。

要求：

- 禁止在 Controller/Service 新增裸数字错误码。
- 动态业务提示使用 `apicode.Newf` 且必须为中文。
- 底层错误用 `apicode.Wrap` 保留 cause，不能暴露给客户端。
- 未识别错误由 `http.ErrorFrom` 统一映射为内部错误。

### 5.2 响应

- 普通成功：`http.Success(ctx, data)`。
- 分页成功：`http.SuccessWithPagination(ctx, list, total, page, pageSize)`。
- 错误：`http.ErrorFrom(ctx, err)` 或已知中间件场景的 `http.ErrorApp`。
- 新业务禁止调用会把 `err.Error()` 写入响应的 `http.WithError`，也禁止启用 `printResponse`、`printErrorResponse` 调试输出。
- 禁止 Controller 自定义另一套 envelope。
- 禁止依赖响应层按字段名格式化时间来掩盖 DTO 语义不清；日期字段应明确表达自然日或营业日。

### 5.3 日志与敏感数据

- 日志 message 必须中文，字段名使用稳定的英文 `snake_case`。
- 可结构化的 ID、单号、状态、耗时、错误必须放在 Zap 字段中，不拼接进 message。
- 密码、Token、Cookie、密钥、身份证、银行卡、完整手机号等不得写入日志或审计明文。
- 不得记录完整请求体，除非已评估并完成敏感字段脱敏。
- 错误日志应保留 `zap.Error(err)`；用户提示与内部诊断必须分离。

## 6. 鉴权、权限与数据范围

### 6.1 接口鉴权

- 用户接口复用 `middleware.AuthMiddleware()`；内部服务 `/api/internal/v1` 使用独立内部服务鉴权，二者不得混用。
- 禁止相信 `X-User-Id`、`X-Store-Id` 或请求参数中的 `role_code/data_scope/user_id` 作为授权依据。
- 前端隐藏按钮不是安全边界，后端路由必须验证权限。

### 6.2 RBAC

- 优先复用菜单和种子数据中已有权限码；新增权限码沿用目标业务现行命名风格，可能为两段（如 `inventory:in`）或三段（如 `store:return:add`），禁止擅自重命名既有权限码。
- 路由层显式挂载最小权限；多个权限择一时使用 `PermissionAny`。
- 角色绕过只能在确有业务依据的局部路由使用 `PermissionOrRoles`。
- 禁止在 Controller 到处手写 `roleCode == "admin"`。

### 6.3 总部与门店

- `super_admin` 可全局访问。
- `admin` 只有 `store_id == 0` 时才是未绑定门店的总部管理员。
- 绑定门店的 admin 仍受门店数据范围约束。
- 是否跨店统一使用既有 `HQUnboundAdmin` 和 `ResolveQueryStoreID` 逻辑。
- 普通门店用户传入的 `store_id` 不能覆盖其真实门店。

### 6.4 行级数据范围

- RBAC 只决定“能否调用接口”，DataScope 决定“能看到哪些行”，两者都必须存在。
- `AuthMiddleware` 只自动注入 Gin Context；调用依赖标准 `context.Context` 数据范围的 Service 前，Controller 必须先调用 `middleware.AttachAuthContextToHTTPRequest(ctx)`，再传入 `ctx.Request.Context()`。
- 列表查询优先复用 `internal/authctx` + `internal/datascope.DataPermission` 及集中 `TablePolicy`；`pkg/datascope` 是部分列表使用的兼容入口，不代表另一套完整权限体系。
- JOIN 查询中的范围列必须带表名或别名，避免歧义。
- 缺失鉴权上下文时必须 fail-closed（例如 `WHERE 1=0`），不能默认全量。
- 列表过滤不能保护详情、更新和删除；单条操作必须在 SQL/Scoped 方法中校验资源归属。
- 后台任务、机器人和内部调用若没有 HTTP AuthContext，必须显式声明数据范围。

## 7. 事务、并发、幂等与审计

### 7.1 事务

以下操作必须按完整业务原子单元设计事务：

- 主表与明细；
- 库存与出入库流水；
- 余额与钱包流水；
- 状态迁移与发送/处理记录；
- 替换关联数据；
- 多表统计口径更新。

要求：

- 事务内全部查询和写入使用 `tx`。
- 金额、库存和状态并发更新使用行锁、条件更新或乐观冲突检测，并检查 `RowsAffected`。
- 失败必须完整回滚，不能留下主单无明细等半成品。
- 远程调用与长耗时计算尽量放到事务之外。

### 7.2 幂等

`middleware/duplicate_request.go` 只提供短期重复请求抑制，不是持久业务幂等。

支付、余额、库存、短信、打印、外部推送和重要单据创建必须结合：

- 业务幂等键；
- 数据库唯一索引；
- 状态条件更新；
- `RowsAffected` 检查；
- 可重试发送状态或 Outbox。

禁止以“先查询是否存在，再插入”替代唯一约束。

### 7.3 审计

- 新增 POST/PUT/PATCH/DELETE 默认应被通用审计中间件覆盖。
- 权限调整、余额/库存人工调整、导出敏感数据等高风险操作应明确审计策略。
- 通用中间件提供请求级审计，但不会自动产生完整 before/after/diff；高风险业务需额外设计业务快照与差异审计。
- 审计请求、Query、before/after/diff 数据时必须脱敏。
- 不得为了减少日志而随意加入审计排除列表。

## 8. 日期、时间、分页与验证

### 8.1 营业日

门店经营日期统一使用 `utils/businessdate`：00:00–04:59 归前一自然日，05:00 起归当天；传入的 `time.Time` 必须具有明确的业务时区。

- 记账、营业额、经营统计、库存经营日等必须使用 `businessdate.Date` 或 `DateString`。
- 禁止在营业日场景直接调用 `time.Now().Format("2006-01-02")`。
- 自然日业务必须明确说明，不得误套营业日。
- 日期请求按 `2006-01-02` 严格解析，无效日期返回统一中文错误码。
- 新逻辑必须考虑部署时区，不得无意混用 UTC 与 Local。

### 8.2 分页

推荐请求字段：

```go
Page     int `form:"page,default=1" binding:"min=1"`
PageSize int `form:"page_size,default=20" binding:"min=1,max=100"`
```

要求：

- 新增普通列表默认页码 1、每页 20、最大 100；兼容旧接口或批量选择场景可保留既有特殊上限，但必须说明原因，不能在无关需求中破坏原契约。
- Count 与数据列表使用完全相同的数据范围和业务筛选。
- 排序必须稳定，通常在业务排序后追加 `id DESC`。
- 使用 `SuccessWithPagination`，禁止另造分页响应格式。

### 8.3 参数验证

- DTO binding 负责必填、长度、范围、枚举和集合元素等结构性验证。
- Service 负责资源存在性、归属、状态和跨字段等业务验证。
- 校验失败必须映射为中文且稳定的 `apicode`。
- 禁止把 Gin/validator 原始英文错误直接返回前端。
- 禁止拼接用户输入生成 SQL；所谓字符串“清洗”不能替代参数绑定。

## 9. 数据库、迁移与配置

- 新增模型或字段时检查 `bootstrap/migrate.go` 的 AutoMigrate 模型及完整性检查。
- 同步维护 `migrations/init.sql`；需要默认数据时同步维护种子 SQL/Seed 版本。
- 唯一性、幂等与关键数据完整性必须由数据库约束兜底。
- 修改索引应考虑存量数据、回填、锁表风险及生产迁移方式。
- 配置统一进入 `config/`，通过环境变量加载；不得散落读取环境变量或硬编码生产值。
- 新配置应提供明确默认值、边界校验和启动时的中文错误提示。
- 存量违规写法不是范例：即使同类旧代码仍使用裸错误、全局 DB、英文日志或直接 `ShouldBind*`，新代码也应优先遵循本文件列出的主链路。
- 生产密钥绝不能写入前端、日志、示例真实值或版本库。

## 10. 前端架构与复用约束

### 10.1 分层

- `web-admin/src/api/`：按业务域定义请求和 TypeScript 类型。
- `views/`：页面编排，不重复实现通用基础组件。
- `components/base/`：按钮、表格、分页、表单、弹窗等基础组件。
- `feedback/`：Toast 与确认弹窗。
- `store/`：Pinia 客户端状态。
- TanStack Vue Query：服务端状态、缓存与刷新。
- `permission/`：路由、指令和动态菜单权限。
- `utils/currentStore.ts`：当前门店/租户解析。

大型页面新增复杂逻辑时应拆成业务组件或 composable，避免继续形成千行级单文件组件。

### 10.2 API 与状态

- 所有请求复用 `src/api/http.ts`，禁止自行创建 Axios 实例或直接 `fetch()` 重复实现认证。
- 请求头 `Authorization`、`x-tenant-id`、`X-Client-Source` 由统一拦截器处理。
- 响应类型复用 `src/api/types.ts` 的 `ApiEnvelope`、`Paginated` 等定义。
- 客户端状态放 Pinia，服务端状态放 Vue Query；禁止将同一服务端列表再复制到全局 Store。
- Mutation 成功后使用一致的 Query Key 精确失效缓存，避免全局粗暴刷新。
- 错误提示使用统一 Toast，并优先展示后端中文业务 message；不要在页面重复实现错误解析器。

### 10.3 权限与租户

- UI 权限复用 `v-permission`、`usePermission()` 和 `BaseTableRowActions`。
- 新页面优先接入后端菜单动态路由；菜单、路由、按钮和后端权限码必须一致。
- 前端权限仅用于体验，不能代替后端授权。
- 当前门店统一通过 user store 和 `currentStore.ts` 获取，禁止从 URL、DOM 或任意 localStorage 字段猜测。
- 不得在 `localStorage` 保存明文密码、API Key 或其他敏感凭证。

### 10.4 组件、表单、日期与分页

- 优先复用 `src/components/base/` 的 `BaseButton`、`BaseTable`、`BasePagination`、`BaseDialog`、`BaseForm` 等。
- 反馈复用 `src/feedback/toast.ts`、`confirm.ts`；不得重复实现 Toast、确认框和 loading。
- 日期选择器使用 `value-format="YYYY-MM-DD"`；不要用 UTC 的 `toISOString().slice(0, 10)` 计算本地“今天”。
- 后端分页参数使用 `page/page_size`，页面大小变化后回到第 1 页。
- 前端校验只改善体验，不能替代后端验证；复杂表单应使用统一 model/rules/validate 机制。
- 用户可见文本、校验提示、Toast 和运行时控制台错误说明必须使用中文。

### 10.5 AI 功能安全

扩展 AI 功能时必须复用：

- `src/components/ai-assistant/AIAssistant.vue`；
- `src/api/aiAssistant.ts`；
- 统一 HTTP、权限、门店上下文、反馈和基础组件；
- `marked` + `DOMPurify` 的安全 Markdown 渲染链路。

禁止：

- 复制出另一套行为漂移的 AI 助手；
- 将 API Key、系统敏感配置、Token、密码或租户隐私拼入前端代码、日志或 prompt；
- 未清洗就通过 `v-html` 渲染模型输出；
- 让模型直接执行删除、记账、扣库存、发短信等副作用；此类操作必须经过后端权限、参数验证、幂等和明确人工确认。

## 11. OpenAPI、文档和生成文件

- 修改 API 的路径、参数、请求体、响应或错误语义后，更新对应 Swagger 注解。
- 执行 `make docs` 生成并校验 `docs/openapi.json`。
- `docs/swagger.*`、`docs/openapi.json` 等生成产物应通过命令更新，不手工局部修改。
- 新权限、新菜单和新页面应同步说明其访问范围。
- README 中命令、目录或环境变量因改动失效时应同步更新。

## 12. 测试与交付检查

### 12.1 后端测试

新增或修改关键逻辑至少覆盖与本次变更相关的场景：

- 正常路径与参数边界；
- 未登录、无权限、门店越权；
- 总部跨店与绑定门店 admin；
- 事务回滚与并发状态冲突；
- 重复提交和幂等；
- 营业日 05:00 边界；
- 分页 Count、筛选和稳定排序；
- 未包装底层错误不泄露；
- 敏感信息不进入日志/响应。

测试风格优先：

- 表驱动测试 + `t.Run`；
- Service 通过窄接口 fake；
- HTTP/中间件使用 `httptest`；
- GORM 查询可使用既有 DryRun 模式，但关键事务与并发行为不能只靠 SQL 字符串推断。

### 12.2 前端检查

当前前端没有正式测试脚本。修改前端至少执行：

```bash
cd web-admin
npm run build
```

并手工验证相关页面的：加载态、空态、错误态、权限隐藏、分页、表单边界和移动端可用性。若引入测试框架，必须同时补充 package scripts 和使用说明，不能只放孤立测试文件。

### 12.3 最低交付命令

按变更范围执行：

```bash
# Go 格式化（仅格式化本次修改的 Go 文件）
gofmt -w <changed-go-files>

# 后端测试（可先运行目标包，再运行全量）
go test ./...

# API 变更
make docs

# 前端变更
cd web-admin && npm run build
```

最终必须：

1. 查看 `git status --short`；
2. 查看全部目标 diff，确认无密钥、调试打印、英文运行时提示和无关改动；
3. 确认未覆盖用户原有修改；
4. 在交付说明中列出修改文件、执行的验证命令和未能验证的风险。

## 13. 明确禁止事项

- 禁止绕过 `pkg/apicode` 新增裸错误码或直接返回 `err.Error()`。
- 禁止新增第三套 HTTP 响应、分页、JWT、权限、日志或日期工具。
- 禁止以普通 admin 角色判断替代总部未绑店判断。
- 禁止裸 `GetByID` 后仅在 Controller 做门店归属判断。
- 禁止列表应用 DataScope 后就忽略详情、更新、删除的数据范围。
- 禁止在事务闭包中混用 `tx` 与原始 DB。
- 禁止只依赖请求防重中间件保护库存、余额、单据和外部发送。
- 禁止在营业日场景自行格式化 `time.Now()`。
- 禁止 SQL 字符串拼接用户输入。
- 禁止在前端重复实现 Axios 客户端、Toast、确认框、权限和租户选择逻辑。
- 禁止仅用前端隐藏按钮保护敏感接口。
- 禁止未经清洗渲染 HTML/Markdown。
- 禁止提交调试打印、英文运行时提示、明文密码或密钥。
- 禁止为追求抽象而重写稳定模块；优先做最小、可验证、与既有风格一致的变更。

## 14. 开始编码前的快速清单

- [ ] 已阅读本文件、`README.md` 和目标业务完整链路。
- [ ] 已用 `rg` 找到同类实现和可复用公共模块。
- [ ] 已确认使用 `pkg/apicode`，没有新造错误/响应体系。
- [ ] 已设计认证、权限、总部/门店和 DataScope。
- [ ] 已区分自然日与营业日。
- [ ] 已设计事务、并发、幂等和审计。
- [ ] 已确认分页、验证、日志和敏感信息策略。
- [ ] 已确认前端复用 API、Store、权限、反馈和基础组件。
- [ ] 所有新增/修改注释、日志和用户提示均为中文。
- [ ] 已确定测试与交付验证命令。
