CREATE TABLE IF NOT EXISTS orders (
    order_uid UUID PRIMARY KEY,
    track_number VARCHAR(50) NOT NULL,
    entry VARCHAR(10) NOT NULL,
    locale VARCHAR(2) NOT NULL,
    internal_signature VARCHAR(100),
    customer_id VARCHAR(50) NOT NULL,
    delivery_service VARCHAR(50) NOT NULL,
    shardkey VARCHAR(10) NOT NULL,
    sm_id INTEGER NOT NULL,
    date_created TIMESTAMPTZ NOT NULL,
    oof_shard VARCHAR(10) NOT NULL
);


CREATE TABLE IF NOT EXISTS delivery (
    order_uid UUID PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(16) NOT NULL,
    zip VARCHAR(10) NOT NULL,
    city VARCHAR(100) NOT NULL,
    address VARCHAR(100) NOT NULL,
    region VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL
);


CREATE TABLE IF NOT EXISTS payment (
    order_uid UUID PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    transaction VARCHAR(100) NOT NULL,
    request_id VARCHAR(100) DEFAULT '',
    currency VARCHAR(3) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    amount INTEGER NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(20) NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total INTEGER NOT NULL,
    custom_fee INTEGER NOT NULL
);


CREATE TABLE IF NOT EXISTS item (
    id SERIAL PRIMARY KEY,
    order_uid UUID REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id INTEGER NOT NULL,
    track_number VARCHAR(50) NOT NULL,
    price INTEGER NOT NULL,
    rid VARCHAR(50) NOT NULL,
    name VARCHAR(50) NOT NULL,
    sale INTEGER NOT NULL,
    size VARCHAR(10) NOT NULL,
    total_price INTEGER NOT NULL,
    nm_id INTEGER NOT NULL,
    brand VARCHAR(50) NOT NULL,
    status INTEGER NOT NULL
);


CREATE INDEX IF NOT EXISTS idx_orders_order_uid ON orders(order_uid);
CREATE INDEX IF NOT EXISTS idx_delivery_order_uid ON delivery(order_uid);
CREATE INDEX IF NOT EXISTS idx_payment_order_uid ON payment(order_uid);
CREATE INDEX IF NOT EXISTS idx_item_order_uid ON item(order_uid);




INSERT INTO orders (
    order_uid,
    track_number,
    entry,
    locale,
    internal_signature,
    customer_id,
    delivery_service,
    shardkey,
    sm_id,
    date_created,
    oof_shard
) VALUES (
    '8864b7f1-c455-4300-bfdc-d339429c2099',
    'WBILMTESTTRACK',
    'WBIL',
    'en',
    '',
    'test',
    'meest',
    '9',
    99,
    '2021-11-26T06:22:19Z',
    '1'
)
ON CONFLICT (order_uid) DO NOTHING;


INSERT INTO delivery (
    order_uid,
    name,
    phone,
    zip,
    city,
    address,
    region,
    email
) VALUES (
    '8864b7f1-c455-4300-bfdc-d339429c2099',
    'Test Testov',
    '+9720000000',
    '2639809',
    'Kiryat Mozkin',
    'Ploshad Mira 15',
    'Kraiot',
    'test@gmail.com'
)
ON CONFLICT DO NOTHING;


INSERT INTO payment (
    order_uid,
    transaction,
    request_id,
    currency,
    provider,
    amount,
    payment_dt,
    bank,
    delivery_cost,
    goods_total,
    custom_fee
) VALUES (
    '8864b7f1-c455-4300-bfdc-d339429c2099',
    '8864b7f1-c455-4300-bfdc-d339429c2099',
    '',
    'USD',
    'wbpay',
    1817,
    1637907727,
    'alpha',
    1500,
    317,
    0
)
ON CONFLICT DO NOTHING;


INSERT INTO item (
    order_uid,
    chrt_id,
    track_number,
    price,
    rid,
    name,
    sale,
    size,
    total_price,
    nm_id,
    brand,
    status
) VALUES (
    '8864b7f1-c455-4300-bfdc-d339429c2099',
    9934930,
    'WBILMTESTTRACK',
    453,
    'ab4219087a764ae0btest',
    'Mascaras',
    30,
    '0',
    317,
    2389212,
    'Vivienne Sabo',
    202
)
ON CONFLICT DO NOTHING;