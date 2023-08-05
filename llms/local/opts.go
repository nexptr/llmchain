package local

import (
	"github.com/nexptr/llmchain/llms"
)

const (
	defaultLLaMAAddr  = `127.0.0.1:50051`
	defaultLLaMAModel = `langchat-16k-v1`
)

// CallOption is a function that configures a LLM.
type ModelOption func(*LLaMA)

func FromYaml(opt llms.ModelOptions) (*LLaMA, error) {

	client := defaultLLaMA()

	err := llms.UnmarshalPlugin(opt.Settings, client)

	if err != nil {
		return nil, err
	}

	client.Model = opt.Name

	return client, nil

}

// WithHosts sets the Host that for  gRPC Server default is 127.0.0.1:50051
func WithHosts(o []string) ModelOption {

	//todo verify o
	return func(p *LLaMA) {
		p.Hosts = o
	}
}

// WithModel sets the model for us.
func WithModel(o string) ModelOption {
	return func(p *LLaMA) {
		p.Model = o
	}
}

func defaultLLaMA() *LLaMA {
	return &LLaMA{
		Hosts: []string{defaultLLaMAAddr},
		Model: "defaultLLaMAModel",
	}
}

var defaultChatRequest = map[string]*GenerationRequest{
	"default": {
		Temperature:       0.8,
		TopP:              1,
		N:                 1,
		MaxNewTokens:      7000,
		TopK:              0,
		RepetitionPenalty: 0,
	},
	"chatglm2-6b": {
		Temperature:       0.7,
		TopP:              2,
		N:                 1,
		MaxNewTokens:      7000,
		TopK:              0,
		RepetitionPenalty: 1,
	},
}
