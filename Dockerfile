# Используем официальный образ Go
FROM golang:alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY ./go.mod ./go.sum ./
# Загружаем зависимости
RUN go mod download

# Копируем весь код приложения
COPY . .

RUN go build -o ./shortener ./cmd/app/main.go

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/shortener ./shortener
COPY --from=builder /app/config/config.yaml ./config/config.yaml
COPY --from=builder /app/internal/database/repo/migrations/ ./internal/database/repo/sql/*.sql

COPY --from=builder /app/web ./web

# Открываем порт
EXPOSE 8082

# Запускаем приложение
CMD ["./shortener"]
