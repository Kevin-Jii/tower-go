package context

import (
	"time"

	"go.uber.org/zap"
	"github.com/Kevin-Jii/tower-go/pkg/logging/logger"
)

// LogRequest 记录 HTTP 请求日志
func LogRequest(l *logger.Logger, method, path, clientIP string, statusCode int, latency time.Duration, userAgent string) {
	fields := []zap.Field{
		zap.String("method", method),
		zap.String("path", path),
		zap.String("client_ip", clientIP),
		zap.Int("status_code", statusCode),
		zap.Duration("latency", latency),
		zap.String("user_agent", userAgent),
	}

	// 根据状态码选择日志级别
	if statusCode >= 500 {
		l.Error("HTTP Server Error", fields...)
	} else if statusCode >= 400 {
		l.Warn("HTTP Client Error", fields...)
	} else {
		l.Info("HTTP Request", fields...)
	}
}

// LogBusinessError 记录业务错误日志
func LogBusinessError(l *logger.Logger, errCode logger.ErrorCode, err error, context map[string]interface{}) {
	fields := []zap.Field{
		zap.Int("error_code", errCode.Code),
		zap.String("error_message", errCode.Message),
	}

	// 添加上下文字段
	for k, v := range context {
		fields = append(fields, zap.Any(k, v))
	}

	// 添加错误详情
	if err != nil {
		fields = append(fields, zap.Error(err))
	}

	l.Error("Business Error", fields...)
}

// LogDatabaseError 记录数据库错误
func LogDatabaseError(l *logger.Logger, operation string, table string, err error, duration time.Duration) {
	fields := []zap.Field{
		zap.String("operation", operation),
		zap.String("table", table),
		zap.Duration("duration", duration),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		l.Error("Database Error", fields...)
	} else {
		l.Debug("Database Operation", fields...)
	}
}

// LogAuthEvent 记录认证事件
func LogAuthEvent(l *logger.Logger, event string, userID uint, username, clientIP string, success bool, reason string) {
	fields := []zap.Field{
		zap.String("event", event),
		zap.Uint("user_id", userID),
		zap.String("username", username),
		zap.String("client_ip", clientIP),
		zap.Bool("success", success),
	}

	if reason != "" {
		fields = append(fields, zap.String("reason", reason))
	}

	if success {
		l.Info("Authentication Success", fields...)
	} else {
		l.Warn("Authentication Failure", fields...)
	}
}

// LogBusinessEvent 记录业务事件
func LogBusinessEvent(l *logger.Logger, event string, userID uint, fields map[string]interface{}) {
	logFields := []zap.Field{
		zap.String("event", event),
		zap.Uint("user_id", userID),
	}

	// 添加额外字段
	for k, v := range fields {
		logFields = append(logFields, zap.Any(k, v))
	}

	l.Info("Business Event", logFields...)
}

// LogPerformance 记录性能日志
func LogPerformance(l *logger.Logger, operation string, duration time.Duration, thresholds map[string]time.Duration, fields map[string]interface{}) {
	logFields := []zap.Field{
		zap.String("operation", operation),
		zap.Duration("duration", duration),
	}

	// 添加阈值信息
	for name, threshold := range thresholds {
		logFields = append(logFields, zap.Duration(name+"_threshold", threshold))
	}

	// 添加额外字段
	for k, v := range fields {
		logFields = append(logFields, zap.Any(k, v))
	}

	// 根据持续时间选择日志级别
	var isSlow bool
	for _, threshold := range thresholds {
		if duration > threshold {
			isSlow = true
			break
		}
	}

	if isSlow {
		logFields = append(logFields, zap.Bool("slow_operation", true))
		l.Warn("Slow Performance", logFields...)
	} else {
		l.Debug("Performance", logFields...)
	}
}

// LogSQLOperation 记录SQL操作
func LogSQLOperation(l *logger.Logger, query string, duration time.Duration, rowsAffected int64, err error) {
	fields := []zap.Field{
		zap.String("query", query),
		zap.Duration("duration", duration),
		zap.Int64("rows_affected", rowsAffected),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		l.Error("SQL Error", fields...)
	} else {
		// 慢查询检查
		if duration > time.Second {
			l.Warn("Slow SQL Query", fields...)
		} else {
			l.Debug("SQL Query", fields...)
		}
	}
}

// LogCacheOperation 记录缓存操作
func LogCacheOperation(l *logger.Logger, operation, key string, hit bool, ttl time.Duration) {
	fields := []zap.Field{
		zap.String("operation", operation),
		zap.String("key", key),
		zap.Bool("hit", hit),
	}

	if ttl > 0 {
		fields = append(fields, zap.Duration("ttl", ttl))
	}

	l.Debug("Cache Operation", fields...)
}

// LogWebSocketEvent 记录WebSocket事件
func LogWebSocketEvent(l *logger.Logger, event string, userID uint, connectionID string, fields map[string]interface{}) {
	logFields := []zap.Field{
		zap.String("event", event),
		zap.Uint("user_id", userID),
		zap.String("connection_id", connectionID),
	}

	// 添加额外字段
	for k, v := range fields {
		logFields = append(logFields, zap.Any(k, v))
	}

	l.Info("WebSocket Event", logFields...)
}

// LogThirdPartyCall 记录第三方服务调用
func LogThirdPartyCall(l *logger.Logger, service, endpoint, method string, statusCode int, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("service", service),
		zap.String("endpoint", endpoint),
		zap.String("method", method),
		zap.Int("status_code", statusCode),
		zap.Duration("duration", duration),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		l.Error("Third Party Service Error", fields...)
	} else {
		// 根据状态码选择日志级别
		if statusCode >= 500 {
			l.Warn("Third Party Service Error", fields...)
		} else if statusCode >= 400 {
			l.Warn("Third Party Client Error", fields...)
		} else {
			l.Info("Third Party Service Success", fields...)
		}
	}
}

// LogSecurityEvent 记录安全事件
func LogSecurityEvent(l *logger.Logger, event string, userID uint, clientIP, userAgent string, severity string, details map[string]interface{}) {
	fields := []zap.Field{
		zap.String("event", event),
		zap.Uint("user_id", userID),
		zap.String("client_ip", clientIP),
		zap.String("user_agent", userAgent),
		zap.String("severity", severity),
	}

	// 添加详细信息
	for k, v := range details {
		fields = append(fields, zap.Any(k, v))
	}

	// 根据严重程度选择日志级别
	switch severity {
	case "high", "critical":
		l.Error("Security Event", fields...)
	case "medium":
		l.Warn("Security Event", fields...)
	case "low":
		l.Info("Security Event", fields...)
	default:
		l.Info("Security Event", fields...)
	}
}

// LogSystemEvent 记录系统事件
func LogSystemEvent(l *logger.Logger, component, event string, level logger.Level, fields map[string]interface{}) {
	logFields := []zap.Field{
		zap.String("component", component),
		zap.String("event", event),
	}

	// 添加额外字段
	for k, v := range fields {
		logFields = append(logFields, zap.Any(k, v))
	}

	// 根据级别记录日志
	switch level {
	case logger.DebugLevel:
		l.Debug("System Event", logFields...)
	case logger.InfoLevel:
		l.Info("System Event", logFields...)
	case logger.WarnLevel:
		l.Warn("System Event", logFields...)
	case logger.ErrorLevel:
		l.Error("System Event", logFields...)
	case logger.FatalLevel:
		l.Fatal("System Event", logFields...)
	default:
		l.Info("System Event", logFields...)
	}
}

// LogMetrics 记录指标日志
func LogMetrics(l *logger.Logger, metrics map[string]interface{}) {
	fields := make([]zap.Field, 0, len(metrics))
	for k, v := range metrics {
		fields = append(fields, zap.Any(k, v))
	}

	l.Debug("Metrics", fields...)
}

// LogAudit 记录审计日志
func LogAudit(l *logger.Logger, action string, resource string, userID uint, details map[string]interface{}) {
	fields := []zap.Field{
		zap.String("action", action),
		zap.String("resource", resource),
		zap.Uint("user_id", userID),
	}

	// 添加详细信息
	for k, v := range details {
		fields = append(fields, zap.Any(k, v))
	}

	l.Info("Audit Log", fields...)
}