# ########################################################## #
# Makefile for Golang Project
# Includes cross-compiling, installation, cleanup
# ########################################################## #

.PHONY: list clean install build build_all glide test vet lint format all

# Check for required command tools to build or stop immediately
EXECUTABLES = go pwd uname date cat
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

# Vars
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BIN_DIR=bin
COVERAGE_DIR=coverage
BIN_RBN_KINESIS=rbn-to-kinesis
COV_PROFILE=${COVERAGE_DIR}/test-coverage.txt
COV_HTML=${COVERAGE_DIR}/test-coverage.html
PKG_RBN_KINESIS=gitlab.disney.com/guys-workspace/rbn-to-kinesis
PLATFORMS=linux 
ARCHITECTURES=amd64
VERSION=$(shell cat version.txt | cut -d' ' -f2)
BUILD=`date +%FT%T%z`
UNAME=$(shell uname)
GOLIST=$(shell go list ./...)

# Binary Build
LDFLAGS_RBN_KINESIS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

# Test Build
LDFLAGS_TEST=-ldflags "-X ${PKG_RBN_KINESIS}.Version=${VERSION} -X ${PKG_RBN_KINESIS}.Build=${BUILD}

default: build

all: clean lint build test install

list:
	@echo "Available GNU make targets..."
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs

vet:
	go vet ${PKG_RBN_KINESIS}

lint:
	go get golang.org/x/lint/golint
	go get github.com/GoASTScanner/gas/cmd/gas
	golint -set_exit_status ${GOLIST}

format:
	go fmt ${PKG_RBN_KINESIS}

build: format vet
	CGO_ENABLED=0 go build -o ${BIN_DIR}/${BIN_RBN_KINESIS} ${LDFLAGS_RBN_KINESIS} ${PKG_RBN_KINESIS}
	docker build -t containerregistry.disney.com/digital/rbn-to-kinesis -f Docker/Dockerfile .

build_all: format vet
    $(shell export GOOS=linux; export GOARCH=amd64; go build -v ${PKG_RBN_KINESIS})

test: build_all
	# tests and code coverage
	mkdir -p $(COVERAGE_DIR)
	go test ${GOLIST} -short -v ${LDFLAGS_TEST} -coverprofile ${COV_PROFILE}
	go tool cover -html=${COV_PROFILE} -o ${COV_HTML}
ifeq ($(UNAME), Darwin)
	open ${COV_HTML}
endif

docs: 
	go get golang.org/x/tools/cmd/godoc
	open http://localhost:6060/${PKG_RBN_KINESIS}
	godoc -http=":6060"
	
install:
	go install -o ${BIN_RBN_KINESIS} ${LDFLAGS_RBN_KINESIS} ${PKG_RBN_KINESIS}

# Remove only what we've created
clean:
	if [ -d ${BIN_DIR} ] ; then rm -rf ${BIN_DIR} ; fi
	if [ -d ${COVERAGE_DIR} ] ; then rm -rf ${COVERAGE_DIR} ; fi
