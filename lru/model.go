package lru

import (
	"container/list"
	"context"
	"math"
	"sync"
	"time"

	"github.com/monaco-io/lib/syncmap"
)

const (
	defaultTTL    = time.Second * 10
	defaultLength = math.MaxInt16
)

// Cache is an LRU cache. It is not safe for concurrent access.
type Cache[K comparable, V any] struct {
	// limit is the maximum number of cache entries before
	// an item is evicted. Zero means no limit.
	limit int

	// expire time
	ttl time.Duration

	cache *list.List
	hash  *syncmap.Map[K, *list.Element]
	lock  sync.Locker
}

type entry[K comparable, V any] struct {
	key    K
	value  V
	expire time.Time
}

type CacheC[K comparable, V any] struct {
	ICache[K, V]
	cb func(context.Context, K) (V, error)
}
