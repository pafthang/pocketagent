FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE=gate
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /service ./cmd/${SERVICE}

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /service .
COPY configs ./configs

ENV CONFIG_DIR=/app/configs

CMD ["./service"]