package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/exppii/llmchain/api/log"
	"github.com/exppii/llmchain/api/model"
	"github.com/exppii/llmchain/llms"
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
		log.D(`current payload: `, payload.String())
		templatedInputs, err := payload.TemplatePromptStrings(manager.GetPrompt())
		log.D(`current templated inputs: `, strings.Join(templatedInputs, `,`))
		if err != nil {

			log.E("err when template input data: ", err.Error())
			//TODO
			c.JSON(http.StatusInternalServerError, ModelNotExistsErr.WithMessage(err.Error()))
			return
		}

		fn := func(s string, c *[]Choice) { *c = append(*c, Choice{Text: s}) }

		payload.TokenCallback = func(s string) bool {
			print(s)
			return true
		}

		var result []Choice
		for _, i := range templatedInputs {
			log.D(`compute choices for inputs: `, i)
			r, err := ComputeChoices(llm, i, payload, fn)
			if err != nil {

				log.E(`error when compute choices: `, err.Error())
				//TODO
				c.JSON(http.StatusInternalServerError, ModelNotExistsErr.WithMessage(err.Error()))
				return
			}
			log.D(`got compute choices result: `, r)
			result = append(result, r...)
		}

		resp := &OpenAIResponse{
			Model:   input.Model, // we have to return what the user sent here, due to OpenAI spec.
			Choices: result,
			Object:  "text_completion",
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

// parseReq parse openai format req
func parseReq(c *gin.Context) (*OpenAIRequest, error) {

	input := &OpenAIRequest{}

	if err := c.Bind(input); err != nil {
		return nil, fmt.Errorf("failed reading parameters from request: ", err.Error())

	}
	return input, nil
}
