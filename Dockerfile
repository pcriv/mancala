FROM golang:1.25

WORKDIR /app

COPY ./ ./

EXPOSE 1323
ENTRYPOINT ["go", "run", "./cmd/server"]
