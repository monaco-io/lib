package lru

import (
	"container/list"
	"math"
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

	ll    *list.List
	cache *syncmap.Map[K, *list.Element]
}

type entry[K comparable, V any] struct {
	key    K
	value  V
	expire time.Time
}

// New creates a new Cache.
// If maxEntries is zero, the cache has no limit and it's assumed
// that eviction is done by the caller.
func New[K comparable, V any](limit int, ttl time.Duration) *Cache[K, V] {
	c := Cache[K, V]{
		limit: defaultLength,
		ttl:   defaultTTL,
		ll:    list.New(),
		cache: syncmap.New[K, *list.Element](),
	}
	if limit != 0 {
		c.limit = limit
	}
	if ttl != 0 {
		c.ttl = ttl
	}
	return &c
}

// Add adds a value to the cache.
func (c *Cache[K, V]) Add(key K, value V) {
	if c.cache == nil {
		c.cache = syncmap.New[K, *list.Element]()
		c.ll = list.New()
	}
	if ee, ok := c.cache.Load(key); ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry[K, V]).value = value
		return
	}
	expire := now().Add(c.ttl)
	ele := c.ll.PushFront(&entry[K, V]{key, value, expire})
	c.cache.Store(key, ele)
	if c.limit != 0 && c.ll.Len() > c.limit {
		c.RemoveOldest()
	}
}

// Get looks up a key's value from the cache.
func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache.Load(key); hit {
		c.ll.MoveToFront(ele)
		data := ele.Value.(*entry[K, V])
		if now().After(data.expire) {
			return
		}
		return data.value, true
	}
	return
}

// Remove removes the provided key from the cache.
func (c *Cache[K, V]) Remove(key K) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache.Load(key); hit {
		c.removeElement(ele)
	}
}

// RemoveOldest removes the oldest item from the cache.
func (c *Cache[K, V]) RemoveOldest() {
	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *Cache[K, V]) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entry[K, V])
	c.cache.Delete(kv.key)
}

// Len returns the number of items in the cache.
func (c *Cache[K, V]) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

// Clear purges all stored items from the cache.
func (c *Cache[K, V]) Clear() {
	c.ll = nil
	c.cache = syncmap.New[K, *list.Element]()
}

func now() time.Time {
	return time.Now()
}
