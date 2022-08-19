package lru

import (
	"container/list"
	"time"
)

func (c *Cache[K, V]) moveToFront(ele *list.Element) {
	c.lock.Lock()
	c.cache.MoveToFront(ele)
	c.lock.Unlock()
}

func (c *Cache[K, V]) pushFront(v any) (ele *list.Element) {
	c.lock.Lock()
	ele = c.cache.PushFront(v)
	c.lock.Unlock()
	return
}

func (c *Cache[K, V]) remove(e *list.Element) {
	c.lock.Lock()
	c.cache.Remove(e)
	c.lock.Unlock()
	ele := e.Value.(*entry[K, V])
	c.hash.Delete(ele.key)
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
