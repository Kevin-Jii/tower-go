# Docker Nginx 图片服务部署指南

## 📦 快速开始

### Windows 用户

```bash
# 1. 一键启动（推荐）
.\start-nginx.bat

# 或手动执行
docker-compose -f docker-compose.nginx.yml up -d
```

### Linux/Mac 用户

```bash
# 1. 赋予执行权限
chmod +x start-nginx.sh

# 2. 一键启动
./start-nginx.sh

# 或手动执行
docker-compose -f docker-compose.nginx.yml up -d
```

---

## 🔧 配置说明

### Docker Compose 配置

文件：`docker-compose.nginx.yml`

```yaml
services:
  nginx-images:
    image: nginx:alpine          # 使用轻量级 Alpine 镜像
    container_name: tower-nginx-images
    ports:
      - "8080:80"                # 宿主机端口:容器端口
    volumes:
      - ./uploads/images:/usr/share/nginx/html/images:rw  # 图片目录
      - ./docker/nginx/nginx.conf:/etc/nginx/conf.d/default.conf:ro  # 配置文件
    restart: unless-stopped      # 自动重启
```

### Nginx 配置

文件：`docker/nginx/nginx.conf`

**核心配置：**
- 监听端口：80 (容器内)
- 图片路径：`/images/`
- CORS 支持：允许跨域访问
- 缓存时间：7天
- 目录浏览：开启（开发环境）

---

## 📁 目录结构

```
tower-go/
├── docker-compose.nginx.yml     # Docker Compose 配置
├── start-nginx.bat              # Windows 启动脚本
├── start-nginx.sh               # Linux/Mac 启动脚本
├── docker/
│   └── nginx/
│       └── nginx.conf           # Nginx 配置文件
└── uploads/
    └── images/                  # 图片存储目录
        ├── test.txt             # 测试文件
        └── 2024/                # 按日期分类
            └── 01/
                └── 15/
                    └── 143052_abc123.png
```

---

## 🚀 部署步骤

### 步骤 1：启动 Nginx 容器

```bash
# Windows
.\start-nginx.bat

# Linux/Mac
./start-nginx.sh
```

### 步骤 2：验证服务

浏览器访问以下 URL：

1. **健康检查**
   ```
   http://localhost:8080/health
   应返回：healthy
   ```

2. **测试文件**
   ```
   http://localhost:8080/images/test.txt
   应返回：Docker Nginx is working!
   ```

3. **目录浏览**
   ```
   http://localhost:8080/images/
   应显示文件列表
   ```

### 步骤 3：配置应用

编辑 `.env` 文件：

```bash
# Windows 路径格式
IMAGE_UPLOAD_PATH=C:/Users/Administrator/Desktop/xdAdmin/tower-go/uploads/images
IMAGE_BASE_URL=http://localhost:8080/images

# Linux 路径格式
# IMAGE_UPLOAD_PATH=/path/to/tower-go/uploads/images
# IMAGE_BASE_URL=http://localhost:8080/images
```

### 步骤 4：重启应用

```bash
.\tower-go.exe
```

### 步骤 5：测试功能

1. 创建报菜记录
2. 检查日志：`tail -f logs/app.log`
3. 查看钉钉群消息
4. 访问图片 URL

---

## 🔍 验证清单

- [ ] Docker 已启动
- [ ] Nginx 容器正在运行
- [ ] `http://localhost:8080/health` 返回 healthy
- [ ] `http://localhost:8080/images/test.txt` 可访问
- [ ] `.env` 配置正确
- [ ] 应用已重启
- [ ] 创建报菜记录成功
- [ ] 图片保存到 uploads/images
- [ ] 图片 URL 可访问
- [ ] 钉钉收到消息

---

## 🛠️ 常用命令

### 查看容器状态
```bash
docker ps | grep tower-nginx
```

### 查看日志
```bash
# 实时查看
docker-compose -f docker-compose.nginx.yml logs -f

# 查看最后 100 行
docker-compose -f docker-compose.nginx.yml logs --tail=100
```

### 重启服务
```bash
docker-compose -f docker-compose.nginx.yml restart
```

### 停止服务
```bash
docker-compose -f docker-compose.nginx.yml down
```

### 进入容器
```bash
docker exec -it tower-nginx-images sh
```

### 查看容器内文件
```bash
docker exec tower-nginx-images ls -la /usr/share/nginx/html/images/
```

---

## 🌐 生产环境部署

### 使用域名访问

1. **配置域名解析**
   ```
   A记录：images.your-domain.com -> 服务器IP
   ```

2. **修改端口映射**
   
   编辑 `docker-compose.nginx.yml`：
   ```yaml
   ports:
     - "80:80"  # 使用标准 HTTP 端口
   ```

3. **更新应用配置**
   ```bash
   IMAGE_BASE_URL=http://images.your-domain.com/images
   ```

### 使用 HTTPS（推荐）

1. **准备 SSL 证书**
   ```bash
   mkdir -p docker/nginx/ssl
   # 将证书放到此目录
   # cert.pem, key.pem
   ```

2. **更新 Nginx 配置**
   
   编辑 `docker/nginx/nginx.conf`：
   ```nginx
   server {
       listen 443 ssl http2;
       server_name images.your-domain.com;
       
       ssl_certificate /etc/nginx/ssl/cert.pem;
       ssl_certificate_key /etc/nginx/ssl/key.pem;
       
       location /images/ {
           alias /usr/share/nginx/html/images/;
           add_header Access-Control-Allow-Origin *;
           expires 7d;
       }
   }
   
   server {
       listen 80;
       server_name images.your-domain.com;
       return 301 https://$server_name$request_uri;
   }
   ```

3. **更新 Docker Compose**
   ```yaml
   ports:
     - "80:80"
     - "443:443"
   volumes:
     - ./docker/nginx/ssl:/etc/nginx/ssl:ro
   ```

4. **更新应用配置**
   ```bash
   IMAGE_BASE_URL=https://images.your-domain.com/images
   ```

---

## 📊 性能优化

### 1. 使用 CDN

```bash
# 配置 CDN 回源到 Nginx
IMAGE_BASE_URL=https://cdn.your-domain.com/images
```

### 2. 限制文件大小

编辑 `docker/nginx/nginx.conf`：
```nginx
client_max_body_size 5M;  # 限制上传大小
```

### 3. Gzip 压缩

```nginx
gzip on;
gzip_types image/png image/jpeg image/gif;
gzip_min_length 1000;
```

### 4. 增加缓存时间

```nginx
expires 30d;  # 缓存 30 天
```

---

## 🐛 故障排查

### 问题 1：容器无法启动

**检查 Docker：**
```bash
docker info
```

**查看错误日志：**
```bash
docker-compose -f docker-compose.nginx.yml logs
```

**常见原因：**
- Docker 未启动
- 端口 8080 被占用
- 配置文件格式错误

**解决方案：**
```bash
# 查看端口占用
netstat -ano | findstr :8080

# 修改端口
# 编辑 docker-compose.nginx.yml，改为 "8081:80"
```

### 问题 2：图片无法访问

**检查容器状态：**
```bash
docker ps | grep tower-nginx
```

**检查文件是否存在：**
```bash
docker exec tower-nginx-images ls -la /usr/share/nginx/html/images/
```

**检查 Nginx 配置：**
```bash
docker exec tower-nginx-images nginx -t
```

### 问题 3：钉钉看不到图片

**原因：** Docker 映射的 `localhost:8080` 只能本机访问

**解决方案：**
1. 使用公网 IP：`http://your-public-ip:8080/images`
2. 使用域名：`http://images.your-domain.com/images`
3. 使用内网穿透工具（测试用）

---

## 📦 备份和恢复

### 备份图片

```bash
# 打包所有图片
tar -czf images-backup-$(date +%Y%m%d).tar.gz uploads/images/

# 备份到远程
scp images-backup-*.tar.gz user@backup-server:/backups/
```

### 恢复图片

```bash
# 解压备份
tar -xzf images-backup-20240115.tar.gz

# 重启 Nginx
docker-compose -f docker-compose.nginx.yml restart
```

---

## 🔐 安全建议

1. **生产环境禁用目录浏览**
   ```nginx
   autoindex off;  # 关闭目录浏览
   ```

2. **限制访问来源**
   ```nginx
   # 只允许特定 IP
   allow 1.2.3.4;
   deny all;
   ```

3. **防盗链**
   ```nginx
   valid_referers none blocked server_names
                  *.your-domain.com;
   if ($invalid_referer) {
       return 403;
   }
   ```

4. **定期更新镜像**
   ```bash
   docker pull nginx:alpine
   docker-compose -f docker-compose.nginx.yml up -d
   ```

---

## 📈 监控

### 查看访问统计

```bash
# 查看访问日志
docker exec tower-nginx-images cat /var/log/nginx/images_access.log

# 实时监控
docker exec tower-nginx-images tail -f /var/log/nginx/images_access.log
```

### 监控容器资源

```bash
docker stats tower-nginx-images
```

---

## 🎯 下一步

1. ✅ **启动 Nginx 容器**
2. ✅ **配置应用 .env**
3. ✅ **测试图片上传**
4. ⏸️ 生产环境配置 HTTPS
5. ⏸️ 配置 CDN 加速

---

**快速命令参考：**

```bash
# 启动
.\start-nginx.bat

# 查看日志
docker-compose -f docker-compose.nginx.yml logs -f

# 停止
docker-compose -f docker-compose.nginx.yml down

# 重启
docker-compose -f docker-compose.nginx.yml restart
```
