package retry

import (
	"fmt"
	"time"

	"github.com/monaco-io/lib/typing/xopt"
)

type config struct {
	RetryTimes int
	Delay      time.Duration
}

func WithRetryTimes(times int) xopt.Option[config] {
	return func(cfg *config) {
		cfg.RetryTimes = times
	}
}

func WithDelay(delay time.Duration) xopt.Option[config] {
	return func(cfg *config) {
		cfg.Delay = delay
	}
}

func Do(f func() error, opts ...xopt.Option[config]) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			err = fmt.Errorf("panic recoverd: %v", r)
		}
	}()
	var cfg config
	xopt.Apply(opts, &cfg)
	for times := 0; times < cfg.RetryTimes; times++ {
		if times != 0 {
			time.Sleep(cfg.Delay)
		}
		if err = f(); err == nil {
			break
		}
	}
	return
}
