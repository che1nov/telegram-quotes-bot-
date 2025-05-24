FROM golang:1.23-alpine

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o bot cmd/main.go

CMD ["./bot"]