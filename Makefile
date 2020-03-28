# Borrowed from: 
# https://github.com/silven/go-example/blob/master/Makefile
# https://vic.demuzere.be/articles/golang-makefile-crosscompile/

BINARY = soundtouch-golang
VET_REPORT = vet.report
TEST_REPORT = tests.xml
GOARCH = amd64


VERSION?=?

ifeq ($(OS), Windows_NT)
        VERSION := $(shell git describe --exact-match --tags 2>nil)
else
        VERSION := $(shell git describe --exact-match --tags 2>/dev/null)
endif


COMMIT=$(shell git rev-parse --short HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Symlink into GOPATH
GITHUB_USERNAME=theovassiliou
BUILD_DIR=${GOPATH}/src/github.com/${GITHUB_USERNAME}/${BINARY}
BIN_DIR=${BUILD_DIR}/bin
CURRENT_DIR=$(shell pwd)
BUILD_DIR_LINK=$(shell readlink ${BUILD_DIR})

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.branch=${BRANCH}"

# Build the project
all: clean build test vet static

build: 
	cd ${BUILD_DIR}; \
	go build volumebutler/volumebutler.go ; \
	cd - >/dev/null

static: linux darwin windows

linux: 
	cd ${BUILD_DIR}; \
	GO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BIN_DIR}/${BINARY}-linux-${GOARCH} . ; \
	cd - >/dev/null

darwin:
	cd ${BUILD_DIR}; \
	GO_ENABLED=0 GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BIN_DIR}/${BINARY}-darwin-${GOARCH} . ; \
	cd - >/dev/null

windows:
	cd ${BUILD_DIR}; \
	GO_ENABLED=0 GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BIN_DIR}/${BINARY}-windows-${GOARCH}.exe . ; \
	cd - >/dev/null

test:
	if ! hash go2xunit 2>/dev/null; then go install github.com/tebeka/go2xunit; fi
	cd ${BUILD_DIR}; \
	go test -v ./... 2>&1 | go2xunit -output ${TEST_REPORT} ; \
	cd - >/dev/null

vet:
	-cd ${BUILD_DIR}; \
	go vet ./... > ${VET_REPORT} 2>&1 ; \
	cd - >/dev/null

fmt:
	cd ${BUILD_DIR}; \
	go fmt $$(go list ./... | grep -v /vendor/) ; \
	cd - >/dev/null

clean:
	-rm -f ${TEST_REPORT}
	-rm -f ${VET_REPORT}
	-rm -f ${BIN_DIR}/${BINARY}-*

.PHONY: static test vet fmt clean build