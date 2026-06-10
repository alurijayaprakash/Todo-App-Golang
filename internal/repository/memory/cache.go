// A lightweight in-memory fallback cache implementing CacheRepository.
package memory

import (
	"context"
	"sync"
	"time"
	"todo-api/internal/handler/domain"
)

type item struct {
	value      interface{}
	expiration int64
}

type MemoryCache struct {
	mu    sync.RWMutex
	store map[string]item
}

func NewMemoryCache() domain.CacheRepository {
	return &MemoryCache{
		store: make(map[string]item),
	}
}

func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = item{
		value:      value,
		expiration: time.Now().Add(ttl).UnixNano(),
	}
	return nil
}

func (c *MemoryCache) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	it, exists := c.store[key]
	if !exists || time.Now().UnixNano() > it.expiration {
		return false, nil
	}

	// Direct assignment since it's an internal pointer copy for in-memory testing
	*(dest.(*interface{})) = it.value
	return true, nil
}

func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
	return nil
}
