package llamacpp

import (
	"strings"

	"github.com/exppii/llmchain/llm"
	"github.com/exppii/llmchain/prompts"
)

//由于llamacpp 支持较多的模型，当时每个模型的 末日prompt可能都不一样。这里通过简单的配置适配基本的模型

const vicunaTemplate = `User: {{.input}}
ASSISTANT:`

const defaultTemplate = `User: {{.input}}
ASSISTANT:
`

// Name implements llmchain.LLM
func loadModelPrompts(model string) *prompts.Template {
	//todo
	switch model {
	case llm.VICUNA_13B:
		tmpl, _ := prompts.New(vicunaTemplate)
		return tmpl
	}

	return nil
}

func getModelNameByModelPath(modelPath string) string {

	if strings.Contains(modelPath, `vincuna`) {
		return llm.VICUNA_13B
	}
	return llm.LLaMA_7B
}
