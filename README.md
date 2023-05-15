# LLMChain

LLMChain is Golang(using cgo for [llama.cpp](https://github.com/ggerganov/llama.cpp) support ) package implementation of [LangChain](https://github.com/hwchase17/langchain).

- OpenAI compatible API
- Supports multiple models
- Support for langchain

**This package is in active mode of building and there are many changes ahead. Hence you can use it with your complete own risk. The package will be considered as stable when version 1.0 is released.**

This project was inspired by heavily inspired by and based on the popular Python [LangChain](https://github.com/hwchase17/langchain). It's also influenced by [LocalAI](https://github.com/go-skynet/LocalAI).


## Basic example


```bash
git clone https://github.com/exppii/llmchain.git
cd llmchain
# make llama.cpp library.
make app

# copy your models to models/
cp your-model.bin models/

#edit config.yaml for you model 
vi ./examples/app/conf.yaml

# start with go run
LIBRARY_PATH=./llms/llamacpp C_INCLUDE_PATH=./llms/llamacpp go run ./examples/app -conf ./examples/app/conf.yaml

# Now API is accessible at localhost:8080
curl http://localhost:8080/v1/models
# {"object":"list","data":[{"id":"ggml-llama-7b","object":"model"}]}


curl http://localhost:8080/v1/completions -H "Content-Type: application/json" -d '{
     "model": "ggml-llama-7b",            
     "prompt": "how many days a week",
     "temperature": 0.7,
     "max_tokens": 256
   }'

```

## Getting Started

## Model compatibility

It is compatible with the models supported by llama.cpp

Tested with:

Vicuna
Alpaca


## Short-term roadmap

- [ ] Binary releases
- [ ] docker releases
- [ ] Have a webUI!

## Acknowledgements

//TODO