# 报菜记录与钉钉推送功能文档

## 功能概述

本模块实现了完整的报菜记录 CRUD 功能,并通过事件驱动架构自动推送消息到钉钉机器人。

## 架构设计

### 1. 事件总线 (Event Bus)
- 位置: `utils/event_bus.go`
- 功能: 实现发布订阅模式,解耦报菜模块与钉钉推送模块
- 特性:
  - 同步发布 (`Publish`)
  - 异步发布 (`PublishAsync`)
  - 支持多订阅者
  - 线程安全

### 2. 钉钉机器人模块
- **数据模型**: `model/dingtalk_bot.go`
  - 支持多个机器人配置
  - 门店级别或全局机器人
  - 启用/禁用状态控制
  - 支持签名密钥安全认证
  
- **持久层**: `module/dingtalk_bot.go`
  - CRUD 操作
  - 按门店查询启用的机器人
  - Webhook 唯一性验证

- **服务层**: `service/dingtalk.go`
  - 文本消息推送
  - Markdown 消息推送
  - 门店广播功能
  - 签名生成

- **控制器**: `controller/dingtalk_bot.go`
  - 机器人配置管理 (仅管理员)
  - 测试连接功能

### 3. 报菜事件监听器
- 位置: `service/menu_report_listener.go`
- 功能: 
  - 订阅 `menu_report.created` 事件
  - 自动构建推送消息
  - 调用钉钉服务广播通知

### 4. 报菜服务增强
- 创建报菜记录时异步发布事件
- 事件包含完整信息:
  - 报菜记录详情
  - 门店名称
  - 菜品名称
  - 操作人员姓名

## API 接口

### 钉钉机器人管理 (管理员专用)

#### 1. 创建机器人配置
```
POST /api/v1/dingtalk-bots
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "name": "测试门店机器人",
  "webhook": "https://oapi.dingtalk.com/robot/send?access_token=xxxxx",
  "secret": "SEC***",  // 可选,签名密钥
  "store_id": 999,     // 可选,null表示全局
  "is_enabled": true,
  "msg_type": "markdown",  // 或 "text"
  "remark": "测试用机器人"
}
```

#### 2. 获取机器人列表
```
GET /api/v1/dingtalk-bots?page=1&page_size=10
Authorization: Bearer <admin_token>
```

#### 3. 更新机器人配置
```
PUT /api/v1/dingtalk-bots/:id
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "is_enabled": false
}
```

#### 4. 删除机器人配置
```
DELETE /api/v1/dingtalk-bots/:id
Authorization: Bearer <admin_token>
```

#### 5. 测试机器人连接
```
POST /api/v1/dingtalk-bots/:id/test
Authorization: Bearer <admin_token>
```

### 报菜记录管理

#### 1. 创建报菜记录 (触发钉钉推送)
```
POST /api/v1/menu-reports
Authorization: Bearer <token>
Content-Type: application/json

{
  "dish_id": 1,
  "quantity": 50,
  "remark": "今日特价菜品"
}
```

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 123,
    "store_id": 999,
    "dish_id": 1,
    "user_id": 10,
    "quantity": 50,
    "remark": "今日特价菜品",
    "created_at": "2025-11-06T17:00:00+08:00",
    "updated_at": "2025-11-06T17:00:00+08:00"
  }
}
```

**自动钉钉推送消息格式:**
```markdown
## 📋 新报菜通知

**菜品名称:** 红烧肉  
**报菜数量:** 50  
**门店名称:** 测试门店  
**操作人员:** 张三  
**报菜时间:** 2025-11-06 17:00:00  
**备注:** 今日特价菜品  

---
*报菜记录ID: 123*
```

#### 2. 查询报菜记录列表
```
GET /api/v1/menu-reports?page=1&page_size=10
Authorization: Bearer <token>
```

#### 3. 按日期范围查询
```
GET /api/v1/menu-reports?start_date=2025-11-01&end_date=2025-11-06
Authorization: Bearer <token>
```

#### 4. 获取单条报菜记录
```
GET /api/v1/menu-reports/:id
Authorization: Bearer <token>
```

#### 5. 更新报菜记录
```
PUT /api/v1/menu-reports/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "quantity": 60,
  "remark": "调整数量"
}
```

#### 6. 删除报菜记录
```
DELETE /api/v1/menu-reports/:id
Authorization: Bearer <token>
```

#### 7. 获取统计数据
```
GET /api/v1/menu-reports/statistics?start_date=2025-11-01&end_date=2025-11-06
Authorization: Bearer <token>
```

## 数据库变更

### 新增表: `ding_talk_bots`
```sql
CREATE TABLE `ding_talk_bots` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '机器人名称',
  `webhook` varchar(500) NOT NULL COMMENT 'Webhook地址',
  `secret` varchar(500) DEFAULT NULL COMMENT '签名密钥',
  `store_id` bigint unsigned DEFAULT NULL COMMENT '所属门店ID',
  `is_enabled` tinyint(1) DEFAULT 1 COMMENT '是否启用',
  `msg_type` varchar(20) DEFAULT 'markdown' COMMENT '消息类型',
  `remark` text COMMENT '备注',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `webhook` (`webhook`),
  KEY `idx_store_id` (`store_id`),
  KEY `idx_is_enabled` (`is_enabled`)
);
```

## 使用流程

### 1. 配置钉钉机器人 (管理员操作)

1. 在钉钉群里创建自定义机器人,获取 Webhook 地址和密钥
2. 调用创建机器人配置 API,将配置保存到数据库
3. 使用测试接口验证机器人连接是否正常

### 2. 创建报菜记录 (门店操作)

1. 调用创建报菜记录 API
2. 系统自动:
   - 保存报菜记录到数据库
   - 发布 `menu_report.created` 事件
   - 事件监听器接收事件
   - 查询该门店的所有启用机器人
   - 向所有机器人推送通知消息

### 3. 查看推送结果

- 钉钉群中会收到格式化的报菜通知
- 日志中会记录推送成功或失败信息

## 错误处理

- 机器人配置错误: 日志记录但不影响报菜创建
- Webhook 重复: 返回 409 Conflict
- 权限不足: 返回 403 Forbidden
- 资源不存在: 返回 404 Not Found

## 扩展性

### 1. 添加更多事件监听器
```go
eventBus.Subscribe("menu_report.created", anotherListener.OnMenuReportCreated)
```

### 2. 支持更多消息类型
在 `model/dingtalk_bot.go` 中添加新的消息结构体,如 Link、FeedCard 等

### 3. 添加更多事件类型
```go
// 报菜更新事件
type MenuReportUpdatedEvent struct {
    Report    *model.MenuReport
    OldQuantity int
    NewQuantity int
}

func (e MenuReportUpdatedEvent) Name() string {
    return "menu_report.updated"
}
```

## 注意事项

1. **异步推送**: 钉钉推送是异步的,不会阻塞报菜创建流程
2. **权限控制**: 钉钉机器人配置管理仅限管理员
3. **门店隔离**: 报菜记录自动关联当前用户门店
4. **多机器人**: 支持同时向多个机器人推送(门店级+全局)
5. **签名安全**: 建议在生产环境使用签名密钥

## 测试建议

1. 先在测试群创建测试机器人
2. 使用测试接口验证连接
3. 创建报菜记录验证自动推送
4. 测试禁用机器人后不再推送
5. 验证多门店机器人隔离
