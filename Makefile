SERVICE:=mancala
BUILD_OPTS:=-ldflags="-s -w"

# Absolute path to this makefile
THIS_MAKEFILE := $(abspath $(lastword $(MAKEFILE_LIST)))

export GOOS ?= $(shell go env GOOS)
export GOARCH ?= $(shell go env GOARCH)

##@ General

## Print this help message
# Parses this Makefile and prints targets that are preceded by "##" comments
help:
	@awk -F : '\
		BEGIN {	\
			in_doc = 0; \
        	printf "\nUsage:\n  make \033[36m<target>\033[0m\n"; \
		} \
		/^##@/ { \
        	printf "\n\033[1m%s\033[0m\n", substr($$0, 5); \
    	} \
		/^## / && in_doc == 0 { \
			in_doc = 1; \
			doc_first_line = $$0; \
			sub(/^## */, "", doc_first_line); \
		} \
		$$0 !~ /^#/ && in_doc == 1 { \
			in_doc = 0; \
			if (NF <= 1) { \
				next; \
			} \
			printf "  \033[36m%-25s\033[0m %s\n", $$1, doc_first_line; \
		} \
	' < "$(THIS_MAKEFILE)"

## Clean all caches and intermediates
clean:
	go clean -i

## Setup local environment
setup: deps-macos .env .env.test

## Install dependencies
deps-macos:
	brew install dbmate
	brew install golangci-lint
	brew install mockery
	brew install yamlfmt


##@ Build

## Build for local platform
build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILD_OPTS) -o bin/$(SERVICE)-connectrpc ./cmd/connectrpc/
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILD_OPTS) -o bin/$(SERVICE)-restapi ./cmd/restapi/
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILD_OPTS) -o bin/$(SERVICE)-web ./cmd/web/

##@ Testing

## Run tests
test: .env.test
	go tool godotenv -f .env.test go tool gotestsum --format=testdox -- -cover ./...

## Run tests with coverage
test.coverage:
	mkdir -p ./coverage
	go tool godotenv -f .env.test go tool gotestsum --format=testdox -- -covermode=count -coverpkg=./... -coverprofile coverage/coverage.out ./...
	grep -v -E -f .covignore ./coverage/coverage.out > ./coverage/coverage.filtered.out
	mv ./coverage/coverage.filtered.out ./coverage/coverage.out
	go tool cover -func coverage/coverage.out -o coverage/coverage.tool
	go tool cover -html coverage/coverage.out -o coverage/coverage.html


##@ Lint

## Run all linters
lint: lint.go lint.yml

## Lint go sources
# Autofix is disabled if CI is set
lint.go:
	go tool golangci-lint run

## Lint yaml files
# Autofix is disabled if CI is set
lint.yml:
	go tool yamlfmt -conf .yamlfmt.yml

##@ Local Development

## Run connectrpc server
connectrpc:
	go tool godotenv -f .env go run ./cmd/connectrpc/

## Run restapi server
rest:
	go tool godotenv -f .env go run ./cmd/restapi/

## Run web frontend
web:
	go tool godotenv -f .env go run ./cmd/web/

# Ensures .env exists
.env:
	cp .env.example .env
	echo ENV=local >> .env

# Ensures .env.test exists
.env.test:
	cp .env.example .env.test
	echo ENV=test >> .env.test

##@ Code Generation

## Generate code
gen:
	go generate ./...
	buf generate proto/

.PRECIOUS: .env .env.test
