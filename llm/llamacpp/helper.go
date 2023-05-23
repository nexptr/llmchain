package llamacpp

// #cgo CXXFLAGS: -I${SRCDIR}/llama.cpp/examples -I${SRCDIR}/llama.cpp
// #cgo LDFLAGS: -L${SRCDIR}/ -lbinding -lm -lstdc++
// #cgo darwin LDFLAGS: -framework Accelerate
// #cgo darwin CXXFLAGS: -std=c++11
// #include "binding.h"
import "C"
import (
	"fmt"
	"strings"
	"sync"
	"unsafe"
)

var (
	m         sync.Mutex
	callbacks = map[uintptr]func(string) bool{}
)

func (l *LLaMACpp) load() error {

	l.tmpl = loadModelPrompts(l.Model)

	mPath := C.CString(l.ModelPath)

	l.state = C.load_model(mPath, C.int(l.ContextSize), C.int(l.Seed), C.bool(l.F16), C.bool(l.MLock), C.bool(l.EnableEmbedding), C.int(l.NGPULayers))
	if l.state == nil {
		return fmt.Errorf("failed loading model")
	}

	return nil

}

// Free model
func (l *LLaMACpp) freeModel() {
	C.llama_free_model(l.state)
}

func (l *LLaMACpp) predictWithOpts(text string, opts *PredictOptions) (string, error) {

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
		C.float(opts.TopP), C.float(opts.Temperature), C.float(opts.Penalty), C.int(opts.Repeat),
		C.bool(opts.IgnoreEOS), C.bool(opts.F16),
		C.int(opts.Batch), C.int(opts.NKeep), pass, C.int(reverseCount),
		C.float(opts.TailFreeSamplingZ), C.float(opts.TypicalP), C.float(opts.FrequencyPenalty), C.float(opts.PresencePenalty),
		C.int(opts.Mirostat), C.float(opts.MirostatETA), C.float(opts.MirostatTAU), C.bool(opts.PenalizeNL), C.CString(opts.LogitBias),
	)
	ret := C.llama_predict(params, l.state, (*C.char)(unsafe.Pointer(&out[0])), C.bool(opts.DebugMode))
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
