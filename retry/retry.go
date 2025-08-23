package retry

import (
	"fmt"
	"time"
)

type Options struct {
	RetryTimes int
	Delay      time.Duration
}

func Do(f func() error, opts Options) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			err = fmt.Errorf("panic recoverd: %v", r)
		}
	}()

	if opts.RetryTimes == 0 {
		opts.RetryTimes = 1
	}

	for times := 0; times < opts.RetryTimes; times++ {
		if times != 0 {
			time.Sleep(opts.Delay)
		}
		if err = f(); err == nil {
			break
		}
	}
	return
}
