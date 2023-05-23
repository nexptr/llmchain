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
func (*BaseChat) Prompt(input string) string {
	//do nothing
	return input
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

// ChatPrompt implements llmchain.Chain
// func (*BaseChat) ChatPrompt(ctx context.Context, req *llmchain.ChatRequest) (*llmchain.ChatRequest, error) {
// 	return req, nil
// }

// ChatPrompt implements llmchain.Chain
func (*BaseChat) ChatPrompt(ctx context.Context, messages []llmchain.Message) ([]llmchain.Message, error) {
	//do nothing
	return messages, nil
}

// PromptArgs implements llmchain.Chain
func (*BaseChat) PromptArgs(args map[string]string) (string, error) {

	input, ok := args[`input`]

	if !ok {
		//TODO
		return "hello", nil
	}

	return input, nil

}
