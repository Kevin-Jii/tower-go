# 错误码快速参考

## 通用错误码（1000-1999）

| 错误码 | 错误名称                  | 错误描述                 |
|--------|---------------------------|--------------------------|
| 200    | ErrSuccess                | 成功                     |
| 1000   | ErrBadRequest             | 请求参数错误             |
| 1001   | ErrInternalServer         | 服务器内部错误           |
| 1002   | ErrNotFound               | 资源不存在               |
| 1003   | ErrMethodNotAllowed       | 请求方法不允许           |
| 1004   | ErrTooManyRequests        | 请求过于频繁             |
| 1005   | ErrServiceBusy            | 服务繁忙，请稍后重试     |
| 1100   | ErrDatabaseQuery          | 数据库查询错误           |
| 1101   | ErrDatabaseInsert         | 数据库插入错误           |
| 1102   | ErrDatabaseUpdate         | 数据库更新错误           |
| 1103   | ErrDatabaseDelete         | 数据库删除错误           |
| 1104   | ErrDuplicateKey           | 数据已存在               |
| 1200   | ErrValidation             | 数据验证失败             |
| 1201   | ErrInvalidEmail           | 邮箱格式不正确           |
| 1202   | ErrInvalidPhone           | 手机号格式不正确         |
| 1203   | ErrPasswordTooWeak        | 密码强度不足             |
| 1204   | ErrInvalidDateFormat      | 日期格式不正确           |
| 1205   | ErrInvalidDataRange       | 数据范围不正确           |
| 1206   | ErrRequiredFieldMissing   | 必填字段缺失             |

## 认证授权错误码（2000-2999）

| 错误码 | 错误名称                | 错误描述               |
|--------|-------------------------|------------------------|
| 2000   | ErrUnauthorized         | 未授权，请先登录       |
| 2001   | ErrTokenInvalid         | Token 无效             |
| 2002   | ErrTokenExpired         | Token 已过期           |
| 2003   | ErrTokenMissing         | Token 缺失             |
| 2004   | ErrTokenMalformed       | Token 格式错误         |
| 2005   | ErrLoginFailed          | 用户名或密码错误       |
| 2006   | ErrAccountLocked        | 账号已被锁定           |
| 2007   | ErrAccountDisabled      | 账号已被禁用           |
| 2008   | ErrPasswordIncorrect    | 密码错误               |
| 2009   | ErrOldPasswordIncorrect | 原密码错误             |
| 2100   | ErrForbidden            | 无权限访问             |
| 2101   | ErrInsufficientPriv     | 权限不足               |
| 2102   | ErrRoleNotFound         | 角色不存在             |
| 2103   | ErrPermissionDenied     | 权限被拒绝             |
| 2104   | ErrStoreAccessDenied    | 无权访问该门店数据     |
| 2105   | ErrCrossStoreOperation  | 不允许跨门店操作       |

## 用户业务错误码（3000-3999）

| 错误码 | 错误名称                  | 错误描述                 |
|--------|---------------------------|--------------------------|
| 3000   | ErrUserNotFound           | 用户不存在               |
| 3001   | ErrUserAlreadyExists      | 用户已存在               |
| 3002   | ErrUsernameAlreadyTaken   | 用户名已被占用           |
| 3003   | ErrPhoneAlreadyTaken      | 手机号已被占用           |
| 3004   | ErrEmailAlreadyTaken      | 邮箱已被占用             |
| 3005   | ErrUserCreateFailed       | 用户创建失败             |
| 3006   | ErrUserUpdateFailed       | 用户更新失败             |
| 3007   | ErrUserDeleteFailed       | 用户删除失败             |
| 3008   | ErrCannotDeleteSelf       | 不能删除自己             |
| 3009   | ErrPasswordResetFailed    | 密码重置失败             |
| 3100   | ErrStoreNotFound          | 门店不存在               |
| 3101   | ErrStoreAlreadyExists     | 门店已存在               |
| 3102   | ErrStoreCreateFailed      | 门店创建失败             |
| 3103   | ErrStoreUpdateFailed      | 门店更新失败             |
| 3104   | ErrStoreDeleteFailed      | 门店删除失败             |
| 3105   | ErrStoreHasUsers          | 门店下还有用户，无法删除 |
| 3200   | ErrRoleAlreadyExists      | 角色已存在               |
| 3201   | ErrRoleCreateFailed       | 角色创建失败             |
| 3202   | ErrRoleUpdateFailed       | 角色更新失败             |
| 3203   | ErrRoleDeleteFailed       | 角色删除失败             |
| 3204   | ErrRoleHasUsers           | 角色下还有用户，无法删除 |

## 菜品业务错误码（4000-4999）

| 错误码 | 错误名称                    | 错误描述                   |
|--------|-----------------------------|----------------------------|
| 4000   | ErrDishNotFound             | 菜品不存在                 |
| 4001   | ErrDishAlreadyExists        | 菜品已存在                 |
| 4002   | ErrDishCreateFailed         | 菜品创建失败               |
| 4003   | ErrDishUpdateFailed         | 菜品更新失败               |
| 4004   | ErrDishDeleteFailed         | 菜品删除失败               |
| 4005   | ErrDishInUse                | 菜品正在使用中，无法删除   |
| 4100   | ErrMenuReportNotFound       | 报菜记录不存在             |
| 4101   | ErrMenuReportCreateFailed   | 报菜创建失败               |
| 4102   | ErrMenuReportUpdateFailed   | 报菜更新失败               |
| 4103   | ErrMenuReportDeleteFailed   | 报菜删除失败               |
| 4104   | ErrMenuReportInvalidDate    | 报菜日期不正确             |
| 4105   | ErrMenuReportDuplicate      | 该日期的报菜已存在         |

## 权限菜单错误码（5000-5999）

| 错误码 | 错误名称                  | 错误描述                 |
|--------|---------------------------|--------------------------|
| 5000   | ErrMenuNotFound           | 菜单不存在               |
| 5001   | ErrMenuAlreadyExists      | 菜单已存在               |
| 5002   | ErrMenuCreateFailed       | 菜单创建失败             |
| 5003   | ErrMenuUpdateFailed       | 菜单更新失败             |
| 5004   | ErrMenuDeleteFailed       | 菜单删除失败             |
| 5005   | ErrMenuHasChildren        | 菜单下还有子菜单，无法删除 |
| 5006   | ErrMenuParentNotFound     | 父菜单不存在             |
| 5007   | ErrMenuCircularReference  | 菜单存在循环引用         |
| 5100   | ErrRoleMenuAssignFailed   | 角色菜单权限分配失败     |
| 5101   | ErrRoleMenuNotFound       | 角色菜单权限不存在       |
| 5102   | ErrStoreRoleMenuFailed    | 门店角色菜单权限分配失败 |
| 5103   | ErrMenuTreeBuildFailed    | 菜单树构建失败           |

## WebSocket 错误码（6000-6999）

| 错误码 | 错误名称              | 错误描述                 |
|--------|-----------------------|--------------------------|
| 6000   | ErrWSConnectionFailed | WebSocket 连接失败       |
| 6001   | ErrWSUpgradeFailed    | WebSocket 升级失败       |
| 6002   | ErrWSMessageInvalid   | WebSocket 消息格式错误   |
| 6003   | ErrWSSessionNotFound  | WebSocket 会话不存在     |
| 6004   | ErrWSSessionConflict  | 账号已在其他设备登录     |
| 6005   | ErrWSBroadcastFailed  | 消息广播失败             |

## 文件上传错误码（7000-7999）

| 错误码 | 错误名称                | 错误描述             |
|--------|-------------------------|----------------------|
| 7000   | ErrFileUploadFailed     | 文件上传失败         |
| 7001   | ErrFileTypeNotAllowed   | 文件类型不允许       |
| 7002   | ErrFileSizeTooLarge     | 文件大小超出限制     |
| 7003   | ErrFileNotFound         | 文件不存在           |
| 7004   | ErrFileDeleteFailed     | 文件删除失败         |
| 7005   | ErrFileDownloadFailed   | 文件下载失败         |
| 7006   | ErrImageProcessFailed   | 图片处理失败         |
| 7007   | ErrFileStorageFull      | 存储空间不足         |
| 7008   | ErrInvalidFileFormat    | 文件格式不正确       |
| 7009   | ErrFileCorrupted        | 文件已损坏           |
| 7010   | ErrFileNameTooLong      | 文件名过长           |
| 7011   | ErrFilePathInvalid      | 文件路径不正确       |
| 7012   | ErrMultipartFormError   | 解析表单数据失败     |
| 7013   | ErrFileReadFailed       | 文件读取失败         |
| 7014   | ErrFileWriteFailed      | 文件写入失败         |
| 7015   | ErrDirectoryCreateFail  | 目录创建失败         |

## 第三方服务错误码（8000-8999）

| 错误码 | 错误名称                | 错误描述             |
|--------|-------------------------|----------------------|
| 8000   | ErrThirdPartyService    | 第三方服务错误       |
| 8001   | ErrSMSServiceFailed     | 短信服务失败         |
| 8002   | ErrEmailServiceFailed   | 邮件服务失败         |
| 8003   | ErrPaymentFailed        | 支付失败             |
| 8004   | ErrPaymentTimeout       | 支付超时             |
| 8005   | ErrOSSUploadFailed      | OSS 上传失败         |
| 8006   | ErrCacheFailed          | 缓存服务失败         |
| 8007   | ErrMessageQueueFailed   | 消息队列失败         |
| 8008   | ErrExternalAPIFailed    | 外部 API 调用失败    |
| 8009   | ErrExternalAPITimeout   | 外部 API 超时        |
| 8010   | ErrWeChatServiceFailed  | 微信服务失败         |
| 8011   | ErrAlipayServiceFailed  | 支付宝服务失败       |
| 8012   | ErrNetworkError         | 网络连接错误         |
| 8013   | ErrServiceUnavailable   | 服务不可用           |
| 8014   | ErrRateLimitExceeded    | 调用频率超限         |
| 8015   | ErrAPIKeyInvalid        | API 密钥无效         |
| 8016   | ErrAuthorizationFailed  | 第三方授权失败       |
| 8017   | ErrCallbackFailed       | 回调处理失败         |
| 8018   | ErrWebhookFailed        | Webhook 处理失败     |
| 8019   | ErrDataSyncFailed       | 数据同步失败         |

## 业务逻辑错误码（9000-9999）

| 错误码 | 错误名称                  | 错误描述             |
|--------|---------------------------|----------------------|
| 9000   | ErrBusinessLogic          | 业务逻辑错误         |
| 9001   | ErrOperationNotAllowed    | 不允许的操作         |
| 9002   | ErrDataInconsistency      | 数据不一致           |
| 9003   | ErrConcurrentConflict     | 并发冲突             |
| 9004   | ErrResourceLocked         | 资源已被锁定         |
| 9005   | ErrQuotaExceeded          | 配额已超限           |
| 9006   | ErrInventoryInsufficient  | 库存不足             |
| 9007   | ErrOrderExpired           | 订单已过期           |
| 9008   | ErrOrderCanceled          | 订单已取消           |
| 9009   | ErrRefundFailed           | 退款失败             |
| 9010   | ErrStatusInvalid          | 状态不正确           |
| 9011   | ErrWorkflowError          | 工作流错误           |
| 9012   | ErrDependencyMissing      | 依赖项缺失           |
| 9013   | ErrConfigError            | 配置错误             |
| 9014   | ErrFeatureDisabled        | 功能已禁用           |
| 9015   | ErrMaintenanceMode        | 系统维护中           |
| 9016   | ErrVersionMismatch        | 版本不匹配           |
| 9017   | ErrDataExpired            | 数据已过期           |
| 9018   | ErrDuplicateOperation     | 重复操作             |
| 9019   | ErrPreconditionFailed     | 前置条件不满足       |

---

## 使用示例

### Controller 层

```go
// 参数验证错误
if err := ctx.ShouldBindJSON(&req); err != nil {
    utils.ErrorWithCode(ctx, utils.ErrBadRequest.WithMessage("请求参数格式错误"))
    return
}

// 用户不存在
user, err := c.service.GetUser(userID)
if err != nil {
    utils.ErrorWithCode(ctx, utils.ErrUserNotFound)
    return
}

// 权限不足
if !hasPermission {
    utils.ErrorWithCode(ctx, utils.ErrForbidden)
    return
}

// 数据库错误
if err := c.service.CreateUser(req); err != nil {
    utils.LogDatabaseError("CreateUser", err)
    utils.ErrorWithCode(ctx, utils.ErrDatabaseInsert)
    return
}

// 成功响应
utils.Success(ctx, user)
```

### Service 层

```go
// 检查用户名是否已存在
if exists {
    return utils.ErrUsernameAlreadyTaken
}

// 检查门店访问权限
if user.StoreID != storeID {
    return utils.ErrStoreAccessDenied
}

// 业务逻辑错误
if !isValid {
    return utils.ErrValidation.WithMessage("手机号必须为11位")
}
```

---

## HTTP 状态码映射

| 错误码范围 | HTTP 状态码              | 说明         |
|------------|--------------------------|--------------|
| 200        | 200 OK                   | 成功         |
| 1000-1999  | 500 Internal Server Error | 服务器错误   |
| 2000-2099  | 401 Unauthorized         | 认证失败     |
| 2100-2999  | 403 Forbidden            | 权限不足     |
| 3000-9999  | 400 Bad Request          | 业务错误     |

---

## 前端处理建议

```javascript
// axios 拦截器示例
axios.interceptors.response.use(
  response => response,
  error => {
    const code = error.response.data.code;
    const message = error.response.data.message;
    
    switch(code) {
      case 2000:
      case 2001:
      case 2002:
        // Token 失效，跳转登录
        router.push('/login');
        break;
      case 2100:
      case 2103:
        // 权限不足
        Message.error('您没有权限执行此操作');
        break;
      default:
        // 显示错误消息
        Message.error(message);
    }
    
    return Promise.reject(error);
  }
);
```
