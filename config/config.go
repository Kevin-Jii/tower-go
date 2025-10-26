package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
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
