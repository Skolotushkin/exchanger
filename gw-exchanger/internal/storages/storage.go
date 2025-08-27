package storages

import "context"

type Storage interface {
	GetAllRates(ctx context.Context) (map[string]float32, error)
	GetRate(ctx context.Context, currency string) (float32, error)
}
