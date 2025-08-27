package cache

import (
	"context"
	"time"
)

type Cache interface {
	// Базовые методы (если вдруг нужны где-то ещё)
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Close() error

	// Специализированные методы для курса валют
	GetRates(ctx context.Context) (map[string]string, error)
	SetRates(ctx context.Context, rates map[string]string, ttl time.Duration) error
}
