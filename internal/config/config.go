package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Структура для определения всех конфигураций
type Config struct {
	HTTPPort      string
    CacheCapacity int
    DBHost        string
    DBPort        string
    DBUser        string
    DBPass        string
    DBName        string
    DBSSLMode     string
    KafkaBrokers  string
    KafkaTopic    string
}

func Load() (*Config, error) {
	// Загрузка файла .env
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

    // Преобразование емкости кэша
    cacheCapacityRaw := os.Getenv("CACHE_CAPACITY")
    cacheCapacity, err := strconv.Atoi(cacheCapacityRaw)
    if err != nil {
        cacheCapacity = 100
    }
	
    // Чтение переменных файла .env
	return &Config{
		HTTPPort:      os.Getenv("HTTP_PORT"),
        CacheCapacity: cacheCapacity,
        DBHost:        os.Getenv("DB_HOST"),
        DBPort:        os.Getenv("DB_PORT"),
        DBUser:        os.Getenv("DB_USER"),
        DBPass:        os.Getenv("DB_PASSWORD"),
        DBName:        os.Getenv("DB_NAME"),
        DBSSLMode:     os.Getenv("DB_SSL_MODE"),
        KafkaBrokers:  os.Getenv("KAFKA_BROKERS"),
        KafkaTopic:    os.Getenv("KAFKA_TOPIC"),
	}, nil
}