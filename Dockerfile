
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go mod файлы
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Если .env не существует, копируем .env.example в .env
RUN [ -f .env ] || cp .env.example .env

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# минимальный образ для запуска
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем собранное приложение
COPY --from=builder /app/main .

# Копируем .env файл
COPY --from=builder /app/.env .env

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]
