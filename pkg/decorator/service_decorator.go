package decorator

import (
	"fmt"
	"time"

	"github.com/Kevin-Jii/tower-go/utils/logging"
)

// ServiceFunc 服务函数类型
type ServiceFunc func() (interface{}, error)

// LoggingDecorator 日志装饰器
func LoggingDecorator(name string, fn ServiceFunc) ServiceFunc {
	return func() (interface{}, error) {
		start := time.Now()
		result, err := fn()
		duration := time.Since(start)

		if err != nil {
			logging.LogError(fmt.Sprintf("[%s] failed in %v: %v", name, duration, err))
		} else {
			logging.LogInfo(fmt.Sprintf("[%s] completed in %v", name, duration))
		}
		return result, err
	}
}

// TimingDecorator 计时装饰器
func TimingDecorator(fn ServiceFunc) (interface{}, error, time.Duration) {
	start := time.Now()
	result, err := fn()
	return result, err, time.Since(start)
}

// RetryDecorator 重试装饰器
func RetryDecorator(maxRetries int, delay time.Duration, fn ServiceFunc) ServiceFunc {
	return func() (interface{}, error) {
		var lastErr error
		for i := 0; i < maxRetries; i++ {
			result, err := fn()
			if err == nil {
				return result, nil
			}
			lastErr = err
			if i < maxRetries-1 {
				time.Sleep(delay)
			}
		}
		return nil, fmt.Errorf("after %d retries: %w", maxRetries, lastErr)
	}
}
