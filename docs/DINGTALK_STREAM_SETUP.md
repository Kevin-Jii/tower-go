# 钉钉 Stream 模式配置指南

## 问题说明

如果看到错误 **"Stream模式连接失败,请参考Stream模式SDK接入指南"**,说明需要在钉钉开放平台配置事件订阅。

## 解决方案

### 方案1: 配置事件订阅(推荐用于接收消息)

如果你需要**接收**钉钉群里的消息(用户@机器人时触发),需要配置事件订阅:

1. **登录钉钉开放平台**
   - 访问: https://open-dev.dingtalk.com
   - 进入你的应用管理

2. **配置事件订阅**
   - 进入 "应用功能" → "事件与回调"
   - 选择 "事件订阅管理"
   - 推送方式选择: **Stream模式推送**
   - 点击 "已完成接入,验证连接通道"

3. **订阅具体事件**
   - 进入 "机器人与消息推送"
   - 订阅 "群会话" 相关事件
   - 保存配置

### 方案2: 仅使用 HTTP API 推送(推荐,无需订阅)

如果你**只需要推送消息到群**,不需要接收用户消息,可以:

1. **不启动 Stream 连接**
   - 创建机器人时选择 `bot_type: "webhook"` 模式
   - 或者不配置 config.yaml 中的 dingtalk.stream 部分

2. **使用 API 推送**
   ```bash
   # 创建 webhook 机器人
   POST /api/v1/dingtalk/robots
   {
     "name": "报菜通知",
     "bot_type": "webhook",
     "webhook": "https://oapi.dingtalk.com/robot/send?access_token=xxx",
     "secret": "SECxxx",
     "is_enabled": true
   }
   ```

3. **获取 Webhook 地址**
   - 钉钉群 → 群设置 → 智能群助手 → 添加机器人
   - 选择"自定义"
   - 复制 Webhook 地址和密钥

## Stream 模式的使用场景

### 适合 Stream 模式:
- ✅ 需要接收用户在群里@机器人的消息
- ✅ 需要实现交互式对话
- ✅ 需要处理群事件(成员加入/退出等)

### 适合 Webhook 模式:
- ✅ 只推送通知消息到群(报菜通知就是这种)
- ✅ 简单易配置,无需后台订阅
- ✅ 无需维持长连接

## 当前系统推荐配置

**对于报菜通知系统,建议使用 Webhook 模式**:

1. 每个门店创建一个 Webhook 机器人
2. 配置对应群的 webhook 地址
3. 报菜时自动推送到对应群

**优势:**
- ✅ 配置简单,5分钟搞定
- ✅ 无需后台订阅事件
- ✅ 无需维护 WebSocket 连接
- ✅ 稳定可靠

## Stream 模式完整配置步骤

如果确实需要 Stream 模式:

### 1. 配置 config.yaml

```yaml
dingtalk:
  stream:
    client_id: dinglsqvmdgbj5xaidwz
    client_secret: sMU4Tr0r4L37bkas0khmyTJ7fuNUYqSXXjH0K8VS7naHazhjAOKE7WCmvFedwMuH
    agent_id: "2974263828"
```

### 2. 钉钉后台配置

1. 登录 https://open-dev.dingtalk.com
2. 进入应用 → 事件与回调 → 事件订阅管理
3. 推送方式: **Stream模式推送**
4. 点击 "已完成接入,验证连接通道"
5. 订阅 "机器人与消息推送" → "群会话"

### 3. 创建 Stream 机器人

```bash
POST /api/v1/dingtalk/robots
{
  "name": "全局Stream机器人",
  "bot_type": "stream",
  "client_id": "dinglsqvmdgbj5xaidwz",
  "client_secret": "sMU4Tr0r4L37bkas0khmyTJ7fuNUYqSXXjH0K8VS7naHazhjAOKE7WCmvFedwMuH",
  "agent_id": "2974263828",
  "is_enabled": true
}
```

### 4. 验证连接

服务器启动后查看日志:
```
INFO Stream bot started successfully {"botID": 1, "botName": "全局Stream机器人"}
```

## 常见问题

### Q: 为什么需要配置事件订阅?
A: Stream 模式是双向通信,钉钉需要知道推送哪些事件给你的应用。如果不配置订阅,连接会因为"没有订阅任何事件"而失败。

### Q: 我只想推送消息,不想接收消息,怎么办?
A: 使用 Webhook 模式!简单快捷,5分钟配置完成,无需后台订阅。

### Q: 可以同时使用两种模式吗?
A: 可以!系统支持混合模式:
- 全局用 Stream 机器人(接收消息)
- 各门店用 Webhook 机器人(推送通知)

### Q: Stream 连接失败怎么办?
A: 检查三点:
1. config.yaml 配置是否正确
2. 钉钉后台是否配置了事件订阅(Stream模式)
3. 是否点击了"验证连接通道"按钮

## 参考链接

- 钉钉开放平台: https://open-dev.dingtalk.com
- Stream 模式文档: https://open.dingtalk.com/document/orgapp/stream-mode-overview
- 机器人开发文档: https://open.dingtalk.com/document/robots/robot-overview
