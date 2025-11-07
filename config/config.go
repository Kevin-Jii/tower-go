package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	DingTalk DingTalkConfig `yaml:"dingtalk"`
}

type AppConfig struct {
	Name string `yaml:"name"`
	Port int    `yaml:"port"`
}

type DatabaseConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	Enabled  bool   `yaml:"enabled"`
}

type DingTalkConfig struct {
	Stream DingTalkStreamConfig `yaml:"stream"`
}

type DingTalkStreamConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	AgentID      string `yaml:"agent_id"`
	MiniAppID    string `yaml:"mini_app_id"`
}

var cfg Config

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	return nil
}

func GetConfig() Config {
	return cfg
}

func GetDatabaseConfig() DatabaseConfig {
	return cfg.Database
}

func GetRedisConfig() RedisConfig {
	return cfg.Redis
}

func GetDingTalkConfig() DingTalkConfig {
	return cfg.DingTalk
}

func GetDingTalkStreamConfig() DingTalkStreamConfig {
	return cfg.DingTalk.Stream
}
