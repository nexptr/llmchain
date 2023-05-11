package llms

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/exppii/llmchain/prompts"
)

// CallOption is a function that configures a LLM.
type ModelOption func(*ModelOptions)

// ModelOptions is a set base of options for models
type ModelOptions struct {
	OpenAIRequest `yaml:"parameters"`
	Name          string `yaml:"name"`
	ModelPath     string `yaml:"model_path"`

	StopWords   []string `yaml:"stopwords"`
	Cutstrings  []string `yaml:"cutstrings"`
	TrimSpace   []string `yaml:"trimspace"`
	ContextSize int      `yaml:"context_size"`

	Parts int  `yaml:"parts"`
	MLock bool `yaml:"mlock"`

	Threads        int                    `yaml:"threads" json:"threads"`
	Debug          bool                   `yaml:"debug"`
	Roles          map[string]string      `yaml:"roles"`
	Embeddings     bool                   `yaml:"embeddings"`
	TemplateConfig prompts.TemplateConfig `yaml:"template"`

	PromptStrings, InputStrings []string

	TokenCallback func(string) bool `yaml:"-" json:"-"`
}

func (m *ModelOptions) Dump() string {
	j, _ := json.Marshal(m)
	return string(j)
}

func (m *ModelOptions) BuildOpts() []ModelOption {

	// Generate the prediction using the language model
	option := []ModelOption{}

	if m.Temperature != 0 {
		option = append(option, WithTemperature(m.Temperature))
	}

	if m.TopP != 0 {
		option = append(option, WithTopP(m.TopP))
	}

	if m.TopK != 0 {
		option = append(option, WithTopK(m.TopK))
	}

	if m.Maxtokens != 0 {
		option = append(option, WithMaxToken(m.Maxtokens))
	}

	if m.Threads != 0 {
		option = append(option, WithThreads(m.Threads))
	}

	if m.Mirostat != 0 {
		option = append(option, WithMirostat(m.Mirostat))
	}

	if m.MirostatETA != 0 {
		option = append(option, WithMirostatETA(m.MirostatETA))
	}

	if m.MirostatTAU != 0 {
		option = append(option, WithMirostatTAU(m.MirostatTAU))
	}

	if m.Debug {
		option = append(option, Debug)
	}

	option = append(option, WithStopWords(m.StopWords...))

	if m.RepeatPenalty != 0 {
		option = append(option, WithPenalty(m.RepeatPenalty))
	}

	if m.Keep != 0 {
		option = append(option, WithKeep(m.Keep))
	}

	if m.Batch != 0 {
		option = append(option, WithBatch(m.Batch))
	}

	if m.F16 {
		option = append(option, EnableF16KV)
	}

	if m.IgnoreEOS {
		option = append(option, IgnoreEOS)
	}

	if m.Seed != 0 {
		option = append(option, WithSeed(m.Seed))
	}

	if m.TemplateConfig.Completion != `` || m.TemplateConfig.Chat != `` || m.TemplateConfig.Edit != `` {
		option = append(option, WithTemplateConfig(m.TemplateConfig))
	}

	return option

}

// Create a new PredictOptions object with the given options.
func (m *ModelOptions) MergeModelOptions(opts ...ModelOption) *ModelOptions {

	for _, opt := range opts {
		opt(m)
	}
	return m
}

var EnableEmbeddings ModelOption = func(p *ModelOptions) { p.Embeddings = true }

var EnableF16Memory ModelOption = func(p *ModelOptions) { p.F16 = true }

// WithContext sets the context size.
func WithContext(c int) ModelOption {
	return func(p *ModelOptions) {
		p.ContextSize = c
	}
}

// WithContext sets the context size.
func WithTemplateConfig(c prompts.TemplateConfig) ModelOption {
	return func(p *ModelOptions) {
		p.TemplateConfig = c
	}
}

func WithModelSeed(c int) ModelOption {
	return func(p *ModelOptions) {
		p.Seed = c
	}
}

func WithParts(c int) ModelOption {
	return func(p *ModelOptions) {
		p.Parts = c
	}
}

var Debug ModelOption = func(p *ModelOptions) {
	p.Debug = true
}

var IgnoreEOS ModelOption = func(p *ModelOptions) {
	p.IgnoreEOS = true
}

var EnableF16KV ModelOption = func(p *ModelOptions) { p.F16 = true }

// WithTokenCallback sets the prompts that will stop predictions.
func WithTokenCallback(fn func(string) bool) ModelOption {
	return func(p *ModelOptions) {
		p.TokenCallback = fn
	}
}

// WithStopWords sets the prompts that will stop predictions.
func WithStopWords(stop ...string) ModelOption {
	return func(p *ModelOptions) {
		p.StopWords = stop
	}
}

// WithSeed sets the random seed for sampling text generation.
func WithSeed(seed int) ModelOption {
	return func(p *ModelOptions) {
		p.Seed = seed
	}
}

// WithThreads sets the number of threads to use for text generation.
func WithThreads(threads int) ModelOption {
	return func(p *ModelOptions) {
		p.Threads = threads
	}
}

// Maxtokens sets the number of tokens to generate.
func WithMaxToken(max int) ModelOption {
	return func(p *ModelOptions) {

		p.Maxtokens = max
	}
}

// WithTopK sets the value for top-K sampling.
func WithTopK(topk int) ModelOption {
	return func(p *ModelOptions) {
		p.TopK = topk
	}
}

// WithTopP sets the value for nucleus sampling.
func WithTopP(topp float64) ModelOption {
	return func(p *ModelOptions) {
		p.TopP = topp
	}
}

// WithTemperature sets the temperature value for text generation.
func WithTemperature(temp float64) ModelOption {
	return func(p *ModelOptions) {
		p.Temperature = temp
	}
}

// WithPenalty sets the repetition penalty for text generation.
func WithPenalty(penalty float64) ModelOption {
	return func(p *ModelOptions) {
		p.RepeatPenalty = penalty
	}
}

// WithRepeat sets the number of times to repeat text generation.
func WithRepeat(repeat int) ModelOption {
	return func(p *ModelOptions) {
		p.Repeat = repeat
	}
}

// WithBatch sets the batch size.
func WithBatch(size int) ModelOption {
	return func(p *ModelOptions) {
		p.Batch = size
	}
}

// WithKeep sets the number of tokens from initial prompt to keep.
func WithKeep(n int) ModelOption {
	return func(p *ModelOptions) {
		p.Keep = n
	}
}

// WithMirostat sets the mirostat parameter.
func WithMirostat(m int) ModelOption {
	return func(p *ModelOptions) {
		p.Mirostat = m
	}
}

// WithMirostatETA sets the mirostat ETA parameter.
func WithMirostatETA(me float64) ModelOption {
	return func(p *ModelOptions) {
		p.MirostatETA = me
	}
}

// WithMirostatTAU sets the mirostat TAU parameter.
func WithMirostatTAU(mt float64) ModelOption {
	return func(p *ModelOptions) {
		p.MirostatTAU = mt
	}
}

// var EnableMLock ModelOption = func(p *ModelOptions) { p.MLock = true }

// WithMaxTokens is an option for LLM.Call.
// func WithMaxTokens(maxTokens int) ModelOption {
// 	return func(o *ModelOptions) {
// 		o.MaxTokens = maxTokens
// 	}
// }

// // WithTemperature is an option for LLM.Call.
// func WithTemperature(temperature float64) ModelOption {
// 	return func(o *ModelOptions) {
// 		o.Temperature = temperature
// 	}
// }

// // WithStopWords is an option for LLM.Call.
// func WithStopWords(stopWords []string) ModelOption {
// 	return func(o *ModelOptions) {
// 		o.StopWords = stopWords
// 	}
// }

// // WithOptions is an option for LLM.Call.
// func WithOptions(options ModelOptions) ModelOption {
// 	return func(o *ModelOptions) {
// 		(*o) = options
// 	}
// }

// Override override input req
func (m *ModelOptions) Override(input *OpenAIRequest) *ModelOptions {

	m.N = input.N
	m.Messages = input.Messages

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

	if input.Maxtokens != 0 && input.Maxtokens < m.Maxtokens {
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

// TemplatePromptStrings TODO replace a function name
func (m *ModelOptions) TemplatePromptStrings(t *prompts.Template) ([]string, error) {

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

func (m *ModelOptions) TemplateMessage(t *prompts.Template) (string, error) {

	mess := []string{}
	for _, i := range m.Messages {
		r := m.Roles[i.Role]
		if r == "" {
			r = i.Role
		}

		content := fmt.Sprint(r, " ", i.Content)
		mess = append(mess, content)
	}

	predInput := strings.Join(mess, "\n")

	templateFile := m.Model

	if m.TemplateConfig.Completion != "" {
		templateFile = m.TemplateConfig.Completion
	}

	return t.Render(templateFile, struct {
		Input string
	}{Input: predInput})

}
