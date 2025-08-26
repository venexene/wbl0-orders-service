package consumer

import (
	"github.com/segmentio/kafka-go"
	"github.com/venexene/wbl0-orders-service/internal/db"
)

type Consumer struct {
	reader *kafka.Reader
	storage *database.Storage
}


