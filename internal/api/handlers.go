package handlers

import (
    "net/http"
	"errors"
    "log"
    "context"
    "github.com/jackc/pgx/v5"
	"github.com/venexene/wbl0-orders-service/internal/db"
    "github.com/gin-gonic/gin"
)

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
        if errors.Is(err, pgx.ErrNoRows) {
            // Заказ не найден
            c.JSON(http.StatusNotFound, gin.H{
                "error": "Failed to find order",
            })
        } else if errors.Is(err, context.DeadlineExceeded) {
            // Таймаут запроса
            c.JSON(http.StatusGatewayTimeout, gin.H{
                "error": "Request timeout",
            })
        } else {
            // Другие ошибки
            log.Printf("Database error: %v", err) // Логирование для дебагинга
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Internal server error",
            })
        }
        return
    }

    c.JSON(http.StatusOK, order)
}