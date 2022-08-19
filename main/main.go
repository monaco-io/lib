package main

import (
	"log"
	"math"
	"time"

	"github.com/monaco-io/lib/lru"
	"golang.org/x/sync/errgroup"
)

func main() {
	var instance lru.ICache[int, int] = lru.New[int, int](100, time.Second*10)

	c1 := time.After(time.Second * 15)
	c2 := time.After(time.Second * 15)
	f1 := func() error {
		for {
			select {
			case <-c1:
				return nil
			default:
				for i := 0; i <= math.MaxInt8; i++ {
					instance.Add(i, i)
					// time.Sleep(time.Second / 100)
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
					if i == 100 {
						instance.Clear()
					}
					v, ok := instance.Get(i)
					log.Println(i, v, ok, instance.Len())
					// instance.Get(i)
				}
			}
		}
	}
	var eg errgroup.Group
	eg.Go(f1)
	eg.Go(f2)
	eg.Wait()
}
