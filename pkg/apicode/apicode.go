package apicode

// Code 统一错误码（与 JSON 响应里的 code 字段一致）及默认中文文案。
// 约定：2xxxx 成功/通用；4xxxx 客户端；5xxxx 服务端。
type Code struct {
	Num int
	Msg string
}

// 常用成功
var OK = Code{200, "success"}

// 认证 / 鉴权（4xxxx）
var (
	AuthHeaderRequired = Code{40101, "缺少 Authorization 请求头"}
	AuthHeaderFormat   = Code{40102, "Authorization 头格式无效"}
	TokenInvalid       = Code{40103, "登录凭证无效或已过期"}

	PermissionDenied     = Code{40301, "无权限"}
	PermissionLoadFailed = Code{50001, "权限数据加载失败"}
)

// 通用请求（仍可与 gin binding 错误文案并存）
var (
	BadRequest       = Code{40001, "请求参数错误"}
	InvalidID        = Code{40002, "无效的 ID"}
	NotFound         = Code{40401, "资源不存在"}
	InternalError    = Code{50000, "服务器内部错误"}
	StoreAccountGone = Code{40402, "记账记录不存在"}

	StoreAccountEditTimeout = Code{40304, "账单创建超过24小时，无法修改"}
)
