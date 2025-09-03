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

// Структура для работы с БД
type Storage struct {
    pool *pgxpool.Pool
}

// Конструктор структуры для БД
func NewStorage(pool *pgxpool.Pool) *Storage {
    return &Storage{pool: pool}
}

// Создание пула соединений к БД
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
        return pool, fmt.Errorf("Failed to create pool: %v", err)
    }

    // Проверка соединения
    if err := pool.Ping(context); err != nil {
        return pool, fmt.Errorf("Failed to ping database: %v", err)
    }

    return pool, nil
}



// Получение заказа по UID из БД
func (s *Storage) GetOrderByUID(context context.Context, orderUID string) (*models.Order, error) {
    // Получение основной информации о заказе
    orderQuery := "SELECT * FROM orders WHERE order_uid = $1"
    var order models.Order
    err := s.pool.QueryRow(context, orderQuery, orderUID).Scan(
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
    err = s.pool.QueryRow(context, deliveryQuery, orderUID).Scan(
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
    paymentQuery := "SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_uid = $1"
    var payment models.Payment
    err = s.pool.QueryRow(context, paymentQuery, orderUID).Scan(
        &payment.Transaction,
        &payment.RequestID,
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
    rows, err := s.pool.Query(context, itemsQuery, orderUID)
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


// Добавление заказа в БД
func (s *Storage) AddOrder(ctx context.Context, order *models.Order) error {
    // Начало транзакции для атомарного добавления данных
    tx, err := s.pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("Failed to begin transaction: %v", err)
    }
    defer tx.Rollback(ctx) // Откат транзакции в случае ошибки

    // Добавление основной информации о заказе
    orderQuery := `
        INSERT INTO orders (
            order_uid, track_number, entry, locale, internal_signature,
            customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `
    _, err = tx.Exec(ctx, orderQuery,
        order.OrderUID,
        order.TrackNumber,
        order.Entry,
        order.Locale,
        order.InternalSignature,
        order.CustomerID,
        order.DeliveryService,
        order.ShardKey,
        order.SMID,
        order.DateCreated,
        order.OOFShard,
    )
    if err != nil {
        return fmt.Errorf("Failed to insert order: %v", err)
    }

    // Добавление информации о доставке
    deliveryQuery := `
        INSERT INTO delivery (
            order_uid, name, phone, zip, city, address, region, email
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
    _, err = tx.Exec(ctx, deliveryQuery,
        order.OrderUID,
        order.Delivery.Name,
        order.Delivery.Phone,
        order.Delivery.Zip,
        order.Delivery.City,
        order.Delivery.Address,
        order.Delivery.Region,
        order.Delivery.Email,
    )
    if err != nil {
        return fmt.Errorf("Failed to insert delivery: %v", err)
    }

    // Добавление информации о платеже
    paymentQuery := `
        INSERT INTO payment (
            order_uid, transaction, request_id, currency, provider, amount, 
            payment_dt, bank, delivery_cost, goods_total, custom_fee
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `
    _, err = tx.Exec(ctx, paymentQuery,
        order.OrderUID,
        order.Payment.Transaction,
        order.Payment.RequestID,
        order.Payment.Currency,
        order.Payment.Provider,
        order.Payment.Amount,
        order.Payment.PaymentDt,
        order.Payment.Bank,
        order.Payment.DeliveryCost,
        order.Payment.GoodsTotal,
        order.Payment.CustomFee,
    )
    if err != nil {
        return fmt.Errorf("Failed to insert payment: %v", err)
    }

    // Добавление информации о товарах
    itemQuery := `
        INSERT INTO item (
            order_uid, chrt_id, track_number, price, rid, name, 
            sale, size, total_price, nm_id, brand, status
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    `
    for _, item := range order.Items {
        _, err = tx.Exec(ctx, itemQuery,
            order.OrderUID,
            item.ChrtID,
            item.TrackNumber,
            item.Price,
            item.Rid,
            item.Name,
            item.Sale,
            item.Size,
            item.TotalPrice,
            item.NmID,
            item.Brand,
            item.Status,
        )
        if err != nil {
            return fmt.Errorf("Failed to insert item: %v", err)
        }
    }

    // Подтверждение транзакции
    if err = tx.Commit(ctx); err != nil {
        return fmt.Errorf("Failed to commit transaction: %v", err)
    }

    return nil
}


// Проверка существования заказа по UID
func (s *Storage) OrderExists(ctx context.Context, orderUID string) (bool, error) {
    query := "SELECT EXISTS(SELECT 1 FROM orders WHERE order_uid = $1)"
    var exists bool

    err := s.pool.QueryRow(ctx, query, orderUID).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("Failed to check order existence: %v", err)
    }

    return exists, nil
}


// Добавление заказа с проверкой на существование
func (s *Storage) AddOrderIfNotExists(ctx context.Context, order *models.Order) error {
    exists, err := s.OrderExists(ctx, order.OrderUID)
    if err != nil {
        return err
    }
    if exists {
        return fmt.Errorf("Order with UID %v already exists", order.OrderUID)
    }
    
    return s.AddOrder(ctx, order)
}