# Makefile.

PROJECT_PATH=$(shell cd "$(dirname "$0" )" &&pwd)
PROJECT_NAME=llmchain
VERSION=$(shell git describe --tags | sed 's/\(.*\)-.*/\1/')
BUILD_DATE=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
BUILD_HASH=$(shell git rev-parse HEAD)
LDFLAGS="-X github.com/cxbooks/cxbooks.buildstamp=${BUILD_DATE} -X github.com/cxbooks/cxbooks.githash=${BUILD_HASH} -X github.com/cxbooks/cxbooks.VERSION=${VERSION} -s -w"
DESTDIR=${PROJECT_PATH}/build
VERSION=v0.0.2


.PHONY: all


all : cxbooks


clean:
	rm -rf ${DESTDIR}
	docker rmi cxbooks:${VERSION}