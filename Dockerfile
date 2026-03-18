FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$(cat VERSION 2>/dev/null || echo dev)" -o remote-manager .

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

COPY --from=builder /build/remote-manager /usr/local/bin/remote-manager

ENTRYPOINT ["remote-manager"]
