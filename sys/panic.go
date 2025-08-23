package sys

import (
	"context"
	"fmt"

	"github.com/monaco-io/lib/xlog"
)

func Recover() (msg string) {
	if e := recover(); e != nil {
		msg = fmt.Sprintf("PANIC RECOVER: [ %+v ]", e)
		xlog.E(context.Background(), msg)
	}
	return
}
