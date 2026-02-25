FROM golang:alpine AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./shortener ./cmd/app/main.go

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/shortener ./shortener
COPY --from=builder /app/config/config.yaml ./config/config.yaml
COPY --from=builder /app/internal/repo/postgres/migrations/ ./internal/repo/postgres/migrations/*.sql

COPY --from=builder /app/front ./front

# Открываем порт
EXPOSE 8081

# Запускаем приложение
CMD ["./shortener"]
