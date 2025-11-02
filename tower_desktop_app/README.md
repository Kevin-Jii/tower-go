# Tower Desktop App (Flutter)

多门店报菜/菜单权限管理系统桌面端（Windows 优先）。

## 功能规划
- 登录，持久化会话
- 获取用户菜单树 & 按钮级权限
- 门店列表查看
- 菜品 CRUD（当前仅列出 + 新增占位）
- 报菜记录列表 & 新增
- 权限控制组件 `PermissionGate`

## 目录结构 (简化)
```
lib/
  core/
    network/api_client.dart
    widgets/permission_gate.dart
  features/
    auth/ (models, api, session)
    menu/
    dish/
    store/
    report/
    common/api_result.dart
  main.dart
```

## 依赖安装
确保已安装 Flutter SDK (>=3.3) 并开启 Windows 支持。
```powershell
flutter config --enable-windows-desktop
flutter pub get
flutter pub run build_runner build --delete-conflicting-outputs
flutter run -d windows
```

## 生成数据模型 (Freezed/Json)
```powershell
flutter pub run build_runner watch --delete-conflicting-outputs
```

## 调整后端地址
修改 `lib/core/network/api_client.dart` 中 `baseUrl` 指向你的后端，例如：
```
static const String baseUrl = 'http://127.0.0.1:8080/api/v1';
```
若为局域网环境：
```
static const String baseUrl = 'http://192.168.x.x:8080/api/v1';
```

## 登录接口契约
POST /auth/login
请求: {"username":"admin","password":"123456"}
响应期望包装：
```
{
  "code":0,
  "message":"OK",
  "data":{
    "token":"<JWT>",
    "user":{ "id":1,"username":"admin","storeId":1,"roleCode":"admin","roleId":1 }
  }
}
```

## 权限与菜单
- 获取用户菜单: GET /menus/user-menus
- 获取权限字符串: GET /menus/user-permissions -> ["dish:add", "store:list", ...]
- 通过 `SessionManager().hasPermission('dish:add')` 或 `PermissionGate(permission: 'dish:add', child: ...)` 控制 UI

## 计划迭代
- 持久化 token (本地文件/SharedPreferences)
- 动态菜单树渲染
- 完整 CRUD 页面
- 国际化 (i18n)
- Windows 安装包 (MSIX / Inno Setup)

## Windows 打包
```powershell
flutter build windows
# 生成的可执行与运行时位于 build/windows/x64/runner/Release
```
MSIX：参见官方文档 https://docs.flutter.dev/deployment/windows

## 开发技巧
- 接口错误封装在 ApiException
- 后端统一返回建议使用 {code,message,data}
- 若为空列表后端返回 []，解析已兼容

## License
内部项目，未指定 License。
