# API 错误码约定

统一使用 `pkg/apicode` 维护错误码和默认提示信息。

错误码格式为 `HTTP 状态码前三位 + 两位业务序号`：

- `400xx`：请求参数或校验错误
- `401xx`：身份认证错误
- `403xx`：权限或业务禁止
- `404xx`：资源不存在
- `409xx`：资源冲突、库存不足、状态冲突
- `422xx`：业务数据校验失败
- `500xx`：服务端错误

业务层返回错误：

```go
return apicode.New(apicode.ProductNotFound)
return apicode.Newf(apicode.InvalidParameter, "商品 %d 的参数无效", productID)
return apicode.Wrap(apicode.InternalError, err)
```

Controller 统一转换响应：

```go
if err != nil {
    http.ErrorFrom(ctx, err)
    return
}
```

`ErrorFrom` 会识别 `apicode` 错误并返回对应 code/message；未包装的底层错误只返回 `50000`，避免把数据库或第三方服务细节暴露给客户端。
