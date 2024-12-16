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