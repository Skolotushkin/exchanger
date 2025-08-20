CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS wallets (
                                       id SERIAL PRIMARY KEY,
                                       user_id INT REFERENCES users(id) ON DELETE CASCADE,
    currency VARCHAR(3) NOT NULL,
    balance NUMERIC(18,2) DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, currency)
    );

CREATE TABLE IF NOT EXISTS transactions (
                                            id SERIAL PRIMARY KEY,
                                            user_id INT REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL, -- deposit, withdraw, exchange
    amount NUMERIC(18,2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
    );

-- Тестовый пользователь
INSERT INTO users (email, password_hash)
VALUES ('test@example.com', '$2a$10$Vb0RrRsmVxv8Y4x/TxyaSehlHl3P0hoxY6gQKz4wz0cMCEc9Uzvve')
    ON CONFLICT DO NOTHING;

-- У него кошелёк в USD
INSERT INTO wallets (user_id, currency, balance)
VALUES (1, 'USD', 10000)
    ON CONFLICT DO NOTHING;