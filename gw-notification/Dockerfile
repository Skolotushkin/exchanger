# Сборочный этап
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
# COPY go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o main ./cmd

# Финальный образ
FROM alpine:3.18

WORKDIR /root/

COPY --from=builder /app/main .
COPY config.env .

EXPOSE 8081

CMD ["./main", "-c", "config.env"]