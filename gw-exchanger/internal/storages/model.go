package storages

type Rate struct {
	Currency string  `db:"currency"`
	Rate     float32 `db:"rate"`
}
