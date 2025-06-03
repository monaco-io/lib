package syncmap

import (
	"sync"
)

type Map[K comparable, V any] struct {
	syncmap *sync.Map
}

func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{syncmap: &sync.Map{}}
}

func (m Map[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.syncmap.Load(key)
	if !ok {
		return
	}
	value, ok = v.(V)
	return
}

func (m Map[K, V]) Store(key K, value V) {
	m.syncmap.Store(key, value)
}

func (m Map[K, V]) Delete(key K) {
	m.syncmap.Delete(key)
}

func (m Map[K, V]) Range(f func(key K, value V) bool) {
	m.syncmap.Range(func(key, value any) bool {
		k, ok := key.(K)
		if !ok {
			return false
		}
		v, ok := value.(V)
		if !ok {
			return false
		}
		return f(k, v)
	})
}
