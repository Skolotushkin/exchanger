package storages

import "time"

type Rate struct {
	ID           int       `db:"id"`
	CurrencyFrom string    `db:"currency_from"`
	CurrencyTo   string    `db:"currency_to"`
	Rate         float32   `db:"rate"`
	UpdatedAt    time.Time `db:"updated_at"`
}
