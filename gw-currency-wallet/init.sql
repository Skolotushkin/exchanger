-- включаем расширение pgcrypto, для генерации UUID и хеширования паролей
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS wallets (
                                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    currency CHAR(3) NOT NULL,
    balance NUMERIC(18,2) DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, currency)
    );

CREATE TABLE IF NOT EXISTS transactions (
                                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL,
    amount NUMERIC(18,2) NOT NULL,
    exchanged_amount NUMERIC(18,2),
    from_currency CHAR(3),
    to_currency CHAR(3),
    currency CHAR(3) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
    );

-- Тестовый пользователь (исправлено название колонки)
INSERT INTO users (email, password)
VALUES ('test@example.com', '$2a$10$Vb0RrRsmVxv8Y4x/TxyaSehlHl3P0hoxY6gQKz4wz0cMCEc9Uzvve')
    ON CONFLICT (email) DO NOTHING;

-- У него кошелёк в USD (исправлен user_id)
INSERT INTO wallets (user_id, currency, balance)
SELECT
    u.id,
    'USD',
    10000
FROM users u
WHERE u.email = 'test@example.com'
    ON CONFLICT (user_id, currency) DO NOTHING;