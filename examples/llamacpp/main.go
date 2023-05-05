package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/exppii/llmchain/llms/llamacpp"
)

var (
	threads   = 4
	tokens    = 128
	modelPath string
)

// LIBRARY_PATH=./llms/llamacpp C_INCLUDE_PATH=./llms/llamacpp go run ./examples/llamacpp -m "./models/7B/ggml-model-q4_1.bin" -t 4
func main() {

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.StringVar(&modelPath, "m", "./models/7B/ggml-model-q4_0.bin", "path to q4_0.bin model file to load")
	flags.IntVar(&threads, "t", runtime.NumCPU(), "number of threads to use during computation")
	flags.IntVar(&tokens, "n", 512, "number of tokens to predict")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("Parsing program arguments failed: %s", err)
		os.Exit(1)
	}

	llm, err := llamacpp.New(modelPath, llamacpp.WithContext(128), llamacpp.WithParts(-1), llamacpp.EnableEmbeddings)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.TODO()

	completion, err := llm.Call(ctx, "How many sides does a square have?")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(`completion: `, completion)

}
