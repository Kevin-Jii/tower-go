# ApiClient 智能封装方法使用指南

## 概述

`ApiClient` 提供了智能封装方法，自动处理：
- ✅ 自动提取响应中的 `data` 字段
- ✅ 自动过滤请求中的 `null` 值
- ✅ 自动处理响应为空的情况
- ✅ 统一的错误处理

## 方法列表

### 1. `postSmart` - 智能 POST 请求

**用途**：创建资源，自动处理 data 提取和 null 过滤

```dart
// 返回创建的对象，如果后端不返回则为 null
Future<Store?> createStore(CreateStoreRequest request) async {
  return await _client.postSmart<Store>(
    path: '/stores',
    data: request.toJson(),
    fromJson: Store.fromJson,
    filterNulls: true, // 自动过滤 null 值
  );
}
```

**参数**：
- `path`: 请求路径
- `data`: 请求数据（Map 或对象的 toJson()）
- `fromJson`: 解析函数，将 Map 转为对象
- `filterNulls`: 是否过滤 null 值（默认 false）
- `queryParameters`: 可选的查询参数

**返回**：`Future<T?>` - 返回 data 中的对象，如果 data 为 null 则返回 null

---

### 2. `putSmart` - 智能 PUT 请求

**用途**：更新资源，自动过滤 null 值（默认开启）

```dart
// 返回更新后的对象，如果后端不返回则为 null
Future<Store?> updateStore(int id, UpdateStoreRequest request) async {
  return await _client.putSmart<Store>(
    path: '/stores/$id',
    data: request.toJson(),
    fromJson: Store.fromJson,
    filterNulls: true, // PUT 默认为 true
  );
}
```

**特点**：
- 默认 `filterNulls = true`，只发送需要更新的字段
- 例如只更新状态：`{"status": 1}` 而不是 `{"name": null, "address": null, "status": 1}`

---

### 3. `getSmart` - 智能 GET 单个对象

**用途**：获取单个资源详情

```dart
Future<Store> getStoreDetail(int id) async {
  return await _client.getSmart<Store>(
    path: '/stores/$id',
    fromJson: Store.fromJson,
  );
}
```

**返回**：`Future<T>` - 返回 data 中的对象（data 为 null 会抛异常）

---

### 4. `getListSmart` - 智能 GET 列表

**用途**：获取资源列表（非分页）

```dart
Future<List<Store>> listStores() async {
  return await _client.getListSmart<Store>(
    path: '/stores/all',
    fromJson: Store.fromJson,
  );
}
```

**返回**：`Future<List<T>>` - 返回 data 中的列表，如果为 null 返回空列表

---

### 5. `deleteSmart` - 智能 DELETE 请求

**用途**：删除资源，不期望返回数据

```dart
Future<void> deleteStore(int id) async {
  await _client.deleteSmart(path: '/stores/$id');
}
```

**返回**：`Future<void>` - 无返回值，成功即完成

---

### 6. `getPage` - 分页请求（原有方法）

**用途**：获取分页数据

```dart
Future<PageResponse<Store>> getStores({int page = 1, int pageSize = 10}) async {
  return _client.getPage<Store>(
    '/stores',
    itemParser: Store.fromJson,
    queryParameters: {
      'page': page,
      'page_size': pageSize,
    },
  );
}
```

## 完整示例：StoreApi

```dart
class StoreApi {
  final ApiClient _client;

  StoreApi(this._client);

  /// 获取分页列表
  Future<PageResponse<Store>> getStores({
    int page = 1,
    int pageSize = 10,
    String? keyword,
  }) async {
    return _client.getPage<Store>(
      '/stores',
      itemParser: Store.fromJson,
      queryParameters: {
        'page': page,
        'page_size': pageSize,
        if (keyword != null && keyword.isNotEmpty) 'keyword': keyword,
      },
    );
  }

  /// 获取所有门店（非分页）
  Future<List<Store>> listStores() async {
    return await _client.getListSmart<Store>(
      path: '/stores/all',
      fromJson: Store.fromJson,
    );
  }

  /// 获取门店详情
  Future<Store> getStoreDetail(int id) async {
    return await _client.getSmart<Store>(
      path: '/stores/$id',
      fromJson: Store.fromJson,
    );
  }

  /// 创建门店
  Future<Store?> createStore(CreateStoreRequest request) async {
    return await _client.postSmart<Store>(
      path: '/stores',
      data: request.toJson(),
      fromJson: Store.fromJson,
      filterNulls: true,
    );
  }

  /// 更新门店
  Future<Store?> updateStore(int id, UpdateStoreRequest request) async {
    return await _client.putSmart<Store>(
      path: '/stores/$id',
      data: request.toJson(),
      fromJson: Store.fromJson,
      filterNulls: true,
    );
  }

  /// 删除门店
  Future<void> deleteStore(int id) async {
    await _client.deleteSmart(path: '/stores/$id');
  }
}
```

## 优势对比

### 使用前（手动处理）：

```dart
Future<Store> updateStore(int id, UpdateStoreRequest request) async {
  // 手动过滤 null
  final data = request.toJson()..removeWhere((key, value) => value == null);
  
  return _client.request<Store>(
    '/stores/$id',
    method: 'PUT',
    data: data,
    converter: (json) {
      // 手动检查 null
      final data = json['data'];
      if (data == null) {
        throw ApiException('更新响应数据为空');
      }
      return Store.fromJson(data as Map<String, dynamic>);
    },
  );
}
```

### 使用后（自动处理）：

```dart
Future<Store?> updateStore(int id, UpdateStoreRequest request) async {
  return await _client.putSmart<Store>(
    path: '/stores/$id',
    data: request.toJson(),
    fromJson: Store.fromJson,
    filterNulls: true, // 自动过滤 null
  );
}
```

## 错误处理

所有方法都会自动：
1. 捕获网络错误并转换为 `ApiException`
2. 验证响应格式
3. 检查业务状态码
4. 提供清晰的错误信息

```dart
try {
  final store = await storeApi.createStore(request);
} on ApiException catch (e) {
  print('错误: ${e.message}');
  print('状态码: ${e.statusCode}');
}
```

## 注意事项

1. **filterNulls**：
   - POST: 默认 `false`，创建时通常需要完整数据
   - PUT: 默认 `true`，更新时只发送修改的字段

2. **返回值为 null**：
   - `postSmart` 和 `putSmart` 返回 `T?`
   - 如果后端不返回 data，返回 null
   - 不影响操作成功性，Provider 中可以刷新列表获取最新数据

3. **列表为空**：
   - `getListSmart` 返回空列表而不是抛异常
   - 适合可选的列表数据

## 迁移指南

### 旧代码：
```dart
await _client.post<void>('/stores', data: request.toJson());
```

### 新代码：
```dart
await _client.postSmart<Store>(
  path: '/stores',
  data: request.toJson(),
  fromJson: Store.fromJson,
  filterNulls: true,
);
```

## 总结

智能封装方法让你：
- 🚀 更少的代码
- 🛡️ 更好的错误处理
- 🎯 统一的调用方式
- 💪 自动处理边缘情况
