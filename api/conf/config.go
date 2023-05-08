package conf

import (
	"fmt"
	"os"

	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
)

// Config 定义 配置结构图
type Config struct {
	APIAddr string `yaml:"api_addr"`

	ModelPath string `yaml:"model_path"`

	LogLevel zapcore.Level `yaml:"verbose"`

	LogDir string `yaml:"log_dir"`
}

func defaultConfig() *Config {
	cf := &Config{
		APIAddr:   "0.0.0.0:8080",
		ModelPath: "./models",
		LogLevel:  zapcore.InfoLevel,
		LogDir:    "./logs",
	}
	return cf
}

func InitConfig(path string) (*Config, error) {

	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	cf := defaultConfig()

	if err := yaml.Unmarshal(f, cf); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config file: %w", err)
	}
	return cf, nil
}
