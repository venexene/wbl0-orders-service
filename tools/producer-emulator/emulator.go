package main

import (
	"context"
	"encoding/json"
	"time"
    "log"

	"github.com/segmentio/kafka-go"
)


//Структура для заказа
type Order struct {
    OrderUID          string    `json:"order_uid"`
    TrackNumber       string    `json:"track_number"`
    Entry             string    `json:"entry"`
    Locale            string    `json:"locale"`
    InternalSignature string    `json:"internal_signature"`
    CustomerID        string    `json:"customer_id"`
    DeliveryService   string    `json:"delivery_service"`
    ShardKey          string    `json:"shardkey"`
    SMID              int       `json:"sm_id"`
    DateCreated       time.Time `json:"date_created"`
    OOFShard          string    `json:"oof_shard"`
    Delivery          Delivery  `json:"delivery"`
    Payment           Payment   `json:"payment"`
    Items             []Item    `json:"items"`
}

//Структура для доставки
type Delivery struct {
    OrderUID string `json:"-"`
    Name     string `json:"name"`
    Phone    string `json:"phone"`
    Zip      string `json:"zip"`
    City     string `json:"city"`
    Address  string `json:"address"`
    Region   string `json:"region"`
    Email    string `json:"email"`
}

//Структура для оплаты
type Payment struct {
    OrderUID     string `json:"-"`
    Transaction  string `json:"transaction"`
    RequestID    string `json:"request_id"`
    Currency     string `json:"currency"`
    Provider     string `json:"provider"`
    Amount       int    `json:"amount"`
    PaymentDt    int64  `json:"payment_dt"`
    Bank         string `json:"bank"`
    DeliveryCost int    `json:"delivery_cost"`
    GoodsTotal   int    `json:"goods_total"`
    CustomFee    int    `json:"custom_fee"`
}

//Структура для предмета заказа
type Item struct {
    ID          int    `json:"-"`
    OrderUID    string `json:"-"`
    ChrtID      int    `json:"chrt_id"`
    TrackNumber string `json:"track_number"`
    Price       int    `json:"price"`
    Rid         string `json:"rid"`
    Name        string `json:"name"`
    Sale        int    `json:"sale"`
    Size        string `json:"size"`
    TotalPrice  int    `json:"total_price"`
    NmID        int    `json:"nm_id"`
    Brand       string `json:"brand"`
    Status      int    `json:"status"`
}


func main() {
	kafkaBrokers := []string{"kafka:9092"}
	topic := "wbl0_orders"

    // Создание райтера Kafka
	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBrokers...), // Преобразование адреса брокера в TCP-формат
		Topic:    topic, // Установка топика
	}
	defer writer.Close() // Отложенное закрытие соединения с райтером

    // Тестовый заказ
	testOrder := Order{
		OrderUID: "6d2a89ac-0ede-40cd-9fed-bd6b88d855b0",
    	TrackNumber: "WBMYTESTTRACKNUMBER",
    	Entry: "WBMY",
    	Locale: "ru",
    	InternalSignature: "",
    	CustomerID: "test_id",
    	DeliveryService: "test_service",
    	ShardKey: "9",
    	SMID: 99,
    	DateCreated: time.Now(),
    	OOFShard: "2",

		Delivery: Delivery{
            OrderUID: "6d2a89ac-0ede-40cd-9fed-bd6b88d855b0",
            Name: "Kafka Test Testov",
            Phone: "+9820000000",
            Zip: "2532712",
            City: "Rostov-on-Don",
            Address: "Lenina street 28",
            Region: "Rostov Region",
            Email: "kafkatest@gmail.com",
        },

        Payment: Payment{
            OrderUID: "6d2a89ac-0ede-40cd-9fed-bd6b88d855b0",
            Transaction: "6d2a89ac-0ede-40cd-9fed-bd6b88d855b0",
            RequestID: "",
            Currency: "RUB",
            Provider: "WBPAY",
            Amount: 2156,
            PaymentDt: 12425325,
            Bank: "Sber",
            DeliveryCost: 2000,
            GoodsTotal: 453,
            CustomFee: 10,
        },

        Items: [] Item{
            {
                OrderUID: "6d2a89ac-0ede-40cd-9fed-bd6b88d855b0",
                ChrtID: 2200232,
                TrackNumber: "WBMYTESTTRACK",
                Price: 559,
                Rid: "cb3416617a723ae0btest",
                Name: "Book",
                Sale: 20,
                Size: "3",
                TotalPrice: 450,
                NmID: 3335551,
                Brand: "AST",
                Status: 404,
            },

            {
                OrderUID: "6d2a89ac-0ede-40cd-9fed-bd6b88d855b0",
                ChrtID: 2221232,
                TrackNumber: "WBMYTESTTRACK",
                Price: 234,
                Rid: "cb3416617a723ae0btest",
                Name: "Pen",
                Sale: 14,
                Size: "1",
                TotalPrice: 200,
                NmID: 3335251,
                Brand: "AST",
                Status: 403,
            },
        },
	}

    // Отправка сообщения в Kafka
    message, err := json.Marshal(testOrder)
    if err != nil {
        log.Fatalf("Failed to marshal order: %v", err)
    }

    // Отправка сообщения в Kafka
    err = writer.WriteMessages(context.Background(),
        kafka.Message{
            Key: []byte(testOrder.OrderUID), // Ключа сообщения
            Value: message, // Тело сообщения
        },
    )

    if err != nil {
        log.Fatalf("Failed to write message: %v", err)
    } else {
        log.Printf("Message successfully sent to Kafka")
    }
}