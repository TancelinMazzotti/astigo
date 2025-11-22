# Build stage
FROM golang:1.24.2-alpine AS builder

RUN apk add --update --no-cache build-base

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o astigo .

# Image definition
FROM alpine:latest
RUN apk update && apk upgrade
RUN apk add curl

RUN adduser -D appuser

WORKDIR /app
RUN mkdir -p /app/config
RUN chown -R appuser:appuser /app

COPY --from=builder /app/astigo .

COPY --chown=appuser:appuser migrations/ migrations/
COPY --chown=appuser:appuser config/config.yaml config/

USER appuser
#HTTP
EXPOSE 8080
#GRPC
EXPOSE 50051

CMD ["./astigo"]