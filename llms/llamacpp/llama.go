package llamacpp

// #cgo CXXFLAGS: -I${SRCDIR}/llama.cpp/examples -I${SRCDIR}/llama.cpp
// #cgo LDFLAGS: -L${SRCDIR}/ -lbinding -lm -lstdc++
// #cgo darwin LDFLAGS: -framework Accelerate
// #cgo darwin CXXFLAGS: -std=c++11
// #include "binding.h"
import "C"

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"unsafe"

	"github.com/exppii/llmchain/llms"
	"github.com/exppii/llmchain/utils"
)

type LLaMACpp struct {
	options llms.ModelOptions

	// Streaming bool
	// ModelPath string
	// Threads   int
	embeddings bool
	state      unsafe.Pointer

	// This is still needed, see: https://github.com/ggerganov/llama.cpp/discussions/784
	mu sync.Mutex
}

var _ llms.LLM = &LLaMACpp{}

func New(modelPath string, opts ...ModelOption) (*LLaMACpp, error) {

	// Check if we already have a loaded model
	if !utils.PathExists(modelPath) {
		return nil, fmt.Errorf("model does not exist")
	}

	mOpts := NewModelOptions(opts...)

	mPath := C.CString(modelPath)

	result := C.load_model(mPath, C.int(mOpts.ContextSize), C.int(mOpts.Parts), C.int(mOpts.Seed), C.bool(mOpts.F16Memory), C.bool(mOpts.MLock), C.bool(mOpts.Embeddings))
	if result == nil {
		return nil, fmt.Errorf("failed loading model")
	}

	ll := &LLaMACpp{state: result, embeddings: mOpts.Embeddings}

	return ll, nil
}

// Free model
func (l *LLaMACpp) Free() {
	C.llama_free_model(l.state)
}

// SupportStream  LLaMACpp support stream
func (l *LLaMACpp) SupportStream() bool {
	return true
}

func (l *LLaMACpp) InferenceFn(input string, payload *llms.Payload, tokenCallback func(string) bool) func() (string, error) {

	return func() (string, error) {

		if tokenCallback != nil {
			l.SetTokenCallback(tokenCallback)
		}

		predictOptions := l.BuildPredictOptions(payload)

		str, er := l.Predict(
			input,
			predictOptions...,
		)
		// Seems that if we don't free the callback explicitly we leave functions registered (that might try to send on closed channels)
		// For instance otherwise the API returns: {"error":{"code":500,"message":"send on closed channel","type":""}}
		// after a stream event has occurred
		l.SetTokenCallback(nil)
		return str, er
	}

}

func (l *LLaMACpp) BuildPredictOptions(c *llms.Payload) []PredictOption {
	// Generate the prediction using the language model
	predictOptions := []PredictOption{
		WithTemperature(c.Temperature),
		WithTopP(c.TopP),
		WithTopK(c.TopK),
		WithTokens(c.Maxtokens),
		WithThreads(c.Threads),
	}

	if c.Mirostat != 0 {
		predictOptions = append(predictOptions, WithMirostat(c.Mirostat))
	}

	if c.MirostatETA != 0 {
		predictOptions = append(predictOptions, WithMirostatETA(c.MirostatETA))
	}

	if c.MirostatTAU != 0 {
		predictOptions = append(predictOptions, WithMirostatTAU(c.MirostatTAU))
	}

	if c.Debug {
		predictOptions = append(predictOptions, Debug)
	}

	predictOptions = append(predictOptions, WithStopWords(c.StopWords...))

	if c.RepeatPenalty != 0 {
		predictOptions = append(predictOptions, WithPenalty(c.RepeatPenalty))
	}

	if c.Keep != 0 {
		predictOptions = append(predictOptions, WithNKeep(c.Keep))
	}

	if c.Batch != 0 {
		predictOptions = append(predictOptions, WithBatch(c.Batch))
	}

	if c.F16 {
		predictOptions = append(predictOptions, EnableF16KV)
	}

	if c.IgnoreEOS {
		predictOptions = append(predictOptions, IgnoreEOS)
	}

	if c.Seed != 0 {
		predictOptions = append(predictOptions, WithSeed(c.Seed))
	}

	return predictOptions
}

// Embeddings
func (l *LLaMACpp) Embeddings(text string, opts ...PredictOption) ([]float32, error) {
	if !l.embeddings {
		return []float32{}, fmt.Errorf("model loaded without embeddings")
	}

	po := NewPredictOptions(opts...)

	input := C.CString(text)
	if po.Tokens == 0 {
		po.Tokens = 99999999
	}
	floats := make([]float32, po.Tokens)
	reverseCount := len(po.StopPrompts)
	reversePrompt := make([]*C.char, reverseCount)
	var pass **C.char
	for i, s := range po.StopPrompts {
		cs := C.CString(s)
		reversePrompt[i] = cs
		pass = &reversePrompt[0]
	}

	params := C.llama_allocate_params(input, C.int(po.Seed), C.int(po.Threads), C.int(po.Tokens), C.int(po.TopK),
		C.float(po.TopP), C.float(po.Temperature), C.float(po.Penalty), C.int(po.Repeat),
		C.bool(po.IgnoreEOS), C.bool(po.F16KV),
		C.int(po.Batch), C.int(po.NKeep), pass, C.int(reverseCount),
		C.float(po.TailFreeSamplingZ), C.float(po.TypicalP), C.float(po.FrequencyPenalty), C.float(po.PresencePenalty),
		C.int(po.Mirostat), C.float(po.MirostatETA), C.float(po.MirostatTAU), C.bool(po.PenalizeNL), C.CString(po.LogitBias),
	)

	ret := C.get_embeddings(params, l.state, (*C.float)(&floats[0]))
	if ret != 0 {
		return floats, fmt.Errorf("embedding inference failed")
	}

	return floats, nil
}

func (l *LLaMACpp) MergePayload(req *llms.OpenAIRequest) *llms.Payload {

	// payload := llms.NewPayload(opt)
	//TODO copy
	m := l.options

	return m.Override(req)

}

func (l *LLaMACpp) Predict(text string, opts ...PredictOption) (string, error) {
	po := NewPredictOptions(opts...)

	if po.TokenCallback != nil {
		setCallback(l.state, po.TokenCallback)
	}

	input := C.CString(text)
	if po.Tokens == 0 {
		po.Tokens = 99999999
	}
	out := make([]byte, po.Tokens)

	reverseCount := len(po.StopPrompts)
	reversePrompt := make([]*C.char, reverseCount)
	var pass **C.char
	for i, s := range po.StopPrompts {
		cs := C.CString(s)
		reversePrompt[i] = cs
		pass = &reversePrompt[0]
	}

	params := C.llama_allocate_params(input, C.int(po.Seed), C.int(po.Threads), C.int(po.Tokens), C.int(po.TopK),
		C.float(po.TopP), C.float(po.Temperature), C.float(po.Penalty), C.int(po.Repeat),
		C.bool(po.IgnoreEOS), C.bool(po.F16KV),
		C.int(po.Batch), C.int(po.NKeep), pass, C.int(reverseCount),
		C.float(po.TailFreeSamplingZ), C.float(po.TypicalP), C.float(po.FrequencyPenalty), C.float(po.PresencePenalty),
		C.int(po.Mirostat), C.float(po.MirostatETA), C.float(po.MirostatTAU), C.bool(po.PenalizeNL), C.CString(po.LogitBias),
	)
	ret := C.llama_predict(params, l.state, (*C.char)(unsafe.Pointer(&out[0])), C.bool(po.DebugMode))
	if ret != 0 {
		return "", fmt.Errorf("inference failed")
	}
	res := C.GoString((*C.char)(unsafe.Pointer(&out[0])))

	res = strings.TrimPrefix(res, " ")
	res = strings.TrimPrefix(res, text)
	res = strings.TrimPrefix(res, "\n")

	for _, s := range po.StopPrompts {
		res = strings.TrimRight(res, s)
	}

	C.llama_free_params(params)

	if po.TokenCallback != nil {
		setCallback(l.state, nil)
	}

	return res, nil
}

// CGo only allows us to use static calls from C to Go, we can't just dynamically pass in func's.
// This is the next best thing, we register the callbacks in this map and call tokenCallback from
// the C code. We also attach a finalizer to LLama, so it will unregister the callback when the
// garbage collection frees it.

// SetTokenCallback registers a callback for the individual tokens created when running Predict. It
// will be called once for each token. The callback shall return true as long as the model should
// continue predicting the next token. When the callback returns false the predictor will return.
// The tokens are just converted into Go strings, they are not trimmed or otherwise changed. Also
// the tokens may not be valid UTF-8.
// Pass in nil to remove a callback.
//
// It is save to call this method while a prediction is running.
func (l *LLaMACpp) SetTokenCallback(callback func(token string) bool) {
	setCallback(l.state, callback)
}

var (
	m         sync.Mutex
	callbacks = map[uintptr]func(string) bool{}
)

//export tokenCallback
func tokenCallback(statePtr unsafe.Pointer, token *C.char) bool {
	m.Lock()
	defer m.Unlock()

	if callback, ok := callbacks[uintptr(statePtr)]; ok {
		return callback(C.GoString(token))
	}

	return true
}

// setCallback can be used to register a token callback for LLama. Pass in a nil callback to
// remove the callback.
func setCallback(statePtr unsafe.Pointer, callback func(string) bool) {
	m.Lock()
	defer m.Unlock()

	if callback == nil {
		delete(callbacks, uintptr(statePtr))
	} else {
		callbacks[uintptr(statePtr)] = callback
	}
}

func (l *LLaMACpp) Call(ctx context.Context, prompt string) (string, error) {

	ret, err := l.Predict(prompt, Debug, WithTokenCallback(func(token string) bool {
		fmt.Print(token)
		return true
	}), WithTokens(128), WithThreads(4), WithTopK(90), WithTopP(0.86), WithStopWords("llama"))

	if err != nil {
		panic(err)
	}
	embeds, err := l.Embeddings(prompt)
	if err != nil {
		fmt.Printf("Embeddings: error %s \n", err.Error())
	}
	fmt.Printf("Embeddings: %v", embeds)
	fmt.Printf("\n\n")

	return ret, err
}
