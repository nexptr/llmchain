package llamacpp

import (
	"github.com/exppii/llmchain"
)

// CallOption is a function that configures a LLM.
type ModelOption func(*LLaMACpp)

func defaultLLaMACpp() *LLaMACpp {

	return &LLaMACpp{
		Model:         "",
		ModelPath:     "",
		Batch:         0,
		F16:           false,
		IgnoreEOS:     false,
		RepeatPenalty: 0,
		Repeat:        0,
		Keep:          0,
		MirostatETA:   0,
		MirostatTAU:   0,
		Mirostat:      0,
	}

}

// WithContext sets the context size.
func WithContext(c int) ModelOption {
	return func(p *LLaMACpp) {
		p.ContextSize = c
	}
}

// WithContext sets the context size.
// func WithTemplateConfig(c prompts.TemplateConfig) ModelOption {
// 	return func(p *LLaMACpp) {
// 		p.TemplateConfig = c
// 	}
// }

func WithModelSeed(c int) ModelOption {
	return func(p *LLaMACpp) {
		p.Seed = c
	}
}

func WithNGPULayers(c int) ModelOption {
	return func(p *LLaMACpp) {
		p.NGPULayers = c
	}
}

var EnableEmbeddings ModelOption = func(p *LLaMACpp) { p.EnableEmbedding = true }

var EnableF16Memory ModelOption = func(p *LLaMACpp) { p.F16 = true }

const (
	tailFreeSamplingZ = 1.0

	typicalP         = 1.0
	frequencyPenalty = 0.0
	presencePenalty  = 0.0
	logitBias        = ``
	penalizeNL       = false
)

type PredictOptions struct {
	Seed, Threads, Maxtokens, TopK, Repeat, Batch, NKeep int
	TopP, Temperature, Penalty                           float32
	F16                                                  bool
	DebugMode                                            bool
	StopWords                                            []string
	IgnoreEOS                                            bool

	TailFreeSamplingZ float32
	TypicalP          float32
	FrequencyPenalty  float32
	PresencePenalty   float32
	Mirostat          int
	MirostatETA       float32
	MirostatTAU       float32
	PenalizeNL        bool
	LogitBias         string
	TokenCallback     func(string) bool
}

type PredictOption func(p *PredictOptions)

func defaultOptions() *PredictOptions {
	return &PredictOptions{
		Seed:              -1,
		Threads:           4,
		Maxtokens:         128,
		Penalty:           1.1,
		Repeat:            64,
		Batch:             8,
		NKeep:             64,
		TopK:              40,
		TopP:              0.95,
		TailFreeSamplingZ: 1.0,
		TypicalP:          1.0,
		Temperature:       0.8,
		FrequencyPenalty:  0.0,
		PresencePenalty:   0.0,
		Mirostat:          0,
		MirostatTAU:       5.0,
		MirostatETA:       0.1,
	}
}

func (m *PredictOptions) buildChatOpts(req *llmchain.ChatRequest) *PredictOptions {
	// Generate the prediction using the language model

	m.Temperature = req.Temperature

	m.TopP = req.TopP
	m.Maxtokens = req.MaxTokens

	// if req.Mirostat != 0 {
	// 	predictOptions = append(predictOptions, llama.SetMirostat(c.Mirostat))
	// }

	// if req.MirostatETA != 0 {
	// 	predictOptions = append(predictOptions, SetMirostatETA(req.MirostatETA))
	// }

	// if req.MirostatTAU != 0 {
	// 	predictOptions = append(predictOptions, SetMirostatTAU(req.MirostatTAU))
	// }

	// if req.Debug {
	// 	predictOptions = append(predictOptions, Debug)
	// }

	// predictOptions = append(predictOptions, SetStopWords(req.StopWords...))

	// if req.RepeatPenalty != 0 {
	// 	predictOptions = append(predictOptions, SetPenalty(req.RepeatPenalty))
	// }

	// if req.Keep != 0 {
	// 	predictOptions = append(predictOptions, SetNKeep(req.Keep))
	// }

	// if req.Batch != 0 {
	// 	predictOptions = append(predictOptions, SetBatch(req.Batch))
	// }

	// if req.F16 {
	// 	predictOptions = append(predictOptions, EnableF16KV)
	// }

	// if req.IgnoreEOS {
	// 	predictOptions = append(predictOptions, IgnoreEOS)
	// }

	// if req.Seed != 0 {
	// 	predictOptions = append(predictOptions, SetSeed(req.Seed))
	// }
	// func chatTokenCallback(model string, responses chan OpenAIResponse) llms.Callback {

	// 	return func(token string) bool {

	// 		resp := OpenAIResponse{
	// 			Model:   model, // we have to return what the user sent here, due to OpenAI spec.
	// 			Choices: []Choice{{Delta: &llms.Message{Role: "assistant", Content: token}}},
	// 			Object:  "chat.completion.chunk",
	// 		}

	// 		responses <- resp

	// 		log.D(`send: `, resp.String())
	// 		return true

	// 	}

	// }

	return m
}

func (m *PredictOptions) buildCompletionOpts(req *llmchain.CompletionRequest) *PredictOptions {
	// Generate the prediction using the language model

	m.Temperature = req.Temperature

	m.TopP = req.TopP
	m.Maxtokens = req.MaxTokens

	// if req.Mirostat != 0 {
	// 	predictOptions = append(predictOptions, llama.SetMirostat(c.Mirostat))
	// }

	// if req.MirostatETA != 0 {
	// 	predictOptions = append(predictOptions, SetMirostatETA(req.MirostatETA))
	// }

	// if req.MirostatTAU != 0 {
	// 	predictOptions = append(predictOptions, SetMirostatTAU(req.MirostatTAU))
	// }

	// if req.Debug {
	// 	predictOptions = append(predictOptions, Debug)
	// }

	// predictOptions = append(predictOptions, SetStopWords(req.StopWords...))

	// if req.RepeatPenalty != 0 {
	// 	predictOptions = append(predictOptions, SetPenalty(req.RepeatPenalty))
	// }

	// if req.Keep != 0 {
	// 	predictOptions = append(predictOptions, SetNKeep(req.Keep))
	// }

	// if req.Batch != 0 {
	// 	predictOptions = append(predictOptions, SetBatch(req.Batch))
	// }

	// if req.F16 {
	// 	predictOptions = append(predictOptions, EnableF16KV)
	// }

	// if req.IgnoreEOS {
	// 	predictOptions = append(predictOptions, IgnoreEOS)
	// }

	// if req.Seed != 0 {
	// 	predictOptions = append(predictOptions, SetSeed(req.Seed))
	// }

	return m
}
