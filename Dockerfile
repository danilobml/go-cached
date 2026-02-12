FROM golang:1.25.6-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o go-cached ./cmd/api

FROM scratch

COPY --from=builder ["/build/go-cached", "/"]

ENTRYPOINT ["/go-cached"]
