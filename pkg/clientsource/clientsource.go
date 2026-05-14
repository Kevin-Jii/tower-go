// Package clientsource 约定请求头 X-Client-Source，用于区分调用端（小程序 / H5 / RN / 管理后台等）。
//
// 与小程序（Taro）工程约定一致时，典型取值：
//   weapp   — 微信小程序（TARO_ENV=weapp）
//   web     — H5（TARO_ENV=h5）
//   app     — React Native 独立 App（TARO_ENV=rn）
//   其他    — 其他 Taro 平台可透传 TARO_ENV 字符串
//   unknown — 未传或空
//
// 本仓库 web-admin（Vite 浏览器端）默认发送 web-admin，可通过环境变量 VITE_CLIENT_SOURCE 覆盖。
package clientsource

import (
	"net/http"
	"strings"
)

// HeaderName 为统一请求头名（大小写不敏感，此处使用规范写法）。
const HeaderName = "X-Client-Source"

// 常用取值（非穷举，服务端接受任意非空字符串）
const (
	SourceWeapp    = "weapp"
	SourceWeb      = "web"
	SourceApp      = "app"
	SourceWebAdmin = "web-admin"
	SourceUnknown  = "unknown"
)

// Parse 从原始头字段解析调用端标识；空或纯空白为 unknown。
func Parse(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return SourceUnknown
	}
	return s
}

// FromRequest 读取 HTTP 请求中的 X-Client-Source。
func FromRequest(r *http.Request) string {
	if r == nil {
		return SourceUnknown
	}
	return Parse(r.Header.Get(HeaderName))
}
