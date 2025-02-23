FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o main cmd/wallet-api/main.go
RUN ls -l /app


FROM postgres:latest
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY --from=builder /app/main .
COPY --from=builder /app/config.env .
RUN ls -l /app
EXPOSE 8080
CMD ["./main"]