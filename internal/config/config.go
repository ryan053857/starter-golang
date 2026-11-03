// Package config 加载进程级配置。
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultConfigFile = "config.yaml"
	dotEnvFile        = ".env"
)

// Config 是进程级配置。
type Config struct {
	Addr        string
	Env         string
	LogLevel    string
	DatabaseURL string
}

// Load 装载配置：默认值、YAML、环境变量；path 为空时使用 config.yaml。
func Load(path string) (Config, error) {
	cfg := Config{
		Addr:     "127.0.0.1:8080",
		Env:      "production",
		LogLevel: "info",
	}

	if path == "" {
		path = defaultConfigFile
	}
	if err := applyYAML(&cfg, path); err != nil {
		return Config{}, err
	}

	if os.Getenv("APP_ENV") == "dev" {
		if err := godotenv.Load(dotEnvFile); err != nil {
			return Config{}, fmt.Errorf("load %s: %w", dotEnvFile, err)
		}
	}

	setIf(&cfg.Env, os.Getenv("APP_ENV"))
	setIf(&cfg.Addr, os.Getenv("ADDR"))
	setIf(&cfg.LogLevel, os.Getenv("LOG_LEVEL"))
	setIf(&cfg.DatabaseURL, os.Getenv("DATABASE_URL"))
	return cfg, nil
}

func setIf(dst *string, src string) {
	if src != "" {
		*dst = src
	}
}
