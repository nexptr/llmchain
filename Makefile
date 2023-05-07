# Makefile.

PROJECT_PATH=$(shell cd "$(dirname "$0" )" &&pwd)
PROJECT_NAME=llmchain
VERSION=$(shell git describe --tags | sed 's/\(.*\)-.*/\1/')
BUILD_DATE=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
BUILD_HASH=$(shell git rev-parse HEAD)
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
	$(MAKE) -C llms/llamacpp $(GENERIC_PREFIX)libbinding.a

clean:
	rm -rf ${DESTDIR}
	docker rmi llmchain:${VERSION}

