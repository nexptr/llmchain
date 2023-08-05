# Makefile.

PROJECT_PATH=$(shell cd "$(dirname "$0" )" &&pwd)
PROJECT_NAME=llmchain
VERSION=$(shell git describe --tags | sed 's/\(.*\)-.*/\1/')
BUILD_DATE=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
BUILD_HASH=$(shell git rev-parse HEAD)
LDFLAGS="-X github.com/nexptr/llmchain.BuildStamp=${BUILD_DATE} -X github.com/nexptr/llmchain.GitHash=${BUILD_HASH} -X github.com/nexptr/llmchain.VERSION=${VERSION} -s -w"

DESTDIR=${PROJECT_PATH}/build
VERSION=v0.0.2

ifeq ($(BUILD_TYPE), "generic")
	GENERIC_PREFIX:=generic-
else
	GENERIC_PREFIX:=
endif


.PHONY: all


all : llmchain


llms/local/chat_grpc.pb.go llms/local/chat.pb.go:
	@echo "===更新生成 chat.proto golang 端代码==="
	@protoc --go_out=${PROJECT_PATH}/ --go-grpc_out=${PROJECT_PATH}/ --go_opt=paths=source_relative  --go-grpc_opt=paths=source_relative ./llms/local/chat.proto

py_proto:
	@cd gen_server && make proto


proto: py_proto llms/local/chat_grpc.pb.go llms/local/chat.pb.go 
	@echo "===更新生成 chat.proto golang 端代码==="

llamacpp/libbinding.a: llamacpp 
	$(MAKE) -C llm/llamacpp $(GENERIC_PREFIX)libbinding.a


app: llamacpp/libbinding.a
	@echo "create llmchain-${VERSION} "
	@mkdir -p ${DESTDIR}/llmchain-${VERSION}/{conf,bin,i18n}

	@echo "copy default configure file"
	@cp -f ${PROJECT_PATH}/examples/app/conf.yaml ${DESTDIR}/llmchain-${VERSION}/conf/conf.yaml

	@echo "build github.com/nexptr/llmchain/examples/app"
	@env  go build -ldflags ${LDFLAGS} -o ${DESTDIR}/llmchain-${VERSION}/bin/app github.com/nexptr/llmchain/examples/app


clean:
	rm -rf ${DESTDIR}
	docker rmi llmchain:${VERSION}

