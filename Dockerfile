FROM golang:1.14

WORKDIR /app

COPY ./ ./

EXPOSE 1323
ENTRYPOINT ["go", "run", "./cmd/api"]