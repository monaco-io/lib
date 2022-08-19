package lru

// package main

// import (
// 	"log"
// 	"math"
// 	"time"

// 	"github.com/monaco-io/lib/lru"
// 	"golang.org/x/sync/errgroup"
// )

// var instance = lru.New[int, int](100, time.Second*5)

// func main() {
// 	c1 := time.After(time.Second * 3)
// 	c2 := time.After(time.Second * 10)
// 	f1 := func() error {
// 		for {
// 			select {
// 			case <-c1:
// 				return nil
// 			default:
// 				for i := 0; i <= math.MaxInt8; i++ {
// 					instance.Set(i, i)
// 				}
// 			}
// 		}
// 	}
// 	f2 := func() error {
// 		for {
// 			select {
// 			case <-c2:
// 				return nil
// 			default:
// 				for i := 0; i <= math.MaxInt8; i++ {
// 					v, ok := instance.Get(i)
// 					log.Println(i, v, ok)
// 				}
// 			}
// 		}
// 	}
// 	var eg errgroup.Group
// 	eg.Go(f1)
// 	eg.Go(f2)
// 	eg.Wait()
// }
