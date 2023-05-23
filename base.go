package llmchain

import "context"

// H is a shortcut for map[string]string
type H = map[string]string

// LLM common interface for lang model
type LLM interface {

	//Name return LLM Name
	Name() string

	//Free free model
	Free()

	//Call 实现最基本的输入输出。
	Call(ctx context.Context, prompt string) (string, error)

	//Chat chatGPT compatible chat/completions input/output
	Chat(ctx context.Context, req *ChatRequest) (ChatResponse, error)

	//Chat chatGPT compatible completions input/output
	Completion(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)

	//Chat chatGPT compatible embeddings input/output
	Embeddings(ctx context.Context, req *EmbeddingsRequest) (*EmbeddingsResponse, error)
}

// Chain 定义基本的LangChain
type Chain interface {

	//Name return LLM Name
	Name() string

	//WithLLM some langchain may using LLM for generate prompt in runtime。
	//这里要考虑历史对话的记录问题实现
	WithLLM(llm LLM)

	//Prompt parse input to prompt using chain prompt template. or maybe using api,db,cached data...
	//这里要考虑历史对话的记录问题实现
	PromptArgs(args H) (string, error)

	//InputPrompt using `input` as user input string
	Prompt(input string) string

	//ChatPrompt 考虑到多轮对话中需要保留现有的角色对话关系，需要经过模版转化的后的输入依然保留角色信息
	ChatPrompt(ctx context.Context, messages []Message) ([]Message, error)
}
