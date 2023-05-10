package llamacpp

import "github.com/exppii/llmchain/llms"

func defaultLLamaModelOptions(opts ...llms.ModelOption) llms.ModelOptions {

	m := llms.ModelOptions{
		OpenAIRequest: llms.OpenAIRequest{
			Seed:          -1,
			F16:           false,
			Maxtokens:     512,
			TopK:          90,
			Repeat:        64,
			TopP:          0.86,
			Temperature:   0.8,
			RepeatPenalty: 1.1,
			MirostatETA:   0.1,
			MirostatTAU:   5,
			Mirostat:      0,
			Batch:         8, //this field must no zero
		},
		StopWords: []string{`llama`},
		// Debug:       true,
		ContextSize: 512,
		MLock:       false,
		Embeddings:  false,
		Threads:     4,
	}

	for _, opt := range opts {
		opt(&m)
	}

	return m

}

const (
	tailFreeSamplingZ = 1.0

	typicalP         = 1.0
	frequencyPenalty = 0.0
	presencePenalty  = 0.0
	logitBias        = ``
	penalizeNL       = false
)

// type PredictOptions struct {
// 	llms.ModelOptions

// }
