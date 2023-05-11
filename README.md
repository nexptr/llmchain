# LLMChain

Go language implementation of LangChain. This project was inspired by <https://github.com/hwchase17/langchain>

//TODO 

**The project has just started, and a working version is expected to be released in June.**


## quick start 


```
git clone https://github.com/exppii/llmchain.git


cd llmchain

# copy your models to models/
cp your-model.bin models/


make app


# start with go run
LIBRARY_PATH=./llms/llamacpp C_INCLUDE_PATH=./llms/llamacpp go run ./examples/app -conf your_conf.yaml


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