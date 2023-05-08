package prompts

type TemplateConfig struct {
	Completion string `yaml:"completion"`
	Chat       string `yaml:"chat"`
	Edit       string `yaml:"edit"`
}
