package auth

import "fmt"

// ErrorCode 错误码结构
type ErrorCode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error 实现 error 接口
func (e ErrorCode) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// WithMessage 返回带自定义消息的错误
func (e ErrorCode) WithMessage(msg string) ErrorCode {
	return ErrorCode{
		Code:    e.Code,
		Message: msg,
	}
}

// WithMessageF 返回带格式化消息的错误
func (e ErrorCode) WithMessageF(format string, args ...interface{}) ErrorCode {
	return ErrorCode{
		Code:    e.Code,
		Message: fmt.Sprintf(format, args...),
	}
}

// ========== 通用错误码 (1000-1999) ==========

var (
	// 成功
	ErrSuccess = ErrorCode{Code: 200, Message: "success"}

	// 通用错误
	ErrBadRequest       = ErrorCode{Code: 1000, Message: "请求参数错误"}
	ErrInternalServer   = ErrorCode{Code: 1001, Message: "服务器内部错误"}
	ErrNotFound         = ErrorCode{Code: 1002, Message: "资源不存在"}
	ErrMethodNotAllowed = ErrorCode{Code: 1003, Message: "请求方法不允许"}
	ErrTooManyRequests  = ErrorCode{Code: 1004, Message: "请求过于频繁"}
	ErrServiceBusy      = ErrorCode{Code: 1005, Message: "服务繁忙，请稍后重试"}

	// 数据库错误
	ErrDatabaseQuery  = ErrorCode{Code: 1100, Message: "数据库查询错误"}
	ErrDatabaseInsert = ErrorCode{Code: 1101, Message: "数据库插入错误"}
	ErrDatabaseUpdate = ErrorCode{Code: 1102, Message: "数据库更新错误"}
	ErrDatabaseDelete = ErrorCode{Code: 1103, Message: "数据库删除错误"}
	ErrDuplicateKey   = ErrorCode{Code: 1104, Message: "数据已存在"}

	// 数据验证错误
	ErrValidation           = ErrorCode{Code: 1200, Message: "数据验证失败"}
	ErrInvalidEmail         = ErrorCode{Code: 1201, Message: "邮箱格式不正确"}
	ErrInvalidPhone         = ErrorCode{Code: 1202, Message: "手机号格式不正确"}
	ErrPasswordTooWeak      = ErrorCode{Code: 1203, Message: "密码强度不足"}
	ErrInvalidDateFormat    = ErrorCode{Code: 1204, Message: "日期格式不正确"}
	ErrInvalidDataRange     = ErrorCode{Code: 1205, Message: "数据范围不正确"}
	ErrRequiredFieldMissing = ErrorCode{Code: 1206, Message: "必填字段缺失"}
)

// ========== 认证授权错误码 (2000-2999) ==========

var (
	// 认证错误
	ErrUnauthorized         = ErrorCode{Code: 2000, Message: "未授权，请先登录"}
	ErrTokenInvalid         = ErrorCode{Code: 2001, Message: "Token 无效"}
	ErrTokenExpired         = ErrorCode{Code: 2002, Message: "Token 已过期"}
	ErrTokenMissing         = ErrorCode{Code: 2003, Message: "Token 缺失"}
	ErrTokenMalformed       = ErrorCode{Code: 2004, Message: "Token 格式错误"}
	ErrLoginFailed          = ErrorCode{Code: 2005, Message: "用户名或密码错误"}
	ErrAccountLocked        = ErrorCode{Code: 2006, Message: "账号已被锁定"}
	ErrAccountDisabled      = ErrorCode{Code: 2007, Message: "账号已被禁用"}
	ErrPasswordIncorrect    = ErrorCode{Code: 2008, Message: "密码错误"}
	ErrOldPasswordIncorrect = ErrorCode{Code: 2009, Message: "原密码错误"}

	// 权限错误
	ErrForbidden           = ErrorCode{Code: 2100, Message: "无权限访问"}
	ErrInsufficientPriv    = ErrorCode{Code: 2101, Message: "权限不足"}
	ErrRoleNotFound        = ErrorCode{Code: 2102, Message: "角色不存在"}
	ErrPermissionDenied    = ErrorCode{Code: 2103, Message: "权限被拒绝"}
	ErrStoreAccessDenied   = ErrorCode{Code: 2104, Message: "无权访问该门店数据"}
	ErrCrossStoreOperation = ErrorCode{Code: 2105, Message: "不允许跨门店操作"}
)

// ========== 用户业务错误码 (3000-3999) ==========

var (
	// 用户相关
	ErrUserNotFound         = ErrorCode{Code: 3000, Message: "用户不存在"}
	ErrUserAlreadyExists    = ErrorCode{Code: 3001, Message: "用户已存在"}
	ErrUsernameAlreadyTaken = ErrorCode{Code: 3002, Message: "用户名已被占用"}
	ErrPhoneAlreadyTaken    = ErrorCode{Code: 3003, Message: "手机号已被占用"}
	ErrEmailAlreadyTaken    = ErrorCode{Code: 3004, Message: "邮箱已被占用"}
	ErrUserCreateFailed     = ErrorCode{Code: 3005, Message: "用户创建失败"}
	ErrUserUpdateFailed     = ErrorCode{Code: 3006, Message: "用户更新失败"}
	ErrUserDeleteFailed     = ErrorCode{Code: 3007, Message: "用户删除失败"}
	ErrCannotDeleteSelf     = ErrorCode{Code: 3008, Message: "不能删除自己"}
	ErrPasswordResetFailed  = ErrorCode{Code: 3009, Message: "密码重置失败"}

	// 门店相关
	ErrStoreNotFound      = ErrorCode{Code: 3100, Message: "门店不存在"}
	ErrStoreAlreadyExists = ErrorCode{Code: 3101, Message: "门店已存在"}
	ErrStoreCreateFailed  = ErrorCode{Code: 3102, Message: "门店创建失败"}
	ErrStoreUpdateFailed  = ErrorCode{Code: 3103, Message: "门店更新失败"}
	ErrStoreDeleteFailed  = ErrorCode{Code: 3104, Message: "门店删除失败"}
	ErrStoreHasUsers      = ErrorCode{Code: 3105, Message: "门店下还有用户，无法删除"}

	// 角色相关
	ErrRoleAlreadyExists = ErrorCode{Code: 3200, Message: "角色已存在"}
	ErrRoleCreateFailed  = ErrorCode{Code: 3201, Message: "角色创建失败"}
	ErrRoleUpdateFailed  = ErrorCode{Code: 3202, Message: "角色更新失败"}
	ErrRoleDeleteFailed  = ErrorCode{Code: 3203, Message: "角色删除失败"}
	ErrRoleHasUsers      = ErrorCode{Code: 3204, Message: "角色下还有用户，无法删除"}
)

// ========== 菜品业务错误码 (4000-4999) ==========

var (
	// 菜品相关
	ErrDishNotFound      = ErrorCode{Code: 4000, Message: "菜品不存在"}
	ErrDishAlreadyExists = ErrorCode{Code: 4001, Message: "菜品已存在"}
	ErrDishCreateFailed  = ErrorCode{Code: 4002, Message: "菜品创建失败"}
	ErrDishUpdateFailed  = ErrorCode{Code: 4003, Message: "菜品更新失败"}
	ErrDishDeleteFailed  = ErrorCode{Code: 4004, Message: "菜品删除失败"}
	ErrDishInUse         = ErrorCode{Code: 4005, Message: "菜品正在使用中，无法删除"}

	// 报菜相关
	ErrMenuReportNotFound     = ErrorCode{Code: 4100, Message: "报菜记录不存在"}
	ErrMenuReportCreateFailed = ErrorCode{Code: 4101, Message: "报菜创建失败"}
	ErrMenuReportUpdateFailed = ErrorCode{Code: 4102, Message: "报菜更新失败"}
	ErrMenuReportDeleteFailed = ErrorCode{Code: 4103, Message: "报菜删除失败"}
	ErrMenuReportInvalidDate  = ErrorCode{Code: 4104, Message: "报菜日期不正确"}
	ErrMenuReportDuplicate    = ErrorCode{Code: 4105, Message: "该日期的报菜已存在"}
)

// ========== 权限菜单错误码 (5000-5999) ==========

var (
	// 菜单相关
	ErrMenuNotFound          = ErrorCode{Code: 5000, Message: "菜单不存在"}
	ErrMenuAlreadyExists     = ErrorCode{Code: 5001, Message: "菜单已存在"}
	ErrMenuCreateFailed      = ErrorCode{Code: 5002, Message: "菜单创建失败"}
	ErrMenuUpdateFailed      = ErrorCode{Code: 5003, Message: "菜单更新失败"}
	ErrMenuDeleteFailed      = ErrorCode{Code: 5004, Message: "菜单删除失败"}
	ErrMenuHasChildren       = ErrorCode{Code: 5005, Message: "菜单下还有子菜单，无法删除"}
	ErrMenuParentNotFound    = ErrorCode{Code: 5006, Message: "父菜单不存在"}
	ErrMenuCircularReference = ErrorCode{Code: 5007, Message: "菜单存在循环引用"}

	// 角色菜单权限相关
	ErrRoleMenuAssignFailed = ErrorCode{Code: 5100, Message: "角色菜单权限分配失败"}
	ErrRoleMenuNotFound     = ErrorCode{Code: 5101, Message: "角色菜单权限不存在"}
	ErrStoreRoleMenuFailed  = ErrorCode{Code: 5102, Message: "门店角色菜单权限分配失败"}
	ErrMenuTreeBuildFailed  = ErrorCode{Code: 5103, Message: "菜单树构建失败"}
)

// ========== WebSocket 错误码 (6000-6999) ==========

var (
	ErrWSConnectionFailed = ErrorCode{Code: 6000, Message: "WebSocket 连接失败"}
	ErrWSUpgradeFailed    = ErrorCode{Code: 6001, Message: "WebSocket 升级失败"}
	ErrWSMessageInvalid   = ErrorCode{Code: 6002, Message: "WebSocket 消息格式错误"}
	ErrWSSessionNotFound  = ErrorCode{Code: 6003, Message: "WebSocket 会话不存在"}
	ErrWSSessionConflict  = ErrorCode{Code: 6004, Message: "账号已在其他设备登录"}
	ErrWSBroadcastFailed  = ErrorCode{Code: 6005, Message: "消息广播失败"}
)

// ========== 文件上传错误码 (7000-7999) ==========

var (
	ErrFileUploadFailed    = ErrorCode{Code: 7000, Message: "文件上传失败"}
	ErrFileTypeNotAllowed  = ErrorCode{Code: 7001, Message: "文件类型不允许"}
	ErrFileSizeTooLarge    = ErrorCode{Code: 7002, Message: "文件大小超出限制"}
	ErrFileNotFound        = ErrorCode{Code: 7003, Message: "文件不存在"}
	ErrFileDeleteFailed    = ErrorCode{Code: 7004, Message: "文件删除失败"}
	ErrFileDownloadFailed  = ErrorCode{Code: 7005, Message: "文件下载失败"}
	ErrImageProcessFailed  = ErrorCode{Code: 7006, Message: "图片处理失败"}
	ErrFileStorageFull     = ErrorCode{Code: 7007, Message: "存储空间不足"}
	ErrInvalidFileFormat   = ErrorCode{Code: 7008, Message: "文件格式不正确"}
	ErrFileCorrupted       = ErrorCode{Code: 7009, Message: "文件已损坏"}
	ErrFileNameTooLong     = ErrorCode{Code: 7010, Message: "文件名过长"}
	ErrFilePathInvalid     = ErrorCode{Code: 7011, Message: "文件路径不正确"}
	ErrMultipartFormError  = ErrorCode{Code: 7012, Message: "解析表单数据失败"}
	ErrFileReadFailed      = ErrorCode{Code: 7013, Message: "文件读取失败"}
	ErrFileWriteFailed     = ErrorCode{Code: 7014, Message: "文件写入失败"}
	ErrDirectoryCreateFail = ErrorCode{Code: 7015, Message: "目录创建失败"}
)

// ========== 第三方服务错误码 (8000-8999) ==========

var (
	ErrThirdPartyService   = ErrorCode{Code: 8000, Message: "第三方服务错误"}
	ErrSMSServiceFailed    = ErrorCode{Code: 8001, Message: "短信服务失败"}
	ErrEmailServiceFailed  = ErrorCode{Code: 8002, Message: "邮件服务失败"}
	ErrPaymentFailed       = ErrorCode{Code: 8003, Message: "支付失败"}
	ErrPaymentTimeout      = ErrorCode{Code: 8004, Message: "支付超时"}
	ErrOSSUploadFailed     = ErrorCode{Code: 8005, Message: "OSS 上传失败"}
	ErrCacheFailed         = ErrorCode{Code: 8006, Message: "缓存服务失败"}
	ErrMessageQueueFailed  = ErrorCode{Code: 8007, Message: "消息队列失败"}
	ErrExternalAPIFailed   = ErrorCode{Code: 8008, Message: "外部 API 调用失败"}
	ErrExternalAPITimeout  = ErrorCode{Code: 8009, Message: "外部 API 超时"}
	ErrWeChatServiceFailed = ErrorCode{Code: 8010, Message: "微信服务失败"}
	ErrAlipayServiceFailed = ErrorCode{Code: 8011, Message: "支付宝服务失败"}
	ErrNetworkError        = ErrorCode{Code: 8012, Message: "网络连接错误"}
	ErrServiceUnavailable  = ErrorCode{Code: 8013, Message: "服务不可用"}
	ErrRateLimitExceeded   = ErrorCode{Code: 8014, Message: "调用频率超限"}
	ErrAPIKeyInvalid       = ErrorCode{Code: 8015, Message: "API 密钥无效"}
	ErrAuthorizationFailed = ErrorCode{Code: 8016, Message: "第三方授权失败"}
	ErrCallbackFailed      = ErrorCode{Code: 8017, Message: "回调处理失败"}
	ErrWebhookFailed       = ErrorCode{Code: 8018, Message: "Webhook 处理失败"}
	ErrDataSyncFailed      = ErrorCode{Code: 8019, Message: "数据同步失败"}
)

// ========== 业务逻辑错误码 (9000-9999) ==========

var (
	ErrBusinessLogic         = ErrorCode{Code: 9000, Message: "业务逻辑错误"}
	ErrOperationNotAllowed   = ErrorCode{Code: 9001, Message: "不允许的操作"}
	ErrDataInconsistency     = ErrorCode{Code: 9002, Message: "数据不一致"}
	ErrConcurrentConflict    = ErrorCode{Code: 9003, Message: "并发冲突"}
	ErrResourceLocked        = ErrorCode{Code: 9004, Message: "资源已被锁定"}
	ErrQuotaExceeded         = ErrorCode{Code: 9005, Message: "配额已超限"}
	ErrInventoryInsufficient = ErrorCode{Code: 9006, Message: "库存不足"}
	ErrOrderExpired          = ErrorCode{Code: 9007, Message: "订单已过期"}
	ErrOrderCanceled         = ErrorCode{Code: 9008, Message: "订单已取消"}
	ErrRefundFailed          = ErrorCode{Code: 9009, Message: "退款失败"}
	ErrStatusInvalid         = ErrorCode{Code: 9010, Message: "状态不正确"}
	ErrWorkflowError         = ErrorCode{Code: 9011, Message: "工作流错误"}
	ErrDependencyMissing     = ErrorCode{Code: 9012, Message: "依赖项缺失"}
	ErrConfigError           = ErrorCode{Code: 9013, Message: "配置错误"}
	ErrFeatureDisabled       = ErrorCode{Code: 9014, Message: "功能已禁用"}
	ErrMaintenanceMode       = ErrorCode{Code: 9015, Message: "系统维护中"}
	ErrVersionMismatch       = ErrorCode{Code: 9016, Message: "版本不匹配"}
	ErrDataExpired           = ErrorCode{Code: 9017, Message: "数据已过期"}
	ErrDuplicateOperation    = ErrorCode{Code: 9018, Message: "重复操作"}
	ErrPreconditionFailed    = ErrorCode{Code: 9019, Message: "前置条件不满足"}
)

// IsErrorCode 判断错误是否为 ErrorCode 类型
func IsErrorCode(err error) bool {
	_, ok := err.(ErrorCode)
	return ok
}

// GetErrorCode 从 error 中提取 ErrorCode
func GetErrorCode(err error) (ErrorCode, bool) {
	ec, ok := err.(ErrorCode)
	return ec, ok
}
