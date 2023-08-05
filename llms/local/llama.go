package local

import (
	"context"
	"fmt"
	"time"

	"github.com/nexptr/llmchain/llms"
	"github.com/nexptr/llmchain/schema"
)

var _ llms.LLM = &LLaMA{}

type LLaMA struct {

	//Model the model name
	Model string `json:"model" yaml:"model"`

	// Hosts of server including the version.
	Hosts []string `json:"hosts" yaml:"hosts"`

	// ChatServiceClient send requsts to gRPC server.
	client ChatServiceClient `json:"-" yaml:"-"`
}

// New return OpenAI compatiable client
func New(opts ...ModelOption) *LLaMA {

	client := defaultLLaMA()

	for _, fn := range opts {
		fn(client)
	}

	return client
}

// Call implements llms.LLM.
func (l *LLaMA) Call(ctx context.Context, prompt string) (string, error) {

	req := l.defaultChatRequest(prompt)

	data, err := l.Chat(ctx, req)

	ret := ``
	if err != nil {
		//TODO
		return ret, err
	}

	for _, v := range data.Choices {
		if v.Message != nil {
			ret += v.Message.Content
		}

	}

	return ret, nil
}

// Chat implements llms.LLM.
func (l *LLaMA) Chat(ctx context.Context, req *schema.ChatRequest) (resp *schema.ChatResponse, err error) {

	if req.StreamCallback != nil {
		req.Stream = true // Nosy ;)
		return call(ctx, l, req, resp, req.StreamCallback)
	}

	vResp := &GenerationReply{}
	vResp, err = call(ctx, l, req, vResp, nil)
	if err != nil {
		return
	}
	if vResp.ErrorCode != ErrorCode_Zero {
		return nil, fmt.Errorf(vResp.Text)
	}

	msg := schema.BuildAIMessage(vResp.Text)

	resp = &schema.ChatResponse{
		ID:      "TODO",
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Choices: []schema.Choice{
			{
				Index:        0,
				FinishReason: vResp.FinishReason,
				Message:      &msg,
			},
		},
		Usage: schema.Usage{
			PromptTokens:     int(vResp.Usage.PromptTokens),
			CompletionTokens: int(vResp.Usage.CompletionTokens),
			TotalTokens:      int(vResp.Usage.TotalTokens),
		},
	}

	return
}

// Completion implements llms.LLM.
func (l *LLaMA) Completion(ctx context.Context, req *schema.CompletionRequest) (resp *schema.CompletionResponse, err error) {
	if err := l.build(ctx); err != nil {
		return nil, err
	}

	in := NewGenerationRequestByCompletionRequest(req)

	reply, err := l.client.Completion(ctx, in)
	if err != nil {
		return nil, err
	}

	if reply.ErrorCode != ErrorCode_Zero {
		return nil, fmt.Errorf(reply.Text)
	} else {
		msg := schema.BuildAIMessage(reply.Text)
		resp = &schema.CompletionResponse{
			Choices: []schema.Choice{
				{
					Message:      &msg,
					FinishReason: reply.FinishReason,
				},
			},
		}
	}
	return
}

// Embeddings implements llms.LLM.
func (l *LLaMA) Embeddings(ctx context.Context, req *schema.EmbeddingsRequest) (resp *schema.EmbeddingsResponse, err error) {
	if err := l.build(ctx); err != nil {
		return nil, err
	}

	in := NewEmbeddingsRequest(req)

	reply, err := l.client.Embedings(ctx, in)
	if err != nil {
		return nil, err
	}

	if reply.ErrorCode != ErrorCode_Zero {
		return nil, fmt.Errorf(reply.ErrorCode.String())
	}

	resp = &schema.EmbeddingsResponse{
		Object: "embedding",
		Data:   []schema.EmbeddingData{},
		Usage:  schema.Usage{},
	}

	for i, v := range reply.Embeddings {
		resp.Data = append(resp.Data, schema.EmbeddingData{
			Object:    "embedding",
			Embedding: v.Embedding,
			Index:     i,
		})
	}
	return
}

// Free implements llms.LLM.
func (l *LLaMA) Free() {
	// TODO
	panic("unimplemented")
}

// Name implements llms.LLM.
func (l *LLaMA) Name() string {
	return l.Model
}
