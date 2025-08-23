package option

type Option[T any] func(*T)

type Options[T any] []Option[T]

func (opts Options[T]) Apply(o *T) {
	for _, opt := range opts {
		opt(o)
	}
}

func Apply[T any](opts []Option[T], o *T) {
	for _, opt := range opts {
		opt(o)
	}
}
