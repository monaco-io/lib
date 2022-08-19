package lru

import (
	"context"

	"github.com/monaco-io/lib/list"
	"github.com/monaco-io/lib/syncmap"
	"github.com/pkg/errors"
)

// Set adds a value to the cache.
func (c *Cache[K, V]) Set(key K, value V) {
	if c.hash == nil {
		c.hash = syncmap.New[K, *list.Element[*entry[K, V]]]()
		c.cache = list.New[*entry[K, V]]()
	}
	if ee, ok := c.hash.Load(key); ok {
		c.moveToFront(ee)
		ee.Value.value = value
		return
	}
	expire := now().Add(c.ttl)
	ele := c.pushFront(&entry[K, V]{key, value, expire})
	c.hash.Store(key, ele)
	if c.limit != 0 && c.cache.Len() > c.limit {
		c.removeOldest()
	}
}

// Get looks up a key's value from the cache.
func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
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
func (c *Cache[K, V]) Remove(key K) {
	if c.hash == nil {
		return
	}
	if ele, hit := c.hash.Load(key); hit {
		c.remove(ele)
	}
}

// Len returns the number of items in the cache.
func (c *Cache[K, V]) Len() int {
	if c.hash == nil {
		return 0
	}
	return c.cache.Len()
}

// Clear purges all stored items from the cache.
func (c *Cache[K, V]) Flush() {
	c.cache = list.New[*entry[K, V]]()
	c.hash = syncmap.New[K, *list.Element[*entry[K, V]]]()
}

// Clear purges all stored items from the cache.
func (c *CacheC[K, V]) GetC(ctx context.Context, key K) (value V, err error) {
	value, ok := c.Get(key)
	if ok {
		return
	}
	value, err = c.cb(ctx, key)
	if err != nil {
		err = errors.Wrap(err, "callback")
		return
	}
	c.Set(key, value)
	return
}
