.PHONY: server build_docs serve_docs

server:
	env REDIS_URL="redis://127.0.0.1:6379/10" go run server.go

test:
	go test ./...

build_docs:
	yarn redoc-cli bundle spec.yml -o public/docs.html

serve_docs:
	yarn redoc-cli serve spec.yml -w