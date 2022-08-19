# lib/lru

thread safe LRU cache.

## example

```go
package main

import (
	"log"

	"github.com/monaco-io/lib/lru"
)

var instance = lru.New[int, int](100, time.Second*5)

func main() {
	i := 1
	//set key value pairs
	instance.Set(i, i)

	//get val from mem
	v, ok := instance.Get(i)
	log.Println(i, v, ok)

}

```

## example concurrent

```go
package main

import (
	"log"
	"math"
	"time"

	"github.com/monaco-io/lib/lru"
	"golang.org/x/sync/errgroup"
)

var instance = lru.New[int, int](100, time.Second*5)

func main() {
	c1 := time.After(time.Second * 3)
	c2 := time.After(time.Second * 10)
	f1 := func() error {
		for {
			select {
			case <-c1:
				return nil
			default:
				for i := 0; i <= math.MaxInt8; i++ {
					instance.Set(i, i)
				}
			}
		}
	}
	f2 := func() error {
		for {
			select {
			case <-c2:
				return nil
			default:
				for i := 0; i <= math.MaxInt8; i++ {
					v, ok := instance.Get(i)
					log.Println(i, v, ok)
				}
			}
		}
	}
	var eg errgroup.Group
	eg.Go(f1)
	eg.Go(f2)
	eg.Wait()
}


```

## example with callback

```go

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/monaco-io/lib/lru"
	"golang.org/x/sync/errgroup"
)

var instance = lru.NewWithCB(100, time.Second*5, func(context.Context, int) (int, error) {
	return 10086, nil
})

func main() {
	c1 := time.After(time.Second * 3)
	c2 := time.After(time.Second * 10)
	f1 := func() error {
		for {
			select {
			case <-c1:
				return nil
			default:
				for i := 0; i <= math.MaxInt8; i++ {
					instance.Set(i, i)
				}
			}
		}
	}
	f2 := func() error {
		for {
			select {
			case <-c2:
				return nil
			default:
				for i := 0; i <= math.MaxInt8; i++ {
					v, ok := instance.GetC(context.Background(), i)
					log.Println(i, v, ok)
				}
			}
		}
	}
	var eg errgroup.Group
	eg.Go(f1)
	eg.Go(f2)
	eg.Wait()
}
```
