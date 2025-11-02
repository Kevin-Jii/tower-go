## WebSocket 会话与多端踢出

### 连接地址
```
ws://<host>:<port>/ws?device_id=<可选设备ID>
```
Header 需包含:
```
Authorization: Bearer <登录获得的JWT>
```

### 消息协议
统一 JSON：
```
{
  "type": "connected|kick|ping|pong|echo|notice",
  "payload": {},
  "ts": 1730570000 // Unix 秒
}
```

服务端发送的主要消息：
- connected: 首次连接成功，payload.session_id
- kick: 被踢下线，payload.reason
- pong: 心跳回应
- notice: 系统广播（预留）

客户端可发送：
- {"type":"ping"}
- 其他结构会被 echo 回显（调试用）

### 会话策略
当前初始化为 `single` 单点登录：
1. 新建 WebSocket 连接时，会踢出该用户的旧连接。
2. 登录（HTTP /auth/login）成功后，也会踢出现有 WebSocket 会话。

可切换到 `multi` 允许多会话并限制最大数量（代码中预留实现，但未在 main 中启用）。

### 登录即时踢出
`/auth/login` 在 single 策略下：
1. 生成新 token
2. 检测该用户是否有现存会话，有则发送 kick 并关闭连接。

### 前端处理建议
1. 建立连接后监听 kick 事件，收到后：
   - 弹出提示 “账号已在其他设备登录”
   - 清理本地 token，跳转登录页
2. 定期发送 ping（如 30s）保持活跃；若多次未收到 pong 则重连。
3. 断线重连需使用最新 token。

### 示例（JavaScript）
```js
const token = localStorage.getItem('token');
const ws = new WebSocket(`ws://localhost:10024/ws?device_id=web-${Date.now()}`);
ws.onopen = () => {
  ws.send(JSON.stringify({type: 'ping'}));
};
ws.onmessage = (e) => {
  const msg = JSON.parse(e.data);
  switch(msg.type){
    case 'connected': console.log('session', msg.payload.session_id); break;
    case 'pong': break;
    case 'kick': alert('被踢下线: ' + (msg.payload?.reason||'')); ws.close(); break;
  }
};
```

### 后续可扩展
- 会话列表 & 管理端踢人 REST 接口
- Redis 存储/分布式广播
- 消息订阅/房间机制
- 服务端定时清理过期 token 对应会话
