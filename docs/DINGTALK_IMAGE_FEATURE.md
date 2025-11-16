# 钉钉图片推送功能实现文档

## 功能概述

通过 **Nginx 托管图片 + Markdown 引用** 的方案，实现报菜记录单带图片推送到钉钉群。

## 技术方案

### 方案对比

| 方案 | 优点 | 缺点 | 状态 |
|------|------|------|------|
| 钉钉 Webhook 图片 | 简单直接 | 不支持图片 | ❌ 不可用 |
| 钉钉 Stream 群消息图片 | 官方支持 | 只支持 sampleText | ❌ 不可用 |
| 钉钉企业公告 | 支持图片 | API 复杂，需要 AgentID | ⚠️ 已实现但备用 |
| **Nginx + Markdown** | 简单高效，兼容性好 | 需要公网访问 | ✅ **当前方案** |

### 最终方案架构

```
报菜创建 
  ↓
生成 PNG 图片
  ↓
保存到 Nginx 托管目录
  ↓
获取图片 URL
  ↓
Markdown 引用图片
  ↓
推送到钉钉群
```

## 实现细节

### 1. 配置管理 (`config/config.go`)

新增应用配置：
```go
type AppConfig struct {
    Name             string
    Port             int
    ImageUploadPath  string // 图片上传目录（绝对路径）
    ImageBaseURL     string // 图片访问基础URL
}
```

环境变量：
```bash
IMAGE_UPLOAD_PATH=/var/www/html/images
IMAGE_BASE_URL=http://your-domain.com/images
```

### 2. 文件管理工具 (`utils/file_helper.go`)

**核心功能：**

#### SaveImageFile()
```go
func SaveImageFile(filename string, imageData []byte) (string, error)
```
- 按日期分类存储：`2024/01/15/143052_a1b2c3d4.png`
- 使用 MD5 + 时间戳生成唯一文件名
- 返回完整访问 URL

#### DeleteImageFile()
```go
func DeleteImageFile(imageURL string) error
```
- 根据 URL 删除对应文件

#### CleanOldImages()
```go
func CleanOldImages(days int) error
```
- 清理指定天数前的旧图片

### 3. 钉钉服务 (`service/dingtalk.go`)

**新增/修改的函数：**

#### sendStreamImageText()
```go
func (s *DingTalkService) sendStreamImageText(
    bot *model.DingTalkBot, 
    title, text string, 
    imageData []byte
) error
```
**流程：**
1. 保存图片到 Nginx 目录
2. 获取图片 URL
3. 在 Markdown 中引用图片：`![报菜明细](url)`
4. 发送消息（优先 Markdown，失败降级为文本）

#### saveImageToNginx()
```go
func (s *DingTalkService) saveImageToNginx(
    imageData []byte, 
    filename string
) (string, error)
```
调用 `utils.SaveImageFile()` 保存图片

#### sendStreamMarkdownWithText()
```go
func (s *DingTalkService) sendStreamMarkdownWithText(
    bot *model.DingTalkBot, 
    title, markdownText, accessToken string
) error
```
**特性：**
- 优先尝试 `sampleMarkdown` 格式
- 失败自动降级为 `sampleText` 格式
- 保留图片链接在文本中

### 4. 消息格式

#### Markdown 格式（优先）
```markdown
## 📋 新报菜通知

**门店名称:** XXX店
**操作人员:** 张三
**报菜时间:** 2024-01-15 14:30:52

**报菜明细:**
- **宫保鸡丁**: 数量 5
- **鱼香肉丝**: 数量 3

![报菜明细](http://your-domain.com/images/2024/01/15/143052_a1b2c3d4.png)
```

#### 纯文本格式（降级）
```
📋 新报菜通知

门店名称: XXX店
操作人员: 张三
报菜时间: 2024-01-15 14:30:52

报菜明细:
• 宫保鸡丁: 数量 5
• 鱼香肉丝: 数量 3

图片链接: http://your-domain.com/images/2024/01/15/143052_a1b2c3d4.png
```

## 配置步骤

### Step 1: 配置 Nginx

创建配置文件 `/etc/nginx/sites-available/images`:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location /images/ {
        alias /var/www/html/images/;
        
        # CORS 配置（钉钉需要）
        add_header Access-Control-Allow-Origin *;
        add_header Access-Control-Allow-Methods 'GET, OPTIONS';
        
        # 缓存配置
        expires 7d;
        add_header Cache-Control "public, immutable";
        
        # 开发环境可开启目录浏览
        autoindex on;
        autoindex_exact_size off;
        autoindex_localtime on;
    }
}
```

启用配置：
```bash
sudo ln -s /etc/nginx/sites-available/images /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### Step 2: 创建图片目录

```bash
sudo mkdir -p /var/www/html/images
sudo chmod 755 /var/www/html/images
sudo chown www-data:www-data /var/www/html/images
```

### Step 3: 配置应用环境变量

编辑 `.env`:

```bash
# 图片上传配置
IMAGE_UPLOAD_PATH=/var/www/html/images
IMAGE_BASE_URL=http://your-domain.com/images

# 生产环境使用 HTTPS（推荐）
# IMAGE_BASE_URL=https://your-domain.com/images
```

### Step 4: 重启应用

```bash
./tower-go.exe
# 或
systemctl restart tower-go
```

## 测试验证

### 1. 手动测试图片访问

```bash
# 创建测试文件
echo "test" > /var/www/html/images/test.txt

# 浏览器访问
curl http://your-domain.com/images/test.txt
```

### 2. 测试报菜功能

1. 创建报菜记录
2. 查看日志确认图片已保存：
```bash
tail -f logs/app.log
# 应看到：
# Image saved successfully imageURL=http://...
```

3. 访问图片 URL 确认可访问
4. 检查钉钉群是否收到消息

### 3. 日志关键信息

成功日志示例：
```
Image saved successfully
  botID=6
  imageURL=http://your-domain.com/images/2024/01/15/143052_a1b2c3d4.png
  imageSize=245678

Sending stream message to DingTalk API
  robotCode=xxx
  msgKey=sampleMarkdown
```

## 故障排查

### 问题 1: 图片保存失败

**症状：**
```
Failed to save image to nginx
error=failed to create directory: permission denied
```

**解决：**
```bash
# 检查目录权限
ls -la /var/www/html/images
sudo chmod 755 /var/www/html/images
sudo chown www-data:www-data /var/www/html/images
```

### 问题 2: 钉钉看不到图片

**可能原因：**
1. URL 不是公网可访问（localhost 不行）
2. 钉钉服务器被防火墙拦截
3. 没有配置 CORS

**解决：**
```bash
# 1. 使用公网域名或IP
IMAGE_BASE_URL=http://your-public-ip/images

# 2. 检查防火墙
sudo ufw allow 80/tcp

# 3. 确认 nginx CORS 配置
add_header Access-Control-Allow-Origin *;
```

### 问题 3: Markdown 不生效

**症状：**
日志显示降级为纯文本：
```
Markdown format not supported, falling back to plain text with link
```

**说明：**
这是正常的，Stream 模式群消息不支持 Markdown。
图片链接会以纯文本形式显示，用户可以点击访问。

## 性能优化

### 1. 磁盘空间管理

定时清理旧图片（30天）：
```bash
# crontab -e
0 3 * * * find /var/www/html/images -type f -mtime +30 -delete
```

或在代码中：
```go
// 每天执行一次
go func() {
    ticker := time.NewTicker(24 * time.Hour)
    for range ticker.C {
        utils.CleanOldImages(30)
    }
}()
```

### 2. CDN 加速

使用 CDN 分发图片：
```bash
IMAGE_BASE_URL=https://cdn.your-domain.com/images
```

### 3. 图片压缩

在生成图片时进行压缩（可选）：
```go
// 在 image_generator.go 中调整 PNG 压缩级别
encoder := png.Encoder{CompressionLevel: png.BestCompression}
```

## 安全建议

1. **限制文件大小**：单个图片不超过 2MB
2. **文件类型验证**：只允许 PNG/JPG/GIF
3. **防盗链**：配置 nginx referer 检查
4. **HTTPS**：生产环境强制使用 HTTPS
5. **定期备份**：重要图片定期备份到对象存储

## 维护检查清单

- [ ] Nginx 服务正常运行
- [ ] 图片目录有足够空间（< 80%）
- [ ] 图片 URL 公网可访问
- [ ] 日志无错误信息
- [ ] 定时清理任务正常执行
- [ ] 钉钉能正常接收图片

## 相关文件

- `config/config.go` - 配置管理
- `utils/file_helper.go` - 文件工具
- `utils/image_generator.go` - 图片生成
- `service/dingtalk.go` - 钉钉服务
- `service/menu_report_listener.go` - 报菜事件监听
- `docs/nginx_image_config.md` - Nginx 详细配置
- `docs/IMAGE_UPLOAD_QUICK_START.md` - 快速开始指南

## 更新日志

- **2024-01-15**: 初始实现，使用 Nginx + Markdown 方案
- 支持图片按日期分类存储
- 支持 Markdown 格式（Stream 模式降级为文本）
- 支持自动清理旧图片
