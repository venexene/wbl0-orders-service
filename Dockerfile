FROM golang:1.25

# Установка рабочей директории
WORKDIR /app

# Копирование зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копирование проекта
COPY . .

# Сборка приложения
RUN go build -o main ./cmd

# Открытие порта
EXPOSE 8080

# Установка возможности запускать main
RUN chmod +x main

# Запуск приложения из папки cmd
CMD ["./main"]