API_SPEC="api/openapi-spec/spec.yml"
REDIS_URL="redis://127.0.0.1:6379/10"

server:
	env REDIS_URL=$(REDIS_URL) go run ./cmd/api

test:
	go test ./...

generate.api:
	oapi-codegen -package openapi -generate "types"  api/openapi-spec/spec.yml > internal/pkg/openapi/types.gen.go
	oapi-codegen -package openapi -generate "server" api/openapi-spec/spec.yml > internal/pkg/openapi/server.gen.go
	oapi-codegen -package openapi -generate "spec"   api/openapi-spec/spec.yml > internal/pkg/openapi/spec.gen.go

.PHONY: build
build:
