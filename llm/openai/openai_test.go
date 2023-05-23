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
		MaxTokens:        1024,
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

func TestOpenAI_Completion(t *testing.T) {

	ai := openai.New(openai.WithAPIHost(`http://192.168.12.213:8000/v1`), openai.WithModel(`vicuna-13b-v1.1`))

	// msg := llmchain.Message{Role: `user`, Content: `怎样计算圆形面积`}
	req := &llmchain.CompletionRequest{
		Model: `vicuna-13b-v1.1`,
		Prompt: `USER: You are given the below API Documentation:
This API endpoint will search the notes for a user.

    Endpoint: https://example.com
    POST /api/notes

    Query parameters:
    q | string | The search term for notes
    size | int | The limit for notes

Using this documentation, generate the full API url to call for answering the user question.
You should build the Request in order to get a response that is as short as possible, while still getting the necessary information to answer the question. Pay attention to deliberately exclude any unnecessary pieces of data in the API call.

Question:Search for notes containing langchain
Request:
ASSISTANT:`, //[]string{`怎样计算圆形面积`}
		// Messages:    []llmchain.Message{msg},
		Temperature: 0.1,
		TopP:        1,
		N:           1,
		Stream:      false,
		// StreamCallback: func(res llmchain.ChatResponse, done bool, err error) {

		// },
		Stop:             []string{},
		MaxTokens:        0,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
		LogitBias:        nil,
		User:             "",
	}

	resp, err := ai.Completion(context.TODO(), req)

	if err != nil {
		t.Errorf(err.Error())
	}

	d, _ := json.Marshal(resp)
	println(string(d))
}

func TestOpenAI_Embeddings(t *testing.T) {

	ai := openai.New(openai.WithAPIHost(`http://192.168.12.213:8000/v1`), openai.WithModel(`vicuna-13b-v1.1`))

	// msg := llmchain.Message{Role: `user`, Content: `怎样计算圆形面积`}
	req := &llmchain.EmbeddingsRequest{
		Model: `vicuna-13b-v1.1`,
		Input: `USER: You are given the below API Documentation:
This API endpoint will search the notes for a user.

    Endpoint: https://example.com
    POST /api/notes

    Query parameters:
    q | string | The search term for notes
    size | int | The limit for notes

Using this documentation, generate the full API url to call for answering the user question.
You should build the Request in order to get a response that is as short as possible, while still getting the necessary information to answer the question. Pay attention to deliberately exclude any unnecessary pieces of data in the API call.

Question:Search for notes containing langchain
Request:
ASSISTANT:`, //[]string{`怎样计算圆形面积`}
		// Messages:    []llmchain.Message{msg},

		User: "",
	}

	resp, err := ai.Embeddings(context.TODO(), req)

	if err != nil {
		t.Errorf(err.Error())
	}

	d, _ := json.Marshal(resp)
	println(string(d))
}
