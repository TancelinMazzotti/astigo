# Build stage
FROM golang:1.24.2-alpine AS debugger

RUN apk add --update --no-cache build-base vim curl htop
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -gcflags="all=-N -l"  \
    -o astigo .

RUN cp $(go env GOPATH)/bin/dlv /usr/local/bin/

EXPOSE 8080
EXPOSE 50051

EXPOSE 2345

CMD ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./astigo"]