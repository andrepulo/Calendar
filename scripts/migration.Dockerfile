# Используем официальный образ Golang для сборки инструмента миграции
FROM golang:1.20 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

# Собираем инструмент миграции
RUN go build -o migrate ./cmd/migrate.go

# Используем минимальный образ Alpine для выполнения инструмента миграции
FROM alpine:3.19.1

# Устанавливаем необходимые пакеты
RUN apk --no-cache add ca-certificates

# Копируем собранный инструмент миграции из образа сборки
COPY --from=builder /app/migrate /migrate

# Устанавливаем переменные окружения для подключения к базе данных
ENV DB_URI=postgres://postgres:password@db:5432/auth?sslmode=disable

# Команда для запуска инструмента миграции
CMD ["/migrate", "-path", "./migrations", "-database", "$(DB_URI)", "up"]