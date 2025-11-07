# 钉钉 Stream 模式集成文档

## 实现状态

### ✅ 已完成

#### 1. 数据模型更新 (`model/dingtalk_bot.go`)
- ✅ 添加 `bot_type` 字段 (webhook/stream)
- ✅ 添加 Stream 模式所需字段:
  - `client_id`: AppKey/SuiteKey
  - `client_secret`: AppSecret/SuiteSecret  
  - `agent_id`: 应用 AgentId (用于消息推送)
- ✅ 更新创建和更新请求模型
- ✅ Webhook 字段改为可选 (Stream 模式不需要)

#### 2. Stream 客户端管理器 (`service/dingtalk_stream_client.go`)
- ✅ 创建 `DingTalkStreamClient` 管理多个 Stream 连接
- ✅ 实现单例模式 (`GetStreamClient`)
- ✅ 支持启动/停止单个机器人连接
- ✅ 支持停止所有连接
- ✅ 线程安全的连接管理
- ✅ 集成钉钉官方 Stream SDK

#### 3. Stream 模式消息推送 (`service/dingtalk.go`)
- ✅ `getStreamAccessToken`: 获取钉钉 access_token
- ✅ `sendStreamText`: Stream 模式发送文本消息
- ✅ `sendStreamMarkdown`: Stream 模式发送 Markdown 消息
- ✅ `sendStreamMessage`: 调用钉钉服务端 API 发送消息

#### 4. 双模式支持
- ✅ `BroadcastToStore` 根据 `bot_type` 自动选择发送方式
- ✅ Webhook 模式: 直接 HTTP POST 到群机器人
- ✅ Stream 模式: 通过钉钉 API 发送企业应用消息

#### 5. 动态连接管理
- ✅ 创建 Stream 机器人时自动启动连接
- ✅ 更新机器人状态时动态启停连接
- ✅ 删除机器人时自动停止连接
- ✅ 切换机器人类型时自动处理连接

#### 6. Bootstrap 集成 (`bootstrap/stream.go`)
- ✅ `InitStreamClients`: 启动时自动连接所有启用的 Stream 机器人
- ✅ `CloseStreamClients`: 关闭时优雅停止所有连接
- ✅ 集成到 `bootstrap/app.go` 启动流程

#### 7. 数据库查询支持
- ✅ `ListEnabledStreamBots`: 查询所有启用的 Stream 机器人

## 核心功能特性

### 消息推送双模式

| 特性 | Webhook 模式 | Stream 模式 |
|-----|------------|------------|
| **实现方式** | 直接 HTTP POST | 钉钉服务端 API |
| **目标** | 群聊机器人 | 企业应用 |
| **认证方式** | access_token + 签名 | AppKey/AppSecret → access_token |
| **消息类型** | text, markdown | text, markdown |
| **@成员** | 支持 @手机号 | 支持 @userid |
| **API 端点** | `oapi.dingtalk.com/robot/send` | `oapi.dingtalk.com/.../asyncsend_v2` |

### 自动连接管理

```go
// 创建 Stream 机器人 → 自动启动 WebSocket 连接
POST /api/v1/dingtalk-bots
{
  "bot_type": "stream",
  "is_enabled": true
}

// 禁用机器人 → 自动停止连接
PUT /api/v1/dingtalk-bots/1
{
  "is_enabled": false
}

// 切换类型 → 自动处理连接
PUT /api/v1/dingtalk-bots/1
{
  "bot_type": "webhook"
}

// 删除机器人 → 自动清理连接
DELETE /api/v1/dingtalk-bots/1
```

### 启动时自动初始化

```go
// bootstrap/app.go
func Run() {
    // ... 初始化数据库、Redis 等
    
    controllers := BuildControllers()
    
    // 自动连接所有启用的 Stream 机器人
    InitStreamClients(controllers.DingTalkBotModule)
    defer CloseStreamClients()
    
    // 启动服务器
    r.Run(addr)
}
```

## API 使用示例

### 创建 Webhook 机器人

```bash
curl -X POST http://localhost:10024/api/v1/dingtalk-bots \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "测试门店通知",
    "bot_type": "webhook",
    "webhook": "https://oapi.dingtalk.com/robot/send?access_token=xxx",
    "secret": "SECxxx",
    "store_id": 1,
    "msg_type": "markdown"
  }'
```

**响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "name": "测试门店通知",
    "bot_type": "webhook",
    "webhook": "https://oapi.dingtalk.com/robot/send?access_token=xxx",
    "is_enabled": true,
    "msg_type": "markdown"
  }
}
```

### 创建 Stream 机器人

```bash
curl -X POST http://localhost:10024/api/v1/dingtalk-bots \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "企业应用通知",
    "bot_type": "stream",
    "client_id": "dingxxxxxx",
    "client_secret": "your_app_secret",
    "agent_id": "123456789",
    "store_id": 1,
    "msg_type": "markdown"
  }'
```

**响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 2,
    "name": "企业应用通知",
    "bot_type": "stream",
    "client_id": "dingxxxxxx",
    "agent_id": "123456789",
    "is_enabled": true,
    "msg_type": "markdown"
  }
}
```

**说明:** Stream 机器人创建成功后会自动启动 WebSocket 连接。

### 测试机器人

```bash
# 测试 Webhook 机器人
curl -X POST http://localhost:10024/api/v1/dingtalk-bots/1/test \
  -H "Authorization: Bearer YOUR_TOKEN"

# 测试 Stream 机器人
curl -X POST http://localhost:10024/api/v1/dingtalk-bots/2/test \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 更新机器人状态

```bash
# 禁用机器人 (Stream 类型会自动停止连接)
curl -X PUT http://localhost:10024/api/v1/dingtalk-bots/2 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "is_enabled": false
  }'

# 启用机器人 (Stream 类型会自动启动连接)
curl -X PUT http://localhost:10024/api/v1/dingtalk-bots/2 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "is_enabled": true
  }'
```

### 切换机器人类型

```bash
# 从 Stream 切换到 Webhook
curl -X PUT http://localhost:10024/api/v1/dingtalk-bots/2 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "bot_type": "webhook",
    "webhook": "https://oapi.dingtalk.com/robot/send?access_token=xxx"
  }'
```

### 触发自动通知

```bash
# 创建报菜记录,自动触发钉钉推送
curl -X POST http://localhost:10024/api/v1/menu-reports \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "dish_id": 1,
    "quantity": 50,
    "unit": "份"
  }'
```

## 配置说明

### Webhook 模式配置

1. **创建群机器人**
   - 钉钉群 → 群设置 → 智能群助手 → 添加机器人
   - 选择"自定义"机器人
   - 获取 Webhook 地址: `https://oapi.dingtalk.com/robot/send?access_token=xxx`

2. **安全设置 (可选)**
   - 开启加签,获取 Secret: `SECxxx`
   - 设置关键词或 IP 白名单

3. **配置到系统**
   ```json
   {
     "bot_type": "webhook",
     "webhook": "完整的 Webhook URL",
     "secret": "加签密钥 (可选)"
   }
   ```

### Stream 模式配置

1. **创建企业内部应用**
   - 钉钉开发者后台 → 应用开发 → 企业内部开发
   - 创建应用,获取:
     - **AppKey (ClientID)**: `dingxxxxxx`
     - **AppSecret (ClientSecret)**: `xxx`
     - **AgentId**: `123456789`

2. **开通权限**
   - 机器人接收消息权限: `qyapi_chat_manage`
   - 发送消息到企业群权限: `qyapi_robot_sendmsg`

3. **配置到系统**
   ```json
   {
     "bot_type": "stream",
     "client_id": "AppKey",
     "client_secret": "AppSecret",
     "agent_id": "AgentId"
   }
   ```

## 工作原理

### Webhook 模式流程

```
报菜记录创建
    ↓
发布事件 (menu_report.created)
    ↓
MenuReportListener 监听
    ↓
BroadcastToStore (bot_type=webhook)
    ↓
直接 HTTP POST 到群机器人
    ↓
群内显示通知
```

### Stream 模式流程

```
启动服务
    ↓
InitStreamClients
    ↓
启动 WebSocket 连接到钉钉平台
    ↓
报菜记录创建
    ↓
发布事件 (menu_report.created)
    ↓
MenuReportListener 监听
    ↓
BroadcastToStore (bot_type=stream)
    ↓
获取 access_token
    ↓
调用钉钉服务端 API 发送消息
    ↓
企业应用推送通知
```

## 数据库迁移

系统会自动迁移 `ding_talk_bots` 表结构:

```sql
ALTER TABLE ding_talk_bots 
  ADD COLUMN bot_type VARCHAR(20) DEFAULT 'webhook',
  ADD COLUMN client_id VARCHAR(200),
  ADD COLUMN client_secret VARCHAR(500),
  ADD COLUMN agent_id VARCHAR(50),
  MODIFY webhook VARCHAR(500) NULL;
```

## 日志监控

### 查看 Stream 连接状态

```bash
# 查看日志
tail -f logs/app.log | grep -i stream

# 关键日志
# ✅ Stream bot started successfully
# ❌ Stream client start failed
# 📊 Stream clients initialized: total=2, success=2
```

### 连接健康检查

```go
streamClient := service.GetStreamClient()

// 检查是否运行
isRunning := streamClient.IsRunning()

// 获取连接数量
botCount := streamClient.GetBotCount()

// 检查特定机器人
client, exists := streamClient.GetClient(botID)
```

## 故障排查

### Stream 连接失败

**症状:** 日志显示 `Stream client start failed`

**可能原因:**
1. ClientID/ClientSecret 错误
2. 网络无法访问钉钉服务器
3. 应用权限未开通

**解决方案:**
```bash
# 1. 验证凭证
curl https://oapi.dingtalk.com/gettoken \
  -d "appkey=YOUR_CLIENT_ID&appsecret=YOUR_CLIENT_SECRET"

# 2. 检查网络连接
ping api.dingtalk.com

# 3. 查看应用权限
# 开发者后台 → 应用详情 → 权限管理
```

### 消息发送失败

**症状:** 日志显示 `dingtalk api error: code=xxx`

**常见错误码:**
- `40001`: access_token 过期或无效
- `40003`: 权限不足
- `60011`: 不在群内或群不存在

**解决方案:**
```bash
# 测试机器人
curl -X POST http://localhost:10024/api/v1/dingtalk-bots/{id}/test

# 查看详细错误
tail -f logs/app.log | grep -A 5 "Failed to send"
```

### Webhook 签名失败

**症状:** 钉钉返回 `sign not match`

**原因:** Secret 配置错误或时间戳不同步

**解决方案:**
1. 检查 Secret 是否正确
2. 确认服务器时间同步: `ntpdate ntp.aliyun.com`

## 性能优化建议

### 1. 连接池管理
当前实现为每个机器人一个连接,适合中小规模场景。

### 2. 消息队列
高并发场景建议添加消息队列:
```go
// 异步发送消息
eventBus.PublishAsync("dingtalk.send", messageData)
```

### 3. 失败重试
添加指数退避重试机制:
```go
for i := 0; i < maxRetries; i++ {
    if err := sendMessage(); err == nil {
        break
    }
    time.Sleep(time.Second * time.Duration(math.Pow(2, float64(i))))
}
```

### 4. 监控告警
集成 Prometheus 监控:
```go
// 消息发送成功率
dingtalk_message_sent_total{type="webhook|stream",status="success|failure"}

// 连接数量
dingtalk_stream_connections_total
```

## 下一步增强建议

### 优先级 1: 接收机器人消息
实现 Stream 模式下接收和处理用户消息。

### 优先级 2: 健康检查 API
```bash
GET /api/v1/dingtalk-bots/{id}/status
{
  "bot_id": 1,
  "bot_type": "stream",
  "is_enabled": true,
  "connection_status": "connected",
  "last_message_time": "2025-11-06T10:30:00Z"
}
```

### 优先级 3: 消息模板
支持自定义消息模板:
```go
type MessageTemplate struct {
    ID       uint
    Name     string
    Type     string // text, markdown
    Template string // 支持变量替换
}
```

### 优先级 4: 批量发送
支持批量发送优化:
```go
func (s *DingTalkService) BatchSend(bots []*model.DingTalkBot, message Message) error
```

## 参考文档

- [Stream 模式介绍](https://open.dingtalk.com/document/development/introduction-to-stream-mode)
- [机器人开发文档](https://open.dingtalk.com/document/orgapp/robot-overview)
- [Go SDK GitHub](https://github.com/open-dingtalk/dingtalk-stream-sdk-go)
- [发送企业消息 API](https://open.dingtalk.com/document/orgapp/asynchronous-sending-of-single-chat-messages-by-robots)
- [获取 access_token](https://open.dingtalk.com/document/orgapp/obtain-orgapp-token)

## 总结

✅ **已实现功能:**
- Webhook 和 Stream 双模式支持
- 自动连接管理 (启动/停止/重连)
- 双模式消息推送 (text/markdown)
- Bootstrap 自动初始化
- 事件驱动架构集成
- 完整的 CRUD API

🎯 **核心优势:**
- 零配置基础设施 (Stream 模式)
- 自动化连接管理
- 统一的消息推送接口
- 灵活的模式切换
- 完善的错误处理和日志
