package lru

import (
	"time"

	"github.com/monaco-io/lib/typing"
)

func (c *lru[K, V]) moveToFront(ele *typing.Element[*entry[K, V]]) {
	c.lock.Lock()
	c.data.MoveToFront(ele)
	c.lock.Unlock()
}

func (c *lru[K, V]) pushFront(v *entry[K, V]) (ele *typing.Element[*entry[K, V]]) {
	c.lock.Lock()
	ele = c.data.PushFront(v)
	c.lock.Unlock()
	return
}

func (c *lru[K, V]) remove(e *typing.Element[*entry[K, V]]) {
	c.lock.Lock()
	c.data.Remove(e)
	c.lock.Unlock()
	c.hash.Delete(e.Value.key)
}

// removeOldest removes the oldest item from the cache.
func (c *lru[K, V]) removeOldest() {
	if c.hash == nil {
		return
	}
	ele := c.data.Back()
	if ele != nil {
		c.remove(ele)
	}
}

func now() time.Time {
	return time.Now()
}
