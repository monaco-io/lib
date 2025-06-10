package objectpool

import (
	"sync"
)

type ObjectPool[T any] struct {
	pool *sync.Pool
}

func New[T any]() *ObjectPool[T] {
	return &ObjectPool[T]{
		pool: &sync.Pool{
			New: func() any {
				return new(T)
			},
		},
	}
}

func (op *ObjectPool[T]) Get() *T {
	return op.pool.Get().(*T)
}

func (op *ObjectPool[T]) Put(obj *T) {
	op.pool.Put(obj)
}
