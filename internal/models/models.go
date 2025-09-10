package models

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
)


//Структура для заказа
type Order struct {
    OrderUID          string    `json:"order_uid" validate:"required,uuid4"`
    TrackNumber       string    `json:"track_number" validate:"required,alphanum,max=50"`
    Entry             string    `json:"entry" validate:"required,alphanum,max=10"`
    Locale            string    `json:"locale" validate:"required,max=2"`
    InternalSignature string    `json:"internal_signature" validate:"omitempty,alphanum,max=100"`
    CustomerID        string    `json:"customer_id" validate:"required,alphanum,max=50"`
    DeliveryService   string    `json:"delivery_service" validate:"required,alphanum,max=50"`
    ShardKey          string    `json:"shardkey" validate:"required,alphanum,max=10"`
    SMID              uint      `json:"sm_id" validate:"required,min=1"`
    DateCreated       time.Time `json:"date_created" validate:"required"`
    OOFShard          string    `json:"oof_shard" validate:"required,alphanum,max=10"`
    Delivery          Delivery  `json:"delivery" validate:"required"`
    Payment           Payment   `json:"payment" validate:"required"`
    Items             []Item    `json:"items" validate:"required,min=1,dive"`
}

//Структура для доставки
type Delivery struct {
    OrderUID string `json:"-"`
    Name     string `json:"name" validate:"required,max=100"`
    Phone    string `json:"phone" validate:"required,e164"`
    Zip      string `json:"zip" validate:"required,max=10"`
    City     string `json:"city" validate:"required,max=100"`
    Address  string `json:"address" validate:"required,max=100"`
    Region   string `json:"region" validate:"required,max=100"`
    Email    string `json:"email" validate:"required,email"`
}

//Структура для оплаты
type Payment struct {
    OrderUID     string `json:"-"`
    Transaction  string `json:"transaction" validate:"required,uuid4"`
    RequestID    string `json:"request_id" validate:"omitempty,max=50"`
    Currency     string `json:"currency" validate:"required,max=3"`
    Provider     string `json:"provider" validate:"required,max=50"`
    Amount       int    `json:"amount" validate:"required,min=1"`
    PaymentDt    uint64 `json:"payment_dt" validate:"required"`
    Bank         string `json:"bank" validate:"required,alphanum,max=20"`
    DeliveryCost uint   `json:"delivery_cost" validate:"required"`
    GoodsTotal   uint   `json:"goods_total" validate:"required"`
    CustomFee    uint   `json:"custom_fee" validate:"required"`
}

//Структура для предмета заказа
type Item struct {
    ID          int    `json:"-"`
    OrderUID    string `json:"-"`
    ChrtID      uint   `json:"chrt_id" validate:"required,min=1"`
    TrackNumber string `json:"track_number" validate:"required,alphanum,max=50"`
    Price       uint   `json:"price" validate:"required,min=1"` 
    Rid         string `json:"rid" validate:"required,alphanum,max=50"`
    Name        string `json:"name" validate:"required,max=50"`
    Sale        uint   `json:"sale" validate:"required"`
    Size        string `json:"size" validate:"required,alphanum,max=10"`
    TotalPrice  uint   `json:"total_price" validate:"required"`
    NmID        uint   `json:"nm_id" validate:"required,min=1"`
    Brand       string `json:"brand" validate:"required,max=50"`
    Status      uint   `json:"status" validate:"required,max=999"`
}

// Загрузка заказа из файла
func LoadOrderFromFile(path string) (*Order, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("Failed to read file %s: %v", path, err)
    }

    var order Order
    if err := json.Unmarshal(data, &order); err != nil {
        return nil, fmt.Errorf("Failed to unmarshal JSON: %v", err)
    }

    validate := validator.New()
    if err := validate.Struct(order); err != nil {
        return nil, fmt.Errorf("Failed to validate order: %v", err)
    }

    return &order, nil
}