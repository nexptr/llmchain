// Package llms defines the types for LLMs.
package llms

type Callback func(token string) bool

// LLM is a langchaingo Large Language Model.
type LLM interface {

	//ModelOptions get current LLM call options
	MergePayload(req *OpenAIRequest) *Payload

	//SupportStream if this LLM support stream
	SupportStream() bool

	InferenceFn(input string, payload *Payload, tokenCallback func(string) bool) func() (string, error)

	// Call(ctx context.Context, prompt string, options ...PredictOption) (string, error)
	// Generate(ctx context.Context, prompts []string, options ...CallOption) ([]*Generation, error)
}

type OpenAIResponse struct {
	Created int      `json:"created,omitempty"`
	Object  string   `json:"object,omitempty"`
	ID      string   `json:"id,omitempty"`
	Model   string   `json:"model,omitempty"`
	Choices []Choice `json:"choices,omitempty"`
	Data    []Item   `json:"data,omitempty"`

	Usage OpenAIUsage `json:"usage"`
}

type OpenAIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Item struct {
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
	Object    string    `json:"object,omitempty"`
}

type Choice struct {
	Index        int      `json:"index,omitempty"`
	FinishReason string   `json:"finish_reason,omitempty"`
	Message      *Message `json:"message,omitempty"`
	Delta        *Message `json:"delta,omitempty"`
	Text         string   `json:"text,omitempty"`
}

type Message struct {
	Role    string `json:"role,omitempty" yaml:"role"`
	Content string `json:"content,omitempty" yaml:"content"`
}

type OpenAIRequest struct {
	Model string `json:"model" yaml:"model"`

	// Prompt is read only by completion API calls
	Prompt interface{} `json:"prompt" yaml:"prompt"`

	// Edit endpoint
	Instruction string      `json:"instruction" yaml:"instruction"`
	Input       interface{} `json:"input" yaml:"input"`

	Stop interface{} `json:"stop" yaml:"stop"`

	// Messages is read only by chat/completion API calls
	Messages []Message `json:"messages" yaml:"messages"`

	Stream bool `json:"stream"`
	Echo   bool `json:"echo"`
	// Common options between all the API calls
	TopP        float64 `json:"top_p" yaml:"top_p"`
	TopK        int     `json:"top_k" yaml:"top_k"`
	Temperature float64 `json:"temperature" yaml:"temperature"`
	Maxtokens   int     `json:"max_tokens" yaml:"max_tokens"`

	N int `json:"n"`

	// Custom parameters - not present in the OpenAI API
	ChainType string `json:"chain_type" yaml:"chain_type"`

	Batch         int     `json:"batch" yaml:"batch"`
	F16           bool    `json:"f16" yaml:"f16"`
	IgnoreEOS     bool    `json:"ignore_eos" yaml:"ignore_eos"`
	RepeatPenalty float64 `json:"repeat_penalty" yaml:"repeat_penalty"`
	Keep          int     `json:"n_keep" yaml:"n_keep"`

	MirostatETA float64 `json:"mirostat_eta" yaml:"mirostat_eta"`
	MirostatTAU float64 `json:"mirostat_tau" yaml:"mirostat_tau"`
	Mirostat    int     `json:"mirostat" yaml:"mirostat"`

	Seed int `json:"seed" yaml:"seed"`
}

func defaultRequest(modelFile string) OpenAIRequest {
	return OpenAIRequest{
		TopP:        0.7,
		TopK:        80,
		Maxtokens:   512,
		Temperature: 0.9,
		Model:       modelFile,
	}
}
