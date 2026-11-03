// Package config 加载进程级配置。
package config

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultConfigFile = "config.yaml"
	dotEnvFile        = ".env"
)

// Config 是进程级配置。
type Config struct {
	Addr     string
	Env      string
	LogLevel string
}

// Load 依次应用默认值、YAML、环境变量后返回配置。
// path 非空时用该路径；为空时用 CONFIG_FILE，再否则用 config.yaml。
// YAML 文件存在时读取。Env 为 dev 时从 .env 补充尚未设置的环境变量。
func Load(path string) (Config, error) {
	cfg := Config{
		Addr: "127.0.0.1:8080",
		Env:  "production",
	}

	if path == "" {
		path = os.Getenv("CONFIG_FILE")
		if path == "" {
			path = defaultConfigFile
		}
	}
	if err := applyYAML(&cfg, path); err != nil {
		return Config{}, err
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = cfg.Env
	}
	if env == "dev" {
		if err := godotenv.Load(dotEnvFile); err != nil && !os.IsNotExist(err) {
			return Config{}, err
		}
	}

	setIf(&cfg.Env, os.Getenv("APP_ENV"))
	setIf(&cfg.Addr, os.Getenv("ADDR"))
	setIf(&cfg.LogLevel, os.Getenv("LOG_LEVEL"))
	if cfg.LogLevel == "" {
		if cfg.Env == "dev" {
			cfg.LogLevel = "debug"
		} else {
			cfg.LogLevel = "info"
		}
	}
	return cfg, nil
}

// setIf 在 src 非空时写入 dst。
func setIf(dst *string, src string) {
	if src != "" {
		*dst = src
	}
}
