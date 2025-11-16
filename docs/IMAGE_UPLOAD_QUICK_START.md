# 图片上传功能快速开始指南

## 快速测试（Windows 本地环境）

### 方案一：使用现有 Web 服务器

如果你已经有 nginx 或其他 Web 服务器：

1. **创建图片目录**
```bash
mkdir C:\nginx\html\images
# 或者根据你的实际路径
mkdir C:\xampp\htdocs\images
```

2. **配置 .env**
```bash
# nginx
IMAGE_UPLOAD_PATH=C:/nginx/html/images
IMAGE_BASE_URL=http://localhost/images

# 或者 Apache/XAMPP
# IMAGE_UPLOAD_PATH=C:/xampp/htdocs/images
# IMAGE_BASE_URL=http://localhost/images
```

3. **重启应用**
```bash
.\tower-go.exe
```

### 方案二：不使用 Web 服务器（临时方案）

如果没有 nginx，可以临时使用本地路径：

1. **创建本地目录**
```bash
mkdir C:\tower-go-images
```

2. **配置 .env（使用应用内置的静态文件服务）**
```bash
IMAGE_UPLOAD_PATH=C:/tower-go-images
IMAGE_BASE_URL=http://localhost:10024/static/images
```

3. **需要在应用中添加静态文件服务**（可选，如果需要自托管）

### 方案三：使用公网可访问的服务器

**最推荐的方案**（钉钉需要公网访问）

1. **在服务器上创建目录**
```bash
# SSH 登录到你的服务器
ssh user@your-server.com

# 创建目录
sudo mkdir -p /var/www/html/images
sudo chmod 755 /var/www/html/images
```

2. **配置 nginx**
```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location /images/ {
        alias /var/www/html/images/;
        add_header Access-Control-Allow-Origin *;
        expires 7d;
        autoindex on;
    }
}
```

3. **配置应用 .env**
```bash
IMAGE_UPLOAD_PATH=/var/www/html/images
IMAGE_BASE_URL=http://your-domain.com/images
# 或使用 HTTPS（推荐）
# IMAGE_BASE_URL=https://your-domain.com/images
```

## 测试图片上传功能

### 1. 手动测试

创建测试图片并访问：

```bash
# 在图片目录创建测试文件
echo "test" > C:\nginx\html\images\test.txt

# 浏览器访问
http://localhost/images/test.txt
```

### 2. 通过应用测试

创建一个报菜记录，查看日志：

```bash
# 查看应用日志
tail -f logs/app.log

# 应该看到类似日志：
# Image saved successfully
# imageURL=http://localhost/images/2024/01/15/143052_a1b2c3d4.png
```

### 3. 在钉钉中查看

创建报菜后，钉钉群应该收到带图片的 Markdown 消息。

## 常见问题

### Q1: 钉钉看不到图片？

**A:** 钉钉机器人需要能够访问图片 URL，确保：
- 图片 URL 是公网可访问的（不能是 localhost）
- 配置了 CORS 跨域头
- 使用 HTTPS（钉钉推荐）

### Q2: 图片保存失败？

**A:** 检查：
```bash
# 确认目录存在
ls C:\nginx\html\images

# 确认应用有写权限
# Windows: 右键目录 -> 属性 -> 安全 -> 编辑权限
# Linux: sudo chmod 755 /var/www/html/images
```

### Q3: Stream 模式不支持 Markdown？

**A:** 根据之前的测试，钉钉 Stream 模式的群消息可能不支持 `sampleMarkdown`。
可以尝试：
- 使用 Webhook 模式（支持 Markdown）
- 或者在文本消息中附加图片链接

## 目录权限设置

### Windows
```bash
# 右键目录 -> 属性 -> 安全
# 确保 Users 组有"修改"权限
```

### Linux
```bash
sudo chown -R www-data:www-data /var/www/html/images
sudo chmod -R 755 /var/www/html/images

# 如果应用以其他用户运行
sudo chown -R your-app-user:your-app-user /var/www/html/images
```

## 生产环境建议

1. **使用 HTTPS**
```bash
IMAGE_BASE_URL=https://your-domain.com/images
```

2. **使用 CDN 加速**
```bash
IMAGE_BASE_URL=https://cdn.your-domain.com/images
```

3. **设置定时清理**
```bash
# crontab -e
0 3 * * * find /var/www/html/images -type f -mtime +30 -delete
```

4. **监控磁盘空间**
```bash
# 设置告警，当磁盘使用超过 80% 时通知
```

## 下一步

配置完成后，测试创建报菜记录，应该会：
1. ✅ 生成 PNG 图片
2. ✅ 保存到 nginx 目录
3. ✅ 获得图片 URL
4. ✅ 在 Markdown 中引用图片
5. ✅ 发送到钉钉群（带图片）

如遇到问题，查看 `logs/app.log` 获取详细信息。
