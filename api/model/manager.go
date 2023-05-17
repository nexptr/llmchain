package model

import (
	"fmt"

	"github.com/exppii/llmchain"
	"github.com/exppii/llmchain/api/conf"
	"github.com/exppii/llmchain/api/log"
	"github.com/exppii/llmchain/chain"
	"github.com/exppii/llmchain/llm"
	"github.com/exppii/llmchain/llm/openai"

	"github.com/exppii/llmchain/prompts"
)

// Manager LLM model manager
type Manager struct {
	cf *conf.Config

	//ModelPath root mode path for llms
	ModelPath string

	loadedModels map[string]llmchain.LLM

	// promptsTemplates map[string]*template.Template

	promptsTemplates *prompts.Render
}

func NewModelManager(cf *conf.Config) *Manager {

	return &Manager{
		cf:           cf,
		loadedModels: map[string]llmchain.LLM{},

		// promptsTemplates: prompts.NewRender(cf.PromptPath),
	}
}

// Load using the configs in config file, load llm models to memory.
func (m *Manager) Load() error {

	for _, o := range m.cf.ModelOptions {

		//load
		modeType := llm.GetModelType(o.Name)

		switch modeType {
		case llm.ModelOpenAI:
			llm, err := openai.FromYaml(o)
			if err != nil {
				return err
			}
			m.loadedModels[o.Name] = llm
			log.I(`loaded llama.cpp from `, o.Name)

		case llm.ModelLLaMACPP:

			// llm, err := m.LoadLLaMACpp(o)
			// if err != nil {

			// 	return err
			// }
			// m.loadedModels[o.Name] = llm
		default:

			return fmt.Errorf(`model: %s not support`, o.Name)
		}

	}

	//Reg all chain

	llmchain.RegChain(chain.NewBaseChatChain())

	return nil

}

// Free clean all model in memory
func (m *Manager) Free() {

	for _, llm := range m.loadedModels {
		log.I(`free model: `, llm.Name())
		llm.Free()
	}

}

func (m *Manager) GetModel(modelName string) (llmchain.LLM, error) {
	panic(`todo`)
	// model, exists := m.loadedModels[modelName]
	// if !exists {
	// 	return nil, fmt.Errorf("model %s not found", modelName)
	// }

	// return model, nil

}

func (m *Manager) LLMChain(modelName string, chainName string) (llmchain.LLM, error) {

	model, exists := m.loadedModels[modelName]
	if !exists {
		return nil, fmt.Errorf("model %s not found", modelName)
	}

	if c, ok := llmchain.GetChain(chainName); ok {
		log.D(`using chain: `, chainName)
		c.WithLLM(model)
		return llmchain.New(model, c), nil
	}

	return model, nil

}

func (m *Manager) GetPrompt() *prompts.Render {

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

// func (m *Manager) LoadLLaMACpp(opts llm.ModelOptions) (*llamacpp.LLaMACpp, error) {
// 	log.I(`loading model: `, opts.Model, ` with path: `, opts.ModelPath, `...`)
// 	return llamacpp.New(opts.ModelPath, opts.BuildOpts()...)

// }

func (m *Manager) LoadOpenAI(modelName string, opts ...openai.ModelOption) (*openai.OpenAI, error) {

	return nil, fmt.Errorf("openai llm todo")
}
