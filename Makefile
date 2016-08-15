GO := $(shell which go)
GO_BUILD = ${GO} build -v -x
GO_BUILD_FLAGS_SHARED := -buildmode=c-shared

ROOT:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BUILD_DIR = ${ROOT}/build

SRC := ${ROOT}/lib/complete
LIBRARY_NAME := $(shell basename $(subst -,,${ROOT}))
PACKAGE_PATH := $(subst $(shell go env GOPATH)/src/,,${SRC})

CC ?= $(shell which clang)
CGO_CFLAGS ?=
CGO_LDFLAGS ?= -L$(shell /opt/llvm/bin/llvm-config --libdir)

all: build

build/golibclang.so:
	CC=${CC} CGO_CFLAGS=${CGO_CFLAGS} CGO_LDFLAGS=${CGO_LDFLAGS} ${GO_BUILD} ${GO_BUILD_FLAGS_SHARED} -o ${BUILD_DIR}/${LIBRARY_NAME}.so ${PACKAGE_PATH}

build/golibclang:
	CC=${CC} CGO_CFLAGS=${CGO_CFLAGS} CGO_LDFLAGS=${CGO_LDFLAGS} ${GO_BUILD} -o ${BUILD_DIR}/golibclang ./cmd/golibclang

run: clean build/golibclang.so
	cd build; python -c "import golibclang; print(golibclang.complete(str('../testdata/boost-asio_server.cpp'), 19, 24))"

clean:
	${RM} -r ./build

.PHONY: build/golibclang.so build/golibclang
