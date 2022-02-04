.PHONY: server test build_docs serve_docs

API_SPEC="api/openapi-spec/spec.yml"
REDIS_URL="redis://127.0.0.1:6379/10"

server:
	env REDIS_URL=$(REDIS_URL) go run ./cmd/api

test:
	go test ./...
