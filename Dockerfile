# Этап сборки
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o cmd/bot .

# Финальный образ
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/bot .
CMD ["./bot"]
