package openai_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/exppii/llmchain"
	"github.com/exppii/llmchain/llm/openai"
)

func TestOpenAI_Chat(t *testing.T) {

	ai := openai.New(openai.WithAPIHost(`http://192.168.12.213:8000/v1`), openai.WithModel(`vicuna-13b-v1.1`))

	resp, err := ai.Chat(context.TODO(), nil)

	if err != nil {
		t.Errorf(err.Error())
	}

	d, _ := json.Marshal(resp)
	println(string(d))
}

func TestOpenAI_Stream(t *testing.T) {

	ai := openai.New(openai.WithAPIHost(`http://192.168.12.213:8000/v1`), openai.WithModel(`vicuna-13b-v1.1`))

	msg := llmchain.Message{Role: `user`, Content: `怎样计算圆形面积`}
	req := &llmchain.ChatRequest{
		Model:       `vicuna-13b-v1.1`,
		Messages:    []llmchain.Message{msg},
		Temperature: 0.8,
		TopP:        1,
		N:           1,
		Stream:      false,
		StreamCallback: func(res llmchain.ChatResponse, done bool, err error) {

		},
		Stop:             []string{},
		MaxTokens:        0,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
		LogitBias:        nil,
		User:             "",
	}

	resp, err := ai.Chat(context.TODO(), req)

	if err != nil {
		t.Errorf(err.Error())
	}

	d, _ := json.Marshal(resp)
	println(string(d))
}
