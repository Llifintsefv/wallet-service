FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /wallet-api ./cmd/wallet-api/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /wallet-api /app/wallet-api
COPY config.env /app/config.env

EXPOSE 8080

CMD ["/app/wallet-api"]