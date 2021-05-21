package sys

import (
	"fmt"

	"github.com/monaco-io/lib/log"
)

func Recover() (msg string) {
	if e := recover(); e != nil {
		msg = fmt.Sprintf("PANIC RECOVER: [ %+v ]", e)
		log.E(msg)
	}
	return
}
