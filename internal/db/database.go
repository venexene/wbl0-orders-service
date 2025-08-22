package database

import (
    "context"
    "fmt"
    "time"
    "errors"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/venexene/wbl0-orders-service/internal/config"
    "github.com/venexene/wbl0-orders-service/internal/models"
)



// Создание пула соединений к базе данных
func CreatePool(cfg *config.Config) (*pgxpool.Pool, error) {
    // Формирование строки подключения
    connectionStr := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=%s",
        cfg.DBUser,
        cfg.DBPass,
        cfg.DBHost,
        cfg.DBPort,
        cfg.DBName,
        cfg.DBSSLMode,
    )

    // Создание контекста с таймаутом для контроля времени выполнения и обработки отмены
    context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Создание пула для соедиения
    pool, err := pgxpool.New(context, connectionStr)
    if err != nil {
        return nil, fmt.Errorf("Failed to create pool: %v", err)
    }

    // Проверка соединения
    if err := pool.Ping(context); err != nil {
        return nil, fmt.Errorf("Failed to ping database: %v", err)
    }

    return pool, nil
}



// Получение заказа по UID из базы данных
func GetOrderByUID(context context.Context, pool *pgxpool.Pool, orderUID string) (*models.Order, error) {
    // Получение основной информации о заказе
    orderQuery := "SELECT * FROM orders WHERE order_uid = $1"
    var order models.Order
    err := pool.QueryRow(context, orderQuery, orderUID).Scan(
        &order.OrderUID,
        &order.TrackNumber,
        &order.Entry,
        &order.Locale,
        &order.InternalSignature,
        &order.CustomerID,
        &order.DeliveryService,
        &order.ShardKey,
        &order.SMID,
        &order.DateCreated,
        &order.OOFShard,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, fmt.Errorf("Failed to find order with UID %v", orderUID)
        }
        return nil, fmt.Errorf("Failed to query database: %v", err)
    }


    // Получение информации о доставке
    deliveryQuery := "SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = $1"
    var delivery models.Delivery
    err = pool.QueryRow(context, deliveryQuery, orderUID).Scan(
        &delivery.Name,
        &delivery.Phone,
        &delivery.Zip,
        &delivery.City,
        &delivery.Address,
        &delivery.Region,
        &delivery.Email,
    )
    if err != nil {
        return nil, fmt.Errorf("Failed to query delivery: %v", err)
    }
    order.Delivery = delivery


    // Получение информации о платеже
    paymentQuery := "SELECT transaction, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_uid = $1"
    var payment models.Payment
    err = pool.QueryRow(context, paymentQuery, orderUID).Scan(
        &payment.Transaction,
        &payment.Currency,
        &payment.Provider,
        &payment.Amount,
        &payment.PaymentDt,
        &payment.Bank,
        &payment.DeliveryCost,
        &payment.GoodsTotal,
        &payment.CustomFee,
    )
    if err != nil {
        return nil, fmt.Errorf("Failed to query payment: %v", err)
    }
    order.Payment = payment


    // Получение информации о товарах
    itemsQuery := "SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM item WHERE order_uid = $1"
    rows, err := pool.Query(context, itemsQuery, orderUID)
    if err != nil {
        return nil, fmt.Errorf("Failed to query items: %v", err)
    }
    defer rows.Close()

    var items []models.Item // Срез для хранения товаров
    for rows.Next() {
        var item models.Item
        err = rows.Scan(
            &item.ChrtID,
            &item.TrackNumber,
            &item.Price,
            &item.Rid,
            &item.Name,
            &item.Sale,
            &item.Size,
            &item.TotalPrice,
            &item.NmID,
            &item.Brand,
            &item.Status,
        )
        if err != nil {
            return nil, fmt.Errorf("Failed to scan item: %v", err)
        }
        items = append(items, item) // Добавление полученного товара в срез
    }

    // Обработка ошибок итерации
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("Failed to iterate items: %v", err)
    }
    order.Items = items


    
    return &order, nil
}