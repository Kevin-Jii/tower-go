# Nginx 图片托管配置指南

## 概述

本项目使用 nginx 托管报菜明细图片，通过 Markdown 格式在钉钉中显示图片。

## 配置步骤

### 1. 创建图片存储目录

```bash
# Linux
sudo mkdir -p /var/www/html/images
sudo chmod 755 /var/www/html/images
sudo chown www-data:www-data /var/www/html/images

# Windows (根据实际路径调整)
mkdir C:\nginx\html\images
```

### 2. 配置 Nginx

编辑 nginx 配置文件 `/etc/nginx/sites-available/default` (Linux) 或 `nginx.conf` (Windows)：

```nginx
server {
    listen 80;
    server_name your-domain.com;  # 替换为你的域名或IP

    # 图片访问路径
    location /images/ {
        alias /var/www/html/images/;  # Linux路径
        # alias C:/nginx/html/images/;  # Windows路径（取消注释并使用此行）
        
        # 允许跨域访问（钉钉需要）
        add_header Access-Control-Allow-Origin *;
        add_header Access-Control-Allow-Methods 'GET, OPTIONS';
        add_header Access-Control-Allow-Headers 'DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type';
        
        # 缓存配置
        expires 7d;
        add_header Cache-Control "public, immutable";
        
        # 自动索引（可选，用于调试）
        autoindex on;
        autoindex_exact_size off;
        autoindex_localtime on;
    }

    # 其他配置...
}
```

### 3. 测试配置并重启 Nginx

```bash
# 测试配置
sudo nginx -t

# 重启 nginx
sudo systemctl restart nginx
# 或
sudo service nginx restart

# Windows
nginx.exe -t
nginx.exe -s reload
```

### 4. 配置应用环境变量

编辑 `.env` 文件：

```bash
# 图片上传配置
IMAGE_UPLOAD_PATH=/var/www/html/images  # nginx 托管的图片目录（绝对路径）
IMAGE_BASE_URL=http://your-domain.com/images  # 图片访问基础URL

# 示例 (本地测试)
# IMAGE_UPLOAD_PATH=C:/nginx/html/images
# IMAGE_BASE_URL=http://localhost/images

# 示例 (生产环境)
# IMAGE_UPLOAD_PATH=/var/www/html/images
# IMAGE_BASE_URL=https://cdn.your-domain.com/images
```

### 5. 验证配置

1. 手动创建一个测试图片：
```bash
echo "test" > /var/www/html/images/test.txt
```

2. 浏览器访问：
```
http://your-domain.com/images/test.txt
```

3. 如果能看到文件内容，说明配置成功。

## 目录结构

图片会按日期自动分类存储：

```
/var/www/html/images/
├── 2024/
│   ├── 01/
│   │   ├── 15/
│   │   │   ├── 143052_a1b2c3d4.png
│   │   │   └── 150230_e5f6g7h8.png
│   │   └── 16/
│   │       └── 091520_i9j0k1l2.png
│   └── 02/
│       └── 01/
│           └── 120000_m3n4o5p6.png
```

访问URL示例：
```
http://your-domain.com/images/2024/01/15/143052_a1b2c3d4.png
```

## HTTPS 配置（推荐）

钉钉推荐使用 HTTPS，配置 SSL 证书：

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;
    
    ssl_certificate /path/to/your/cert.pem;
    ssl_certificate_key /path/to/your/key.pem;
    
    location /images/ {
        alias /var/www/html/images/;
        add_header Access-Control-Allow-Origin *;
        expires 7d;
    }
}

# HTTP 重定向到 HTTPS
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}
```

对应的环境变量：
```bash
IMAGE_BASE_URL=https://your-domain.com/images
```

## 磁盘空间管理

### 自动清理旧图片

可以设置定时任务清理30天前的图片：

```bash
# 编辑 crontab
crontab -e

# 添加定时任务（每天凌晨3点执行）
0 3 * * * find /var/www/html/images -type f -mtime +30 -delete
```

### 程序内清理

也可以在应用中调用清理函数：

```go
// 清理30天前的图片
err := utils.CleanOldImages(30)
```

## 故障排查

### 1. 图片保存失败

检查目录权限：
```bash
ls -la /var/www/html/images
# 确保应用进程用户有写权限
sudo chown -R www-data:www-data /var/www/html/images
sudo chmod -R 755 /var/www/html/images
```

### 2. 图片无法访问

检查 nginx 配置：
```bash
sudo nginx -t
curl http://your-domain.com/images/test.txt
```

查看 nginx 日志：
```bash
sudo tail -f /var/log/nginx/error.log
sudo tail -f /var/log/nginx/access.log
```

### 3. 钉钉图片不显示

- 确认图片 URL 可以公网访问（钉钉服务器需要能访问）
- 检查跨域配置是否正确
- 使用 HTTPS（钉钉推荐）
- 图片大小不要超过 2MB

## 安全建议

1. **限制上传大小**：在应用层限制图片大小（建议 < 2MB）
2. **文件类型验证**：只允许图片格式（PNG, JPG, GIF）
3. **防盗链**：配置 nginx 防止外部盗链
4. **定期备份**：重要图片定期备份
5. **监控磁盘**：设置磁盘空间告警

## CDN 加速（可选）

对于高并发场景，建议使用 CDN：

```bash
# 使用 CDN 后的配置
IMAGE_BASE_URL=https://cdn.your-domain.com/images
```

配置 CDN 回源到你的 nginx 服务器即可。
