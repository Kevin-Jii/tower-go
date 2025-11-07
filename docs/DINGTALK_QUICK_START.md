# 钉钉机器人双模式快速入门

## 一、选择合适的模式

### Webhook 模式
**适用场景:**
- ✅ 消息发送到群聊
- ✅ 生产环境部署
- ✅ 简单配置

**限制:**
- ❌ 需要公网 IP 和域名 (接收消息)
- ❌ 需要配置 TLS 证书 (接收消息)

### Stream 模式
**适用场景:**
- ✅ 本地开发测试
- ✅ 企业应用通知
- ✅ 无公网 IP/域名
- ✅ 需要接收用户消息

**限制:**
- ❌ 需要企业内部应用权限

## 二、Webhook 模式配置 (5分钟)

### 1. 创建群机器人

1. 打开钉钉群 → **群设置**
2. 点击 **智能群助手**
3. 点击 **添加机器人**
4. 选择 **自定义** → 点击 **添加**
5. 填写机器人名称: `报菜通知`
6. **安全设置** → 选择 **加签**
7. 复制:
   - **Webhook 地址**: `https://oapi.dingtalk.com/robot/send?access_token=xxx`
   - **加签密钥**: `SECxxx`

### 2. 添加到系统

```bash
curl -X POST http://localhost:10024/api/v1/dingtalk-bots \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "报菜通知",
    "bot_type": "webhook",
    "webhook": "https://oapi.dingtalk.com/robot/send?access_token=你的token",
    "secret": "SECxxx你的密钥",
    "store_id": 1,
    "msg_type": "markdown"
  }'
```

### 3. 测试发送

```bash
curl -X POST http://localhost:10024/api/v1/dingtalk-bots/1/test \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

✅ 群内收到测试消息即配置成功!

## 三、Stream 模式配置 (10分钟)

### 1. 创建企业内部应用

1. 访问 [钉钉开发者后台](https://open-dev.dingtalk.com/)
2. 登录企业管理员账号
3. 点击 **应用开发** → **企业内部开发**
4. 点击 **创建应用**
5. 填写应用信息:
   - 应用名称: `报菜通知系统`
   - 应用描述: `自动推送报菜记录`
   - 应用图标: 上传图标
6. 创建成功后,在 **应用信息** 页面获取:
   - **AppKey (ClientID)**: `dingxxxxxx`
   - **AppSecret (ClientSecret)**: `xxx`
   - **AgentId**: `123456789`

### 2. 开通应用权限

1. 点击 **权限管理**
2. 搜索并开通以下权限:
   - ✅ **企业会话消息发送** (`qyapi_robot_sendmsg`)
   - ✅ **通讯录只读权限** (`Contact.User.Read`)
3. 点击 **保存**

### 3. 发布应用

1. 点击 **版本管理与发布**
2. 点击 **创建新版本**
3. 填写版本说明
4. 点击 **确认发布**
5. 等待审核通过 (一般几分钟)

### 4. 添加到系统

```bash
curl -X POST http://localhost:10024/api/v1/dingtalk-bots \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "企业应用通知",
    "bot_type": "stream",
    "client_id": "dingxxxxxx",
    "client_secret": "你的AppSecret",
    "agent_id": "123456789",
    "store_id": 1,
    "msg_type": "markdown"
  }'
```

### 5. 验证连接

创建成功后,查看日志:

```bash
tail -f logs/app.log | grep -i stream
```

看到以下日志表示连接成功:
```
✅ Stream bot started successfully | botID=2 | botName=企业应用通知
```

### 6. 测试发送

```bash
curl -X POST http://localhost:10024/api/v1/dingtalk-bots/2/test \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

✅ 工作通知收到测试消息即配置成功!

## 四、触发自动通知

### 创建报菜记录

```bash
curl -X POST http://localhost:10024/api/v1/menu-reports \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "dish_id": 1,
    "quantity": 50,
    "unit": "份",
    "remark": "晚市使用"
  }'
```

**自动触发流程:**
```
创建报菜记录
    ↓
发布 menu_report.created 事件
    ↓
MenuReportListener 监听到事件
    ↓
BroadcastToStore 发送通知
    ↓
根据 bot_type 选择发送方式
    ↓
Webhook → 群消息
Stream  → 工作通知
```

### 通知消息示例

```markdown
📝 报菜记录通知

菜品: 宫保鸡丁
数量: 50 份
门店: 总店
操作人: 张三
时间: 2025-11-06 14:30:15
```

## 五、常见问题

### Q1: Webhook 签名失败?
**A:** 检查 Secret 是否正确复制,确保服务器时间同步。

### Q2: Stream 连接失败?
**A:** 检查:
1. ClientID/ClientSecret 是否正确
2. 应用权限是否开通
3. 网络是否可访问 `api.dingtalk.com`

### Q3: 没有收到通知?
**A:** 检查:
1. 机器人是否启用: `"is_enabled": true`
2. StoreID 是否匹配
3. 查看日志: `tail -f logs/app.log`

### Q4: 如何禁用机器人?
**A:** 
```bash
curl -X PUT http://localhost:10024/api/v1/dingtalk-bots/1 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"is_enabled": false}'
```

Stream 机器人会自动断开连接。

### Q5: 如何切换模式?
**A:**
```bash
# 从 Webhook 切换到 Stream
curl -X PUT http://localhost:10024/api/v1/dingtalk-bots/1 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "bot_type": "stream",
    "client_id": "dingxxxxxx",
    "client_secret": "xxx",
    "agent_id": "123456789"
  }'
```

系统会自动处理连接管理。

## 六、管理多个机器人

### 不同门店配置不同机器人

```bash
# 门店1 - Webhook
curl -X POST http://localhost:10024/api/v1/dingtalk-bots \
  -d '{"name":"门店1通知","bot_type":"webhook","store_id":1,...}'

# 门店2 - Stream
curl -X POST http://localhost:10024/api/v1/dingtalk-bots \
  -d '{"name":"门店2通知","bot_type":"stream","store_id":2,...}'

# 门店3 - 双模式
curl -X POST http://localhost:10024/api/v1/dingtalk-bots \
  -d '{"name":"门店3-群通知","bot_type":"webhook","store_id":3,...}'
  
curl -X POST http://localhost:10024/api/v1/dingtalk-bots \
  -d '{"name":"门店3-应用通知","bot_type":"stream","store_id":3,...}'
```

### 全局机器人

不设置 `store_id` 或设为 `null`,机器人会接收所有门店的通知:

```bash
curl -X POST http://localhost:10024/api/v1/dingtalk-bots \
  -d '{"name":"全局通知","bot_type":"stream",...}'
```

## 七、监控与维护

### 查看机器人列表

```bash
curl http://localhost:10024/api/v1/dingtalk-bots?page=1&page_size=10 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### 查看连接状态

```bash
# 方式1: 查看日志
tail -f logs/app.log | grep -E "(Stream|DingTalk)"

# 方式2: 测试连接
curl -X POST http://localhost:10024/api/v1/dingtalk-bots/{id}/test \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### 重启 Stream 连接

```bash
# 禁用再启用
curl -X PUT http://localhost:10024/api/v1/dingtalk-bots/2 \
  -d '{"is_enabled": false}'

curl -X PUT http://localhost:10024/api/v1/dingtalk-bots/2 \
  -d '{"is_enabled": true}'
```

## 八、最佳实践

### 1. 开发环境使用 Stream
- ✅ 无需内网穿透
- ✅ 快速调试

### 2. 生产环境推荐 Webhook
- ✅ 简单稳定
- ✅ 群聊通知直观

### 3. 关键通知使用双模式
- ✅ Webhook → 群聊提醒
- ✅ Stream → 工作通知 (带@)

### 4. 定期测试连接
```bash
# 每日健康检查脚本
#!/bin/bash
for id in 1 2 3; do
  curl -X POST http://localhost:10024/api/v1/dingtalk-bots/$id/test
  sleep 1
done
```

## 九、下一步

- 📚 查看完整文档: [DINGTALK_STREAM_MODE.md](./DINGTALK_STREAM_MODE.md)
- 🔧 自定义消息模板
- 📊 集成监控告警
- 🤖 实现机器人自动回复

---

**遇到问题?** 查看日志: `tail -f logs/app.log | grep -i dingtalk`
