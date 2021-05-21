package main

import (
	"log"
	"time"

	"github.com/monaco-io/lib"
)

func main() {
	exit1 := func() {
		time.Sleep(time.Second * 1)
		log.Println("exit1 ...")
	}
	exit2 := func() {
		time.Sleep(time.Second * 2)
		log.Println("exit2 ...")
	}
	exit3 := func() {
		time.Sleep(time.Second * 3)
		log.Println("exit3 ...")
	}
	lib.ExitGrace(exit1, exit2, exit3)
}
