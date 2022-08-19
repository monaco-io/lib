package lru

import (
	"time"

	"github.com/monaco-io/lib/list"
)

func (c *Cache[K, V]) moveToFront(ele *list.Element[*entry[K, V]]) {
	c.lock.Lock()
	c.cache.MoveToFront(ele)
	c.lock.Unlock()
}

func (c *Cache[K, V]) pushFront(v *entry[K, V]) (ele *list.Element[*entry[K, V]]) {
	c.lock.Lock()
	ele = c.cache.PushFront(v)
	c.lock.Unlock()
	return
}

func (c *Cache[K, V]) remove(e *list.Element[*entry[K, V]]) {
	c.lock.Lock()
	c.cache.Remove(e)
	c.lock.Unlock()
	c.hash.Delete(e.Value.key)
}

// removeOldest removes the oldest item from the cache.
func (c *Cache[K, V]) removeOldest() {
	if c.hash == nil {
		return
	}
	ele := c.cache.Back()
	if ele != nil {
		c.remove(ele)
	}
}

func now() time.Time {
	return time.Now()
}
