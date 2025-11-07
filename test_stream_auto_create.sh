#!/bin/bash
# 测试 Stream 机器人自动创建

echo "=== 测试 Stream 机器人自动创建 ==="
echo ""

echo "1. 检查配置文件..."
if grep -q "client_id: dinglsqvmdgbj5xaidwz" config/config.yaml; then
    echo "✅ config.yaml 配置正确"
else
    echo "❌ config.yaml 缺少 Stream 配置"
    exit 1
fi

echo ""
echo "2. 启动服务器(请在另一个终端运行)..."
echo "   命令: go run cmd/main.go"
echo ""
echo "3. 查看日志,应该看到:"
echo "   - Stream bot created from config"
echo "   - Stream bot started successfully"
echo ""
echo "4. 验证数据库记录:"
echo "   SELECT * FROM dingtalk_bots WHERE bot_type='stream';"
echo ""
echo "5. 钉钉后台配置事件订阅:"
echo "   - 访问: https://open-dev.dingtalk.com"
echo "   - 进入应用 → 事件与回调 → 事件订阅管理"
echo "   - 推送方式: Stream模式推送"
echo "   - 订阅: 机器人接收消息"
echo "   - 点击: 已完成接入,验证连接通道"
echo ""
echo "=== 配置完成后,Stream 机器人将自动连接! ==="
