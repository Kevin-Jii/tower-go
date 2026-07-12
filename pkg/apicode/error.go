package apicode

import (
	"errors"
	"fmt"
)

// ErrorCode 约定：错误码使用 HTTP 状态码前三位 + 两位业务序号。
// 例如 40410 表示“资源不存在”类错误中的第 10 个错误。
// 这样既方便前端按 HTTP 语义处理，也避免业务层散落裸数字。
type ErrorCode struct {
	Code  Code
	Cause error
}

// Error 允许简单场景直接返回 Code 作为 error。
func (c Code) Error() string {
	return fmt.Sprintf("[%d] %s", c.Num, c.Msg)
}

// Error 实现 error 接口。底层原因通过 Unwrap 保留给日志和 errors.Is 使用。
func (e *ErrorCode) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("[%d] %s", e.Code.Num, e.Code.Msg)
}

// Unwrap 返回底层错误。
func (e *ErrorCode) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// New 创建一个不带底层原因的业务错误。
func New(code Code) error {
	return &ErrorCode{Code: code}
}

// Newf 创建一个使用动态提示信息的业务错误，错误码保持不变。
func Newf(code Code, format string, args ...interface{}) error {
	return &ErrorCode{Code: code.WithMessage(fmt.Sprintf(format, args...))}
}

// Wrap 使用业务错误码包装底层错误。
func Wrap(code Code, cause error) error {
	if cause == nil {
		return New(code)
	}
	return &ErrorCode{Code: code, Cause: cause}
}

// WithMessage 返回带自定义提示信息的错误码。
func (c Code) WithMessage(message string) Code {
	c.Msg = message
	return c
}

// WithMessageF 返回带格式化提示信息的错误码。
func (c Code) WithMessageF(format string, args ...interface{}) Code {
	return c.WithMessage(fmt.Sprintf(format, args...))
}

// Resolve 从 error 中解析统一错误码。第二个返回值表示是否为已识别的业务错误。
func Resolve(err error) (Code, bool) {
	if err == nil {
		return OK, false
	}

	var appErr *ErrorCode
	if errors.As(err, &appErr) && appErr != nil {
		return appErr.Code, true
	}

	// 允许直接把 Code 作为 error 返回，便于简单场景使用。
	var code Code
	if errors.As(err, &code) {
		return code, true
	}
	return InternalError, false
}

// CodeOf 获取错误码。未识别的底层错误统一归类为服务器内部错误。
func CodeOf(err error) Code {
	code, _ := Resolve(err)
	return code
}

// Is 判断错误是否属于指定业务错误码。
func Is(err error, code Code) bool {
	actual, ok := Resolve(err)
	return ok && actual.Num == code.Num
}

// HTTPStatus 将规范错误码映射为 HTTP 状态码。
// 旧的 400/401/500 等 HTTP 状态码仍由响应层原逻辑处理。
func HTTPStatus(code int) int {
	if code >= 10000 {
		status := code / 100
		if status >= 400 && status <= 599 {
			return status
		}
	}
	if code >= 100 && code <= 599 {
		return code
	}
	return 0
}
