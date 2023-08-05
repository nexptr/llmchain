package chains

import (
	"context"

	"github.com/nexptr/llmchain/llms"
	"github.com/nexptr/llmchain/schema"
)

const (
	BaseChatChain = `base_chat_chain`
)

// BaseChat base chat Lang chain,this chain just do nothing
type BaseChat struct {
	name string
}

// Chat implements Chain.
func (*BaseChat) Chat(ctx context.Context, inputs map[string]any, options ...ChainCallOption) (map[string]any, error) {
	panic("unimplemented")
}

// GetInputKeys implements Chain.
func (*BaseChat) GetInputKeys() []string {
	panic("unimplemented")
}

// GetMemory implements Chain.
func (*BaseChat) GetMemory() schema.Memory {
	panic("unimplemented")
}

// GetName implements Chain.
func (*BaseChat) GetName() string {
	panic("unimplemented")
}

// GetOutputKeys implements Chain.
func (*BaseChat) GetOutputKeys() []string {
	panic("unimplemented")
}

var _ Chain = &BaseChat{}

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
func (*BaseChat) WithLLM(llm llms.LLM) {
	//do nothing
}

// ChatPrompt implements llmchain.Chain
// func (*BaseChat) ChatPrompt(ctx context.Context, req *llmchain.ChatRequest) (*llmchain.ChatRequest, error) {
// 	return req, nil
// }

// ChatPrompt implements llmchain.Chain
func (*BaseChat) ChatPrompt(ctx context.Context, messages []schema.Message) ([]schema.Message, error) {
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
