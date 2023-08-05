package local

import (
	context "context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/nexptr/llmchain/llms"
	"github.com/nexptr/llmchain/schema"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (l *LLaMA) defaultChatRequest(prompt string, options ...llms.CallOption) *schema.ChatRequest {

	opts := llms.InitCallOptions(options...)

	msg := schema.Message{Role: `user`, Content: prompt}

	req := &schema.ChatRequest{
		Model:       l.Model,
		Messages:    []schema.Message{msg},
		Temperature: float32(opts.Temperature),
		// TopP:        1,
		// N:           1,
		// Stream: false,
		// StreamCallback: ,
		Stop:      opts.StopWords,
		MaxTokens: opts.MaxTokens,
		// PresencePenalty:  0,
		// FrequencyPenalty: 0,
		// LogitBias:        nil,
		// User:             "",
	}

	if opts.CallBackFn != nil {
		req.Stream = true
		req.StreamCallback = opts.CallBackFn
	}

	return req
}

func (l *LLaMA) build(ctx context.Context) error {
	if l.client == nil {
		endpoint, err := l.endpoint()
		if err != nil {
			return err
		}
		client_conn, err := grpc.DialContext(ctx, endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		l.client = NewChatServiceClient(client_conn)
	}
	return nil
}

// endpoint 默认轮训后段模型
func (l *LLaMA) endpoint() (string, error) {
	if len(l.Hosts) == 1 {
		return l.Hosts[0], nil
	}
	// 生成随机索引
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Intn(len(l.Hosts))

	return l.Hosts[index], nil
}

func NewGenerationRequestByChatRequest(model_name string, req *schema.ChatRequest) *GenerationRequest {
	request, ok := defaultChatRequest[model_name]
	if !ok {
		request = defaultChatRequest["default"]
	}

	if req.Temperature > 0 {
		request.Temperature = req.Temperature
	}

	if req.TopP > 0 {
		request.TopP = req.TopP
	}

	if req.N > 0 {
		request.N = int32(req.N)
	}

	if req.MaxTokens > 0 {
		request.MaxNewTokens = int32(req.MaxTokens)
	}

	if req.PresencePenalty > 0 {
		request.RepetitionPenalty = req.PresencePenalty
	}

	request.Prompt = PromptMessage(model_name, req.Messages)

	return request
}

func NewGenerationRequestByCompletionRequest(req *schema.CompletionRequest) *GenerationRequest {
	request := &GenerationRequest{
		Temperature:       req.Temperature,
		TopP:              req.TopP,
		N:                 int32(req.N),
		MaxNewTokens:      int32(req.MaxTokens),
		Prompt:            req.Prompt,
		RepetitionPenalty: req.PresencePenalty,
		Echo:              req.Echo,

		// TopK:              0,
		// Stop:              req.Stop,
		// StopTokenIds:      []string{},
	}

	return request
}

func NewEmbeddingsRequest(req *schema.EmbeddingsRequest) *EmbedingsMessage {
	request := &EmbedingsMessage{}

	if value, ok := req.Input.(string); ok {
		request.Prompt = []string{value}
	} else if value, ok := req.Input.([]string); ok {
		request.Prompt = value
	}

	return request
}

func call[T any](ctx context.Context, l *LLaMA, req *schema.ChatRequest, resp T, cb schema.SreamCallBack) (T, error) {
	if err := l.build(ctx); err != nil {
		return resp, err
	}

	in := NewGenerationRequestByChatRequest(l.Model, req)
	if cb != nil {
		c, err := l.client.Chat(ctx, in)
		if err != nil {
			return resp, err
		}
		go listen(c, cb)
		return resp, nil
	}

	reply, err := l.client.Completion(ctx, in)
	if err != nil {
		return resp, err
	}
	if reply.Usage == nil {
		reply.Usage = &TokenUsage{}
	}

	if reply.ErrorCode != ErrorCode_Zero {
		return resp, fmt.Errorf(reply.Text)
	} else {
		by, _ := json.Marshal(reply)
		if err := json.Unmarshal(by, resp); err != nil {
			return resp, err
		}
	}

	return resp, nil
}

func listen(client ChatService_ChatClient, cb schema.SreamCallBack) {

	var rErr error
	defer func() {
		cb(nil, true, rErr)
	}()

	for {
		reply, err := client.Recv()
		if err != nil {
			if err != io.EOF {
				rErr = err
			}
			return
		}
		if reply.ErrorCode != ErrorCode_Zero {
			rErr = fmt.Errorf(reply.Text)
			return
		}
		if reply.Usage == nil {
			reply.Usage = &TokenUsage{}
		}

		msg := schema.BuildAIMessage(reply.Text)
		retD := &schema.ChatResponse{
			ID: "todo",
			Choices: []schema.Choice{{
				FinishReason: reply.FinishReason,
			}},
			Usage: schema.Usage{
				PromptTokens:     int(reply.Usage.PromptTokens),
				CompletionTokens: int(reply.Usage.CompletionTokens),
				TotalTokens:      int(reply.Usage.TotalTokens),
			},
		}

		if reply.FinishReason != `` {
			retD.Choices[0].Message = &msg
			cb(retD, false, nil)
		} else {
			retD.Choices[0].Delta = &msg
			cb(retD, false, nil)
		}
	}
}
