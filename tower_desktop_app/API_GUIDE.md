# API 响应格式说明

## 后端统一响应格式

### 1. 普通响应（无分页，meta 为 null）
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "phone": "13012345678"
  },
  "meta": null
}
```

### 2. 分页响应（meta 包裹分页信息）
```json
{
  "code": 200,
  "message": "success",
  "data": [ {"id": 1, "username": "user1"}, {"id": 2, "username": "user2"} ],
  "meta": {
    "total": 100,
    "page": 1,
    "page_size": 10,
    "page_count": 10,
    "has_more": true
  }
}
```

**约定**：
- `meta == null`：非分页接口，`data` 可能是对象、数组或其他类型
- `meta != null`：分页接口，`data` 必须是数组，`meta` 包含分页元信息

### 3. 错误响应
```json
{
  "code": 400,
  "message": "参数错误",
  "data": null,
  "meta": null
}
```

## ApiClient 使用方式

### 按 HTTP 方法对 data 的期望
| 方法 | 是否必须有 data | 典型场景 | 说明 |
| ---- | -------------- | -------- | ---- |
| GET | 必须（除非仅检查状态，不推荐） | 获取详情/列表 | 若缺少 `data` 将抛出 `data 字段为空` 异常 |
| POST | 可选 | 新增/执行动作 | 成功可返回对象或不返回（`data: null`） |
| PUT | 可选 | 更新资源 | 同 POST，返回最新对象或省略 |
| DELETE | 可选 | 删除资源 | 通常无 `data` |

> 设计理由：读取型接口通常需要 `data`；写入型接口只关心是否成功（`code==200`），减少后端负担。

### 普通请求（返回单个对象）
```dart
// 使用 converter 自动处理 data 字段
final user = await apiClient.get<User>(
  '/users/1',
  converter: (json) => User.fromJson(json['data']),
);
```

### 分页请求（推荐使用 getPage/postPage 辅助方法）
```dart
// 自动校验 meta 并解析分页信息
final pageResp = await apiClient.getPage<User>(
  '/users',
  queryParameters: {'page': 1, 'page_size': 10},
  itemParser: (json) => User.fromJson(json),
);

// pageResp 是 PageResponse<User>，包含 list/total/page/pageSize
print('总数: ${pageResp.total}, 当前页: ${pageResp.page}');
for (var user in pageResp.list) {
  print(user.username);
}
```

**优势**：
- 自动校验 `meta` 存在且格式正确，缺失或为 null 抛异常
- 严格类型检查 `total/page/page_size` 必须为 int 且非 null
- 可选字段 `page_count/has_more` 存在时也会校验类型
- 无需手动编写 converter 解析逻辑

### 手动 converter（复杂场景或非标准格式）
```dart
// converter 接收完整响应体（包含 data、meta）
final response = await apiClient.get<UserListResponse>(
  '/users',
  queryParameters: {'page': 1, 'page_size': 10},
  converter: (json) {
    return UserListResponse(
      list: (json['data'] as List)
          .map((e) => User.fromJson(e))
          .toList(),
      total: json['meta']['total'],
      page: json['meta']['page'],
      pageSize: json['meta']['page_size'],
    );
  },
);
```

### 不使用 converter（直接返回 data）
```dart
// 返回 Map<String, dynamic>（data 字段内容）
final data = await apiClient.get<Map<String, dynamic>>('/users/1');

// 返回 List（data 字段内容）
final list = await apiClient.get<List<dynamic>>('/users');
```

## 优势说明

1. **统一错误处理**：ApiClient 自动检查 `code` 字段，非 200 抛出 `ApiException`
2. **灵活类型转换**：通过 `converter` 可处理任意复杂响应结构
3. **分页信息封装**：分页数据与元信息（total、page）一起返回
4. **类型安全**：编译期检查，避免运行时类型错误
5. **避免重复代码**：不再需要在每个 API 方法中手动提取 `data` 字段
6. **严格校验**：`getPage/postPage` 强制分页接口必须包含有效的 `meta`

## 其他模块快速迁移

### 角色管理示例（分页）
```dart
class RoleApi {
  final ApiClient _client;
  
  Future<PageResponse<Role>> getRoles({int page = 1}) async {
    return await _client.getPage<Role>(
      '/roles',
      queryParameters: {'page': page, 'page_size': 10},
      itemParser: (json) => Role.fromJson(json),
    );
  }
  
  Future<Role> createRole(CreateRoleRequest req) async {
    return await _client.post<Role>(
      '/roles',
      data: req,
      converter: (json) => Role.fromJson(json['data']),
    );
  }
}
```

### 门店管理示例（非分页列表）
```dart
class StoreApi {
  final ApiClient _client;
  
  // 后端返回 {code, message, data: [...], meta: null}
  Future<List<Store>> getStores() async {
    return await _client.get<List<Store>>(
      '/stores',
      converter: (json) => (json['data'] as List)
          .map((e) => Store.fromJson(e))
          .toList(),
    );
  }
}
```

### 用户管理示例（使用 getPage）
```dart
class UserApi {
  final ApiClient _client;
  
  Future<UserListResponse> getUsers({
    int page = 1,
    int pageSize = 10,
    String? keyword,
  }) async {
    final pageResp = await _client.getPage<User>(
      '/users',
      queryParameters: {
        'page': page,
        'page_size': pageSize,
        if (keyword != null && keyword.isNotEmpty) 'keyword': keyword,
      },
      itemParser: (json) => User.fromJson(json),
    );
    return UserListResponse(
      list: pageResp.list,
      total: pageResp.total,
      page: pageResp.page,
      pageSize: pageResp.pageSize,
    );
  }
}
```
