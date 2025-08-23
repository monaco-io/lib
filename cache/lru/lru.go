package lru

import (
	"context"
	"math"
	"sync"
	"time"

	x "github.com/monaco-io/lib/typing"
)

const (
	defaultTTL    = time.Second * 1 << 6
	defaultLength = math.MaxInt16
)

// lru is an LRU lru. It is not safe for concurrent access.
type lru[K comparable, V any] struct {
	// size is the maximum number of cache entries before
	// an item is evicted. Zero means no size.
	size int

	// expire time
	ttl time.Duration
	cb  func(context.Context, K) (V, error)

	data *x.LinkedList[*entry[K, V]]
	hash *x.SyncMap[K, *x.Element[*entry[K, V]]]
	lock sync.Locker
}

type entry[K comparable, V any] struct {
	key    K
	value  V
	expire time.Time
}

// Set adds a value to the cache.
func (c *lru[K, V]) Set(key K, value V) {
	if c.hash == nil {
		c.hash = x.NewSyncMap[K, *x.Element[*entry[K, V]]]()
		c.data = x.NewLinkedList[*entry[K, V]]()
	}
	if ee, ok := c.hash.Load(key); ok {
		c.moveToFront(ee)
		ee.Value.value = value
		return
	}
	expire := now().Add(c.ttl)
	ele := c.pushFront(&entry[K, V]{key, value, expire})
	c.hash.Store(key, ele)
	if c.size != 0 && c.data.Len() > c.size {
		c.removeOldest()
	}
}

// Get looks up a key's value from the cache.
func (c *lru[K, V]) get(key K) (value V, ok bool) {
	if c.hash == nil {
		return
	}
	if ele, hit := c.hash.Load(key); hit {
		if now().After(ele.Value.expire) {
			return
		}
		c.moveToFront(ele)
		return ele.Value.value, true
	}
	return
}

// Remove removes the provided key from the cache.
func (c *lru[K, V]) Remove(key K) {
	if c.hash == nil {
		return
	}
	if ele, hit := c.hash.Load(key); hit {
		c.remove(ele)
	}
}

// Len returns the number of items in the cache.
func (c *lru[K, V]) Len() int {
	if c.hash == nil {
		return 0
	}
	return c.data.Len()
}

// Clear purges all stored items from the cache.
func (c *lru[K, V]) Flush() {
	c.data = x.NewLinkedList[*entry[K, V]]()
	c.hash = x.NewSyncMap[K, *x.Element[*entry[K, V]]]()
}

// Clear purges all stored items from the cache.
func (c *lru[K, V]) Get(ctx context.Context, key K) (dv V, _ error) {
	if v, ok := c.get(key); ok {
		return v, nil
	}
	if c.cb != nil {
		if v, err := c.cb(ctx, key); err != nil {
			return dv, err
		} else {
			c.Set(key, v)
			return v, nil
		}
	}
	return dv, ErrMiss
}
