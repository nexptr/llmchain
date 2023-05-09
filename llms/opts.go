package llms

// CallOption is a function that configures a LLM.
type ModelOption func(*ModelOptions)

// BaseOptions is a set base of options for LLM.Call.
type ModelOptions = Payload

// WithMaxTokens is an option for LLM.Call.
// func WithMaxTokens(maxTokens int) ModelOption {
// 	return func(o *ModelOptions) {
// 		o.MaxTokens = maxTokens
// 	}
// }

// // WithTemperature is an option for LLM.Call.
// func WithTemperature(temperature float64) ModelOption {
// 	return func(o *ModelOptions) {
// 		o.Temperature = temperature
// 	}
// }

// // WithStopWords is an option for LLM.Call.
// func WithStopWords(stopWords []string) ModelOption {
// 	return func(o *ModelOptions) {
// 		o.StopWords = stopWords
// 	}
// }

// // WithOptions is an option for LLM.Call.
// func WithOptions(options ModelOptions) ModelOption {
// 	return func(o *ModelOptions) {
// 		(*o) = options
// 	}
// }
