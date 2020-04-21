FROM golang:1.14.1-alpine3.11

WORKDIR /app

COPY ./go.mod ./go.sum /app/
RUN go mod download

COPY . /app/

CMD go run ./cmd/loader/main.go
