package llmchain

import (
	"encoding/json"
	"fmt"
)

type ObjectType string

const (
	OTModel           ObjectType = "model"
	OTModelPermission ObjectType = "model_permission"
	OTList            ObjectType = "list"
	OTEdit            ObjectType = "edit"
	OTTextCompletion  ObjectType = "text_completion"
	OTEEmbedding      ObjectType = "embedding"
	OTFile            ObjectType = "file"
	OTFineTune        ObjectType = "fine-tune"
	OTFineTuneEvent   ObjectType = "fine-tune-event"
)

type SreamCallBack func(res ChatResponse, done bool, err error)

// ChatRequest:
// https://platform.openai.com/docs/guides/chat/chat-completions-beta
// https://platform.openai.com/docs/api-reference/chat
type ChatRequest struct {

	// Model: ID of the model to use.
	// Currently, only gpt-3.5-turbo and gpt-3.5-turbo-0301 are supported.
	Model string `json:"model" yaml:"model"`

	// Messages: The messages to generate chat completions for, in the chat format.
	// https://platform.openai.com/docs/guides/chat/introduction
	// Including the conversation history helps when user instructions refer to prior messages.
	// In the example above, the user's final question of "Where was it played?" only makes sense in the context of the prior messages about the World Series of 2020.
	// Because the models have no memory of past requests, all relevant information must be supplied via the conversation.
	// If a conversation cannot fit within the model's token limit, it will need to be shortened in some way.
	Messages []Message `json:"messages" yaml:"messages"`

	// Temperature: What sampling temperature to use, between 0 and 2.
	// Higher values like 0.8 will make the output more random, while lower values like 0.2 will make it more focused and deterministic.
	// We generally recommend altering this or top_p but not both.
	// Defaults to 1.
	Temperature float32 `json:"temperature,omitempty" yaml:"temperature"`

	// TopP: An alternative to sampling with temperature, called nucleus sampling,
	// where the model considers the results of the tokens with top_p probability mass.
	// So 0.1 means only the tokens comprising the top 10% probability mass are considered.
	// We generally recommend altering this or temperature but not both.
	// Defaults to 1.
	TopP float32 `json:"top_p,omitempty" yaml:"top_p"`

	// N: How many chat completion choices to generate for each input message.
	// Defaults to 1.
	N int `json:"n,omitempty"`

	// Stream: If set, partial message deltas will be sent, like in ChatGPT.
	// Tokens will be sent as data-only server-sent events as they become available,
	// with the stream terminated by a data: [DONE] message.
	Stream bool `json:"stream,omitempty" `

	// StreamCallback is a callback funciton to handle stream response.
	// If provided, this library automatically set `Stream` `true`.
	// This field is added by github.com/otiai10/openaigo only to handle Stream.
	// Thus, it is omitted when the client excute HTTP request.
	StreamCallback func(res ChatResponse, done bool, err error) `json:"-" yaml:"-"`

	// Stop: Up to 4 sequences where the API will stop generating further tokens.
	// Defaults to null.
	Stop []string `json:"stop,omitempty" yaml:"stop"`

	// MaxTokens: The maximum number of tokens allowed for the generated answer.
	// By default, the number of tokens the model can return will be (4096 - prompt tokens).
	MaxTokens int `json:"max_tokens,omitempty" yaml:"max_tokens"`

	// PresencePenalty: Number between -2.0 and 2.0.
	// Positive values penalize new tokens based on whether they appear in the text so far,
	// increasing the model's likelihood to talk about new topics.
	// See more information about frequency and presence penalties.
	// https://platform.openai.com/docs/api-reference/parameter-details
	PresencePenalty float32 `json:"presence_penalty,omitempty"`

	// FrequencyPenalty: Number between -2.0 and 2.0.
	// Positive values penalize new tokens based on their existing frequency in the text so far,
	// decreasing the model's likelihood to repeat the same line verbatim.
	// See more information about frequency and presence penalties.
	// https://platform.openai.com/docs/api-reference/parameter-details
	FrequencyPenalty float32 `json:"frequency_penalty,omitempty"`

	// LogitBias: Modify the likelihood of specified tokens appearing in the completion.
	// Accepts a json object that maps tokens (specified by their token ID in the tokenizer)
	// to an associated bias value from -100 to 100.
	// Mathematically, the bias is added to the logits generated by the model prior to sampling.
	// The exact effect will vary per model, but values between -1 and 1 should decrease or increase likelihood of selection;
	// values like -100 or 100 should result in a ban or exclusive selection of the relevant token.
	LogitBias map[string]int `json:"logit_bias,omitempty"`

	// User: A unique identifier representing your end-user, which can help OpenAI to monitor and detect abuse. Learn more.
	// https://platform.openai.com/docs/guides/safety-best-practices/end-user-ids
	User string `json:"user,omitempty"`

	// Custom parameters - not present in the OpenAI API

	// Langchain using Langchain default: baseChat
	Langchain string `json:"langchain,omitempty" yaml:"langchain"`
}

func (r *ChatRequest) String() string {
	j, _ := json.Marshal(r)
	return string(j)
}

type CompletionRequest struct {

	// Model: ID of the model to use.
	// You can use the List models API to see all of your available models, or see our Model overview for descriptions of them.
	// See https://beta.openai.com/docs/api-reference/completions/create#completions/create-model
	Model string `json:"model"`

	// Prompt: The prompt(s) to generate completions for, encoded as a string, array of strings, array of tokens, or array of token arrays.
	// Note that <|endoftext|> is the document separator that the model sees during training, so if a prompt is not specified the model will generate as if from the beginning of a new document.
	// See https://beta.openai.com/docs/api-reference/completions/create#completions/create-prompt
	// Prompt []string `json:"prompt"`

	Prompt interface{} `json:"prompt"`
	// MaxTokens: The maximum number of tokens to generate in the completion.
	// The token count of your prompt plus max_tokens cannot exceed the model's context length. Most models have a context length of 2048 tokens (except for the newest models, which support 4096).
	// See https://beta.openai.com/docs/api-reference/completions/create#completions/create-max_tokens
	MaxTokens int `json:"max_tokens,omitempty"`

	// Temperature: What sampling temperature to use. Higher values means the model will take more risks. Try 0.9 for more creative applications, and 0 (argmax sampling) for ones with a well-defined answer.
	// We generally recommend altering this or top_p but not both.
	// See https://beta.openai.com/docs/api-reference/completions/create#completions/create-temperature
	Temperature float32 `json:"temperature,omitempty"`

	// Suffix: The suffix that comes after a completion of inserted text.
	// See https://beta.openai.com/docs/api-reference/completions/create#completions/create-suffix
	Suffix string `json:"suffix,omitempty"`

	// TopP: An alternative to sampling with temperature, called nucleus sampling,
	// where the model considers the results of the tokens with top_p probability mass.
	// So 0.1 means only the tokens comprising the top 10% probability mass are considered.
	// We generally recommend altering this or temperature but not both.
	// See https://beta.openai.com/docs/api-reference/completions/create#completions/create-top_p
	TopP float32 `json:"top_p,omitempty"`

	// N: How many completions to generate for each prompt.
	// Note: Because this parameter generates many completions, it can quickly consume your token quota.
	// Use carefully and ensure that you have reasonable settings for max_tokens and stop.
	// See https://beta.openai.com/docs/api-reference/completions/create#completions/create-n
	N int `json:"n,omitempty"`

	// Stream: Whether to stream back partial progress.
	// If set, tokens will be sent as data-only server-sent events as they become available,
	// with the stream terminated by a data: [DONE] message.
	// See https://beta.openai.com/docs/api-reference/completions/create#completions/create-stream
	Stream bool `json:"stream,omitempty"`

	// LogProbs: Include the log probabilities on the logprobs most likely tokens, as well the chosen tokens.
	// For example, if logprobs is 5, the API will return a list of the 5 most likely tokens. The API will always return the logprob of the sampled token, so there may be up to logprobs+1 elements in the response.
	// The maximum value for logprobs is 5. If you need more than this, please contact us through our Help center and describe your use case.
	// See https://beta.openai.com/docs/api-reference/completions/create#completions/create-logprobs
	LogProbs int `json:"logprobs,omitempty"`

	// Echo: Echo back the prompt in addition to the completion.
	// See https://beta.openai.com/docs/api-reference/completions/create#completions/create-echo
	Echo bool `json:"echo,omitempty"`

	// Stop: Up to 4 sequences where the API will stop generating further tokens. The returned text will not contain the stop sequence.
	// See https://beta.openai.com/docs/api-reference/completions/create#completions/create-stop
	Stop []string `json:"stop,omitempty"`

	// PresencePenalty: Number between -2.0 and 2.0.
	// Positive values penalize new tokens based on whether they appear in the text so far, increasing the model's likelihood to talk about new topics.
	// See more information about frequency and presence penalties.
	// See https://beta.openai.com/docs/api-reference/completions/create#completions/create-presence_penalty
	PresencePenalty float32 `json:"presence_penalty,omitempty"`

	// FrequencyPenalty: Number between -2.0 and 2.0.
	// Positive values penalize new tokens based on their existing frequency in the text so far, decreasing the model's likelihood to repeat the same line verbatim.
	// See more information about frequency and presence penalties.
	// See https://beta.openai.com/docs/api-reference/completions/create#completions/create-frequency_penalty
	FrequencyPenalty float32 `json:"frequency_penalty,omitempty"`

	// BestOf: Generates best_of completions server-side and returns the "best" (the one with the highest log probability per token). Results cannot be streamed.
	// When used with n, best_of controls the number of candidate completions and n specifies how many to return - best_of must be greater than n.
	// Note: Because this parameter generates many completions, it can quickly consume your token quota. Use carefully and ensure that you have reasonable settings for max_tokens and stop.
	// See https://beta.openai.com/docs/api-reference/completions/create#completions/create-best_of
	BestOf int `json:"best_of,omitempty"`

	// LogitBias: Modify the likelihood of specified tokens appearing in the completion.
	// Accepts a json object that maps tokens (specified by their token ID in the GPT tokenizer) to an associated bias value from -100 to 100. You can use this tokenizer tool (which works for both GPT-2 and GPT-3) to convert text to token IDs. Mathematically, the bias is added to the logits generated by the model prior to sampling. The exact effect will vary per model, but values between -1 and 1 should decrease or increase likelihood of selection; values like -100 or 100 should result in a ban or exclusive selection of the relevant token.
	// As an example, you can pass {"50256": -100} to prevent the <|endoftext|> token from being generated.
	// See https://beta.openai.com/docs/api-reference/completions/create#completions/create-logit_bias
	LogitBias map[string]int `json:"logit_bias,omitempty"`

	// User: A unique identifier representing your end-user, which can help OpenAI to monitor and detect abuse. Learn more.
	// See https://beta.openai.com/docs/api-reference/completions/create#completions/create-user
	User string `json:"user,omitempty"`

	// Custom parameters - not present in the OpenAI API

	// Langchain using Langchain default: baseChat
	Langchain string `json:"langchain,omitempty" yaml:"langchain"`

	//prompt strings
	PromptStrings []string `json:"-" yaml:"-"`
}

func (r *CompletionRequest) String() string {
	j, _ := json.Marshal(r)
	return string(j)
}

// ChatMessage: An element of messages parameter.
// The main input is the messages parameter. Messages must be an array of message objects,
// where each object has a role (either "system", "user", or "assistant")
// and content (the content of the message).
// Conversations can be as short as 1 message or fill many pages.
type Message struct {

	// Role: Either of "system", "user", "assistant".
	// Typically, a conversation is formatted with a system message first, followed by alternating user and assistant messages.
	// The system message helps set the behavior of the assistant. In the example above, the assistant was instructed with "You are a helpful assistant."
	// The user messages help instruct the assistant. They can be generated by the end users of an application, or set by a developer as an instruction.
	// The assistant messages help store prior responses. They can also be written by a developer to help give examples of desired behavior.
	Role string `json:"role,omitempty" yaml:"role"`

	// Content: A content of the message.
	Content string `json:"content,omitempty" yaml:"content"`
}

type CompletionResponse struct {
	ID      string     `json:"id,omitempty"`
	Object  ObjectType `json:"object,omitempty"`
	Created int64      `json:"created,omitempty"`
	Model   string     `json:"model,omitempty"`
	Choices []Choice   `json:"choices,omitempty"`
	Usage   Usage      `json:"usage"`
}

type ChatResponse struct {
	ID      string   `json:"id,omitempty"`
	Object  string   `json:"object,omitempty"`
	Created int64    `json:"created,omitempty"`
	Choices []Choice `json:"choices,omitempty"`
	Data    []Item   `json:"data,omitempty"`

	Usage Usage `json:"usage"`
}

func (r *ChatResponse) String() string {
	j, _ := json.Marshal(r)
	return string(j)
}

type Choice struct {
	Index        int      `json:"index,omitempty"`
	FinishReason string   `json:"finish_reason,omitempty"`
	Message      *Message `json:"message,omitempty"`
	Delta        *Message `json:"delta,omitempty"`
	Text         string   `json:"text,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Item struct {
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
	Object    string    `json:"object,omitempty"`
}

type EmbeddingsRequest struct {
	Model string      `json:"model"`
	Input interface{} `json:"input"` //string or []string
	User  string      `json:"user,omitempty"`
}

func (r *EmbeddingsRequest) Verify() error {
	switch r.Input.(type) {
	case string, []string:
		return nil
	}
	return fmt.Errorf(`input field must be string or []string`)
}

func (r *EmbeddingsRequest) String() string {
	j, _ := json.Marshal(r)
	return string(j)
}

type EmbeddingsResponse struct {
	Object string          `json:"object"`
	Data   []EmbeddingData `json:"data"`
	Usage  Usage           `json:"usage"`
}

func (r *EmbeddingsResponse) String() string {
	j, _ := json.Marshal(r)
	return string(j)
}

type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}
