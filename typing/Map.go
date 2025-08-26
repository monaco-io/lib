package typing

import (
	"github.com/samber/lo"
)

type Map[K comparable, V any] map[K]V

type MapX Map[any, any]

func (m Map[K, V]) Get(key K) (V, bool) {
	value, exists := m[key]
	return value, exists
}

func (m Map[K, V]) Set(key K, value V) {
	m[key] = value
}

func (m Map[K, V]) Keys() []K {
	return lo.Keys(m)
}

func (m Map[K, V]) HasKey(key K) bool {
	return lo.HasKey(m, key)
}

func (m Map[K, V]) Values() []V {
	return lo.Values(m)
}
