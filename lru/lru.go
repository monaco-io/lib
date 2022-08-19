package lru

import (
	"container/list"
	"time"

	"github.com/monaco-io/lib/syncmap"
)

// Add adds a value to the cache.
func (c *Cache[K, V]) Add(key K, value V) {
	if c.hash == nil {
		c.hash = syncmap.New[K, *list.Element]()
		c.cache = list.New()
	}
	if ee, ok := c.hash.Load(key); ok {
		c.moveToFront(ee)
		ee.Value.(*entry[K, V]).value = value
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
		data := ele.Value.(*entry[K, V])
		if now().After(data.expire) {
			return
		}
		c.moveToFront(ele)
		return data.value, true
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

// Flush purges all stored items from the cache.
func (c *Cache[K, V]) Flush() {
	c.cache = list.New()
	c.hash = syncmap.New[K, *list.Element]()
}

func now() time.Time {
	return time.Now()
}
