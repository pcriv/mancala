.PHONY: server test build_docs serve_docs

API_SPEC="api/openapi.yml"
REDIS_URL="redis://127.0.0.1:6379/10"

server:
	env REDIS_URL=$(REDIS_URL) go run ./cmd/api

test:
	go test ./...

generate.api:
	oapi-codegen -package openapi -generate "types"  $(API_SPEC) > internal/openapi/types.gen.go
	oapi-codegen -package openapi -generate "server" $(API_SPEC) > internal/openapi/server.gen.go
	oapi-codegen -package openapi -generate "spec"   $(API_SPEC) > internal/openapi/spec.gen.go
