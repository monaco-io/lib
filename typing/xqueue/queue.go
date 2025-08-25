package xqueue

import (
	"log"
	"sync"

	"github.com/monaco-io/lib/typing/xopt"
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
	ch chan T

	consumerHandler func(data T) error
	errHandler      func(err error)

	maxProps int
	once     sync.Once
	closed   bool
	mu       sync.RWMutex
	wg       sync.WaitGroup
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
		consumerHandler: consumer,
		errHandler:      cfg.errHandler,

		maxProps: cfg.maxProps,
	}
	return &q
}

func (q *queue[T]) Input(item ...T) {
	q.once.Do(q.dequeue)

	for _, v := range item {
		q.mu.RLock()
		if q.closed {
			q.mu.RUnlock()
			return
		}
		q.mu.RUnlock()
		q.ch <- v
	}
}

func (q *queue[T]) Close() {
	go q.CloseSync()
}

func (q *queue[T]) CloseSync() {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return
	}
	q.closed = true

	close(q.ch) // 关闭数据通道
	q.wg.Wait() // 等待错误处理goroutine完成
}

func (q *queue[T]) dequeue() {
	for range q.maxProps {
		q.wg.Go(func() {
			for data := range q.ch {
				if err := q.consumerHandler(data); err != nil {
					q.errHandler(err)
				}
			}
		})
	}
}
