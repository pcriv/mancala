SERVICE:=mancala
BUILD_OPTS:=-ldflags="-s -w" -mod=vendor

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
setup: deps git-hooks .env .env.test

## Install git hooks
git-hooks:
	lefthook install

## Install dependencies
deps:
	go install github.com/joho/godotenv/cmd/godotenv@latest
	go install gotest.tools/gotestsum@latest
	go install github.com/evilmartians/lefthook@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/google/yamlfmt/cmd/yamlfmt@latest


##@ Build

## Build for local platform
build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILD_OPTS) -o $(SERVICE) main.go

##@ Testing

## Run tests
test: .env.test
	godotenv -f .env.test gotestsum --format=testdox -- -cover ./...

## Run tests with coverage
test.coverage:
	mkdir -p ./coverage
	godotenv -f .env.test gotestsum --format=testdox -- -covermode=count -coverpkg=./... -coverprofile coverage/coverage.out ./...
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
	golangci-lint run $(if $(CI),,--fix) -c .golangci.yml --timeout 5m

## Lint yaml files
# Autofix is disabled if CI is set
lint.yml:
	yamlfmt -conf .yamlfmt.yml

##@ Local Development

## Run grpc-server
grpc-server:
	godotenv -f .env go run ./cmd/grpc-server/

## Run connect-server
connect-server:
	godotenv -f .env go run ./cmd/connect-server/

## Run rest-server
rest-server:
	godotenv -f .env go run ./cmd/rest-server/

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
