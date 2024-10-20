# Используем официальный образ Golang
FROM golang:1.20 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем остальные файлы
COPY . .

# Собираем приложение
RUN go build -o app ./cmd/main.go

# Используем официальный образ Alpine для выполнения приложения
FROM alpine:3.19.1

# Устанавливаем необходимые пакеты
RUN apk --no-cache add ca-certificates

# Копируем собранное приложение из образа builder
COPY --from=builder /app/app .

# Устанавливаем переменные окружения для подключения к базе данных
ENV DB_URI=postgres://user:password@db:5432/dbname?sslmode=disable
ENV HTTP_PORT=8089

# Запускаем приложение
CMD ["./app"]