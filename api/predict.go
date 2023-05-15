package api

import (
	"github.com/exppii/llmchain/api/log"

	"github.com/exppii/llmchain/llms"
)

func ComputeChoices(llm llms.LLM, predInput string, payload *llms.ModelOptions, cb func(string, *[]Choice)) ([]Choice, error) {
	result := []Choice{}

	n := payload.N

	if payload.N == 0 {
		n = 1
	}

	// get the model function to call for the result
	log.D(`LLMInference -----`)
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

func LLMEmbedding(llm llms.LLM, s string, tokens []int, payload *llms.ModelOptions) (func() ([]float32, error), error) {

	// fn := llm.Embeddings

	return func() ([]float32, error) {
		// This is still needed, see: https://github.com/ggerganov/llama.cpp/discussions/784
		// Embeddings(input string, tokens []int, payload *ModelOptions) func() ([]float32, error)
		embeds, err := llm.Embeddings(s, tokens, payload)
		if err != nil {
			return embeds, err
		}
		// Remove trailing 0s
		for i := len(embeds) - 1; i >= 0; i-- {
			if embeds[i] == 0.0 {
				embeds = embeds[:i]
			} else {
				break
			}
		}
		return embeds, nil
	}, nil
}
