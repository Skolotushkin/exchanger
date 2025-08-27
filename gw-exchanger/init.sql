CREATE TABLE rates (
    currency TEXT PRIMARY KEY,
    rate     FLOAT NOT NULL
);

-- Добавим данные
INSERT INTO rates (currency, rate) VALUES
    ('USD', 1.0),        -- базовая валюта
    ('RUB', 93.5),
    ('EUR', 0.92);