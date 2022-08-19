package lru

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
	Clear()
}
