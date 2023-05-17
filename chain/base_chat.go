package chain

import (
	"context"

	"github.com/exppii/llmchain"
)

const (
	BaseChatChain = `base_chat_chain`
)

// BaseChat base chat Lang chain,this chain just do nothing
type BaseChat struct {
	name string
}

var _ llmchain.Chain = &BaseChat{}

func NewBaseChatChain() *BaseChat {

	return &BaseChat{name: `base_chat_chain`}
}

// Prompt implements llmchain.Chain args key:input
func (*BaseChat) Prompt(args map[string]string) (string, error) {

	input, ok := args[`input`]

	if !ok {
		//TODO
		return "hello", nil
	}

	return input, nil
}

// WithLLM implements llmchain.Chain
func (c *BaseChat) Name() string {
	//do nothing
	return c.name
}

// WithLLM implements llmchain.Chain
func (*BaseChat) WithLLM(llm llmchain.LLM) {
	//do nothing
}

// InputPrompt implements llmchain.Chain
func (*BaseChat) InputPrompt(input string) (string, error) {
	return input, nil
}

// ChatPrompt implements llmchain.Chain
func (*BaseChat) ChatPrompt(ctx context.Context, req *llmchain.ChatRequest) (*llmchain.ChatRequest, error) {
	return req, nil
}

// CompletionPrompt implements llmchain.Chain
func (*BaseChat) CompletionPrompt(ctx context.Context, req *llmchain.CompletionRequest) (*llmchain.CompletionRequest, error) {
	return req, nil
}
