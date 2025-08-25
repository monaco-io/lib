package xopt

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

func Boolean(v ...bool) bool {
	if len(v) > 0 {
		return v[0]
	}
	return false
}
