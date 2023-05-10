package llms

type ModelType int

const (
	ModelUnknown  ModelType = 0   //
	ModelOpenAI   ModelType = 100 //
	ModelLLaMA    ModelType = 200 //
	ModelLLaMACPP ModelType = 300 //
	ModelGPT4All  ModelType = 400 //
)

// modelMap store the model name -> model type maps
var modelMap map[string]ModelType

func init() {
	modelMap = map[string]ModelType{
		"gpt-4":              ModelOpenAI,
		"gpt-4-0314":         ModelOpenAI,
		"gpt-4-32k":          ModelOpenAI,
		"gpt-4-32k-0314":     ModelOpenAI,
		"gpt-3.5-turbo":      ModelOpenAI,
		"gpt-3.5-turbo-0301": ModelOpenAI,
		"text-ada-001":       ModelOpenAI,
		"ada":                ModelOpenAI,
		"text-babbage-001":   ModelOpenAI,
		"babbage":            ModelOpenAI,
		"text-curie-001":     ModelOpenAI,
		"curie":              ModelOpenAI,
		"davinci":            ModelOpenAI,
		"text-davinci-003":   ModelOpenAI,
		"text-davinci-002":   ModelOpenAI,
		"code-davinci-002":   ModelOpenAI,
		"code-davinci-001":   ModelOpenAI,
		"code-cushman-002":   ModelOpenAI,
		"code-cushman-001":   ModelOpenAI,
		"ggml-llama-7b":      ModelLLaMACPP,
		"ggml-llama-13b":     ModelLLaMACPP,
		"ggml-vicuna-13b":    ModelLLaMACPP,
	}
}

func GetModelType(model string) ModelType {

	t, ok := modelMap[model]

	if ok {
		return t
	}
	return ModelUnknown
}
