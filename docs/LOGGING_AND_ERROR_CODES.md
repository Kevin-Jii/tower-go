# 结构化日志与错误码体系使用指南

## 概述

本项目集成了 **uber-go/zap** 结构化日志系统和统一的错误码体系，提供高性能的日志记录和规范的错误处理机制。

---

## 一、结构化日志系统

### 1.1 日志配置

日志系统在 `cmd/main.go` 中初始化，默认配置：

```go
logConfig := &utils.LogConfig{
    Level:      "info",           // 日志级别
    FilePath:   "logs/app.log",   // 日志文件路径
    MaxSize:    100,              // 单个文件最大 100MB
    MaxBackups: 10,               // 保留 10 个旧文件
    MaxAge:     30,               // 保留 30 天
    Compress:   true,             // 压缩旧日志
    Console:    true,             // 同时输出到控制台
}
```

### 1.2 日志级别

支持以下日志级别（按严重程度递增）：

| 级别    | 说明                   | 使用场景                       |
|---------|------------------------|--------------------------------|
| Debug   | 调试信息               | 开发调试、详细的执行流程       |
| Info    | 一般信息               | 正常业务操作、系统状态         |
| Warn    | 警告信息               | 潜在问题、性能问题             |
| Error   | 错误信息               | 业务错误、异常情况             |
| Fatal   | 致命错误（会退出程序） | 无法恢复的严重错误             |

### 1.3 基础日志方法

#### 结构化日志（推荐）

```go
import (
    "tower-go/utils"
    "go.uber.org/zap"
)

// Debug 日志
utils.LogDebug("用户查询", 
    zap.String("username", "admin"),
    zap.Int("page", 1),
)

// Info 日志
utils.LogInfo("用户登录成功", 
    zap.Uint("user_id", 123),
    zap.String("ip", "192.168.1.1"),
)

// Warn 日志
utils.LogWarn("查询性能慢", 
    zap.Duration("duration", time.Second*2),
    zap.String("sql", "SELECT * FROM users"),
)

// Error 日志
utils.LogError("数据库查询失败", 
    zap.Error(err),
    zap.String("operation", "GetUser"),
)

// Fatal 日志（会终止程序）
utils.LogFatal("数据库连接失败", zap.Error(err))
```

#### 格式化日志（Printf 风格）

```go
// 简单场景使用
utils.Infof("用户 %s 登录成功", username)
utils.Warnf("查询耗时 %v，超过阈值", duration)
utils.Errorf("解析 JSON 失败: %v", err)
```

### 1.4 业务日志快捷方法

#### HTTP 请求日志

```go
utils.LogRequest(
    ctx.Request.Method,
    ctx.Request.URL.Path,
    ctx.ClientIP(),
    200,
    time.Since(startTime),
)
```

**输出示例：**
```json
{
  "timestamp": "2025-11-03 18:30:45.123",
  "level": "info",
  "msg": "HTTP Request",
  "method": "POST",
  "path": "/api/v1/users",
  "ip": "192.168.1.100",
  "status": 200,
  "latency": "0.125s"
}
```

#### 业务错误日志

```go
utils.LogBusinessError(
    utils.ErrUserNotFound,
    err,
    zap.String("user_id", userID),
)
```

#### 数据库错误日志

```go
utils.LogDatabaseError(
    "CreateUser",
    err,
    zap.String("username", username),
)
```

#### 认证错误日志

```go
utils.LogAuthError("login", userID, "密码错误")
```

#### 性能日志

```go
startTime := time.Now()
// ... 业务逻辑 ...
utils.LogPerformance("BuildMenuTree", time.Since(startTime))
```

**自动识别慢操作（>1秒）：**
```json
{
  "level": "warn",
  "msg": "Slow Operation",
  "operation": "BuildMenuTree",
  "duration": "2.345s"
}
```

#### WebSocket 事件日志

```go
utils.LogWebSocket("user_connected", userID,
    zap.String("session_id", sessionID),
)
```

#### 第三方服务调用日志

```go
utils.LogThirdParty(
    "wechat_pay",
    "create_order",
    true,
    time.Millisecond*500,
    nil,
)
```

---

## 二、错误码体系

### 2.1 错误码结构

所有错误码定义在 `utils/errors.go` 中，遵循统一格式：

```go
type ErrorCode struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}
```

### 2.2 错误码分类

| 错误码范围 | 分类               | 说明                          |
|------------|--------------------|-------------------------------|
| 200        | 成功               | 请求成功                      |
| 1000-1999  | 通用错误           | 系统级错误、数据库错误        |
| 2000-2999  | 认证授权错误       | 登录、Token、权限相关         |
| 3000-3999  | 用户业务错误       | 用户、门店、角色管理          |
| 4000-4999  | 菜品业务错误       | 菜品、报菜相关                |
| 5000-5999  | 权限菜单错误       | 菜单管理、权限分配            |
| 6000-6999  | WebSocket 错误     | 连接、会话相关                |
| 7000-7999  | 文件上传错误       | 文件操作、存储相关            |
| 8000-8999  | 第三方服务错误     | 外部 API、支付、短信等        |
| 9000-9999  | 业务逻辑错误       | 工作流、状态、并发冲突        |

### 2.3 常用错误码

#### 通用错误（1000-1999）

```go
utils.ErrBadRequest       // 1000: 请求参数错误
utils.ErrInternalServer   // 1001: 服务器内部错误
utils.ErrNotFound         // 1002: 资源不存在
utils.ErrDatabaseQuery    // 1100: 数据库查询错误
utils.ErrDuplicateKey     // 1104: 数据已存在
utils.ErrValidation       // 1200: 数据验证失败
utils.ErrInvalidEmail     // 1201: 邮箱格式不正确
utils.ErrInvalidPhone     // 1202: 手机号格式不正确
```

#### 认证授权错误（2000-2999）

```go
utils.ErrUnauthorized     // 2000: 未授权，请先登录
utils.ErrTokenInvalid     // 2001: Token 无效
utils.ErrTokenExpired     // 2002: Token 已过期
utils.ErrLoginFailed      // 2005: 用户名或密码错误
utils.ErrForbidden        // 2100: 无权限访问
utils.ErrStoreAccessDenied // 2104: 无权访问该门店数据
```

#### 用户业务错误（3000-3999）

```go
utils.ErrUserNotFound        // 3000: 用户不存在
utils.ErrUsernameAlreadyTaken // 3002: 用户名已被占用
utils.ErrPhoneAlreadyTaken   // 3003: 手机号已被占用
utils.ErrStoreNotFound       // 3100: 门店不存在
```

#### 菜品业务错误（4000-4999）

```go
utils.ErrDishNotFound        // 4000: 菜品不存在
utils.ErrMenuReportNotFound  // 4100: 报菜记录不存在
```

### 2.4 错误码使用方式

#### 方式一：直接使用响应方法（推荐）

```go
import "tower-go/utils"

// 使用错误码响应
utils.ErrorWithCode(ctx, utils.ErrUserNotFound)

// 使用错误码响应（带数据）
utils.ErrorWithCodeAndData(ctx, utils.ErrValidation, gin.H{
    "field": "phone",
    "reason": "长度必须为11位",
})
```

#### 方式二：自定义错误消息

```go
// 使用自定义消息
customErr := utils.ErrUserNotFound.WithMessage("用户 ID 123 不存在")
utils.ErrorWithCode(ctx, customErr)

// 使用格式化消息
customErr := utils.ErrStoreNotFound.WithMessageF("门店 %d 不存在", storeID)
utils.ErrorWithCode(ctx, customErr)
```

#### 方式三：兼容旧代码

```go
// 仍然支持原有的 Error 方法
utils.Error(ctx, http.StatusBadRequest, "请求参数错误")
```

### 2.5 Controller 层示例

```go
func (c *UserController) GetUser(ctx *gin.Context) {
    userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
    if err != nil {
        utils.ErrorWithCode(ctx, utils.ErrBadRequest.WithMessage("用户ID格式错误"))
        return
    }

    user, err := c.service.GetByID(uint(userID))
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            utils.ErrorWithCode(ctx, utils.ErrUserNotFound)
            return
        }
        utils.LogDatabaseError("GetUser", err, zap.Uint64("user_id", userID))
        utils.ErrorWithCode(ctx, utils.ErrDatabaseQuery)
        return
    }

    utils.Success(ctx, user)
}
```

### 2.6 Service 层示例

```go
func (s *UserService) CreateUser(req *CreateUserRequest) error {
    // 检查用户名是否已存在
    existing, _ := s.module.GetByUsername(req.Username)
    if existing != nil {
        return utils.ErrUsernameAlreadyTaken
    }

    // 检查手机号是否已存在
    existing, _ = s.module.GetByPhone(req.Phone)
    if existing != nil {
        return utils.ErrPhoneAlreadyTaken
    }

    // 创建用户
    if err := s.module.Create(user); err != nil {
        utils.LogDatabaseError("CreateUser", err, 
            zap.String("username", req.Username),
        )
        return utils.ErrUserCreateFailed
    }

    utils.LogBusinessEvent("user_created", 
        zap.String("username", req.Username),
        zap.Uint("store_id", req.StoreID),
    )
    
    return nil
}
```

---

## 三、日志文件管理

### 3.1 日志文件结构

```
logs/
├── app.log           # 当前日志文件
├── app-2025110301.log.gz  # 压缩的旧日志
├── app-2025110302.log.gz
└── ...
```

### 3.2 日志轮转规则

- **大小轮转**: 单个文件超过 100MB 自动切分
- **时间轮转**: 保留最近 30 天的日志
- **数量轮转**: 最多保留 10 个备份文件
- **自动压缩**: 旧日志自动 gzip 压缩

### 3.3 日志查看命令

```bash
# 查看实时日志
tail -f logs/app.log

# 查看最近 100 行
tail -n 100 logs/app.log

# 查看错误日志
grep '"level":"error"' logs/app.log

# 查看特定用户的操作
grep '"user_id":123' logs/app.log

# 查看慢查询
grep '"msg":"Slow Operation"' logs/app.log

# 查看特定时间段
grep '2025-11-03 18:' logs/app.log
```

---

## 四、最佳实践

### 4.1 日志记录原则

✅ **应该记录的日志：**
- 关键业务操作（登录、创建、删除）
- 错误和异常情况
- 性能问题（慢查询、慢接口）
- 外部服务调用
- 权限拒绝

❌ **不应该记录的日志：**
- 密码、Token 等敏感信息
- 过于频繁的操作（每秒数千次）
- 正常的读取操作（可用 Debug 级别）

### 4.2 错误处理原则

1. **使用标准错误码**：优先使用预定义的错误码
2. **记录详细日志**：错误发生时记录上下文信息
3. **避免泄露敏感信息**：返回给前端的错误消息不应包含内部实现细节
4. **统一响应格式**：使用 `ErrorWithCode` 确保响应格式一致

### 4.3 性能优化建议

1. **生产环境使用 Info 级别**：Debug 日志开销较大
2. **避免过度日志**：高频操作使用采样记录
3. **异步日志**：zap 已默认使用异步写入
4. **合理轮转**：根据磁盘空间调整保留策略

---

## 五、配置调整

### 5.1 修改日志级别

编辑 `cmd/main.go`：

```go
logConfig := &utils.LogConfig{
    Level: "debug",  // 改为 debug 查看详细日志
    // ...
}
```

### 5.2 关闭控制台输出

```go
logConfig := &utils.LogConfig{
    Console: false,  // 生产环境可关闭控制台输出
    // ...
}
```

### 5.3 调整文件大小

```go
logConfig := &utils.LogConfig{
    MaxSize: 200,     // 单文件 200MB
    MaxBackups: 30,   // 保留 30 个文件
    MaxAge: 90,       // 保留 90 天
    // ...
}
```

---

## 六、监控与告警

### 6.1 错误日志监控

```bash
# 统计错误数量
grep '"level":"error"' logs/app.log | wc -l

# 按错误码分组统计
grep '"error_code"' logs/app.log | jq '.error_code' | sort | uniq -c

# 查找最频繁的错误
grep '"level":"error"' logs/app.log | jq '.error_msg' | sort | uniq -c | sort -rn
```

### 6.2 性能监控

```bash
# 统计慢操作
grep '"msg":"Slow Operation"' logs/app.log | jq '{operation, duration}'

# 统计接口响应时间
grep '"msg":"HTTP Request"' logs/app.log | jq '{path, latency}' | grep -v '0.0'
```

### 6.3 业务监控

```bash
# 统计今日登录次数
grep '"msg":"Business Event"' logs/app.log | grep '"event":"user_login"' | wc -l

# 查看异常登录
grep '"level":"warn"' logs/app.log | grep 'Authentication Error'
```

---

## 七、常见问题

### Q1: 日志文件过大怎么办？

**A:** 调整 `MaxSize` 参数，或减少 `MaxBackups` 和 `MaxAge`。

### Q2: 如何集成 ELK 或其他日志平台？

**A:** 日志已是 JSON 格式，可直接使用 Filebeat 采集：

```yaml
# filebeat.yml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /path/to/logs/app.log
  json.keys_under_root: true
```

### Q3: 如何添加自定义错误码？

**A:** 在 `utils/errors.go` 中按分类添加：

```go
// 自定义业务错误码（9000-9999）
var (
    ErrCustomBusiness = ErrorCode{Code: 9100, Message: "自定义业务错误"}
)
```

### Q4: 如何在中间件中使用日志？

**A:** 直接调用日志方法：

```go
func LoggerMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        startTime := time.Now()
        ctx.Next()
        
        utils.LogRequest(
            ctx.Request.Method,
            ctx.Request.URL.Path,
            ctx.ClientIP(),
            ctx.Writer.Status(),
            time.Since(startTime),
        )
    }
}
```

---

## 八、迁移指南

### 从标准 log 迁移

**旧代码：**
```go
log.Printf("用户登录: %s", username)
log.Fatalf("数据库连接失败: %v", err)
```

**新代码：**
```go
utils.Infof("用户登录: %s", username)
utils.LogFatal("数据库连接失败", zap.Error(err))
```

### 从自定义错误响应迁移

**旧代码：**
```go
ctx.JSON(http.StatusBadRequest, gin.H{
    "code": 400,
    "message": "用户不存在",
})
```

**新代码：**
```go
utils.ErrorWithCode(ctx, utils.ErrUserNotFound)
```

---

## 总结

✅ **已实现功能：**
- 结构化日志（JSON 格式）
- 日志分级（Debug/Info/Warn/Error/Fatal）
- 日志轮转（大小、时间、数量）
- 统一错误码体系（200+ 错误码）
- 业务日志快捷方法
- 自动日志记录（HTTP 请求、错误、性能）

🎯 **效果：**
- 日志可搜索、可分析
- 错误码规范统一
- 便于问题定位和性能监控
- 支持对接日志平台（ELK、Grafana）
