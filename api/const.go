package api

import "github.com/exppii/llmchain/llms"

var (
	ReqArgsErr        = Response{Code: 400}
	ModelNotExistsErr = Response{Code: 1000}
)

type Choice = llms.Choice

type OpenAIResponse = llms.OpenAIResponse

type OpenAIRequest = llms.OpenAIRequest
