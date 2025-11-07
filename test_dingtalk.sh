#!/bin/bash
# 钉钉双模式功能测试脚本

BASE_URL="http://localhost:10024/api/v1"
ADMIN_TOKEN="YOUR_ADMIN_TOKEN_HERE"

echo "================================"
echo "钉钉双模式功能测试"
echo "================================"
echo ""

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

function test_api() {
    local name=$1
    local method=$2
    local url=$3
    local data=$4
    
    echo -e "${YELLOW}测试: ${name}${NC}"
    
    if [ -z "$data" ]; then
        response=$(curl -s -X $method "${BASE_URL}${url}" \
            -H "Authorization: Bearer ${ADMIN_TOKEN}")
    else
        response=$(curl -s -X $method "${BASE_URL}${url}" \
            -H "Authorization: Bearer ${ADMIN_TOKEN}" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi
    
    echo "$response" | jq '.'
    
    if echo "$response" | jq -e '.code == 0' > /dev/null; then
        echo -e "${GREEN}✅ 成功${NC}"
    else
        echo -e "${RED}❌ 失败${NC}"
    fi
    echo ""
}

# 1. 创建 Webhook 机器人
echo "================================"
echo "1. 创建 Webhook 机器人"
echo "================================"
test_api "创建 Webhook 机器人" "POST" "/dingtalk-bots" '{
  "name": "测试Webhook机器人",
  "bot_type": "webhook",
  "webhook": "https://oapi.dingtalk.com/robot/send?access_token=test_webhook_token",
  "secret": "SECtest_secret",
  "store_id": 1,
  "msg_type": "markdown",
  "remark": "自动化测试创建"
}'

# 2. 创建 Stream 机器人
echo "================================"
echo "2. 创建 Stream 机器人"
echo "================================"
test_api "创建 Stream 机器人" "POST" "/dingtalk-bots" '{
  "name": "测试Stream机器人",
  "bot_type": "stream",
  "client_id": "dingtest123456",
  "client_secret": "test_app_secret",
  "agent_id": "999999999",
  "store_id": 1,
  "msg_type": "markdown",
  "remark": "自动化测试创建"
}'

# 3. 列出所有机器人
echo "================================"
echo "3. 查询机器人列表"
echo "================================"
test_api "获取机器人列表" "GET" "/dingtalk-bots?page=1&page_size=10"

# 4. 获取机器人详情
echo "================================"
echo "4. 获取机器人详情"
echo "================================"
test_api "获取 Webhook 机器人详情" "GET" "/dingtalk-bots/1"

# 5. 测试机器人连接
echo "================================"
echo "5. 测试机器人连接"
echo "================================"
echo -e "${YELLOW}注意: 测试可能失败因为使用的是测试凭证${NC}"
test_api "测试 Webhook 机器人" "POST" "/dingtalk-bots/1/test"

# 6. 更新机器人
echo "================================"
echo "6. 更新机器人配置"
echo "================================"
test_api "禁用机器人" "PUT" "/dingtalk-bots/1" '{
  "is_enabled": false
}'

sleep 1

test_api "重新启用机器人" "PUT" "/dingtalk-bots/1" '{
  "is_enabled": true
}'

# 7. 切换机器人类型
echo "================================"
echo "7. 切换机器人类型"
echo "================================"
echo -e "${YELLOW}注意: 实际环境需要提供真实的凭证${NC}"
test_api "从 Webhook 切换到 Stream" "PUT" "/dingtalk-bots/1" '{
  "bot_type": "stream",
  "client_id": "dingtest789",
  "client_secret": "new_secret",
  "agent_id": "888888888"
}'

sleep 1

test_api "从 Stream 切换回 Webhook" "PUT" "/dingtalk-bots/1" '{
  "bot_type": "webhook",
  "webhook": "https://oapi.dingtalk.com/robot/send?access_token=new_token"
}'

# 8. 删除机器人
echo "================================"
echo "8. 删除机器人"
echo "================================"
echo -e "${YELLOW}提示: 删除前请确认机器人ID${NC}"
read -p "是否继续删除测试机器人? (y/n): " confirm

if [ "$confirm" == "y" ]; then
    test_api "删除第一个机器人" "DELETE" "/dingtalk-bots/1"
    test_api "删除第二个机器人" "DELETE" "/dingtalk-bots/2"
    echo -e "${GREEN}清理完成${NC}"
else
    echo -e "${YELLOW}跳过删除${NC}"
fi

echo ""
echo "================================"
echo "测试完成"
echo "================================"
echo ""
echo "下一步:"
echo "1. 查看日志: tail -f logs/app.log | grep -i dingtalk"
echo "2. 使用真实凭证创建机器人"
echo "3. 创建报菜记录测试自动通知"
echo ""
