package main

import (
	"context"
    mrand "math/rand/v2"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"github.com/segmentio/kafka-go"
    "github.com/google/uuid"
)

//Структура для заказа
type Order struct {
    OrderUID string `json:"order_uid"`
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
    
    // Отправка 10 сообщений о добавлении случайных заказов в БД
    for i := 0; i < 10; i++ {
        order := generateRandomOrder()

        // Отправка сообщения в Kafka
        message, err := json.Marshal(order)
        if err != nil {
            log.Printf("Failed to marshal order: %v", err)
			continue
        }

        // Отправка сообщения в Kafka
        err = writer.WriteMessages(context.Background(),
            kafka.Message{
                Key:   []byte(order.OrderUID), // Ключа сообщения
                Value: message, // Тело сообщения
            },
        )

        if err != nil {
            log.Fatalf("Failed to write message: %v", err)
        } else {
            log.Printf("Message successfully sent to Kafka")
        }

        time.Sleep(1 * time.Second)
    }
}

// Генерация случайного заказа
func generateRandomOrder() Order {
    orderUID := uuid.New().String()
    currentTime := time.Now()

    delivery := Delivery {
        OrderUID: orderUID,
        Name:    generateRandomString(100),
        Phone:   generateRandomPhone(),
        Zip:     fmt.Sprintf("%d", mrand.IntN(100000)),
        City:    generateRandomString(20),
        Address: generateRandomString(20),
        Region:  generateRandomString(20),
        Email:   generateRandomEmail(),
    }

    payment := Payment {
        OrderUID:     orderUID,
        Transaction:  orderUID,
        RequestID:    generateRandomString(10),
        Currency:     "USD",
        Provider:     "WBPAY",
        Amount:       mrand.IntN(100000) + 10,
        PaymentDt:    currentTime.Unix() - int64(mrand.IntN(1000000)),
        Bank:         generateRandomString(15),
        DeliveryCost: mrand.IntN(1000),
        GoodsTotal:   mrand.IntN(500),
        CustomFee:    mrand.IntN(100),
    }

    itemsCount := mrand.IntN(5) + 1 // Случайное число предметов
    var items []Item
    for i := 0; i < itemsCount; i++ {
        item := Item {
            OrderUID:    orderUID,
            ChrtID:      mrand.IntN(100000) + 1,
            TrackNumber: generateRandomString(10), 
            Price:       mrand.IntN(100000),
            Rid:         generateRandomString(20),
            Name:        generateRandomString(15),
            Sale:        mrand.IntN(99),
            Size:        fmt.Sprintf("%d", mrand.IntN(5) + 1),
            TotalPrice:  mrand.IntN(100000),
            NmID:        mrand.IntN(100000) + 1,
            Brand:       generateRandomString(12),
            Status:      mrand.IntN(999),
        }
        items = append(items, item)
    }

    order := Order {
        OrderUID:          orderUID,
        TrackNumber:       generateRandomString(10),
        Entry:             generateRandomString(4),
        Locale:            "ru",
        InternalSignature: generateRandomString(10),
        CustomerID:        generateRandomString(8),
        DeliveryService:   generateRandomString(7),
        ShardKey:          fmt.Sprintf("%d", mrand.IntN(10)),
        SMID:              mrand.IntN(100) + 1,
        DateCreated:       currentTime,
        OOFShard:          fmt.Sprintf("%d", mrand.IntN(10)),
        Delivery:          delivery,
        Payment:           payment,
        Items:             items,
    }
    
    return order
}


// Генерация случайной строки
func generateRandomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[mrand.IntN(len(charset))]
    }
    return string(b)
}


// Генерация слайного номера телефона
func generateRandomPhone() string {
    return fmt.Sprintf("+7%d", mrand.IntN(1000000000))
}

// Генерация случайного email
func generateRandomEmail() string {
    return fmt.Sprintf("%s@%s.com", generateRandomString(10), generateRandomString(5))
}