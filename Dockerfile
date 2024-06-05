# Use a imagem oficial do Go como base
FROM golang:latest AS builder

WORKDIR /app

COPY . . 

RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]
