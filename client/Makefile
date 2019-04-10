
PROJECT_NAME= domaintest
VERSION = 0.0.1
BUILD_ENV := CGO_ENABLED=0

LDFLAGS=-ldflags "-s -w"
TARGET_EXEC = domainTest
GO_FILE = $(ls *.go|grep -v _test)
SOC:=$(shell uname -m)

.PHONY : all

setup:
	mkdir -p build/aarch64
	mkdir -p build/x86_64
	mkdir -p build/armv7l

all: build-aarch64 build-arm build-amd64

build-aarch64: setup
	${BUILD_ENV} GOARCH=arm64 GOOS=linux  go build ${LDFLAGS} -o build/aarch64/${TARGET_EXEC} ${GO_FILE}

build-arm: setup
	${BUILD_ENV} GOARCH=arm GOOS=linux go build ${LDFLAGS} -o build/armv7l/${TARGET_EXEC} ${GO_FILE}

build-amd64: setup
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build ${LDFLAGS} -o build/x86_64/${TARGET_EXEC} ${GO_FILE}

docker: build-aarch64 build-arm build-amd64
	echo ${SOC}
	sed -i 's/FLAG/${SOC}/1' Dockerfile
	docker build -t ${PROJECT_NAME}:${VERSION} .
	sed -i 's/${SOC}/FLAG/1' Dockerfile

clean:
	rm -rf build
