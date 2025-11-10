package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Level 日志级别
type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	DPanicLevel
	PanicLevel
	FatalLevel
)

// String 返回日志级别字符串
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case DPanicLevel:
		return "dpanic"
	case PanicLevel:
		return "panic"
	case FatalLevel:
		return "fatal"
	default:
		return "info"
	}
}

// ParseLevel 解析日志级别字符串
func ParseLevel(s string) Level {
	switch s {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn", "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "dpanic":
		return DPanicLevel
	case "panic":
		return PanicLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}

// ToZapLevel 转换为zap日志级别
func (l Level) ToZapLevel() zapcore.Level {
	return zapcore.Level(l)
}

// Encoder 编码器类型
type Encoder string

const (
	JSONEncoder    Encoder = "json"
	ConsoleEncoder Encoder = "console"
)

// Config 日志配置
type Config struct {
	Level              Level    `yaml:"level" json:"level"`
	Encoder            Encoder  `yaml:"encoder" json:"encoder"`
	FilePath           string   `yaml:"file_path" json:"file_path"`
	MaxSize            int      `yaml:"max_size" json:"max_size"`
	MaxBackups         int      `yaml:"max_backups" json:"max_backups"`
	MaxAge             int      `yaml:"max_age" json:"max_age"`
	Compress           bool     `yaml:"compress" json:"compress"`
	Console            bool     `yaml:"console" json:"console"`
	EnableCaller       bool     `yaml:"enable_caller" json:"enable_caller"`
	EnableStackTrace   bool     `yaml:"enable_stack_trace" json:"enable_stack_trace"`
	CallerSkip         int      `yaml:"caller_skip" json:"caller_skip"`
	TimeFormat         string   `yaml:"time_format" json:"time_format"`
	OutputPaths        []string `yaml:"output_paths" json:"output_paths"`
	ErrorOutputPaths   []string `yaml:"error_output_paths" json:"error_output_paths"`
	InitialFields      map[string]interface{} `yaml:"initial_fields" json:"initial_fields"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:            InfoLevel,
		Encoder:          JSONEncoder,
		FilePath:         "logs/app.log",
		MaxSize:          100,
		MaxBackups:       10,
		MaxAge:           30,
		Compress:         true,
		Console:          true,
		EnableCaller:     true,
		EnableStackTrace: true,
		CallerSkip:       1,
		TimeFormat:       "2006-01-02 15:04:05.000",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    make(map[string]interface{}),
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.MaxSize <= 0 {
		return fmt.Errorf("max_size must be positive")
	}
	if c.MaxBackups < 0 {
		return fmt.Errorf("max_backups cannot be negative")
	}
	if c.MaxAge < 0 {
		return fmt.Errorf("max_age cannot be negative")
	}
	if c.CallerSkip < 0 {
		return fmt.Errorf("caller_skip cannot be negative")
	}
	if c.TimeFormat == "" {
		c.TimeFormat = "2006-01-02 15:04:05.000"
	}
	return nil
}

// ErrorCode 错误码结构
type ErrorCode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Logger 日志管理器
type Logger struct {
	config       *Config
	zapLogger    *zap.Logger
	sugaredLogger *zap.SugaredLogger
}

// New 创建日志管理器
func New(config *Config) (*Logger, error) {
	if config == nil {
		config = DefaultConfig()
	}
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid logger config: %w", err)
	}

	zapLogger, err := buildZapLogger(config)
	if err != nil {
		return nil, fmt.Errorf("failed to build zap logger: %w", err)
	}

	return &Logger{
		config:        config,
		zapLogger:     zapLogger,
		sugaredLogger: zapLogger.Sugar(),
	}, nil
}

// NewWithDefaults 使用默认配置创建日志管理器
func NewWithDefaults() (*Logger, error) {
	return New(DefaultConfig())
}

// NewDevelopment 创建开发环境日志管理器
func NewDevelopment() (*Logger, error) {
	config := DefaultConfig()
	config.Level = DebugLevel
	config.Encoder = ConsoleEncoder
	config.Console = true
	config.EnableCaller = true
	config.EnableStackTrace = true
	return New(config)
}

// NewProduction 创建生产环境日志管理器
func NewProduction() (*Logger, error) {
	config := DefaultConfig()
	config.Level = InfoLevel
	config.Encoder = JSONEncoder
	config.Console = false
	config.EnableCaller = false
	config.EnableStackTrace = false
	return New(config)
}

// buildZapLogger 构建zap日志器
func buildZapLogger(config *Config) (*zap.Logger, error) {
	// 解析日志级别
	level := config.Level.ToZapLevel()

	// 创建日志目录
	if config.FilePath != "" {
		logDir := filepath.Dir(config.FilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}
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
		EncodeTime:     customTimeEncoder(config.TimeFormat),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 根据编码器类型配置编码器
	var encoder zapcore.Encoder
	switch config.Encoder {
	case JSONEncoder:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case ConsoleEncoder:
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 创建输出写入器
	var writers []zapcore.Core

	// 文件输出
	if config.FilePath != "" {
		fileWriter := &lumberjack.Logger{
			Filename:   config.FilePath,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
			LocalTime:  true,
		}
		writers = append(writers, zapcore.NewCore(encoder, zapcore.AddSync(fileWriter), level))
	}

	// 控制台输出
	if config.Console {
		writers = append(writers, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
	}

	// 自定义输出路径
	for _, path := range config.OutputPaths {
		if path == "stdout" {
			writers = append(writers, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
		} else if path == "stderr" {
			writers = append(writers, zapcore.NewCore(encoder, zapcore.AddSync(os.Stderr), level))
		} else {
			file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err == nil {
				writers = append(writers, zapcore.NewCore(encoder, zapcore.AddSync(file), level))
			}
		}
	}

	if len(writers) == 0 {
		writers = append(writers, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
	}

	// 组合核心
	core := zapcore.NewTee(writers...)

	// 构建选项
	var options []zap.Option
	if config.EnableCaller {
		options = append(options, zap.AddCaller())
		if config.CallerSkip > 0 {
			options = append(options, zap.AddCallerSkip(config.CallerSkip))
		}
	}
	if config.EnableStackTrace {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	// 添加初始字段
	if len(config.InitialFields) > 0 {
		options = append(options, zap.Fields(convertToZapFields(config.InitialFields)...))
	}

	// 创建日志器
	return zap.New(core, options...), nil
}

// customTimeEncoder 自定义时间编码器
func customTimeEncoder(format string) zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(format))
	}
}

// convertToZapFields 转换字段为zap字段
func convertToZapFields(fields map[string]interface{}) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}

// ========== 基础日志方法 ==========

// Debug 调试日志
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.zapLogger.Debug(msg, fields...)
}

// Info 信息日志
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.zapLogger.Info(msg, fields...)
}

// Warn 警告日志
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.zapLogger.Warn(msg, fields...)
}

// Error 错误日志
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.zapLogger.Error(msg, fields...)
}

// Fatal 致命错误日志
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.zapLogger.Fatal(msg, fields...)
}

// Panic 恐慌日志
func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.zapLogger.Panic(msg, fields...)
}

// ========== 格式化日志方法 ==========

// Debugf 格式化调试日志
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.sugaredLogger.Debugf(template, args...)
}

// Infof 格式化信息日志
func (l *Logger) Infof(template string, args ...interface{}) {
	l.sugaredLogger.Infof(template, args...)
}

// Warnf 格式化警告日志
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugaredLogger.Warnf(template, args...)
}

// Errorf 格式化错误日志
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugaredLogger.Errorf(template, args...)
}

// Fatalf 格式化致命错误日志
func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.sugaredLogger.Fatalf(template, args...)
}

// Panicf 格式化恐慌日志
func (l *Logger) Panicf(template string, args ...interface{}) {
	l.sugaredLogger.Panicf(template, args...)
}

// ========== 结构化日志方法 ==========

// Debugw 带字段的调试日志
func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Debugw(msg, keysAndValues...)
}

// Infow 带字段的信息日志
func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Infow(msg, keysAndValues...)
}

// Warnw 带字段的警告日志
func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Warnw(msg, keysAndValues...)
}

// Errorw 带字段的错误日志
func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Errorw(msg, keysAndValues...)
}

// Fatalw 带字段的致命错误日志
func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Fatalw(msg, keysAndValues...)
}

// Panicw 带字段的恐慌日志
func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Panicw(msg, keysAndValues...)
}

// ========== 工具方法 ==========

// Sync 同步缓冲区
func (l *Logger) Sync() error {
	return l.zapLogger.Sync()
}

// GetZapLogger 获取原生zap日志器
func (l *Logger) GetZapLogger() *zap.Logger {
	return l.zapLogger
}

// GetSugaredLogger 获取sugared日志器
func (l *Logger) GetSugaredLogger() *zap.SugaredLogger {
	return l.sugaredLogger
}

// GetConfig 获取配置
func (l *Logger) GetConfig() *Config {
	return l.config
}

// SetLevel 动态设置日志级别
func (l *Logger) SetLevel(level Level) {
	l.config.Level = level
	// 注意：zap不支持动态修改日志级别，需要重新创建logger
}

// With 添加字段返回新的日志器
func (l *Logger) With(fields ...zap.Field) *Logger {
	newZapLogger := l.zapLogger.With(fields...)
	return &Logger{
		config:         l.config,
		zapLogger:      newZapLogger,
		sugaredLogger:  newZapLogger.Sugar(),
	}
}

// Named 创建命名的子日志器
func (l *Logger) Named(name string) *Logger {
	newZapLogger := l.zapLogger.Named(name)
	return &Logger{
		config:         l.config,
		zapLogger:      newZapLogger,
		sugaredLogger:  newZapLogger.Sugar(),
	}
}

// 全局默认日志器
var defaultLogger *Logger

// InitDefault 初始化默认日志器
func InitDefault(config *Config) error {
	var err error
	defaultLogger, err = New(config)
	return err
}

// GetDefault 获取默认日志器
func GetDefault() *Logger {
	if defaultLogger == nil {
		defaultLogger, _ = NewWithDefaults()
	}
	return defaultLogger
}

// 全局便利函数
func Debug(msg string, fields ...zap.Field) {
	GetDefault().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	GetDefault().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	GetDefault().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	GetDefault().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	GetDefault().Fatal(msg, fields...)
}

func Debugf(template string, args ...interface{}) {
	GetDefault().Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	GetDefault().Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	GetDefault().Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	GetDefault().Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	GetDefault().Fatalf(template, args...)
}