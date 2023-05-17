# Makefile.

PROJECT_PATH=$(shell cd "$(dirname "$0" )" &&pwd)
PROJECT_NAME=llmchain
VERSION=$(shell git describe --tags | sed 's/\(.*\)-.*/\1/')
BUILD_DATE=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
BUILD_HASH=$(shell git rev-parse HEAD)
LDFLAGS="-X github.com/exppii/llmchain.BuildStamp=${BUILD_DATE} -X github.com/exppii/llmchain.GitHash=${BUILD_HASH} -X github.com/exppii/llmchain.VERSION=${VERSION} -s -w"

DESTDIR=${PROJECT_PATH}/build
VERSION=v0.0.2

ifeq ($(BUILD_TYPE), "generic")
	GENERIC_PREFIX:=generic-
else
	GENERIC_PREFIX:=
endif


.PHONY: all


all : llmchain


llamacpp:
	git submodule update --init --recursive --depth 1

llamacpp/libbinding.a: llamacpp 
	$(MAKE) -C llm/llamacpp $(GENERIC_PREFIX)libbinding.a


app: llamacpp/libbinding.a
	@echo "create llmchain-${VERSION} "
	@mkdir -p ${DESTDIR}/llmchain-${VERSION}/{conf,bin,i18n}

	@echo "copy default configure file"
	@cp -f ${PROJECT_PATH}/examples/app/conf.yaml ${DESTDIR}/llmchain-${VERSION}/conf/conf.yaml

	@echo "build github.com/exppii/llmchain/examples/app"
	@env  go build -ldflags ${LDFLAGS} -o ${DESTDIR}/llmchain-${VERSION}/bin/app github.com/exppii/llmchain/examples/app


clean:
	rm -rf ${DESTDIR}
	docker rmi llmchain:${VERSION}

