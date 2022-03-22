# Go related variables.
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOFILES := $(wildcard *.go)
GOPATH := $(shell go env GOPATH)

# for goimports, comma separated
LOCAL_PACKAGES="coinbase"

# don't test commands.  they should have minimal code.
TEST_PACKAGES := $(shell go list ./... | grep -v /cmd/)

# external dependencies
GOIMPORTS := $(shell command -v goimports )
GOLANGCI_LINT := $(shell command -v golangci-lint )


.PHONY: all
all: test build

build:
	@echo ">>> building..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -race -o bin/vwap main.go

goimports:
ifndef GOIMPORTS
	@echo ">>> goimports missing, installing"
	@go install golang.org/x/tools/cmd/goimports@latest > /dev/null
	@GOIMPORTS := $(shell command -v goimports )
endif
	@echo ">>> goimports"
	@${GOIMPORTS} -w -local ${LOCAL_PACKAGES} ./cmd/** ./internal/** # ./pkg/**

golangci_lint:
ifndef GOLANGCI_LINT
	@echo ">>> golangci-lint missing, installing"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOPATH}/bin v1.42.0
	@GOLANGCI_LINT := $(shell command -v golangci-lint )
endif
	@echo ">>> golangci-lint"
	@${GOLANGCI_LINT} run

tidy: goimports golangci_lint
	@echo ">>> mod tidying"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go mod tidy

test:
	@echo ">>> testing"
	@go test -covermode=count -coverprofile=count.out $(TEST_PACKAGES)

coverage: test
	@go tool cover -html=count.out
