// This defines the Cache Interface,
// allowing you to swap between an In-Memory cache and an external store like Redis seamlessly.

package domain

import (
	"context"
	"time"
)

// CacheRepository defines behavior for caching operations.
type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) (bool, error)
	Delete(ctx context.Context, key string) error
}
