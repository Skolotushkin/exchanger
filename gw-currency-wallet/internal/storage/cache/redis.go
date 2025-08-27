package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const cacheKeyRates = "exchange:rates"

type RedisCache struct {
	client *redis.Client
	logger *zap.Logger
}

func NewRedisCache(addr /*, password*/ string, db int, logger *zap.Logger) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		// Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second, // таймаут подключения
		ReadTimeout:  3 * time.Second, // таймаут чтения
		WriteTimeout: 3 * time.Second, // таймаут записи
		PoolSize:     10,              // размер пула соединений
		MinIdleConns: 2,               // минимальное количество idle соединений
	})
	// Проверяем соединение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	logger.Info("Redis cache initialized",
		zap.String("addr", addr),
		zap.Int("db", db))

	return &RedisCache{client: client, logger: logger}, nil
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}

func (r *RedisCache) GetRates(ctx context.Context) (map[string]string, error) {
	raw, err := r.client.Get(ctx, cacheKeyRates).Result()
	if err != nil {
		return nil, err
	}

	var rates map[string]string
	if err := json.Unmarshal([]byte(raw), &rates); err != nil {
		return nil, err
	}
	return rates, nil
}

func (r *RedisCache) SetRates(ctx context.Context, rates map[string]string, ttl time.Duration) error {
	data, err := json.Marshal(rates)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, cacheKeyRates, string(data), ttl).Err()
}
