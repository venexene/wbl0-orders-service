FROM golang:1.24

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем ВЕСЬ проект
COPY . .

# Собираем приложение
RUN go build -o main ./cmd

# Экспонируем порт
EXPOSE 8080

RUN chmod +x main

# Запускаем приложение из папки cmd
CMD ["./main"]