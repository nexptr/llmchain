package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/exppii/llmchain"
	"github.com/exppii/llmchain/api/log"
	"github.com/exppii/llmchain/api/model"
	"github.com/gin-gonic/gin"
)

func chatTokenCallback(model string, responses chan OpenAIResponse) llms.Callback {

	return func(token string) bool {

		resp := OpenAIResponse{
			Model:   model, // we have to return what the user sent here, due to OpenAI spec.
			Choices: []Choice{{Delta: &llms.Message{Role: "assistant", Content: token}}},
			Object:  "chat.completion.chunk",
		}

		responses <- resp

		log.D(`send: `, resp.String())
		return true

	}

}

func chatEndpointHandler(manager *model.Manager) gin.HandlerFunc {

	process := func(llm llms.LLM, s string, opts *llms.ModelOptions, responses chan OpenAIResponse) {
		ComputeChoices(llm, s, opts, func(s string, c *[]Choice) {})

		close(responses)

		log.D(`exit ComputeChoices process`)
	}

	return func(c *gin.Context) {

		input, err := parseReq(c)
		if err != nil {
			log.E("failed reading parameters from request: ", err.Error())
			//todo 从中间件拿取语言类型
			c.JSON(http.StatusBadRequest, ReqArgsErr.WithMessage(err.Error()))
			return
		}

		log.D("input: ", input.String())

		// if input.Model == `gpt-3.5-turbo` {
		// 	input.Model = `ggml-llama-7b`
		// }

		llm, err := manager.GetModel(input.Model)

		if err != nil {

			log.E("model not found or loaded")

			c.JSON(http.StatusInternalServerError, ModelNotExistsErr.WithMessage(err.Error()))
			return
		}

		payload := llm.MergeModelOptions(input)

		log.D(`current payload: `, payload.String())

		templatedInput, err := payload.TemplateMessage(manager.GetPrompt())

		if err != nil {
			//TODO handle error
			log.E("templatedInput format failed")

			c.JSON(http.StatusInternalServerError, ModelNotExistsErr.WithMessage(err.Error()))
			return
		}
		log.D(`current templated inputs: `, templatedInput)

		if input.Stream {
			log.D("Stream request received")

			c.Header("Content-Type", "text/event-stream")
			c.Header("Cache-Control", "no-cache")
			c.Header("Connection", "keep-alive")
			c.Header("Transfer-Encoding", "chunked")
		}

		if input.Stream {

			responses := make(chan OpenAIResponse)

			payload.TokenCallback = chatTokenCallback(llm.Name(), responses)

			go process(llm, templatedInput, payload, responses)

			c.Stream(func(w io.Writer) bool {

				ev, ok := <-responses

				if ok {
					var buf bytes.Buffer
					enc := json.NewEncoder(&buf)
					enc.Encode(ev)
					io.WriteString(w, "event: data\n\n")
					io.WriteString(w, fmt.Sprintf("data: %s\n\n", buf.String()))
					// log.D(`send: `, buf.String())
					//continue
					return true
				}

				io.WriteString(w, "event: data\n\n")

				resp := &OpenAIResponse{
					Model:   input.Model, // we have to return what the user sent here, due to OpenAI spec.
					Choices: []Choice{{FinishReason: "stop"}},
				}
				respData := resp.String()

				io.WriteString(w, fmt.Sprintf("data: %s\n\n", resp.String()))
				log.D("Sending chunk: ", respData)
				//close stream
				return false

			})

			log.D(`finish chat...`)
			return

		}

		result, err := ComputeChoices(llm, templatedInput, payload, func(s string, c *[]Choice) {
			*c = append(*c, Choice{Message: &llms.Message{Role: "assistant", Content: s}})
		})
		if err != nil {

			log.E(`error when compute choices: `, err.Error())
			//TODO
			c.JSON(http.StatusInternalServerError, ModelNotExistsErr.WithMessage(err.Error()))
			return
		}
		log.D(`got compute choices result: `, result)

		resp := &OpenAIResponse{
			Model:   input.Model, // we have to return what the user sent here, due to OpenAI spec.
			Choices: result,
			Object:  "text_completion",
		}

		// Return the prediction in the response body
		c.JSON(http.StatusOK, resp)

		println(llm)

	}
}

func editEndpointHandler(manager *model.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// https://platform.openai.com/docs/api-reference/completions
func completionEndpointHandler(manager *model.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {

		log.I(`parse input...`)
		input := &llmchain.CompletionRequest{}

		if err := c.Bind(input); err != nil {
			// return nil, fmt.Errorf("failed reading parameters from request: ", err.Error())
			log.E("failed reading parameters from request: ", err.Error())
			//todo 从中间件拿取语言类型
			c.JSON(http.StatusBadRequest, ReqArgsErr.WithMessage(err.Error()))
			return
		}

		log.D(`current input:`, input.String())

		llm, err := manager.LLMChain(input.Model, input.Langchain)

		if err != nil {

			log.E("model not found or loaded:", err)

			c.JSON(http.StatusInternalServerError, ModelNotExistsErr.WithMessage(err.Error()))
			return
		}

		resp, err := llm.Completion(context.TODO(), input)

		if err != nil {

			log.E("model not found or loaded:", err)

			c.JSON(http.StatusInternalServerError, ModelNotExistsErr.WithMessage(err.Error()))
			return
		}
		// jsonResult, _ := json.Marshal(resp)
		// log.Debug().Msgf("Response: %s", jsonResult)

		// Return the prediction in the response body
		c.JSON(http.StatusOK, resp)

	}
}

// https://platform.openai.com/docs/api-reference/embeddings
func embeddingsEndpointHandler(manager *model.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.I(`parse embeddings input...`)
		input, err := parseReq(c)
		if err != nil {
			log.E("failed reading parameters from request: ", err.Error())
			//todo 从中间件拿取语言类型
			c.JSON(http.StatusBadRequest, ReqArgsErr.WithMessage(err.Error()))
			return
		}

		llm, err := manager.GetModel(input.Model)

		if err != nil {

			log.E("model not found or loaded")

			c.JSON(http.StatusInternalServerError, ModelNotExistsErr.WithMessage(err.Error()))
			return
		}

		payload := llm.MergeModelOptions(input)

		items := []llms.Item{}

		for i, s := range payload.InputToken {
			// get the model function to call for the result
			embedFn, err := LLMEmbedding(llm, "", s, payload)

			if err != nil {

				log.E("model not found or loaded")

				c.JSON(http.StatusInternalServerError, ModelNotExistsErr.WithMessage(err.Error()))
				return
			}

			embeddings, err := embedFn()
			if err != nil {

				log.E("model not found or loaded")

				c.JSON(http.StatusInternalServerError, ModelNotExistsErr.WithMessage(err.Error()))
				return
			}
			items = append(items, llms.Item{Embedding: embeddings, Index: i, Object: "embedding"})
		}

		for i, s := range payload.InputStrings {
			// get the model function to call for the result
			embedFn, err := LLMEmbedding(llm, s, []int{}, payload)
			if err != nil {

				log.E("model not found or loaded")

				c.JSON(http.StatusInternalServerError, ModelNotExistsErr.WithMessage(err.Error()))
				return
			}

			embeddings, err := embedFn()
			if err != nil {

				log.E("model not found or loaded")

				c.JSON(http.StatusInternalServerError, ModelNotExistsErr.WithMessage(err.Error()))
				return
			}
			items = append(items, llms.Item{Embedding: embeddings, Index: i, Object: "embedding"})
		}

		resp := &OpenAIResponse{
			Model:  input.Model, // we have to return what the user sent here, due to OpenAI spec.
			Data:   items,
			Object: "list",
		}

		log.D(`embeddings resp:`, resp.String())

		// Return the prediction in the response body
		c.JSON(http.StatusOK, resp)
	}
}

func listModelsHandler(manager *model.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {

		log.I(`received list model req`)
		list := manager.ListModels()

		models := []OpenAIModel{}
		for _, m := range list {
			models = append(models, OpenAIModel{ID: m, Object: "model"})
		}

		resp := struct {
			Object string        `json:"object"`
			Data   []OpenAIModel `json:"data"`
		}{
			Object: "list",
			Data:   models,
		}

		c.JSON(http.StatusOK, resp)

	}
}

func chatCallback(resp chan llmchain.ChatResponse, done chan error) llmchain.SreamCallBack {
	return func(r llmchain.ChatResponse, d bool, e error) {
		if d {
			done <- e
		} else {
			resp <- r
		}
	}
}

func chatEndpointHandler(manager *model.Manager) gin.HandlerFunc {

	return func(c *gin.Context) {

		input := &llmchain.ChatRequest{}

		if err := c.Bind(input); err != nil {
			// return nil, fmt.Errorf("failed reading parameters from request: ", err.Error())
			log.E("failed reading parameters from request: ", err.Error())
			//todo 从中间件拿取语言类型
			c.JSON(http.StatusBadRequest, ReqArgsErr.WithMessage(err.Error()))
			return
		}

		log.D(`current input:`, input.String())

		llm, err := manager.LLMChain(input.Model, input.Langchain)

		if err != nil {

			log.E("model not found or loaded:", err)

			c.JSON(http.StatusInternalServerError, ModelNotExistsErr.WithMessage(err.Error()))
			return
		}

		if input.Stream {

			data := make(chan llmchain.ChatResponse)
			done := make(chan error)
			defer close(data)
			defer close(done)

			input.StreamCallback = chatCallback(data, done)

			resp, err := llm.Chat(context.TODO(), input)

			if err != nil {
				log.E("run stream completion failed: ", err.Error(), resp.String())
				// return resp, err
				c.JSON(http.StatusInternalServerError, ModelNotExistsErr.WithMessage(err.Error()))
				return
			}

			c.Header("Content-Type", "text/event-stream")
			c.Header("Cache-Control", "no-cache")
			c.Header("Connection", "keep-alive")
			c.Header("Transfer-Encoding", "chunked")

			c.Stream(func(w io.Writer) bool {

				for {
					select {
					case payload := <-data:

						var buf bytes.Buffer
						enc := json.NewEncoder(&buf)
						enc.Encode(payload)
						io.WriteString(w, "event: data\n\n")
						io.WriteString(w, fmt.Sprintf("data: %s\n\n", buf.String()))
						// log.D(`send: `, buf.String())
						//continue
						return true

						// fmt.Print(payload.Choices[0].Delta.Content)
					case err = <-done:

						io.WriteString(w, "event: data\n\n")

						resp := &llmchain.ChatResponse{
							// Model:   input.Model, // we have to return what the user sent here, due to OpenAI spec.
							Choices: []llmchain.Choice{{FinishReason: "stop"}},
						}
						respData := resp.String()

						io.WriteString(w, fmt.Sprintf("data: %s\n\n", resp.String()))
						log.D("Sending chunk: ", respData)
						//close stream
						return false

						// fmt.Print("\n")
						// return res, err
					}
				}

			})

			log.D(`finish chat...`)
			return

		}

		//completion
		resp, err := llm.Chat(context.TODO(), input)

		if err != nil {
			log.E(`run chat without stream: `, err.Error())
			c.JSON(http.StatusInternalServerError, ReqArgsErr.WithMessage(err.Error()))

		}
		c.JSON(http.StatusOK, resp)

	}

}
