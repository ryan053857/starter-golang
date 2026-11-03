package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type yamlFile struct {
	Env      string       `yaml:"env"`
	HTTP     yamlHTTP     `yaml:"http"`
	Log      yamlLog      `yaml:"log"`
	Database yamlDatabase `yaml:"database"`
}

type yamlDatabase struct {
	URL string `yaml:"url"`
}

type yamlHTTP struct {
	Addr string `yaml:"addr"`
}

type yamlLog struct {
	Level string `yaml:"level"`
}

// applyYAML 读取 path 的 YAML 并写入 cfg。
func applyYAML(cfg *Config, path string) error {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("read config %s: %w", path, err)
	}

	var file yamlFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return fmt.Errorf("parse config %s: %w", path, err)
	}

	setIf(&cfg.Env, file.Env)
	setIf(&cfg.Addr, file.HTTP.Addr)
	setIf(&cfg.LogLevel, file.Log.Level)
	setIf(&cfg.DatabaseURL, file.Database.URL)
	return nil
}
