# Tower-Go 部署指南

## 📦 环境配置

### 前置条件

在开始部署之前，请确保已安装以下软件：

- **Go**: 1.20 或更高版本
- **MySQL**: 8.0 或更高版本
- **Redis**: 5.0 或更高版本
- **Git**: 用于代码克隆

### 运行环境要求

| 资源项 | 开发环境 | 生产环境 |
|--------|---------|---------|
| CPU | 1 核 | 2 核+ |
| 内存 | 1 GB | 2 GB+ |
| 磁盘空间 | 500 MB | 1 GB+ |
| 操作系统 | Windows/Linux/macOS | Linux (推荐 Ubuntu/CentOS) |

## 🔧 配置文件详解

### 配置文件模板

项目使用 `.env` 文件管理配置，模板文件为 `.env.example`：

```env
# ======================================
# Tower Go 应用配置
# ======================================

# 应用配置
APP_NAME=tower-go
APP_PORT=10024

# 数据库配置
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=your_secure_password_here
DB_NAME=tower

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password_here
REDIS_DB=0
REDIS_ENABLED=true

# 钉钉配置
DINGTALK_CLIENT_ID=your_client_id_here
DINGTALK_CLIENT_SECRET=your_client_secret_here
DINGTALK_AGENT_ID=your_agent_id_here
DINGTALK_MINI_APP_ID=your_mini_app_id_here

# JWT配置 - 请使用强密码（至少32位随机字符串）
JWT_SECRET=your_jwt_secret_here_at_least_32_characters

# 日志配置
LOG_LEVEL=info

# ======================================
# 安全说明
# ======================================
# 1. 请修改所有密码和密钥
# 2. 生产环境请使用环境变量或密钥管理服务
# 3. 不要将此文件提交到版本控制
# 4. 定期轮换密码和密钥
# ======================================
```

### 关键配置项说明

#### 应用配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `APP_NAME` | 应用名称 | tower-go |
| `APP_PORT` | 服务监听端口 | 10024 |

#### 数据库配置

| 配置项 | 说明 | 示例 |
|--------|------|------|
| `DB_HOST` | 数据库主机地址 | localhost / 10.0.1.10 |
| `DB_PORT` | 数据库端口 | 3306 |
| `DB_USERNAME` | 数据库用户名 | root |
| `DB_PASSWORD` | 数据库密码 | 强密码 |
| `DB_NAME` | 数据库名称 | tower |

**安全建议**：
- 生产环境使用内网地址，避免暴露在公网
- 为应用创建独立的数据库用户，只授予最小权限
- 使用 12 位以上复杂密码

```sql
-- 创建数据库用户示例（生产环境）
CREATE USER 'tower_app'@'10.0.1.%' IDENTIFIED BY 'YourStrongPassword123!';
GRANT SELECT, INSERT, UPDATE, DELETE ON tower.* TO 'tower_app'@'10.0.1.%';
FLUSH PRIVILEGES;
```

#### Redis 配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `REDIS_HOST` | Redis 服务器地址 | localhost |
| `REDIS_PORT` | Redis 端口 | 6379 |
| `REDIS_PASSWORD` | Redis 密码 | 空 |
| `REDIS_DB` | Redis 数据库编号 | 0 |
| `REDIS_ENABLED` | 是否启用 Redis | true |

**安全建议**：
- 生产环境设置强密码
- 使用内网地址
- 配置 `requirepass` 和 `bind` 参数

#### JWT 配置

| 配置项 | 说明 | 示例 |
|--------|------|------|
| `JWT_SECRET` | JWT 签名密钥 | 32 位以上随机字符串 |

**安全建议：**
- 必须设置 32 位以上的随机字符串
- 生产环境定期更换（需同步更新所有已登录用户的 Token）
- 可以使用工具生成：

```bash
# Linux/macOS
openssl rand -base64 64

# 或
head -c 32 /dev/random | base64
```

#### 钉钉配置（可选）

钉钉集成用于发送通知：

| 配置项 | 说明 |
|--------|------|
| `DINGTALK_CLIENT_ID` | 钉钉应用的 Client ID |
| `DINGTALK_CLIENT_SECRET` | 钉钉应用的 Client Secret |
| `DINGTALK_AGENT_ID` | 应用 Agent ID |

[如何获取钉钉配置](https://open.dingtalk.com/document/orgapp-server/getappinfo)

#### 日志配置

| 配置项 | 说明 | 可选值 |
|--------|------|--------|
| `LOG_LEVEL` | 日志级别 | debug, info, warn, error |

**环境建议**：
- **开发环境**: debug
- **测试环境**: info
- **生产环境**: warn 或 error

## 🚀 部署方式

### 方式一：本地部署（推荐给开发/测试）

#### 1. 克隆项目

```bash
git clone https://github.com/your-org/tower-go.git
cd tower-go
```

#### 2. 安装依赖

```bash
# 安装 Go 依赖
go mod tidy

# 安装 Swagger
# Windows
go install github.com/swaggo/swag/cmd/swag@latest

# Linux/macOS
go install github.com/swaggo/swag/cmd/swag@latest

# 验证安装
swag --version
```

#### 3. 配置环境

```bash
# 复制配置文件
cp .env.example .env

# Windows
copy .env.example .env
```

编辑 `.env` 文件，配置数据库等信息。

#### 4. 创建数据库

```bash
# 使用 MySQL 客户端连接
mysql -u root -p
```

```sql
-- 创建数据库
CREATE DATABASE tower CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 初始化数据（可选）
-- 运行初始化 SQL 脚本
-- source /path/to/init.sql;
```

#### 5. 生成 API 文档

```bash
swag init -g cmd/main.go
```

#### 6. 启动服务

**Windows:**

```bash
# 直接运行
go run cmd/main.go

# 或编译后运行
go build -o tower-go.exe cmd/main.go
./tower-go.exe
```

**Linux/macOS:**

```bash
# 使用 Makefile
make run

# 或直接运行
go run cmd/main.go
```

#### 7. 验证部署

```bash
# 查看服务状态
curl http://localhost:10024/api/v1/stores

# 或使用浏览器访问 Swagger
# http://localhost:10024/api/v1/swagger/index.html
```

### 方式二：编译部署（推荐生产环境）

#### 1. 交叉编译

**编译 Linux 可执行文件（在 Windows 上）**

```bash
# Windows PowerShell
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -ldflags="-s -w" -o tower-go-linux cmd/main.go

# Linux/macOS
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o tower-go-linux cmd/main.go
```

**编译参数说明**

- `GOOS`: 目标操作系统 (linux, windows, darwin)
- `GOARCH`: 目标架构 (amd64, arm64)
- `-ldflags="-s -w"`: 去除调试信息，减小文件体积

**常见编译组合**

| 编译命令 | 说明 |
|---------|------|
| `GOOS=linux GOARCH=amd64` | Linux 64位 |
| `GOOS=linux GOARCH=arm64` | Linux ARM64 (如树莓派) |
| `GOOS=windows GOARCH=amd64` | Windows 64位 |
| `GOOS=darwin GOARCH=amd64` | macOS Intel |
| `GOOS=darwin GOARCH=arm64` | macOS M1/M2 |

#### 2. 上传到服务器

```bash
# 使用 scp 上传
scp tower-go-linux user@your-server:/opt/tower-go/
scp .env user@your-server:/opt/tower-go/

# 或使用 FTP 工具
# FileZilla, WinSCP 等
```

#### 3. 配置系统服务

**Linux systemd 服务**

创建服务文件 `/etc/systemd/system/tower-go.service`:

```ini
[Unit]
Description=Tower Go Application
After=network.target mysqld.service redis.service
Wants=mysqld.service redis.service

[Service]
Type=simple
User=tower
Group=tower

# 工作目录
WorkingDirectory=/opt/tower-go

# 启动命令
ExecStart=/opt/tower-go/tower-go-linux

# 环境变量
Environment="APP_PORT=10024"
Environment="DB_HOST=localhost"

# 重启策略
Restart=on-failure
RestartSec=5s

# 进程数限制
LimitNOFILE=65536

# 环境文件
EnvironmentFile=/opt/tower-go/.env

[Install]
WantedBy=multi-user.target
```

**启动服务**

```bash
# 重新加载 systemd
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start tower-go

# 设置开机自启
sudo systemctl enable tower-go

# 查看服务状态
sudo systemctl status tower-go

# 查看日志
sudo journalctl -u tower-go -f
```

### 方式三：Docker 部署（推荐）

#### 1. 构建 Docker 镜像

**Dockerfile**

```dockerfile
# Build stage
FROM golang:1.20-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 生成 Swagger 文档
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/main.go

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o tower-go cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 创建非 root 用户
RUN addgroup -g 1001 -S tower && \
    adduser -S tower -u 1001 -G tower

# 从 builder 阶段复制文件
COPY --from=builder /app/tower-go .
COPY --from=builder /app/docs ./docs

# 创建日志目录
RUN mkdir -p /app/logs && chown -R tower:tower /app

# 切换用户
USER tower

EXPOSE 10024

CMD ["./tower-go"]
```

**构建镜像**

```bash
# 构建
docker build -t tower-go:latest .

# 查看镜像
docker images | grep tower-go
```

#### 2. 使用 Docker Compose

**docker-compose.yml**

```yaml
version: '3.8'

services:
  # Tower Go 应用
  app:
    build: .
    image: tower-go:latest
    container_name: tower-go-app
    restart: unless-stopped
    ports:
      - "10024:10024"
    environment:
      # 数据库配置
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USERNAME=tower_app
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=tower

      # Redis 配置
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - REDIS_DB=0
      - REDIS_ENABLED=true

      # JWT 配置
      - JWT_SECRET=${JWT_SECRET}

      # 钉钉配置
      - DINGTALK_CLIENT_ID=${DINGTALK_CLIENT_ID}
      - DINGTALK_CLIENT_SECRET=${DINGTALK_CLIENT_SECRET}

      # 日志配置
      - LOG_LEVEL=info
    depends_on:
      - mysql
      - redis
    volumes:
      - ./logs:/app/logs
    networks:
      - tower-network

  # MySQL 数据库
  mysql:
    image: mysql:8.0
    container_name: tower-mysql
    restart: unless-stopped
    environment:
      - MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
      - MYSQL_DATABASE=tower
      - MYSQL_USER=tower_app
      - MYSQL_PASSWORD=${DB_PASSWORD}
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql  # 可选：初始化数据
    command:
      - --character-set-server=utf8mb4
      - --collation-server=utf8mb4_unicode_ci
    networks:
      - tower-network

  # Redis 缓存
  redis:
    image: redis:7-alpine
    container_name: tower-redis
    restart: unless-stopped
    command:
      - redis-server
      - --appendonly yes
      - --requirepass ${REDIS_PASSWORD}
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - tower-network

volumes:
  mysql_data:
    driver: local
  redis_data:
    driver: local

networks:
  tower-network:
    driver: bridge
```

**环境变量文件 `.env`**

```env
# MySQL Root 密码
MYSQL_ROOT_PASSWORD=YourStrongRootPassword123!

# 应用数据库密码
DB_PASSWORD=YourAppPassword123!

# Redis 密码
REDIS_PASSWORD=YourRedisPassword123!

# JWT 密钥
JWT_SECRET=Your32+CharacterRandomStringHere

# 钉钉配置（可选）
DINGTALK_CLIENT_ID=your_client_id
DINGTALK_CLIENT_SECRET=your_client_secret
```

#### 3. 启动服务

```bash
# 1. 创建环境变量文件
cp .env.example .env
# 编辑 .env，填写各密码

# 2. 启动服务
docker-compose up -d

# 3. 查看日志
docker-compose logs -f app

# 4. 验证服务
curl http://localhost:10024/api/v1/stores

# 5. 停止服务
docker-compose down

# 6. 停止并删除数据卷
docker-compose down -v
```

#### 4. 服务管理命令

```bash
# 查看服务状态
docker-compose ps

# 重启应用
docker-compose restart app

# 查看应用日志
docker logs -f tower-go-app

# 进入容器调试
docker exec -it tower-go-app /bin/sh

# 查看 MySQL 日志
docker logs -f tower-mysql

# 查看 Redis 日志
docker logs -f tower-redis
```

### 方式四：Kubernetes 部署（生产推荐）

#### 1. 创建 ConfigMap

**tower-config.yaml**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: tower-config
  namespace: default
data:
  APP_NAME: "tower-go"
  APP_PORT: "10024"
  DB_DRIVER: "mysql"
  DB_HOST: "mysql-service"
  DB_PORT: "3306"
  DB_USERNAME: "tower_app"
  DB_NAME: "tower"
  REDIS_HOST: "redis-service"
  REDIS_PORT: "6379"
  REDIS_DB: "0"
  REDIS_ENABLED: "true"
  LOG_LEVEL: "info"
```

#### 2. 创建 Secret

**tower-secret.yaml**

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: tower-secret
  namespace: default
type: Opaque
data:
  # echo -n "password" | base64
  DB_PASSWORD: cGFzc3dvcmQ=
  REDIS_PASSWORD: cmVkaXNfcGFzc3dvcmQ=
  JWT_SECRET: eW91cl9qd3Rfc2VjcmV0X2tleQ==
  DINGTALK_CLIENT_ID: eW91cl9jbGllbnRfaWQ=
  DINGTALK_CLIENT_SECRET: eW91cl9jbGllbnRfc2VjcmV0
```

应用配置：

```bash
kubectl apply -f tower-config.yaml
kubectl apply -f tower-secret.yaml
```

#### 3. 创建 Deployment

**tower-deployment.yaml**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: tower-go
  namespace: default
  labels:
    app: tower-go
spec:
  replicas: 3  # 生产环境建议 3 个副本
  selector:
    matchLabels:
      app: tower-go
  template:
    metadata:
      labels:
        app: tower-go
    spec:
      containers:
      - name: tower-go
        image: your-registry.com/tower-go:latest
        ports:
        - containerPort: 10024
        envFrom:
        - configMapRef:
            name: tower-config
        - secretRef:
            name: tower-secret
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 10024
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 10024
          initialDelaySeconds: 5
          periodSeconds: 5
```

#### 4. 创建 Service

**tower-service.yaml**

```yaml
apiVersion: v1
kind: Service
metadata:
  name: tower-go-service
  namespace: default
spec:
  selector:
    app: tower-go
  ports:
    - protocol: TCP
      port: 80
      targetPort: 10024
  type: LoadBalancer  # 或使用 NodePort/ClusterIP
```

#### 5. 部署应用

```bash
# 应用所有配置
kubectl apply -f tower-config.yaml
kubectl apply -f tower-secret.yaml
kubectl apply -f tower-deployment.yaml
kubectl apply -f tower-service.yaml

# 查看部署状态
kubectl get deployments
kubectl get pods
kubectl get services

# 查看日志
kubectl logs -f deployment/tower-go

# 伸缩副本数
kubectl scale deployment tower-go --replicas=5
```

## 🏥 健康检查

### HTTP 健康检查接口

```bash
# 访问健康检查接口
curl http://localhost:10024/health

# 期望响应
{
  "status": "ok",
  "timestamp": "2025-11-11T10:00:00+08:00"
}
```

### 依赖服务健康检查

脚本：`health-check.sh`

```bash
#!/bin/bash

# 配置
APP_URL="http://localhost:10024"
DB_HOST="localhost"
DB_PORT="3306"
REDIS_HOST="localhost"
REDIS_PORT="6379"

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "========== Tower-Go 健康检查 =========="

# 1. 检查服务端口
echo -n "检查应用端口... "
if nc -z localhost 10024 2>/dev/null; then
    echo -e "${GREEN}✓ 正常${NC}"
else
    echo -e "${RED}✗ 异常${NC}"
    exit 1
fi

# 2. 检查 HTTP 接口
echo -n "检查 HTTP 接口... "
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" $APP_URL/health)
if [ "$HTTP_STATUS" == "200" ]; then
    echo -e "${GREEN}✓ 正常${NC}"
else
    echo -e "${RED}✗ 异常 (状态码: $HTTP_STATUS)${NC}"
    exit 1
fi

# 3. 检查数据库
echo -n "检查 MySQL... "
if nc -z $DB_HOST $DB_PORT 2>/dev/null; then
    echo -e "${GREEN}✓ 正常{NC}"
else
    echo -e "${RED}✗ 异常{NC}"
    exit 1
fi

# 4. 检查 Redis
echo -n "检查 Redis... "
if nc -z $REDIS_HOST $REDIS_PORT 2>/dev/null; then
    echo -e "${GREEN}✓ 正常${NC}"
else
    echo -e "${RED}✗ 异常${NC}"
    exit 1
fi

echo "======================================"
echo -e "${GREEN}✓ 所有服务正常${NC}"
```

**使用方法**

```bash
# 添加执行权限
chmod +x health-check.sh

# 运行检查
./health-check.sh
```

## 🚨 常见问题

### Q1: 服务启动失败，提示数据库连接错误

**症状**

```
Error 1045: Access denied for user 'root'@'localhost'
```

**解决方案**

1. 检查数据库账号密码是否正确
2. 检查数据库是否运行
3. 检查网络连通性
4. 确认用户是否有远程访问权限

```sql
-- 创建远程访问用户
CREATE USER 'tower_app'@'%' IDENTIFIED BY 'StrongPassword';
GRANT ALL PRIVILEGES ON tower.* TO 'tower_app'@'%';
FLUSH PRIVILEGES;
```

### Q2: JWT Token 认证失败

**症状**

```
Error: signature is invalid
```

**解决方案**

1. 检查 JWT_SECRET 是否配置
2. 确认前后端使用相同的密钥
3. Token 是否过期

### Q3: Redis 连接失败

**症状**

```
Error: connection refused
```

**解决方案**

1. 检查 Redis 服务是否运行
2. 检查密码是否正确
3. 检查 Redis 配置 `bind` 和 `requirepass`

### Q4: 端口被占用

**症状**

```
listen tcp :10024: bind: address already in use
```

**解决方案**

```bash
# 查看占用端口的进程
# Linux
netstat -tulnp | grep 10024

# Windows
netstat -ano | findstr :10024

# macOS
lsof -i :10024

# 终止进程（Linux）
kill -9 <PID>

# 或修改配置
# .env 文件中修改 APP_PORT
```

### Q5: 权限不足

**症状**

```
Permission denied
```

**解决方案**

```bash
# 修改日志目录权限
sudo chown -R tower:tower /path/to/logs
sudo chmod 755 /path/to/logs
```

### Q6: 钉钉通知发送失败

**症状**

```
DingTalk API error: invalid credentials
```

**解决方案**

1. 检查钉钉配置是否正确
2. 确认应用权限是否开启
3. 检查网络是否能访问钉钉 API

## 📊 性能调优

### 1. MySQL 优化

**my.cnf 配置**

```ini
[mysqld]
# 基本配置
character-set-server=utf8mb4
collation-server=utf8mb4_unicode_ci

# 连接数
max_connections = 500
max_connect_errors = 1000

# InnoDB 配置
innodb_buffer_pool_size = 1G
innodb_log_file_size = 256M
innodb_flush_log_at_trx_commit = 2

# 查询缓存
query_cache_type = 1
query_cache_size = 64M
query_cache_limit = 2M

# 慢查询日志
slow_query_log = 1
slow_query_log_file = /var/log/mysql/slow.log
long_query_time = 1
log_queries_not_using_indexes = 1
```

### 2. Redis 优化

**redis.conf 配置**

```conf
# 内存限制
maxmemory 512mb
maxmemory-policy allkeys-lru

# 持久化
save 900 1
save 300 10
save 60 10000

# TCP 连接
tcp-keepalive 300
timeout 300
```

### 3. 系统调优

**Linux 内核参数**

```bash
# 增加文件描述符限制
echo "* soft nofile 65536" >> /etc/security/limits.conf
echo "* hard nofile 65536" >> /etc/security/limits.conf

# 增加端口范围
echo "net.ipv4.ip_local_port_range = 1024 65535" >> /etc/sysctl.conf

# 启用 TCP 快速打开
echo "net.ipv4.tcp_fastopen = 3" >> /etc/sysctl.conf

# 优化 TCP 连接
echo "net.ipv4.tcp_tw_reuse = 1" >> /etc/sysctl.conf
echo "net.ipv4.tcp_fin_timeout = 15" >> /etc/sysctl.conf

# 生效
sysctl -p
```

## 📈 监控告警

### 1. 应用监控

**Prometheus 集成**（推荐）

```go
// 添加 Prometheus 监控
go get github.com/prometheus/client_golang/prometheus/promhttp

// 在路由中添加
router.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

**关键指标**

```
- tower_http_requests_total
- tower_http_request_duration_seconds
- tower_db_query_duration_seconds
- tower_redis_operation_duration_seconds
```

### 2. 日志告警

**ELK Stack 集成**

```bash
# Filebeat 配置
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /opt/tower-go/logs/app.log

output.elasticsearch:
  hosts: ["localhost:9200"]
```

### 3. 告警规则（Prometheus）

**tower-alerts.yml**

```yaml
groups:
- name: tower-app
  rules:
  # 服务不可用
  - alert: TowerAppDown
    expr: up{job="tower-go"} == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "Tower Go 应用服务不可用"
      description: "{{ $labels.instance }} 已经停止运行超过 1 分钟"

  # 错误率过高
  - alert: TowerHighErrorRate
    expr: rate(tower_http_requests_total{status=~"5.."}[5m]) > 0.1
    for: 2m
    labels:
      severity: warning
    annotations:
      summary: "Tower Go 错误率过高"
      description: "最近5分钟错误率超过 10%"

  # 响应时间过长
  - alert: TowerHighLatency
    expr: histogram_quantile(0.95, rate(tower_http_request_duration_seconds_bucket[5m])) > 0.5
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Tower Go 响应时间过长"
      description: "95% 的请求响应时间超过 500ms"
```

## 📋 部署检查清单

### 环境准备

- [ ] 服务器配置满足最低要求
- [ ] Go 环境已安装（1.20+）
- [ ] MySQL 8.0+ 已安装并运行
- [ ] Redis 5.0+ 已安装并运行
- [ ] 防火墙放行端口（10024, 3306, 6379）

### 配置检查

- [ ] `.env` 文件已创建
- [ ] 数据库密码已修改为强密码
- [ ] JWT_SECRET 已设置为 32 位以上随机字符串
- [ ] 钉钉配置已填写（如需要）
- [ ] 日志级别设置为适当级别

### 安全加固

- [ ] MySQL 使用内网地址
- [ ] Redis 设置密码并限制访问
- [ ] 生产环境使用非 root 用户运行
- [ ] 敏感配置使用环境变量或 Secret
- [ ] 启用 HTTPS（生产环境）

### 服务验证

- [ ] 服务成功启动无错误
- [ ] 可以访问 Swagger 文档
- [ ] 可以正常登录
- [ ] 可以创建门店/菜品/报菜
- [ ] 钉钉通知正常（如启用）
- [ ] 日志正常记录

### 监控告警

- [ ] 健康检查接口正常
- [ ] 日志收集已配置
- [ ] 监控已部署（可选）
- [ ] 告警规则已配置（可选）

## 🔄 版本升级

### 升级步骤

1. **备份数据**

```bash
# 备份数据库
mysqldump -u root -p tower > tower_backup_$(date +%Y%m%d%H%M%S).sql

# 备份配置文件
cp .env .env.backup
```

2. **拉取新代码**

```bash
git pull origin main
```

3. **更新依赖**

```bash
go mod tidy
```

4. **重新生成 Swagger**

```bash
swag init -g cmd/main.go
```

5. **编译重启**

```bash
# 重新编译
go build -o tower-go cmd/main.go

# 重启服务
sudo systemctl restart tower-go
```

6. **验证升级**

```bash
# 查看服务状态
sudo systemctl status tower-go

# 查看日志
journalctl -u tower-go -f -n 100
```

### 版本兼容说明

| 版本 | 说明 | 升级注意事项 |
|------|------|-------------|
| v1.0.x | 初始版本 | - |
| v1.1.x | 新增功能 | 需执行迁移脚本 |
| v2.0.x | 重大更新 | 不兼容升级，需重新配置 |

## 🆘 故障排查

### 查看系统资源

```bash
# 查看内存使用
free -h

# 查看磁盘使用
df -h

# 查看进程
top

# 查看网络连接
netstat -tunlp | grep 10024
```

### 查看应用日志

```bash
# 系统日志
journalctl -u tower-go -f

# 查看特定时间日志
journalctl -u tower-go --since "2025-11-11 10:00:00" --until "2025-11-11 11:00:00"

# 应用日志
tail -f logs/app.log
```

### 数据库排查

```bash
# 查看 MySQL 进程
mysqladmin processlist

# 慢查询分析
mysqldumpslow /var/log/mysql/slow.log

# 查看表大小
SELECT table_name, ROUND(((data_length + index_length) / 1024 / 1024), 2) AS size_mb
FROM information_schema.TABLES
WHERE table_schema = 'tower'
ORDER BY (data_length + index_length) DESC;
```

### 获取帮助

1. 查看 [FAQ](#常见问题)
2. 查看 GitHub Issues
3. 提交 Issue 反馈

---

**文档版本**: v1.0.0
**最后更新**: 2025-11-11
**维护团队**: Tower-Go Team
