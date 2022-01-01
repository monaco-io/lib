package retry

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Options struct {
	Context    context.Context
	RetryTimes int
	Delay      time.Duration
}

func Do(f func(context.Context) error, opts Options) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			// find out exactly what the error was and set err
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				// Fallback err (per specs, error strings should be lowercase w/o punctuation
				err = fmt.Errorf("unknown panic: %v", r)
			}
		}
	}()

	if opts.RetryTimes == 0 {
		opts.RetryTimes = 1
	}
	if opts.Context == nil {
		opts.Context = context.Background()
	}

	for times := 0; times < opts.RetryTimes; times++ {
		if times != 0 {
			time.Sleep(opts.Delay)
		}
		err = f(opts.Context)
		if err == nil {
			break
		}
	}
	return
}
