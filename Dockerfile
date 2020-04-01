FROM golang:1.14.1-alpine3.11

WORKDIR /app

RUN apk add build-base
RUN apk add --no-cache git mercurial
RUN go get github.com/cespare/reflex

COPY ./go.mod ./go.sum /app/
RUN go mod download

COPY . /app/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -installsuffix cgo -o /app/loader ./cmd/loader/main.go

ENTRYPOINT reflex -d none -s -R "^client|node_modules" -r \.go$ -- go run cmd/server/main.go
