package api

import (
	"github.com/exppii/llmchain/llms"
)

func ComputeChoices(llm llms.LLM, predInput string, payload *llms.Payload, cb func(string, *[]Choice), tokenCallback func(string) bool) ([]Choice, error) {
	result := []Choice{}

	n := payload.N

	if payload.N == 0 {
		n = 1
	}

	// get the model function to call for the result
	predFunc := LLMInference(llm, predInput, payload, tokenCallback)

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

func LLMInference(llm llms.LLM, predInput string, payload *llms.Payload, tokenCallback func(string) bool) func() (string, error) {

	// get the model function to call for the result
	fn := llm.InferenceFn(predInput, payload, tokenCallback)

	return func() (string, error) {
		// This is still needed, see: https://github.com/ggerganov/llama.cpp/discussions/784
		// mutexMap.Lock()
		// l, ok := mutexes[modelFile]
		// if !ok {
		// 	m := &sync.Mutex{}
		// 	mutexes[modelFile] = m
		// 	l = m
		// }
		// mutexMap.Unlock()
		// l.Lock()
		// defer l.Unlock()
		//TODO multithread lock

		res, err := fn()
		if tokenCallback != nil && !llm.SupportStream() {
			tokenCallback(res)
		}
		return res, err
	}

}
