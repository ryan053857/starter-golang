// Package config 加载进程级配置。
package config

import (
	"os"
)

const defaultConfigFile = "config.yaml"

// Config 是进程级配置。
type Config struct {
	Addr     string
	Env      string
	LogLevel string
}

// Load 依次应用默认值、环境变量后返回配置。
// path 非空时用该路径；为空时用 CONFIG_FILE，再否则用 config.yaml。
// 当前脚手架仅使用环境变量；path 预留给后续 YAML 加载。
func Load(path string) (Config, error) {
	if path == "" {
		path = os.Getenv("CONFIG_FILE")
		if path == "" {
			path = defaultConfigFile
		}
	}
	_ = path // 预留给后续 YAML 加载

	cfg := Config{
		Addr: "127.0.0.1:8080",
		Env:  "production",
	}
	setIf(&cfg.Addr, os.Getenv("ADDR"))
	setIf(&cfg.Env, os.Getenv("APP_ENV"))
	setIf(&cfg.LogLevel, os.Getenv("LOG_LEVEL"))
	return cfg, nil
}

func setIf(dst *string, v string) {
	if v != "" {
		*dst = v
	}
}
