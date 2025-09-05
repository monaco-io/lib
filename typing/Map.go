package typing

import (
	"fmt"
	"net/url"

	"github.com/samber/lo"
)

type (
	Map[K comparable, V any] map[K]V
	MapX                     Map[string, any]
)

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

func (m Map[K, V]) URLValues() url.Values {
	values := url.Values{}
	for k, v := range m {
		values.Add(fmt.Sprintf("%v", k), fmt.Sprintf("%v", v))
	}
	return values
}

func (m Map[K, V]) URLEncode() string {
	return m.URLValues().Encode()
}

func (m MapX) URLValues() url.Values {
	values := url.Values{}
	for k, v := range m {
		values.Add(k, fmt.Sprintf("%v", v))
	}
	return values
}

func (m MapX) URLEncode() string {
	return m.URLValues().Encode()
}
