FROM golang:1.26

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENTRYPOINT ["air"]
