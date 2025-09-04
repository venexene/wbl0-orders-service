package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/venexene/wbl0-orders-service/internal/config"
	"github.com/venexene/wbl0-orders-service/internal/db"
	"github.com/venexene/wbl0-orders-service/internal/api"
	"github.com/venexene/wbl0-orders-service/internal/kafka"
)

func main() {
	// Получение конфигураций
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("Loaded config")


	// Создание контекста для получения сигнала о завершении
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()


	// Подключение к БД через создание пула соединений
    pool, err := database.CreatePool(cfg)
    if err != nil {
        log.Fatalf("Failed to connect database: %v", err)
    }
    defer pool.Close()
	storage := database.NewStorage(pool)
	log.Println("Connected database")

	
	// Создание консьюмера Kafka
	kafkaConsumer := consumer.NewConsumer(
		strings.Split(cfg.KafkaBrokers, ","),
		cfg.KafkaTopic,
		storage,
	)
	defer kafkaConsumer.Close()
	log.Println("Created Kafka consumer")

	//Запуск консьюмера в горутине
	go func() {
		kafkaConsumer.Consume(context.Background())
	} ()
	log.Printf("Started consume proccess for topic %s", cfg.KafkaTopic)
	

	// Создание роутера
	router := gin.Default()
	log.Printf("Created GIN router")


	// Создание хендлера
	handler := handlers.NewHandler(storage, cfg)


    //Тестовый эндпоинт для проверки работы сервера
    router.GET("/server_check", func(c *gin.Context) {
		handler.TestServerHandle(c)
    })

    //Тестовый эндпоинт для проверки подключения к базе
    router.GET("/db_check", func(c *gin.Context) {
		handler.TestDBHandle(c)
    })

	//Тестовый эндпоинт для проверки работы Kafka
	router.GET("/kafka_check", func(c *gin.Context) {
		handler.TestKafkaHandle(c)
	})

	//Эндпоинт для получения информации о заказе по UID
	router.GET("/orders/:uid", func(c *gin.Context) {
    	handler.GetOrderByUIDHandle(c)
	})

	//Эндпоинт для получения UID всех заказов
	router.GET("/all_orders_uids", func(c *gin.Context) {
    	handler.GetAllOrdersUIDHandle(c)
	})


	// Создание сервера
	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}
	log.Printf("Created server")


	// Запуск сервера в горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()
	log.Printf("Started HTTP server on port %s", cfg.HTTPPort)


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