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
