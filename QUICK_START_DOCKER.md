# 🚀 Docker Nginx 快速开始（3分钟）

## 一键启动

### Windows 用户

```bash
# 双击运行或命令行执行
.\start-nginx.bat
```

### Linux/Mac 用户

```bash
chmod +x start-nginx.sh
./start-nginx.sh
```

---

## 验证部署

启动完成后，浏览器访问以下 URL 验证：

### ✅ 1. 健康检查
```
http://localhost:8080/health
```
应该看到：`healthy`

### ✅ 2. 测试文件
```
http://localhost:8080/images/test.txt
```
应该看到：`Docker Nginx is working!`

### ✅ 3. 目录浏览
```
http://localhost:8080/images/
```
应该看到文件列表

---

## 配置已完成

`.env` 文件已自动配置：

```bash
IMAGE_UPLOAD_PATH=C:/Users/Administrator/Desktop/xdAdmin/tower-go/uploads/images
IMAGE_BASE_URL=http://localhost:8080/images
```

---

## 开始测试

### 1. 重启应用
```bash
.\tower-go.exe
```

### 2. 创建报菜记录
通过 API 或管理界面创建报菜

### 3. 检查结果

**查看日志：**
```bash
tail -f logs/app.log
# 查找：Image saved successfully
```

**访问图片：**
```
http://localhost:8080/images/2024/01/15/xxxxxx.png
```

**检查钉钉群：**
应该收到带图片的消息

---

## 常用命令

```bash
# 查看容器状态
docker ps | grep tower-nginx

# 查看日志
docker-compose -f docker-compose.nginx.yml logs -f

# 重启服务
docker-compose -f docker-compose.nginx.yml restart

# 停止服务
docker-compose -f docker-compose.nginx.yml down
```

---

## ⚠️ 注意事项

### 本地测试限制

Docker 映射的 `localhost:8080` **只能本机访问**，钉钉服务器无法访问。

**解决方案：**

1. **开发测试：** 使用内网穿透工具
   - ngrok: `ngrok http 8080`
   - frp
   
2. **生产环境：** 部署到有公网 IP 的服务器
   ```bash
   IMAGE_BASE_URL=http://your-domain.com:8080/images
   ```

---

## 📚 详细文档

- `docs/DOCKER_NGINX_SETUP.md` - 完整部署指南
- `IMAGE_FEATURE_README.md` - 功能使用指南
- `SUMMARY.md` - 技术实现总结

---

**下一步：** 启动应用并测试报菜功能！
