-- Таблица курсов валют
CREATE TABLE IF NOT EXISTS exchange_rates (
    id SERIAL PRIMARY KEY,
    from_currency VARCHAR(3) NOT NULL,
    to_currency   VARCHAR(3) NOT NULL,
    rate          NUMERIC(12,6) NOT NULL,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (from_currency, to_currency)
    );

-- Начальные значения (USD, EUR, RUB)
INSERT INTO exchange_rates (from_currency, to_currency, rate) VALUES
    ('USD', 'EUR', 0.92),
    ('EUR', 'USD', 1.08),
    ('USD', 'RUB', 94.30),
    ('RUB', 'USD', 0.01),
    ('EUR', 'RUB', 102.50),
    ('RUB', 'EUR', 0.25)
    ON CONFLICT (from_currency, to_currency) DO UPDATE SET rate = EXCLUDED.rate;
