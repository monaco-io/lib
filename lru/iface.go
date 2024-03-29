package lru

import (
	"context"
	"sync"
	"time"

	"github.com/monaco-io/lib/list"
	"github.com/monaco-io/lib/syncmap"
)

// ICache is the interface for thread safe LRU cache.
type ICache[K comparable, V any] interface {
	// Sets a value to the cache, returns true if an eviction occurred and
	// updates the "recently used"-ness of the key.
	Set(K, V)

	// Returns key's value from the cache and
	// updates the "recently used"-ness of the key. #value, isFound
	Get(K) (V, bool)

	// Removes a key from the cache.
	Remove(K)

	Len() int
	Flush()
}

// New creates a new Cache.
// If maxEntries is zero, the cache has no limit and it's assumed
// that eviction is done by the caller.
func New[K comparable, V any](limit int, ttl time.Duration) ICache[K, V] {
	c := Cache[K, V]{
		limit: defaultLength,
		ttl:   defaultTTL,
		cache: list.New[*entry[K, V]](),
		hash:  syncmap.New[K, *list.Element[*entry[K, V]]](),
		lock:  new(sync.Mutex),
	}
	if limit != 0 {
		c.limit = limit
	}
	if ttl != 0 {
		c.ttl = ttl
	}
	return &c
}

type ICacheC[K comparable, V any] interface {
	ICache[K, V]
	GetC(context.Context, K) (V, error)
}

func NewC[K comparable, V any](limit int, ttl time.Duration, cb func(context.Context, K) (V, error)) ICacheC[K, V] {
	c := CacheC[K, V]{
		ICache: New[K, V](limit, ttl),
		cb:     cb,
	}

	return &c
}
