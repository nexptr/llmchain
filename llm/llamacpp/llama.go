package llamacpp

import (
	"context"
	"fmt"
	"sync"
	"unsafe"

	"github.com/exppii/llmchain"
	"github.com/exppii/llmchain/api/log"
	"github.com/exppii/llmchain/llm"
	"github.com/exppii/llmchain/prompts"
	"github.com/exppii/llmchain/utils"
)

var _ llmchain.LLM = &LLaMACpp{}

type LLaMACpp struct {
	Model string `json:"model" yaml:"model"`

	ModelPath string `json:"model_path" yaml:"model_path"`

	ContextSize int `yaml:"context_size"`
	// Parts       int  `yaml:"parts"`
	Seed            int  `json:"seed" yaml:"seed"`
	F16             bool `json:"f16" yaml:"f16"`
	MLock           bool `yaml:"mlock"`
	EnableEmbedding bool `yaml:"embeddings"`
	NGPULayers      int  `yaml:"ngpu_layers"`

	Threads int `yaml:"threads" json:"threads"`

	Batch int `json:"batch" yaml:"batch"`

	IgnoreEOS     bool    `json:"ignore_eos" yaml:"ignore_eos"`
	RepeatPenalty float64 `json:"repeat_penalty" yaml:"repeat_penalty"`
	Repeat        int     `json:"repeat" yaml:"repeat"`
	Keep          int     `json:"n_keep" yaml:"n_keep"`

	MirostatETA float64 `json:"mirostat_eta" yaml:"mirostat_eta"`
	MirostatTAU float64 `json:"mirostat_tau" yaml:"mirostat_tau"`
	Mirostat    int     `json:"mirostat" yaml:"mirostat"`

	state unsafe.Pointer `json:"-"`

	tmpl *prompts.Template
	// This is still needed, see: https://github.com/ggerganov/llama.cpp/discussions/784
	sync.Mutex
}

func FromYaml(opt llm.ModelOptions) (*LLaMACpp, error) {

	client := defaultLLaMACpp()

	err := llm.UnmarshalPlugin(opt.Settings, client)

	if err != nil {
		return nil, err
	}

	// client.Model = opt.Name

	// return New()
	err = client.load()

	return client, err

}

func New(modelPath string, opts ...ModelOption) (*LLaMACpp, error) {

	// Check if we already have a loaded model
	if !utils.PathExists(modelPath) {
		return nil, fmt.Errorf("model does not exist: %s", modelPath)
	}

	client := defaultLLaMACpp()

	client.ModelPath = modelPath
	for _, fn := range opts {
		fn(client)
	}

	//如果model 名字没有设置，尝试根据加载的文件名设定model 名
	if client.Model == `` {
		client.Model = getModelNameByModelPath(client.ModelPath)
	}

	err := client.load()

	return client, err
}

// Free implements llmchain.LLM
func (l *LLaMACpp) Free() {
	l.freeModel()
}

// Name implements llmchain.LLM
func (l *LLaMACpp) Name() string {
	return l.Model
}

// Call implements llmchain.LLM
func (l *LLaMACpp) Call(ctx context.Context, prompt string) (string, error) {

	opts := defaultOptions()

	templateInput, err := l.tmpl.Render(prompts.H{`input`: prompt})

	if err != nil {
		return "", err
	}

	prediction, err := l.predictWithOpts(templateInput, opts)
	return prediction, err
}

func chatTokenCallback(id string, responses chan llmchain.ChatResponse) func(token string) bool {

	return func(token string) bool {

		resp := llmchain.ChatResponse{
			ID:      id, // we have to return what the user sent here, due to OpenAI spec.
			Choices: []llmchain.Choice{{Delta: &llmchain.Message{Role: "assistant", Content: token}}},
			Object:  "chat.completion.chunk",
		}

		responses <- resp

		log.D(`send: `, resp.String())
		return true

	}

}

// Chat implements llmchain.LLM
func (l *LLaMACpp) Chat(ctx context.Context, req *llmchain.ChatRequest) (llmchain.ChatResponse, error) {

	opts := defaultOptions().buildChatOpts(req)

	n := req.N

	if req.N == 0 {
		n = 1
	}

	//TODO
	s := req.Messages[0].Content

	if req.Stream {

		responses := make(chan llmchain.ChatResponse)
		//todo gen id
		opts.TokenCallback = chatTokenCallback(`id`, responses)

		go func() {
			l.ComputeChoices(s, n, opts, func(s string, c *[]llmchain.Choice) {})

			//done
			close(responses)

			log.D(`exit ComputeChoices process`)
		}()

		for ev := range responses {

			req.StreamCallback(ev, false, nil)

		}

		req.StreamCallback(llmchain.ChatResponse{}, true, nil)

		return llmchain.ChatResponse{}, nil

	}

	result, err := l.ComputeChoices(s, n, opts, func(s string, c *[]llmchain.Choice) {
		*c = append(*c, llmchain.Choice{Message: &llmchain.Message{Role: "assistant", Content: s}})
	})

	if err != nil {

		log.E(`error when compute choices: `, err.Error())
		return llmchain.ChatResponse{}, nil
		//TODO
		// c.JSON(http.StatusInternalServerError, ModelNotExistsErr.WithMessage(err.Error()))
		// return
	}
	log.D(`got compute choices result: `, result)

	resp := llmchain.ChatResponse{
		// Model:   input.Model, // we have to return what the user sent here, due to OpenAI spec.
		Choices: result,
		Object:  "text_completion",
	}

	return resp, nil
}

// Completion implements llmchain.LLM
func (l *LLaMACpp) Completion(ctx context.Context, req *llmchain.CompletionRequest) (*llmchain.CompletionResponse, error) {

	// opts := defaultOptions().buildCompletionOpts(req)

	// n := req.N

	// if req.N == 0 {
	// 	n = 1
	// }

	var result []llmchain.Choice
	// for _, i := range req.Prompt {
	// 	// A model can have a "file.bin.tmpl" file associated with a prompt template prefix

	// 	r, err := l.ComputeChoices(i, n, opts, func(s string, c *[]llmchain.Choice) {
	// 		*c = append(*c, llmchain.Choice{Text: s})
	// 	})

	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	result = append(result, r...)
	// }

	resp := &llmchain.CompletionResponse{
		Model:   req.Model, // we have to return what the user sent here, due to OpenAI spec.
		Choices: result,
		Object:  "text_completion",
	}

	// jsonResult, _ := json.Marshal(resp)
	// log.Debug().Msgf("Response: %s", jsonResult)

	return resp, nil

}

// Embeddings implements LLM
func (l *LLaMACpp) Embeddings(ctx context.Context, req *llmchain.EmbeddingsRequest) (resp *llmchain.EmbeddingsResponse, err error) {

	// p := "/embeddings"
	// return call(ctx, l, http.MethodPost, p, req, resp, nil)
	panic(`TODO`)

}

func (l *LLaMACpp) ComputeChoices(input string, N int, predict *PredictOptions, cb func(string, *[]llmchain.Choice)) ([]llmchain.Choice, error) {

	result := []llmchain.Choice{}

	// get the model function to call for the result

	for i := 0; i < N; i++ {

		prediction, err := l.predictWithOpts(input, predict)

		if err != nil {
			return result, err
		}

		// prediction = Finetune(*config, predInput, prediction)
		cb(prediction, &result)
		//result = append(result, Choice{Text: prediction})

	}
	return result, nil
}
