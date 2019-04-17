PROJECT_NAME= domaintest-server
VERSION = 0.0.1
BUILD_ENV := CGO_ENABLED=0

LDFLAGS=-ldflags "-s -w"
TARGET_EXEC = server
GO_FILE = $(ls *.go|grep -v _test)

.PHONY : all

all: build-amd64
setup:
	mkdir -p build/x86_64

build-amd64: setup
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build ${LDFLAGS} -o build/x86_64/${TARGET_EXEC} ${GO_FILE}

docker: build-amd64
	docker build -t ${PROJECT_NAME}:${VERSION} .
	docker tag ${PROJECT_NAME}:${VERSION} ${PROJECT_NAME}:latest
	
clean:
	rm -rf build