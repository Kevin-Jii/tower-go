-- 创建 dingtalk_bot 表
CREATE TABLE IF NOT EXISTS `tower_ding_talk_bot` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '机器人名称',
  `bot_type` varchar(20) DEFAULT 'webhook' COMMENT '机器人类型: webhook, stream',
  `webhook` varchar(500) DEFAULT NULL COMMENT 'Webhook 地址（webhook 模式）',
  `secret` varchar(500) DEFAULT NULL COMMENT '签名密钥（webhook 模式）',
  `client_id` varchar(200) DEFAULT NULL COMMENT 'AppKey/SuiteKey (stream 模式)',
  `client_secret` varchar(500) DEFAULT NULL COMMENT 'AppSecret/SuiteSecret (stream 模式)',
  `agent_id` varchar(50) DEFAULT NULL COMMENT '应用 AgentId (stream 模式推送消息用)',
  `store_id` int unsigned DEFAULT NULL COMMENT '所属门店（null 表示全局）',
  `is_enabled` tinyint(1) DEFAULT TRUE COMMENT '是否启用',
  `msg_type` varchar(20) DEFAULT 'markdown' COMMENT '消息类型: text, markdown',
  `remark` text COMMENT '备注',
  `robot_code` varchar(100) DEFAULT NULL COMMENT '钉钉机器人编码(robotCode)',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_tower_ding_talk_bot_store_id` (`store_id`),
  KEY `idx_tower_ding_talk_bot_is_enabled` (`is_enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='钉钉机器人配置表';

-- 插入一个示例机器人（可选）
INSERT INTO `tower_ding_talk_bot` (`name`, `bot_type`, `is_enabled`) VALUES ('示例机器人', 'webhook', FALSE);
