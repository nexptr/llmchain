package llamacpp

type ModelOption func(*ModelOptions)

type ModelOptions struct {
	ContextSize int
	Parts       int
	Seed        int
	F16Memory   bool
	MLock       bool
	Embeddings  bool
}

// Create a new PredictOptions object with the given options.
func NewModelOptions(opts ...ModelOption) ModelOptions {
	p := defaultModelOptions
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithContext sets the context size.
func WithContext(c int) ModelOption {
	return func(p *ModelOptions) {
		p.ContextSize = c
	}
}

func WithModelSeed(c int) ModelOption {
	return func(p *ModelOptions) {
		p.Seed = c
	}
}

func WithParts(c int) ModelOption {
	return func(p *ModelOptions) {
		p.Parts = c
	}
}

var defaultModelOptions ModelOptions = ModelOptions{
	ContextSize: 512,
	Seed:        0,
	F16Memory:   false,
	MLock:       false,
	Embeddings:  false,
}

var EnableEmbeddings ModelOption = func(p *ModelOptions) { p.Embeddings = true }

var EnableF16Memory ModelOption = func(p *ModelOptions) { p.F16Memory = true }

var EnableF16KV PredictOption = func(p *PredictOptions) { p.F16KV = true }

var EnableMLock ModelOption = func(p *ModelOptions) { p.MLock = true }

type PredictOption func(p *PredictOptions)

type PredictOptions struct {
	Seed, Threads, Tokens, TopK, Repeat, Batch, NKeep int
	TopP, Temperature, Penalty                        float64
	F16KV                                             bool
	DebugMode                                         bool
	StopPrompts                                       []string
	IgnoreEOS                                         bool

	TailFreeSamplingZ float64
	TypicalP          float64
	FrequencyPenalty  float64
	PresencePenalty   float64
	Mirostat          int
	MirostatETA       float64
	MirostatTAU       float64
	PenalizeNL        bool
	LogitBias         string
	TokenCallback     func(string) bool
}

// Create a new PredictOptions object with the given options.
func NewPredictOptions(opts ...PredictOption) PredictOptions {
	p := defaultPredictOptions
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

var defaultPredictOptions PredictOptions = PredictOptions{
	Seed:              -1,
	Threads:           4,
	Tokens:            128,
	Penalty:           1.1,
	Repeat:            64,
	Batch:             8,
	NKeep:             64,
	TopK:              40,
	TopP:              0.95,
	TailFreeSamplingZ: 1.0,
	TypicalP:          1.0,
	Temperature:       0.8,
	FrequencyPenalty:  0.0,
	PresencePenalty:   0.0,
	Mirostat:          0,
	MirostatTAU:       5.0,
	MirostatETA:       0.1,
}

var Debug PredictOption = func(p *PredictOptions) {
	p.DebugMode = true
}

var IgnoreEOS PredictOption = func(p *PredictOptions) {
	p.IgnoreEOS = true
}

// WithTokenCallback sets the prompts that will stop predictions.
func WithTokenCallback(fn func(string) bool) PredictOption {
	return func(p *PredictOptions) {
		p.TokenCallback = fn
	}
}

// WithStopWords sets the prompts that will stop predictions.
func WithStopWords(stop ...string) PredictOption {
	return func(p *PredictOptions) {
		p.StopPrompts = stop
	}
}

// WithSeed sets the random seed for sampling text generation.
func WithSeed(seed int) PredictOption {
	return func(p *PredictOptions) {
		p.Seed = seed
	}
}

// WithThreads sets the number of threads to use for text generation.
func WithThreads(threads int) PredictOption {
	return func(p *PredictOptions) {
		p.Threads = threads
	}
}

// WithTokens sets the number of tokens to generate.
func WithTokens(tokens int) PredictOption {
	return func(p *PredictOptions) {
		p.Tokens = tokens
	}
}

// WithTopK sets the value for top-K sampling.
func WithTopK(topk int) PredictOption {
	return func(p *PredictOptions) {
		p.TopK = topk
	}
}

// WithTopP sets the value for nucleus sampling.
func WithTopP(topp float64) PredictOption {
	return func(p *PredictOptions) {
		p.TopP = topp
	}
}

// WithTemperature sets the temperature value for text generation.
func WithTemperature(temp float64) PredictOption {
	return func(p *PredictOptions) {
		p.Temperature = temp
	}
}

// WithPenalty sets the repetition penalty for text generation.
func WithPenalty(penalty float64) PredictOption {
	return func(p *PredictOptions) {
		p.Penalty = penalty
	}
}

// WithRepeat sets the number of times to repeat text generation.
func WithRepeat(repeat int) PredictOption {
	return func(p *PredictOptions) {
		p.Repeat = repeat
	}
}

// WithBatch sets the batch size.
func WithBatch(size int) PredictOption {
	return func(p *PredictOptions) {
		p.Batch = size
	}
}

// WithKeep sets the number of tokens from initial prompt to keep.
func WithNKeep(n int) PredictOption {
	return func(p *PredictOptions) {
		p.NKeep = n
	}
}

// WithTailFreeSamplingZ sets the tail free sampling, parameter z.
func WithTailFreeSamplingZ(tfz float64) PredictOption {
	return func(p *PredictOptions) {
		p.TailFreeSamplingZ = tfz
	}
}

// WithTypicalP sets the typicality parameter, p_typical.
func WithTypicalP(tp float64) PredictOption {
	return func(p *PredictOptions) {
		p.TypicalP = tp
	}
}

// WithFrequencyPenalty sets the frequency penalty parameter, freq_penalty.
func WithFrequencyPenalty(fp float64) PredictOption {
	return func(p *PredictOptions) {
		p.FrequencyPenalty = fp
	}
}

// WithPresencePenalty sets the presence penalty parameter, presence_penalty.
func WithPresencePenalty(pp float64) PredictOption {
	return func(p *PredictOptions) {
		p.PresencePenalty = pp
	}
}

// WithMirostat sets the mirostat parameter.
func WithMirostat(m int) PredictOption {
	return func(p *PredictOptions) {
		p.Mirostat = m
	}
}

// WithMirostatETA sets the mirostat ETA parameter.
func WithMirostatETA(me float64) PredictOption {
	return func(p *PredictOptions) {
		p.MirostatETA = me
	}
}

// WithMirostatTAU sets the mirostat TAU parameter.
func WithMirostatTAU(mt float64) PredictOption {
	return func(p *PredictOptions) {
		p.MirostatTAU = mt
	}
}

// WithPenalizeNL sets whether to penalize newlines or not.
func WithPenalizeNL(pnl bool) PredictOption {
	return func(p *PredictOptions) {
		p.PenalizeNL = pnl
	}
}

// WithLogitBias sets the logit bias parameter.
func WithLogitBias(lb string) PredictOption {
	return func(p *PredictOptions) {
		p.LogitBias = lb
	}
}
