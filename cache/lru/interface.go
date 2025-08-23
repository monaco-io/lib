package lru

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/monaco-io/lib/typing"
	"github.com/monaco-io/lib/typing/option"
)

var ErrMiss = errors.New("lib.cache.lru:miss")

func IsErrMiss(err error) bool {
	return errors.Is(err, ErrMiss)
}

// ICache is the interface for thread safe LRU cache.
type ICache[K comparable, V any] interface {
	// Sets a value to the cache, returns true if an eviction occurred and
	// updates the "recently used"-ness of the key.
	Set(K, V)

	Get(context.Context, K) (V, error)

	// Removes a key from the cache.
	Remove(K)

	Len() int
	Flush()
}
type config[K comparable, V any] struct {
	limit      uint
	ttl        time.Duration
	sourceFunc func(context.Context, K) (V, error)
}

func WithLimit[K comparable, V any](limit uint) option.Option[config[K, V]] {
	return func(cfg *config[K, V]) {
		cfg.limit = limit
	}
}

func WithTTL[K comparable, V any](ttl time.Duration) option.Option[config[K, V]] {
	return func(cfg *config[K, V]) {
		cfg.ttl = ttl
	}
}

func WithSourceFunc[K comparable, V any](sf func(context.Context, K) (V, error)) option.Option[config[K, V]] {
	return func(cfg *config[K, V]) {
		cfg.sourceFunc = sf
	}
}

// New creates a new Cache.
// If maxEntries is zero, the cache has no limit and it's assumed
// that eviction is done by the caller.
func New[K comparable, V any](opts ...option.Option[config[K, V]]) ICache[K, V] {
	cfg := config[K, V]{
		limit: defaultLength,
		ttl:   defaultTTL,
	}
	option.Apply(opts, &cfg)
	c := lru[K, V]{
		size: int(cfg.limit),
		ttl:  cfg.ttl,
		cb:   cfg.sourceFunc,

		data: typing.NewLinkedList[*entry[K, V]](),
		hash: typing.NewSyncMap[K, *typing.Element[*entry[K, V]]](),
		lock: new(sync.Mutex),
	}
	return &c
}
