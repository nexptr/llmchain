package chain

import (
	"context"
	"fmt"
	"net/http"

	"github.com/exppii/llmchain"
	"github.com/exppii/llmchain/prompts"
)

// BaseChat base chat Lang chain,this chain just do nothing
type APIChain struct {
	name string

	reqTempl *prompts.Template

	respTempl *prompts.Template

	client *http.Client
	l      llmchain.LLM
}

var _ llmchain.Chain = &APIChain{}

func NewAPIChain(name, docs string) *APIChain {

	reqTempl := fmt.Sprintf(apiReqTemplate, docs)
	respTempl := fmt.Sprintf(apiRespTemplate, docs)

	q, _ := prompts.New(reqTempl, `Input`)

	p, _ := prompts.New(respTempl, `Request`, `APIResp`)

	return &APIChain{
		name:      name,
		reqTempl:  q,
		respTempl: p,
	}
}

// Name implements llmchain.Chain
func (c *APIChain) Name() string {
	return c.name
}

// ChatPrompt implements llmchain.Chain,考虑到message是多轮对话，我们这里应该只是对最后一个User的问题进行转化。但是这样可能，导致之前对历史问题丢失上下文。
// 这里当前处置手段是，当前的默认策略为：
// 将当前对话转化为：历史原始问题+历史答案，+封装后的最后的问题
func (*APIChain) ChatPrompt(ctx context.Context, messages []llmchain.Message) ([]llmchain.Message, error) {
	panic("unimplemented")
}

// Prompt implements llmchain.Chain
func (*APIChain) Prompt(input string) string {
	panic("unimplemented")
}

// Prompt implements llmchain.Chain
func (c *APIChain) PromptArgs(args map[string]string) (string, error) {

	//input should be APIdocs and Question

	p, err := c.reqTempl.Render(args)
	if err != nil {

		return "", err
	}

	//since we have api req prompt, send to llm to get truely req
	reqStr, err := c.l.Call(context.Background(), p)

	if err != nil {

		return "", err
	}

	//send real req with http client
	respStr, err := c.sendReqWithHttpClient(reqStr)

	if err != nil {

		return "", err
	}

	sp, err := c.respTempl.Render(prompts.H{`Request`: reqStr, `APIResp`: respStr})
	if err != nil {

		return "", err
	}

	//since we have api req prompt, send to llm to get truely req
	return c.l.Call(context.Background(), sp)
}

// WithLLM implements llmchain.Chain
func (c *APIChain) sendReqWithHttpClient(req string) (string, error) {
	if c.client == nil {
		c.client = http.DefaultClient
	}

	//todo parse reqstring to real url and data
	return `ok`, nil

}

// WithLLM implements llmchain.Chain
func (c *APIChain) WithLLM(llm llmchain.LLM) {
	c.l = llm
}

const apiReqTemplate = `You are given the below API Documentation:
%s
Using this documentation, generate the full request to call for answering the user question.
You should build the API request in order to get a response that is as short as possible, while still getting the necessary information to answer the question. Pay attention to deliberately exclude any unnecessary pieces of data in the API call.

Question:{{.Input}}
API request:`

const apiRespTemplate = apiReqTemplate + ` {{.Request}}

Here is the response from the API:

{{.APIResp}}

Summarize this response to answer the original question.

Summary:`
