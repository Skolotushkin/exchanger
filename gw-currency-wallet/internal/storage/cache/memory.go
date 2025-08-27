package cache

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

type InMemoryCache struct {
	data  map[string]item
	mutex sync.RWMutex
}

type item struct {
	value      string
	expiration int64
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]item),
	}
}

func (m *InMemoryCache) Get(ctx context.Context, key string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	it, ok := m.data[key]
	if !ok || (it.expiration > 0 && time.Now().UnixNano() > it.expiration) {
		return "", errors.New("cache: key not found")
	}
	return it.value, nil
}

func (m *InMemoryCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	exp := int64(0)
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}
	m.data[key] = item{value: value, expiration: exp}
	return nil
}

func (m *InMemoryCache) Close() error { return nil }

func (m *InMemoryCache) GetRates(ctx context.Context) (map[string]string, error) {
	raw, err := m.Get(ctx, cacheKeyRates)
	if err != nil {
		return nil, err
	}
	var rates map[string]string
	if err := json.Unmarshal([]byte(raw), &rates); err != nil {
		return nil, err
	}
	return rates, nil
}

func (m *InMemoryCache) SetRates(ctx context.Context, rates map[string]string, ttl time.Duration) error {
	data, err := json.Marshal(rates)
	if err != nil {
		return err
	}
	return m.Set(ctx, cacheKeyRates, string(data), ttl)
}
