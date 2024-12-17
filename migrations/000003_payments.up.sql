CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    transaction VARCHAR(255),
    request_id VARCHAR(255),
    currency VARCHAR(10),
    provider VARCHAR(50),
    amount DECIMAL(12, 2),
    payment_dt BIGINT,
    bank VARCHAR(100),
    delivery_cost DECIMAL(12, 2),
    goods_total DECIMAL(12, 2),
    custom_fee DECIMAL(12, 2), 
    order_uid VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE
);
