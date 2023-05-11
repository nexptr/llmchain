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
	sync.Mutex
}

var _ llms.LLM = &LLaMACpp{}

func New(modelPath string, opts ...llms.ModelOption) (*LLaMACpp, error) {

	// Check if we already have a loaded model
	if !utils.PathExists(modelPath) {
		return nil, fmt.Errorf("model does not exist: %s", modelPath)
	}

	mOpts := defaultLLamaModelOptions(opts...)

	mPath := C.CString(modelPath)

	result := C.load_model(mPath, C.int(mOpts.ContextSize), C.int(mOpts.Parts), C.int(mOpts.Seed), C.bool(mOpts.F16), C.bool(mOpts.MLock), C.bool(mOpts.Embeddings))
	if result == nil {
		return nil, fmt.Errorf("failed loading model")
	}

	ll := &LLaMACpp{state: result, embeddings: mOpts.Embeddings, options: mOpts}

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

// SupportStream  LLaMACpp support stream
func (l *LLaMACpp) Name() string {
	return l.options.Name
}
func (l *LLaMACpp) InferenceFn(input string, predict *llms.ModelOptions) func() (string, error) {

	return func() (string, error) {

		str, er := l.PredictWithOpts(input, predict)
		// Seems that if we don't free the callback explicitly we leave functions registered (that might try to send on closed channels)
		// For instance otherwise the API returns: {"error":{"code":500,"message":"send on closed channel","type":""}}
		// after a stream event has occurred
		return str, er
	}

}

// Embeddings
func (l *LLaMACpp) Embeddings(text string, opts ...llms.ModelOption) ([]float32, error) {
	if !l.embeddings {
		return []float32{}, fmt.Errorf("model loaded without embeddings")
	}

	//copy from base
	po := l.options

	for _, f := range opts {
		f(&po)
	}

	input := C.CString(text)
	if po.Maxtokens == 0 {
		po.Maxtokens = 99999999
	}
	floats := make([]float32, po.Maxtokens)
	reverseCount := len(po.StopWords)
	reversePrompt := make([]*C.char, reverseCount)
	var pass **C.char
	for i, s := range po.StopWords {
		cs := C.CString(s)
		reversePrompt[i] = cs
		pass = &reversePrompt[0]
	}

	params := C.llama_allocate_params(input, C.int(po.Seed), C.int(po.Threads), C.int(po.Maxtokens), C.int(po.TopK),
		C.float(po.TopP), C.float(po.Temperature), C.float(po.RepeatPenalty), C.int(po.Repeat),
		C.bool(po.IgnoreEOS), C.bool(po.F16),
		C.int(po.Batch), C.int(po.Keep), pass, C.int(reverseCount),
		C.float(tailFreeSamplingZ), C.float(typicalP), C.float(frequencyPenalty), C.float(presencePenalty),
		C.int(po.Mirostat), C.float(po.MirostatETA), C.float(po.MirostatTAU), C.bool(penalizeNL), C.CString(logitBias),
	)

	ret := C.get_embeddings(params, l.state, (*C.float)(&floats[0]))
	if ret != 0 {
		return floats, fmt.Errorf("embedding inference failed")
	}

	return floats, nil
}

func (l *LLaMACpp) MergeModelOptions(req *llms.OpenAIRequest) *llms.ModelOptions {
	//copy options from base config
	m := l.options

	return m.Override(req)

}

func (l *LLaMACpp) Predict(text string, opts ...llms.ModelOption) (string, error) {

	//copy from base
	op := l.options

	for _, f := range opts {
		f(&op)
	}

	return l.PredictWithOpts(text, &op)

}

func (l *LLaMACpp) PredictWithOpts(text string, opts *llms.ModelOptions) (string, error) {

	// This is still needed, see: https://github.com/ggerganov/llama.cpp/discussions/784
	l.Lock()
	defer l.Unlock()

	if opts.TokenCallback != nil {
		println(`update TokenCallback`)
		setCallback(l.state, opts.TokenCallback)
	}

	input := C.CString(text)
	if opts.Maxtokens == 0 {
		opts.Maxtokens = 99999999
	}

	out := make([]byte, opts.Maxtokens)

	reverseCount := len(opts.StopWords)
	reversePrompt := make([]*C.char, reverseCount)
	var pass **C.char
	for i, s := range opts.StopWords {
		cs := C.CString(s)
		reversePrompt[i] = cs
		pass = &reversePrompt[0]
	}

	params := C.llama_allocate_params(input, C.int(opts.Seed), C.int(opts.Threads), C.int(opts.Maxtokens), C.int(opts.TopK),
		C.float(opts.TopP), C.float(opts.Temperature), C.float(opts.RepeatPenalty), C.int(opts.Repeat),
		C.bool(opts.IgnoreEOS), C.bool(opts.F16),
		C.int(opts.Batch), C.int(opts.Keep), pass, C.int(reverseCount),
		C.float(tailFreeSamplingZ), C.float(typicalP), C.float(frequencyPenalty), C.float(presencePenalty),
		C.int(opts.Mirostat), C.float(opts.MirostatETA), C.float(opts.MirostatTAU), C.bool(penalizeNL), C.CString(logitBias),
	)
	ret := C.llama_predict(params, l.state, (*C.char)(unsafe.Pointer(&out[0])), C.bool(opts.Debug))
	if ret != 0 {
		return "", fmt.Errorf("inference failed")
	}
	res := C.GoString((*C.char)(unsafe.Pointer(&out[0])))

	res = strings.TrimPrefix(res, " ")
	res = strings.TrimPrefix(res, text)
	res = strings.TrimPrefix(res, "\n")

	for _, s := range opts.StopWords {
		res = strings.TrimRight(res, s)
	}

	C.llama_free_params(params)

	if opts.TokenCallback != nil {
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

	ret, err := l.Predict(prompt, llms.WithTokenCallback(func(token string) bool {
		fmt.Print(token)
		return true
	}), llms.WithMaxToken(128), llms.WithThreads(4), llms.WithTopK(90), llms.WithTopP(0.86), llms.WithStopWords("llama"))

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
