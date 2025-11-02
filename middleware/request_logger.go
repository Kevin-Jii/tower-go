package middleware

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLoggerMiddleware 打印所有进入的请求参数信息
// 包括: method, path, query, path params, body(限长), clientIP, latency
// 注意：读取 body 后会重新写回以保证后续 handler 能继续绑定 JSON
func RequestLoggerMiddleware(maxBody int) gin.HandlerFunc {
	if maxBody <= 0 {
		maxBody = 2048
	}
	return func(c *gin.Context) {
		start := time.Now()

		// 读取 Body（仅在可能有内容时）
		var bodySnippet string
		if c.Request.Body != nil {
			data, _ := io.ReadAll(c.Request.Body)
			// 还原 Body 供后续使用
			c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
			if len(data) > maxBody {
				bodySnippet = string(data[:maxBody]) + "..."
			} else {
				bodySnippet = string(data)
			}
		}

		// 组装日志
		log.Printf("[REQ] %s %s query=%q pathParams=%v body=%s client=%s", c.Request.Method, c.Request.URL.Path, c.Request.URL.RawQuery, c.Params, bodySnippet, c.ClientIP())

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		log.Printf("[RESP] %d %s %s latency=%s", status, c.Request.Method, c.Request.URL.Path, latency)
	}
}
