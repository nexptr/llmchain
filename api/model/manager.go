package model

import (
	"fmt"

	"github.com/exppii/llmchain/api/conf"
	"github.com/exppii/llmchain/api/log"
	"github.com/exppii/llmchain/llms"
	"github.com/exppii/llmchain/llms/llamacpp"
	"github.com/exppii/llmchain/llms/openai"
	"github.com/exppii/llmchain/prompts"
)

// Manager LLM model manager
type Manager struct {
	cf *conf.Config

	//ModelPath root mode path for llms
	ModelPath string

	loadedModels map[string]llms.LLM

	// promptsTemplates map[string]*template.Template

	promptsTemplates *prompts.Template
}

func NewModelManager(cf *conf.Config) *Manager {

	return &Manager{
		cf:           cf,
		loadedModels: map[string]llms.LLM{},

		promptsTemplates: prompts.NewTemplate(cf.PromptPath),
	}
}

// Load using the configs in config file, load llm models to memory.
func (m *Manager) Load() error {

	for _, o := range m.cf.ModelOptions {

		//load
		modeType := llms.GetModelType(o.Name)

		switch modeType {
		case llms.ModelOpenAI:

		case llms.ModelLLaMACPP:

			llm, err := m.LoadLLaMACpp(o)
			if err != nil {

				return err
			}
			m.loadedModels[o.Name] = llm
			log.I(`loaded llama.cpp from `, o.ModelPath)
		default:

			return fmt.Errorf(`model: %s not support`, o.Name)
		}

	}

	return nil

}

// Free clean all model in memory
func (m *Manager) Free() {
	//TODO free all model
}

func (m *Manager) GetModel(modelName string) (llms.LLM, error) {

	model, exists := m.loadedModels[modelName]
	if !exists {
		return nil, fmt.Errorf("model %s not found", modelName)
	}

	return model, nil

}

func (m *Manager) GetPrompt() *prompts.Template {

	return m.promptsTemplates

}

func (m *Manager) ListModels() []string {
	ret := make([]string, len(m.loadedModels))

	i := 0
	for k := range m.loadedModels {
		ret[i] = k
		i++
	}
	return ret
}

func (m *Manager) LoadLLaMACpp(opts llms.ModelOptions) (*llamacpp.LLaMACpp, error) {
	log.I(`loading model: `, opts.Model, ` with path: `, opts.ModelPath, `...`)
	return llamacpp.New(opts.ModelPath, opts.BuildOpts()...)

}

func (m *Manager) LoadOpenAI(modelName string, opts ...openai.ModelOption) (*openai.OpenAI, error) {

	return nil, fmt.Errorf("openai llm todo")
}

// func (m *Manager) GreedyLoad(modelFile string, llamaOpts []llamacpp.ModelOption, threads uint32) (llms.LLM, error) {

// 	model, exists := m.loadedModels[modelFile]
// 	if exists {
// 		// muModels.Unlock()
// 		return model, nil
// 	}

// 	//try

// 	if model, err := m.LoadLLaMACpp(modelFile, llamaOpts...); err == nil {
// 		// updateModels(model)
// 		return model, nil
// 	} else {
// 		fmt.Printf(`could not load llama model: `, err.Error())
// 	}

// 	if _, err := m.LoadOpenAI(modelFile); err == nil {
// 		// updateModels(model)
// 		//TODO
// 		// return model, nil
// 		return nil, fmt.Errorf("no avail  model - all backends returned")
// 	} else {
// 		fmt.Printf(`could not load openai model: `, err.Error())
// 	}

// 	return nil, fmt.Errorf("no avail  model - all backends returned")
// }
