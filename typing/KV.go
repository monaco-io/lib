package typing

type KV[K comparable, V any] struct {
	k K
	v V
}

func NewKV[K comparable, V any](key K, value V) KV[K, V] {
	return KV[K, V]{
		k: key,
		v: value,
	}
}

func (kv KV[K, V]) Get() (K, V) {
	return kv.k, kv.v
}
