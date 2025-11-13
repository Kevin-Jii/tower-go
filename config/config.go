package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config 应用配置
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	DingTalk DingTalkConfig
}

// AppConfig 应用配置
type AppConfig struct {
	Name string
	Port int
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	Enabled  bool
}

// DingTalkWebhookConfig 钉钉Webhook配置
type DingTalkWebhookConfig struct {
	MenuReportURL string // 报菜记录通知地址
}

// DingTalkConfig 钉钉配置
type DingTalkConfig struct {
	Webhook DingTalkWebhookConfig
	Stream  DingTalkStreamConfig
}

// DingTalkStreamConfig 钉钉Stream配置
type DingTalkStreamConfig struct {
	ClientID     string
	ClientSecret string
	AgentID      string
	MiniAppID    string
}

// 全局配置实例
var cfg *Config

// InitConfig 初始化配置
func InitConfig() {
	// 首先加载.env文件
	loadEnvFile()

	cfg = &Config{
		App:      loadAppConfig(),
		Database: loadDatabaseConfig(),
		Redis:    loadRedisConfig(),
		DingTalk: loadDingTalkConfig(),
	}
}

// loadEnvFile 加载.env文件
func loadEnvFile() {
	envPath := filepath.Join(".", ".env")

	// 检查项目根目录下的.env文件
	if _, err := os.Stat(envPath); err == nil {
		loadEnvFromFile(envPath)
		return
	}

	// 检查config目录下的.env文件
	envPath = filepath.Join("config", ".env")
	if _, err := os.Stat(envPath); err == nil {
		loadEnvFromFile(envPath)
		return
	}
}

// loadEnvFromFile 从文件加载环境变量
func loadEnvFromFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 去除值的双引号
		if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
			value = value[1 : len(value)-1]
		}

		// 设置环境变量
		os.Setenv(key, value)
	}
}

// GetConfig 获取配置
func GetConfig() *Config {
	if cfg == nil {
		InitConfig()
	}
	return cfg
}

// GetDatabaseConfig 获取数据库配置
func GetDatabaseConfig() DatabaseConfig {
	return GetConfig().Database
}

// GetRedisConfig 获取Redis配置
func GetRedisConfig() RedisConfig {
	return GetConfig().Redis
}

// GetDingTalkConfig 获取钉钉配置
func GetDingTalkConfig() DingTalkConfig {
	return GetConfig().DingTalk
}

// GetDingTalkStreamConfig 获取钉钉Stream配置
func GetDingTalkStreamConfig() DingTalkStreamConfig {
	return GetConfig().DingTalk.Stream
}

// GetDingTalkMenuReportURL 获取钉钉报菜记录通知地址
func GetDingTalkMenuReportURL() string {
	return GetConfig().DingTalk.Webhook.MenuReportURL
}

// loadAppConfig 加载应用配置
func loadAppConfig() AppConfig {
	return AppConfig{
		Name: getAppString("APP_NAME", "tower-go"),
		Port: getAppInt("APP_PORT", 10024),
	}
}

// DebugConfig 打印当前配置（仅用于调试）
func DebugConfig() {
	// 这里我们直接从环境变量读取，避免循环依赖
	dbPassword := os.Getenv("DB_PASSWORD")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	if dbPassword != "" {
		// 不打印真实密码，只确认是否读取到
		println("✅ DB_PASSWORD: loaded")
	} else {
		println("❌ DB_PASSWORD: not found")
	}

	if redisPassword != "" {
		// 不打印真实密码，只确认是否读取到
		println("✅ REDIS_PASSWORD: loaded")
	} else {
		println("❌ REDIS_PASSWORD: not found")
	}
}

// loadDatabaseConfig 加载数据库配置
func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Driver:   getAppString("DB_DRIVER", "mysql"),
		Host:     getAppString("DB_HOST", "localhost"),
		Port:     getAppInt("DB_PORT", 3306),
		Username: getAppString("DB_USERNAME", "root"),
		Password: getAppString("DB_PASSWORD", ""),
		Database: getAppString("DB_NAME", "tower"),
	}
}

// loadRedisConfig 加载Redis配置
func loadRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     getAppString("REDIS_HOST", "localhost"),
		Port:     getAppInt("REDIS_PORT", 6379),
		Password: getAppString("REDIS_PASSWORD", ""),
		DB:       getAppInt("REDIS_DB", 0),
		Enabled:  getAppBool("REDIS_ENABLED", true),
	}
}

// loadDingTalkConfig 加载钉钉配置
func loadDingTalkConfig() DingTalkConfig {
	return DingTalkConfig{
		Webhook: DingTalkWebhookConfig{
			MenuReportURL: getAppString("DINGTALK_MENU_REPORT_WEBHOOK_URL", ""),
		},
		Stream: DingTalkStreamConfig{
			ClientID:     getAppString("DINGTALK_CLIENT_ID", ""),
			ClientSecret: getAppString("DINGTALK_CLIENT_SECRET", ""),
			AgentID:      getAppString("DINGTALK_AGENT_ID", ""),
			MiniAppID:    getAppString("DINGTALK_MINI_APP_ID", ""),
		},
	}
}

// getAppString 获取字符串类型配置
func getAppString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getAppInt 获取整数类型配置
func getAppInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getAppBool 获取布尔类型配置
func getAppBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
