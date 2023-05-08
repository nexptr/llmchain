package model

import (
	"fmt"
	"text/template"

	"github.com/exppii/llmchain/llms"
	"github.com/exppii/llmchain/llms/llamacpp"
	"github.com/exppii/llmchain/llms/openai"
)

// Manager LLM model manager
type Manager struct {
	//ModelPath root mode path for llms
	ModelPath string

	loadedModels map[string]llms.LLM

	promptsTemplates map[string]*template.Template
}

func NewModelManager(modelPath string) *Manager {
	return &Manager{
		ModelPath: modelPath,

		loadedModels: map[string]llms.LLM{},

		promptsTemplates: make(map[string]*template.Template),
	}
}

func (m *Manager) GetModel(modelName string) (llms.LLM, bool) {

	model, exists := m.loadedModels[modelName]

	return model, exists
}

func (m *Manager) LoadLLaMACpp(modelName string, opts ...llamacpp.ModelOption) (*llamacpp.LLaMACpp, error) {

	panic(`TODO`)

}

func (m *Manager) LoadOpenAI(modelName string, opts ...openai.ModelOption) (*openai.OpenAI, error) {

	return nil, fmt.Errorf("openai llm todo")
}

func (m *Manager) GreedyLoad(modelFile string, llamaOpts []llamacpp.ModelOption, threads uint32) (llms.LLM, error) {

	model, exists := m.loadedModels[modelFile]
	if exists {
		// muModels.Unlock()
		return model, nil
	}

	//try

	if model, err := m.LoadLLaMACpp(modelFile, llamaOpts...); err == nil {
		// updateModels(model)
		return model, nil
	} else {
		fmt.Printf(`could not load llama model: `, err.Error())
	}

	if model, err := m.LoadOpenAI(modelFile); err == nil {
		// updateModels(model)
		return model, nil
	} else {
		fmt.Printf(`could not load openai model: `, err.Error())
	}

	return nil, fmt.Errorf("no avail  model - all backends returned")
}
