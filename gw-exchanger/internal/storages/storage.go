package storages

import "context"

type Storage interface {
	GetRate(ctx context.Context, from, to string) (float32, error)
	GetAllRates(ctx context.Context) (map[string]float32, error)
	Close() error
}
