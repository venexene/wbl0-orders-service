package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/venexene/wbl0-orders-service/internal/config"
	"github.com/venexene/wbl0-orders-service/internal/db"
	"github.com/venexene/wbl0-orders-service/internal/api"
)

func main() {

	// Получение конфигураций
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}


	// Создание контекста для получения сигнала о завершении
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()


	// Подключение к базе данных через создание пула соединений
    pool, err := database.CreatePool(cfg)
    if err != nil {
        log.Fatalf("Failed to connect database: %v", err)
    }
    defer pool.Close()
    log.Println("DataBase connected")

	storage := database.NewStorage(pool)

    log.Println("DataBase connected")
	// Создание роутера
	router := gin.Default()


    //Тестовый эндпоинт для проверки работы сервера
    router.GET("/server_check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "Server definetly works",
        })
    })


    //Тестовый эндпоинт для проверки подключения к базе
    router.GET("/db_check", func(c *gin.Context) {
		var result string
		err := pool.QueryRow(context.Background(), "SELECT 'DataBase definetly works'").Scan(&result)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to connect database",
			})
			return
		}
    
		c.JSON(http.StatusOK, gin.H{
			"status": result,
        })
    })


	//Эндпоинт для получения информации о заказе
	router.GET("/orders/:uid", func(c *gin.Context) {
    	handlers.GetOrderByUIDHandler(c, storage)
	})


	// Создание сервера
	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	// Запуск сервера в горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()
	log.Printf("HTTP server started on port %s", cfg.HTTPPort)


	// Ожидание сигнала завершения
	<-ctx.Done()
	stop() // Отмена подписки на сигнал
	log.Println("Shutting down server...")
	
	//Создание контекста с таймаутом для корректного завершения
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Закрытие сервера
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}
	log.Println("Shutdown server")
}