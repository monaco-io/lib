package typing

import (
	"sync"
)

type SyncMap[K comparable, V any] struct {
	syncmap *sync.Map
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{syncmap: &sync.Map{}}
}

func (m SyncMap[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.syncmap.Load(key)
	if !ok {
		return
	}
	value, ok = v.(V)
	return
}

func (m SyncMap[K, V]) Store(key K, value V) {
	m.syncmap.Store(key, value)
}

func (m SyncMap[K, V]) Delete(key K) {
	m.syncmap.Delete(key)
}

func (m SyncMap[K, V]) Range(f func(key K, value V) bool) {
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
