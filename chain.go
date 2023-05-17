package llmchain

import "context"

type SreamCallBack func(res ChatResponse, done bool, err error)

// LLM common interface for lang model
type LLM interface {

	//Name return LLM Name
	Name() string

	//Free free model
	Free()

	Call(ctx context.Context, prompt string) (string, error)

	Chat(ctx context.Context, req *ChatRequest) (ChatResponse, error)

	Completion(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)

	//Embeddings
}

type Chain interface {

	//Name return LLM Name
	Name() string
	//WithLLM some langchain may using LLM for generate prompt in runtime。
	//这里要考虑历史对话的记录问题实现
	WithLLM(llm LLM)
	//Prompt parse input to prompt using chain prompt template. or maybe using api,db,cached data...
	//这里要考虑历史对话的记录问题实现
	Prompt(args H) (string, error)

	//InputPrompt using `input` as user input string
	InputPrompt(input string) (string, error)

	ChatPrompt(ctx context.Context, req *ChatRequest) (*ChatRequest, error)

	CompletionPrompt(ctx context.Context, req *CompletionRequest) (*CompletionRequest, error)
}

// H is a shortcut for map[string]string
type H = map[string]string

type LLMChain struct {
	llm   LLM
	chain Chain
}

var _ LLM = &LLMChain{}

// New return LLMChain
func New(llm LLM, chain Chain) *LLMChain {
	return &LLMChain{
		llm, chain,
	}
}

// Free implements LLM
func (*LLMChain) Free() {
	//do nothing
}

// Name implements LLM
func (l *LLMChain) Name() string {
	//
	return l.llm.Name()
}

// Run Simple warp for RunCompletion,
func (l *LLMChain) Call(ctx context.Context, input string) (string, error) {

	ret := ""

	prompt, err := l.chain.InputPrompt(input)

	if err != nil {
		//TODO
		return ret, err
	}

	return l.llm.Call(ctx, prompt)

}

func (l *LLMChain) Completion(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {

	r, err := l.chain.CompletionPrompt(ctx, req)

	if err != nil {
		//TODO
		return nil, err
	}

	return l.llm.Completion(ctx, r)

}

// func(res ChatCompletionResponse, done bool, err error)

func (l *LLMChain) Chat(ctx context.Context, req *ChatRequest) (ChatResponse, error) {

	//
	r, err := l.chain.ChatPrompt(ctx, req)

	if err != nil {
		//TODO
		return ChatResponse{}, err
	}

	return l.llm.Chat(ctx, r)
}
