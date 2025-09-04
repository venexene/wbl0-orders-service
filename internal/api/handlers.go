package handlers

import (
    "net/http"
	"errors"
    "context"
    "time"
    "log"
    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5"
    "github.com/segmentio/kafka-go"
	"github.com/venexene/wbl0-orders-service/internal/db"
    "github.com/venexene/wbl0-orders-service/internal/config"
)


// Хендлер для обработки тестового запроса к серверу
func TestServerHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status": "Server definetly works",
    })
}


// Хендлер для обработки тестового запроса к БД
func TestDBHandler(c *gin.Context, storage *database.Storage) {
    res, err := storage.TestDB()
    if err != nil {
        log.Printf("Failed to test database: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to connect database",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status": res,
    })
}


// Хендлер для обработки тестового запроса к Kafka
func TestKafkaHandler(c *gin.Context, cfg *config.Config) {
    kafkaBrokers := cfg.KafkaBrokers

    // Контекст с таймаутом для ограничения времени выполнения операции с Kafka
    ctxKafka, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel() // Отложенный вызов отмены контекса для освобождения ресурсов

    // Установка соединений с Kafka
    connKafka, err := kafka.DialContext(ctxKafka, "tcp", kafkaBrokers)
    if err != nil {
        log.Printf("Failed to test Kafka: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error" : "Failed to connect Kafka",
        })
        return
    }
    defer connKafka.Close() // Отложенное закрытие соединения с Kafka

    broker := connKafka.Broker() // Получение информации о брокере
    c.JSON(http.StatusOK, gin.H {
        "status":    "Kafka definetly works",
        "brokers":   kafkaBrokers,
        "broker_id": broker.ID,
    })
}


// Хендлер для обработки запроса на получение всей информации о заказе по UID
func GetOrderByUIDHandler(c *gin.Context, storage *database.Storage) {
    orderUID := c.Param("uid") // Извлечение UID из URL
    
    // Проверка что передан не пустой UID
    if orderUID == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "No UID recieved",
        })
        return
    }

    // Вызов функции получения всей информации о заказе по UID
    order, err := storage.GetOrderByUID(c.Request.Context(), orderUID)
    
    //Обработка ошибок получения заказа
    if err != nil {
        log.Printf("Failed to get info by UID: %v", err)
        if errors.Is(err, pgx.ErrNoRows) {
            c.JSON(http.StatusNotFound, gin.H{
                "error": "Failed to find order",
            })
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Internal server error",
            })
        }
        return
    }

    c.JSON(http.StatusOK, order)
}


// Хендлер для получения UID всех заказов
func GetAllOrdersUIDHandler(c *gin.Context, storage *database.Storage) {
    orderUIDs, err := storage.GetAllOrdersUID(c)

    if err != nil {
        log.Printf("Failed to get UIDs: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to get order UIDs",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "order_uids": orderUIDs,
    })
}