package chains

import (
	"context"

	"github.com/nexptr/llmchain/schema"
)

// H is a shortcut for map[string]string
type H = map[string]string

// Chain is the interface all chains must implement.
type Chain interface {
	GetName() string
	// Chat runs the logic of the chain and returns the output. This method should
	// not be called directly. Use rather the Chat function that handles the memory
	// of the chain.
	Chat(ctx context.Context, inputs map[string]any, options ...ChainCallOption) (map[string]any, error)
	// GetMemory gets the memory of the chain.
	GetMemory() schema.Memory
	// InputKeys returns the input keys the chain expects.
	GetInputKeys() []string
	// OutputKeys returns the output keys the chain expects.
	GetOutputKeys() []string
}

// ChainCallOption is a function that can be used to modify the behavior of the Call function.
type ChainCallOption func(*chainCallOptions)

type chainCallOptions struct {
	StopWords []string
}

// WithStopWords is a ChainCallOption that can be used to set the stop words of the chain.
func WithStopWords(stopWords []string) ChainCallOption {
	return func(options *chainCallOptions) {
		options.StopWords = stopWords
	}
}
