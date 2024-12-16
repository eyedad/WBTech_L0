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
