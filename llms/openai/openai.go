package openai

import "github.com/exppii/llmchain/llms"

type OpenAI struct {
}

func (l *OpenAI) Options() llms.ModelOptions {
	panic(`TODO`)
}

// TODO
type ModelOption func(*ModelOptions)

type ModelOptions struct {
	Token string
	Model string
	// Seed       int
	// F16Memory  bool
	// MLock      bool
	// Embeddings bool
}

// WithToken sets the openai token.
func WithToken(c string) ModelOption {
	return func(p *ModelOptions) {
		p.Token = c
	}
}

func WithModel(c string) ModelOption {
	return func(p *ModelOptions) {
		p.Model = c
	}
}
