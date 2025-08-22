package lru

// package main

// import (
// 	"log"

// 	"github.com/monaco-io/lib/lru"
// )

// var instance = lru.New[int, int](100, time.Second*5)

// func main() {
// 	i := 1
// 	//set key value pairs
// 	instance.Set(i, i)

// 	//get val from mem
// 	v, ok := instance.Get(i)
// 	log.Println(i, v, ok)

// }

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
// 	set := func() error {
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
// 	get := func() error {
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
// 	eg.Go(set)
// 	eg.Go(get)
// 	eg.Go(get)
// 	eg.Wait()
// }
