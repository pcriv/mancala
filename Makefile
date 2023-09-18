API_SPEC="openapi/spec.yml"
REDIS_URL="redis://127.0.0.1:6379/10"

server:
	env REDIS_URL=$(REDIS_URL) go run ./cmd/server

test:
	go test ./...

generate.api:
	oapi-codegen -package openapi -generate "types"  $(API_SPEC) > internal/web/openapi/types.gen.go
	oapi-codegen -package openapi -generate "server" $(API_SPEC) > internal/web/openapi/server.gen.go
	oapi-codegen -package openapi -generate "spec"   $(API_SPEC) > internal/web/openapi/spec.gen.go
