package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger
var SugaredLogger *zap.SugaredLogger

// ErrorCode 错误码结构
type ErrorCode struct {
	Code    int
	Message string
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level"`       // 日志级别: debug, info, warn, error
	FilePath   string `yaml:"file_path"`   // 日志文件路径
	MaxSize    int    `yaml:"max_size"`    // 单个日志文件最大大小（MB）
	MaxBackups int    `yaml:"max_backups"` // 保留的旧日志文件最大数量
	MaxAge     int    `yaml:"max_age"`     // 保留旧日志文件的最大天数
	Compress   bool   `yaml:"compress"`    // 是否压缩旧日志
	Console    bool   `yaml:"console"`     // 是否同时输出到控制台
}

// InitLogger 初始化日志系统
func InitLogger(config *LogConfig) error {
	if config == nil {
		config = &LogConfig{
			Level:      "info",
			FilePath:   "logs/app.log",
			MaxSize:    100,
			MaxBackups: 10,
			MaxAge:     30,
			Compress:   true,
			Console:    true,
		}
	}

	// 解析日志级别
	level := parseLogLevel(config.Level)

	// 创建日志目录
	logDir := filepath.Dir(config.FilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 配置日志轮转
	fileWriter := &lumberjack.Logger{
		Filename:   config.FilePath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
		LocalTime:  true,
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// JSON 编码器（文件输出）
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// 控制台编码器（带颜色）
	consoleEncoderConfig := encoderConfig
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	// 创建 Core
	var cores []zapcore.Core

	// 文件输出
	cores = append(cores, zapcore.NewCore(
		jsonEncoder,
		zapcore.AddSync(fileWriter),
		level,
	))

	// 控制台输出
	if config.Console {
		cores = append(cores, zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			level,
		))
	}

	// 组合 Core
	core := zapcore.NewTee(cores...)

	// 创建 Logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	SugaredLogger = Logger.Sugar()

	return nil
}

// customTimeEncoder 自定义时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// parseLogLevel 解析日志级别
func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// CloseLogger 关闭日志系统（优雅退出时调用）
func CloseLogger() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

// ========== 便捷的日志方法 ==========

// LogDebug 调试日志
func LogDebug(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(msg, fields...)
	}
}

// LogInfo 信息日志
func LogInfo(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(msg, fields...)
	}
}

// LogWarn 警告日志
func LogWarn(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(msg, fields...)
	}
}

// LogError 错误日志
func LogError(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(msg, fields...)
	}
}

// LogFatal 致命错误日志（会终止程序）
func LogFatal(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Fatal(msg, fields...)
	}
}

// LogPanic 恐慌日志（会触发 panic）
func LogPanic(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Panic(msg, fields...)
	}
}

// ========== 格式化日志方法（支持 Printf 风格） ==========

// Debugf 格式化调试日志
func Debugf(template string, args ...interface{}) {
	if SugaredLogger != nil {
		SugaredLogger.Debugf(template, args...)
	}
}

// Infof 格式化信息日志
func Infof(template string, args ...interface{}) {
	if SugaredLogger != nil {
		SugaredLogger.Infof(template, args...)
	}
}

// Warnf 格式化警告日志
func Warnf(template string, args ...interface{}) {
	if SugaredLogger != nil {
		SugaredLogger.Warnf(template, args...)
	}
}

// Errorf 格式化错误日志
func Errorf(template string, args ...interface{}) {
	if SugaredLogger != nil {
		SugaredLogger.Errorf(template, args...)
	}
}

// Fatalf 格式化致命错误日志
func Fatalf(template string, args ...interface{}) {
	if SugaredLogger != nil {
		SugaredLogger.Fatalf(template, args...)
	}
}

// ========== 带上下文的日志方法 ==========

// InfoWithFields 带字段的信息日志
func InfoWithFields(msg string, fields map[string]interface{}) {
	if SugaredLogger != nil {
		SugaredLogger.Infow(msg, convertFields(fields)...)
	}
}

// ErrorWithFields 带字段的错误日志
func ErrorWithFields(msg string, fields map[string]interface{}) {
	if SugaredLogger != nil {
		SugaredLogger.Errorw(msg, convertFields(fields)...)
	}
}

// WarnWithFields 带字段的警告日志
func WarnWithFields(msg string, fields map[string]interface{}) {
	if SugaredLogger != nil {
		SugaredLogger.Warnw(msg, convertFields(fields)...)
	}
}

// convertFields 转换字段格式
func convertFields(fields map[string]interface{}) []interface{} {
	result := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		result = append(result, k, v)
	}
	return result
}

// ========== 业务日志快捷方法 ==========

// LogRequest 记录 HTTP 请求日志
func LogRequest(method, path, ip string, statusCode int, latency time.Duration) {
	LogInfo("HTTP Request",
		zap.String("method", method),
		zap.String("path", path),
		zap.String("ip", ip),
		zap.Int("status", statusCode),
		zap.Duration("latency", latency),
	)
}

// LogBusinessError 记录业务错误日志（带错误码）
func LogBusinessError(errCode ErrorCode, err error, fields ...zap.Field) {
	allFields := append([]zap.Field{
		zap.Int("error_code", errCode.Code),
		zap.String("error_msg", errCode.Message),
	}, fields...)

	if err != nil {
		allFields = append(allFields, zap.Error(err))
	}

	LogError("Business Error", allFields...)
}

// LogDatabaseError 记录数据库错误
func LogDatabaseError(operation string, err error, fields ...zap.Field) {
	allFields := append([]zap.Field{
		zap.String("operation", operation),
		zap.Error(err),
	}, fields...)

	LogError("Database Error", allFields...)
}

// LogAuthError 记录认证错误
func LogAuthError(action string, userID uint, reason string) {
	LogWarn("Authentication Error",
		zap.String("action", action),
		zap.Uint("user_id", userID),
		zap.String("reason", reason),
	)
}

// LogBusinessEvent 记录业务事件
func LogBusinessEvent(event string, fields ...zap.Field) {
	LogInfo("Business Event",
		append([]zap.Field{zap.String("event", event)}, fields...)...,
	)
}

// LogPerformance 记录性能日志
func LogPerformance(operation string, duration time.Duration, fields ...zap.Field) {
	allFields := append([]zap.Field{
		zap.String("operation", operation),
		zap.Duration("duration", duration),
	}, fields...)

	if duration > time.Second {
		LogWarn("Slow Operation", allFields...)
	} else {
		LogDebug("Performance", allFields...)
	}
}

// LogSQL 记录 SQL 执行日志
func LogSQL(query string, duration time.Duration, rows int64) {
	LogDebug("SQL Execution",
		zap.String("query", query),
		zap.Duration("duration", duration),
		zap.Int64("rows", rows),
	)
}

// LogCacheOperation 记录缓存操作
func LogCacheOperation(operation, key string, hit bool) {
	LogDebug("Cache Operation",
		zap.String("operation", operation),
		zap.String("key", key),
		zap.Bool("hit", hit),
	)
}

// LogWebSocket 记录 WebSocket 事件
func LogWebSocket(event string, userID uint, fields ...zap.Field) {
	allFields := append([]zap.Field{
		zap.String("event", event),
		zap.Uint("user_id", userID),
	}, fields...)

	LogInfo("WebSocket Event", allFields...)
}

// LogThirdParty 记录第三方服务调用
func LogThirdParty(service, action string, success bool, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("service", service),
		zap.String("action", action),
		zap.Bool("success", success),
		zap.Duration("duration", duration),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		LogError("Third Party Service", fields...)
	} else {
		LogInfo("Third Party Service", fields...)
	}
}
