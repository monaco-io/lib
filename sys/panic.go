package sys

import (
	"fmt"
	"log"
)

func Recover() (msg string) {
	if e := recover(); e != nil {
		msg = fmt.Sprintf("PANIC RECOVER: [ %+v ]", e)
		log.Println(msg)
	}
	return
}
