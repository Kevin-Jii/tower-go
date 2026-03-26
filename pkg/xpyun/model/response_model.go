package model

// XPYunResp HTTP响应
type XPYunResp struct {
	HttpStatusCode int
	Content        *XPYunRespContent
}

// XPYunRespContent 响应内容
type XPYunRespContent struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data,omitempty"`
	OrderId string      `json:"orderId,omitempty"`
}

// IsSuccess 检查请求是否成功
func (r *XPYunRespContent) IsSuccess() bool {
	return r.Code == 0
}