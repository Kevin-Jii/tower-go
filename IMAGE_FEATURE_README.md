# 📸 钉钉图片推送功能使用指南

> 报菜记录自动生成图片并推送到钉钉群

---

## 🚀 5分钟快速开始

### 步骤 1️⃣：选择部署方式

<table>
<tr>
<td width="50%">

#### 方案 A：使用 Nginx（推荐）

适用于：生产环境、已有 Nginx 服务器

优势：性能好、稳定、支持 CDN

</td>
<td width="50%">

#### 方案 B：本地测试

适用于：本地开发、快速测试

优势：简单快速、无需外部依赖

</td>
</tr>
</table>

---

### 步骤 2️⃣：配置环境

#### 方案 A：Nginx 配置

**1. 创建图片目录**
```bash
# Linux
sudo mkdir -p /var/www/html/images
sudo chmod 755 /var/www/html/images

# Windows
mkdir C:\nginx\html\images
```

**2. 配置 Nginx**
```nginx
# 编辑 nginx.conf 或 sites-available/default
location /images/ {
    alias /var/www/html/images/;
    add_header Access-Control-Allow-Origin *;
    expires 7d;
}
```

**3. 重启 Nginx**
```bash
# Linux
sudo systemctl restart nginx

# Windows
nginx.exe -s reload
```

**4. 配置应用 .env**
```bash
IMAGE_UPLOAD_PATH=/var/www/html/images
IMAGE_BASE_URL=http://your-domain.com/images
```

#### 方案 B：本地测试配置

**Windows 用户：**
```bash
# 以管理员身份运行
.\setup_local_images.bat
```

**Linux 用户：**
```bash
mkdir -p ./uploads/images
chmod 755 ./uploads/images
```

**配置 .env：**
```bash
IMAGE_UPLOAD_PATH=./uploads/images
IMAGE_BASE_URL=http://localhost:10024/static/images
```

---

### 步骤 3️⃣：测试验证

**1. 测试文件访问**
```bash
# 创建测试文件
echo "test" > /var/www/html/images/test.txt

# 浏览器访问
http://your-domain.com/images/test.txt
```

**2. 创建报菜记录**
- 通过 API 或管理界面创建报菜
- 等待几秒

**3. 检查钉钉群**
- 应该收到带图片的消息
- 或收到带图片链接的消息

**4. 查看日志确认**
```bash
tail -f logs/app.log
# 查找：Image saved successfully
```

---

## 📋 完整文档导航

```
📁 项目根目录
│
├── 📄 SUMMARY.md                          ⭐ 完整实现总结
├── 📄 IMPLEMENTATION_CHECKLIST.md         ✅ 实现和测试清单
│
├── 📁 docs/
│   ├── 📄 IMAGE_UPLOAD_QUICK_START.md    🚀 快速开始（推荐先看）
│   ├── 📄 nginx_image_config.md          🔧 Nginx 详细配置
│   └── 📄 DINGTALK_IMAGE_FEATURE.md      📚 完整技术文档
│
└── 📁 utils/
    └── 📄 file_helper.go                  💻 文件管理工具代码
```

**建议阅读顺序：**
1. 本文件 - 了解概况
2. `docs/IMAGE_UPLOAD_QUICK_START.md` - 快速开始
3. `docs/nginx_image_config.md` - Nginx 配置
4. `SUMMARY.md` - 完整总结

---

## ❓ 常见问题

### Q1: 钉钉看不到图片怎么办？

**A:** 确保图片 URL 可以公网访问

```bash
# ❌ 错误示例（钉钉无法访问）
IMAGE_BASE_URL=http://localhost/images
IMAGE_BASE_URL=http://192.168.1.100/images

# ✅ 正确示例
IMAGE_BASE_URL=http://your-domain.com/images
IMAGE_BASE_URL=https://cdn.example.com/images
```

### Q2: 图片保存失败？

**A:** 检查目录权限

```bash
# Linux
sudo chmod 755 /var/www/html/images
sudo chown www-data:www-data /var/www/html/images

# 查看日志
tail -f logs/app.log | grep "Failed to save"
```

### Q3: Markdown 格式不生效？

**A:** 正常现象，钉钉 Stream 模式群消息可能不支持 Markdown

系统会自动降级为纯文本 + 图片链接，用户点击链接即可查看图片。

### Q4: 需要清理旧图片吗？

**A:** 建议设置定时清理（30天）

```bash
# Linux 定时任务
crontab -e
# 添加：每天凌晨3点清理30天前的图片
0 3 * * * find /var/www/html/images -type f -mtime +30 -delete
```

---

## 🎯 核心功能

- ✅ 自动生成 PNG 格式图片
- ✅ 按日期分类存储（2024/01/15/...）
- ✅ 唯一文件名（时间戳+MD5）
- ✅ Markdown 格式支持
- ✅ 自动降级机制
- ✅ 详细日志记录
- ✅ 支持 CDN 加速
- ✅ 定期清理功能

---

## 🔍 技术架构

```
报菜创建 → 生成图片 → 保存文件 → 生成URL → 推送钉钉
   ↓          ↓          ↓          ↓          ↓
 API      PNG格式    Nginx目录   HTTP链接   Markdown
 请求     800px      按日期分类   公网可访问  自动降级
```

---

## 📞 获取帮助

### 查看日志
```bash
# 实时查看
tail -f logs/app.log

# 搜索错误
grep "error\|Error\|ERROR" logs/app.log

# 搜索图片相关
grep "image\|Image" logs/app.log
```

### 验证配置
```bash
# 检查环境变量
cat .env | grep IMAGE

# 测试目录权限
ls -la /var/www/html/images

# 测试 Nginx
curl http://your-domain.com/images/test.txt
```

### 调试模式
```bash
# 设置日志级别为 debug
LOG_LEVEL=debug
```

---

## 🌟 生产环境建议

### 安全配置
- [ ] 使用 HTTPS（推荐）
- [ ] 配置防盗链
- [ ] 限制文件大小
- [ ] 定期备份

### 性能优化
- [ ] 配置 CDN
- [ ] 启用 Gzip 压缩
- [ ] 设置缓存策略
- [ ] 图片质量优化

### 监控告警
- [ ] 磁盘空间监控
- [ ] 上传失败告警
- [ ] 访问日志分析
- [ ] 性能指标统计

---

## ✅ 验收清单

部署完成后，请确认：

- [ ] 图片目录已创建并有写权限
- [ ] Nginx 配置正确并已重启
- [ ] .env 配置正确
- [ ] 应用已重启
- [ ] 测试文件可访问
- [ ] 报菜功能正常
- [ ] 钉钉收到消息
- [ ] 图片链接可访问
- [ ] 日志无错误

---

## 📊 当前状态

| 项目 | 状态 |
|------|------|
| 代码开发 | ✅ 完成 |
| 编译测试 | ✅ 通过 |
| 文档编写 | ✅ 完成 |
| 环境配置 | 🔄 待完成 |
| 功能测试 | 🔄 待完成 |
| 生产部署 | ⏸️ 待定 |

---

**下一步：** 按照本文档配置环境，然后进行功能测试

**预计耗时：** 5-10 分钟（首次配置）

**如有问题，请查看详细文档：** `docs/IMAGE_UPLOAD_QUICK_START.md`
