package conf

import (
	"fmt"
	"os"

	"github.com/exppii/llmchain/llms"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
)

// Config 定义 配置结构图
type Config struct {

	//APIAddr , default: 0.0.0.0:8080
	APIAddr string `yaml:"api_addr"`

	//default is `./models`
	ModelPath string `yaml:"model_path"`

	//default is `./prompts/*.tmpl`
	PromptPath string `yaml:"prompt_path"`

	LogLevel zapcore.Level `yaml:"log_level"`

	LogDir string `yaml:"log_dir"`

	ModelOptions map[string]llms.ModelOption `yaml:"model_options"`
}

func defaultConfig() *Config {
	cf := &Config{

		APIAddr:    "0.0.0.0:8080",
		ModelPath:  "./models",
		PromptPath: "./prompts/*.tmpl",
		LogLevel:   zapcore.InfoLevel,
		LogDir:     "./logs",
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
