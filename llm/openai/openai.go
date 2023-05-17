package openai

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/exppii/llmchain"
	"github.com/exppii/llmchain/llm"
)

var _ llmchain.LLM = &OpenAI{}

type OpenAI struct {
	Model string `json:"model" yaml:"model"`
	// APIKey issued by OpenAI console.
	// See https://beta.openai.com/account/api-keys
	APIKey string `json:"api_key" yaml:"api_key"`

	// BaseURL of API including the version.
	// e.g., https://api.openai.com/v1
	APIHost string `json:"api_host" yaml:"api_host"`

	// Organization
	Organization string `json:"organization" yaml:"organization"`

	// HTTPClient (optional) to proxy HTTP request.
	// If nil, *http.DefaultClient will be used.
	HTTPClient *http.Client `json:"-" yaml:"-"`
}

func FromYaml(opt llm.ModelOptions) (*OpenAI, error) {

	client := defaultOpenAI()

	err := llm.UnmarshalPlugin(opt.Settings, client)

	if err != nil {
		return nil, err
	}

	client.Model = opt.Name

	return client, nil

}

// New return OpenAI compatiable client
func New(opts ...ModelOption) *OpenAI {

	client := defaultOpenAI()

	for _, fn := range opts {
		fn(client)
	}

	return client
}

// String dump openAI
func (l *OpenAI) Name() string {
	return l.Model
}

// Free implements llmchain.LLM
func (l *OpenAI) Free() {
	//do nothing
	// if l.HTTPClient != nil {
	// 	l.HTTPClient.CloseIdleConnections()
	// }
}

// String dump openAI
func (l *OpenAI) String() string {
	j, _ := json.Marshal(l)
	return string(j)
}

func (l *OpenAI) Call(ctx context.Context, prompt string) (string, error) {

	req := l.defaultChatRequest(prompt)

	resp, err := l.Chat(ctx, req)

	if err != nil {
		//TODO
		return "", err
	}
	ret := ``
	for _, v := range resp.Choices {
		if v.Message != nil {
			ret += v.Message.Content
		}

	}

	return ret, nil

}

// Chat implements llmchain.LLM
func (l *OpenAI) Chat(ctx context.Context, rawReq *llmchain.ChatRequest) (resp llmchain.ChatResponse, err error) {

	p := "/chat/completions"

	// if rawReq == nil {
	// 	rawReq = l.defaultChatRequest(prompt)
	// }

	if rawReq.StreamCallback != nil {
		rawReq.Stream = true // Nosy ;)
		return call(ctx, l, http.MethodPost, p, rawReq, resp, rawReq.StreamCallback)
	}
	return call(ctx, l, http.MethodPost, p, rawReq, resp, nil)

}

// Completion implements llmchain.LLM
func (l *OpenAI) Completion(ctx context.Context, rawReq *llmchain.CompletionRequest) (resp *llmchain.CompletionResponse, err error) {

	//todo create req for prompt
	p := "/completions"

	// if rawReq == nil {
	// 	rawReq = defaultCompletionRequest(prompt)
	// }

	return call(ctx, l, http.MethodPost, p, rawReq, resp, nil)
}

func (l *OpenAI) defaultChatRequest(prompt string) *llmchain.ChatRequest {

	msg := llmchain.Message{Role: `user`, Content: prompt}

	return &llmchain.ChatRequest{
		Model:       l.Model,
		Messages:    []llmchain.Message{msg},
		Temperature: 0.8,
		TopP:        1,
		N:           1,
		Stream:      false,
		// StreamCallback: ,
		Stop:             []string{},
		MaxTokens:        0,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
		LogitBias:        nil,
		User:             "",
	}
}
