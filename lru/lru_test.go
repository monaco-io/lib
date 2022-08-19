package lru

import (
	"math"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

var instance ICache[int, int] = New[int, int](100, time.Second*10)

func TestAdd(t *testing.T) {
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
					t.Log(i, v, ok)
				}
			}
		}
	}
	var eg errgroup.Group
	eg.Go(f1)
	eg.Go(f2)
	eg.Wait()
}
