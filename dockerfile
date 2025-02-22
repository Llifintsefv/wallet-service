FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o main cmd/wallet-api/main.go
RUN ls -l /app


FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache ca-certificates
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
RUN ls -l /app
EXPOSE 8080
CMD ["./main"]