package config

import (
	"os"
    "github.com/joho/godotenv"
)

// Структура для определения всех конфигураций
type Config struct {
	HTTPPort  string
    DBHost    string
    DBPort    string
    DBUser    string
    DBPass    string
    DBName    string
    DBSSLMode string
    KafkaBrokers string
    KafkaTopic   string
    KafkaGroupID string
}

func Load() (*Config, error) {
	// Загрузка файла .env
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	
    // Чтение переменных файла .env
	return &Config{
		HTTPPort:os.Getenv("HTTP_PORT"),
        DBHost:os.Getenv("DB_HOST"),
        DBPort:os.Getenv("DB_PORT"),
        DBUser:os.Getenv("DB_USER"),
        DBPass:os.Getenv("DB_PASSWORD"),
        DBName:os.Getenv("DB_NAME"),
        DBSSLMode:os.Getenv("DB_SSL_MODE"),
        KafkaBrokers:os.Getenv("KAFKA_BROKERS"),
        KafkaTopic:os.Getenv("KAFKA_TOPIC"),
        KafkaGroupID:os.Getenv("KAFKA_GROUP_ID"),
	}, nil
}