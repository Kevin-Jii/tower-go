package middleware

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/utils/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryMiddleware 统一错误恢复中间件
// 自动捕获 panic，记录日志，返回统一错误响应
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取堆栈信息
				stack := getStackTrace(3)

				// 记录错误日志
				logging.LogError("Panic Recovered",
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.Any("error", err),
					zap.String("stack", stack),
				)

				// 返回统一错误响应
				c.JSON(500, gin.H{
					"code":    500,
					"message": "系统内部错误，请稍后重试",
				})

				// 中断请求
				c.Abort()
			}
		}()
		c.Next()
	}
}

// getStackTrace 获取堆栈信息
func getStackTrace(skip int) string {
	var builder strings.Builder
	builder.WriteString("\nStack trace:\n")

	for i := skip; i < skip+10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		if fn != nil {
			builder.WriteString(fmt.Sprintf("  %s:%d (%s)\n", file, line, fn.Name()))
		}
	}

	return builder.String()
}

// TimeoutMiddleware 请求超时控制（简化版）
// 注意：实际超时控制建议在 Nginx 或负载均衡器层面实现
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 排除健康检查和 WebSocket 接口
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/ready" || c.Request.URL.Path == "/live" || c.Request.URL.Path == "/ws" {
			c.Next()
			return
		}

		// Go 1.21+ 支持 ContextWithTimeout
		// ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		// defer cancel()
		// c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// RateLimitMiddleware 请求频率限制（简化版）
// 实际生产环境建议使用 Redis 实现分布式限流
func RateLimitMiddleware(requestsPerSecond int) gin.HandlerFunc {
	// 这里使用简单的令牌桶简化实现
	// 建议生产环境使用 https://github.com/go-redis/redis 结合 Lua 脚本
	return func(c *gin.Context) {
		// 暂时跳过，后续可集成 Redis 实现
		c.Next()
	}
}

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// SecurityHeadersMiddleware 安全头中间件
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止 XSS 攻击
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")

		// 内容安全策略
		c.Header("Content-Security-Policy", "default-src 'self'")

		c.Next()
	}
}

// RequestIDMiddleware 请求追踪ID
// 为每个请求生成唯一ID，便于日志追踪
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Set("requestID", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), randomString(8))
}

// randomString 生成随机字符串
func randomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
