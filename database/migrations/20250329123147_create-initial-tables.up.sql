CREATE TABLE offers (
    id SERIAL PRIMARY KEY,
    final_price DECIMAL NOT NULL,
    carrier VARCHAR(255) NOT NULL,
    service VARCHAR(255) NOT NULL,
    delivery_time INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CHECK (delivery_time > 0)
);