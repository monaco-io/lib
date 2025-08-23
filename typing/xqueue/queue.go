package xqueue

import (
	"log"
	"sync"

	"github.com/monaco-io/lib/typing/xopt"
	"golang.org/x/sync/errgroup"
)

type Queue[T any] interface {
	Input(item ...T)
	Close()
	CloseSync()
}

const (
	defaultBufferSize = 1 << 10
	defaultMaxProps   = 1 << 4
)

type Config struct {
	maxProps   int
	errHandler func(err error)
}

func WithMaxProps(n int) xopt.Option[Config] {
	return func(o *Config) {
		o.maxProps = n
	}
}

func WithErrorHandler(fn func(err error)) xopt.Option[Config] {
	return func(o *Config) {
		o.errHandler = fn
	}
}

type queue[T any] struct {
	ch    chan T
	errCh chan error

	consumerHandler func(data T) error
	errHandler      func(err error)

	eg       errgroup.Group
	maxProps int
	once     sync.Once
}

func New[T any](consumer func(data T) error, opts ...xopt.Option[Config]) Queue[T] {
	cfg := Config{
		maxProps: defaultMaxProps,
		errHandler: func(err error) {
			log.Printf("lib.queue Error occurred: %v\n", err)
		},
	}
	xopt.Apply(opts, &cfg)
	q := queue[T]{
		ch:              make(chan T, defaultBufferSize),
		errCh:           make(chan error),
		consumerHandler: consumer,
		errHandler:      cfg.errHandler,

		maxProps: cfg.maxProps,
	}
	q.eg.SetLimit(q.maxProps)
	return &q
}

func (q *queue[T]) SetMaxProps(n int) {
	q.maxProps = n
}

func (q *queue[T]) Input(item ...T) {
	q.once.Do(q.dequeue)
	for _, v := range item {
		q.ch <- v
	}
}

func (q *queue[T]) Close() {
	go q.CloseSync()
}

func (q *queue[T]) CloseSync() {
	_ = q.eg.Wait()
	close(q.ch)
	close(q.errCh)
}

func (q *queue[T]) dequeue() {
	q.eg.Go(func() error {
		for err := range q.errCh {
			q.errHandler(err)
		}
		return nil
	})
	q.eg.Go(func() error {
		for data := range q.ch {
			if err := q.consumerHandler(data); err != nil {
				q.errCh <- err
			}
		}
		return nil
	})
}
