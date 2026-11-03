package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// yamlFile 是 YAML 配置文件结构。
type yamlFile struct {
	Env  string   `yaml:"env"`
	HTTP yamlHTTP `yaml:"http"`
	Log  yamlLog  `yaml:"log"`
}

// yamlHTTP 是 HTTP 段。
type yamlHTTP struct {
	Addr string `yaml:"addr"`
}

// yamlLog 是日志段。
type yamlLog struct {
	Level string `yaml:"level"`
}

// applyYAML 将 path 中非空字段写入 cfg；文件不存在时返回 nil。
func applyYAML(cfg *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read config %s: %w", path, err)
	}

	var file yamlFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return fmt.Errorf("parse config %s: %w", path, err)
	}

	setIf(&cfg.Env, file.Env)
	setIf(&cfg.Addr, file.HTTP.Addr)
	setIf(&cfg.LogLevel, file.Log.Level)
	return nil
}
