package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/venexene/wbl0-orders-service/internal/db"
	"github.com/venexene/wbl0-orders-service/internal/models"
)

// Структура консьюмера
type Consumer struct {
	reader  *kafka.Reader
	storage *database.Storage
}

// Конструктор консьюмера
func NewConsumer(brokers []string, topic string, storage *database.Storage) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic: topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
		MaxWait: time.Second,
		Dialer: &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
		},
		MaxAttempts: 3,

	})

	return &Consumer{reader: reader, storage: storage}
}

// Основной метод для получения сообщений
func (c *Consumer) Consume(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(ctx) // Чтение сообщений из Kafka
		if err != nil {
			log.Printf("Kafka failed to consume: %v", err)
			continue
		}
		log.Printf("Received message: %s", string(msg.Value))

		var order models.Order
		// Десериализация JSON
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}
		
		// Сохраниение в БД
		if err := c.storage.AddOrderIfNotExists(ctx, &order); err != nil {
			log.Printf("Failed to add order: %v", err)
		} else {
			log.Printf("Order saved with UID %s", order.OrderUID)
		}
	}
}

// Функция для закрытия соединения с Kafka
func (c *Consumer) Close() error {
	return c.reader.Close()
}