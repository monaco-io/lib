package lru

import (
	"container/list"
	"sync"
	"time"

	"github.com/monaco-io/lib/syncmap"
)

// ICache is the interface for thread safe LRU cache.
type ICache[K comparable, V any] interface {
	// Adds a value to the cache, returns true if an eviction occurred and
	// updates the "recently used"-ness of the key.
	Add(key K, value V)

	// Returns key's value from the cache and
	// updates the "recently used"-ness of the key. #value, isFound
	Get(key K) (value V, ok bool)

	// Removes a key from the cache.
	Remove(key K)

	Len() int
	Flush()
}

// New creates a new Cache.
// If maxEntries is zero, the cache has no limit and it's assumed
// that eviction is done by the caller.
func New[K comparable, V any](limit int, ttl time.Duration) *Cache[K, V] {
	c := Cache[K, V]{
		limit: defaultLength,
		ttl:   defaultTTL,
		cache: list.New(),
		hash:  syncmap.New[K, *list.Element](),
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
