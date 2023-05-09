package llms

import (
	"github.com/exppii/llmchain/prompts"
)

type Payload struct {
	OpenAIRequest `yaml:"parameters"`
	Name          string   `yaml:"name"`
	StopWords     []string `yaml:"stopwords"`
	Cutstrings    []string `yaml:"cutstrings"`
	TrimSpace     []string `yaml:"trimspace"`
	ContextSize   int      `yaml:"context_size"`

	Threads        int                    `yaml:"threads"`
	Debug          bool                   `yaml:"debug"`
	Roles          map[string]string      `yaml:"roles"`
	Embeddings     bool                   `yaml:"embeddings"`
	TemplateConfig prompts.TemplateConfig `yaml:"template"`

	PromptStrings, InputStrings []string
}

type TemplateConfig struct {
	Completion string `yaml:"completion"`
	Chat       string `yaml:"chat"`
	Edit       string `yaml:"edit"`
}

func NewPayload(*ModelOptions) *Payload {

	return &Payload{
		// F16:         true,
		Debug:       false,
		Threads:     4,
		ContextSize: 512,
	}

}

func (m *Payload) TemplatePromptStrings(t *prompts.Template) ([]string, error) {

	templateFile := m.Model

	if m.TemplateConfig.Completion != "" {
		templateFile = m.TemplateConfig.Completion
	}

	templatedInputs := []string{}

	for _, i := range m.PromptStrings {
		templatedInput, err := t.Render(templateFile, struct {
			Input string
		}{Input: i})

		if err != nil {
			return templatedInputs, err
		}

		templatedInputs = append(templatedInputs, templatedInput)
	}

	return templatedInputs, nil

}

func (m *Payload) Override(input *OpenAIRequest) *Payload {
	if input.Echo {
		m.Echo = input.Echo
	}
	if input.TopK != 0 {
		m.TopK = input.TopK
	}
	if input.TopP != 0 {
		m.TopP = input.TopP
	}

	if input.Temperature != 0 {
		m.Temperature = input.Temperature
	}

	if input.Maxtokens != 0 {
		m.Maxtokens = input.Maxtokens
	}

	switch stop := input.Stop.(type) {
	case string:
		if stop != "" {
			m.StopWords = append(m.StopWords, stop)
		}
	case []interface{}:
		for _, pp := range stop {
			if s, ok := pp.(string); ok {
				m.StopWords = append(m.StopWords, s)
			}
		}
	}

	if input.RepeatPenalty != 0 {
		m.RepeatPenalty = input.RepeatPenalty
	}

	if input.Keep != 0 {
		m.Keep = input.Keep
	}

	if input.Batch != 0 {
		m.Batch = input.Batch
	}

	if input.F16 {
		m.F16 = input.F16
	}

	if input.IgnoreEOS {
		m.IgnoreEOS = input.IgnoreEOS
	}

	if input.Seed != 0 {
		m.Seed = input.Seed
	}

	if input.Mirostat != 0 {
		m.Mirostat = input.Mirostat
	}

	if input.MirostatETA != 0 {
		m.MirostatETA = input.MirostatETA
	}

	if input.MirostatTAU != 0 {
		m.MirostatTAU = input.MirostatTAU
	}

	switch inputs := input.Input.(type) {
	case string:
		if inputs != "" {
			m.InputStrings = append(m.InputStrings, inputs)
		}
	case []interface{}:
		for _, pp := range inputs {
			if s, ok := pp.(string); ok {
				m.InputStrings = append(m.InputStrings, s)
			}
		}
	}

	switch p := input.Prompt.(type) {
	case string:
		m.PromptStrings = append(m.PromptStrings, p)
	case []interface{}:
		for _, pp := range p {
			if s, ok := pp.(string); ok {
				m.PromptStrings = append(m.PromptStrings, s)
			}
		}
	}
	return m
}
