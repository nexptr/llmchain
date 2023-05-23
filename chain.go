package llmchain

import "context"

// LLMChain 实现基本的LLMChain
type LLMChain struct {
	llm   LLM
	chain Chain
}

// 这里让 LLMChain也实现LLM接口。如此，也可以把LLMChain当成是一种特殊的LLM
var _ LLM = &LLMChain{}

// New return LLMChain
func New(llm LLM, chain Chain) *LLMChain {
	return &LLMChain{
		llm, chain,
	}
}

// Free implements LLM
func (l *LLMChain) Free() {
	//do nothing
	if l.llm != nil {
		l.Free()
	}
	//TODO maybe chain also need free()
}

// Name implements LLM
func (l *LLMChain) Name() string {
	//
	return l.llm.Name()
}

// Run Simple warp for RunCompletion,
func (l *LLMChain) Call(ctx context.Context, input string) (string, error) {

	//using chain to get prompted text
	prompt := l.chain.Prompt(input)

	//then use l.llm do real call
	return l.llm.Call(ctx, prompt)

}

func (l *LLMChain) Completion(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {

	// In ChatGPT Completion Request Prompt: The prompt(s) to generate completions for, encoded as a string, array of strings, array of tokens, or array of token arrays.
	// Note that <|endoftext|> is the document separator that the model sees during training, so if a prompt is not specified the model will generate as if from the beginning of a new document.
	// See https://beta.openai.com/docs/api-reference/completions/create#completions/create-prompt
	// So,Prompt may be string or []string
	switch p := req.Prompt.(type) {
	case string:
		req.Prompt = l.chain.Prompt(p)
		// req.PromptStrings = append(req.PromptStrings, l.chain.Prompt(p))
	case []string:
		prompts := []string{}
		for _, pp := range p {
			prompts = append(prompts, l.chain.Prompt(pp))
		}
		req.Prompt = prompts
	}

	return l.llm.Completion(ctx, req)

}

// func(res ChatCompletionResponse, done bool, err error)

func (l *LLMChain) Chat(ctx context.Context, req *ChatRequest) (ChatResponse, error) {

	//由于请求中间可能包括多轮对话。
	messages, err := l.chain.ChatPrompt(ctx, req.Messages)

	if err != nil {
		//TODO
		return ChatResponse{}, err
	}

	req.Messages = messages

	return l.llm.Chat(ctx, req)
}

// Embeddings implements LLM
func (l *LLMChain) Embeddings(ctx context.Context, req *EmbeddingsRequest) (*EmbeddingsResponse, error) {
	return l.llm.Embeddings(ctx, req)
}
