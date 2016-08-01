GO := $(shell which go)
GO_BUILD = ${GO} build ${GO_BUILD_FLAGS}
GO_BUILD_FLAGS := -v -x -buildmode=c-shared

ROOT:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BUILD_DIR = ${ROOT}/build

SRC := ${ROOT}/lib/complete
LIBRARY_NAME := $(shell basename $(subst -,,${ROOT}))
PACKAGE_PATH := $(subst $(shell go env GOPATH)/src/,,${SRC})

CC ?= $(shell which clang)
CGO_CFLAGS ?=
CGO_LDFLAGS ?= -L$(shell llvm-config --libdir)

all: build

build/golibclang.so:
	CC=${CC} CGO_CFLAGS=${CGO_CFLAGS} CGO_LDFLAGS=${CGO_LDFLAGS} ${GO_BUILD} -o ${BUILD_DIR}/${LIBRARY_NAME}.so ${PACKAGE_PATH}

run: clean build/golibclang.so
	cd build; python -c "import golibclang; print(golibclang.complete(19, 24))"

clean:
	${RM} -r ./build

.PHONY: build/golibclang.so
