package api

import (
	"github.com/exppii/llmchain/llms"
)

func ComputeChoices(llm llms.LLM, predInput string, payload *llms.ModelOptions, cb func(string, *[]Choice)) ([]Choice, error) {
	result := []Choice{}

	n := payload.N

	if payload.N == 0 {
		n = 1
	}

	// get the model function to call for the result
	predFunc := LLMInference(llm, predInput, payload)

	for i := 0; i < n; i++ {
		prediction, err := predFunc()
		if err != nil {
			return result, err
		}

		// prediction = Finetune(*config, predInput, prediction)
		cb(prediction, &result)
		//result = append(result, Choice{Text: prediction})

	}
	return result, nil
}

func LLMInference(llm llms.LLM, predInput string, payload *llms.ModelOptions) func() (string, error) {

	// get the model function to call for the result
	fn := llm.InferenceFn(predInput, payload)

	return func() (string, error) {

		res, err := fn()
		if payload.TokenCallback != nil && !llm.SupportStream() {
			payload.TokenCallback(res)
		}
		return res, err
	}

}
