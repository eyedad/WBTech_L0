CREATE TABLE orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(100),
    entry VARCHAR(50),
    customer_id VARCHAR(100),
    delivery_service VARCHAR(100),
    shardkey VARCHAR(20),                           
    sm_id INT,                             
    date_created TIMESTAMP WITH TIME ZONE,  
    oof_shard VARCHAR(20),                          
    locale VARCHAR(10),                     
    internal_signature TEXT                 
);

CREATE TABLE deliveries (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200),
    phone VARCHAR(20),
    zip VARCHAR(20),
    city VARCHAR(100),
    address VARCHAR(250),
    region VARCHAR(100),
    email VARCHAR(100),
    order_uid VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE
);

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

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    chrt_id INT,
    track_number VARCHAR(100),
    price DECIMAL(12, 2),
    rid VARCHAR(255),
    name VARCHAR(255),
    sale INT,
    size VARCHAR(50),
    total_price DECIMAL(12, 2),
    nm_id INT,
    brand VARCHAR(255),
	status INT,
	order_uid VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE
);