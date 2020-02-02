FROM golang:1.13

WORKDIR /app

COPY ./ ./

EXPOSE 1323
ENTRYPOINT ["go", "run", "server.go"]