package sys

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/monaco-io/lib/log"
)

// ExitGrace exec callback functions before exit (SIGINT/SIGQUIT/SIGTERM)
func ExitGrace(callback ...func()) {
	ch := make(chan os.Signal, 1)
	signal.Notify(
		ch,

		// Ctrl-Z 发送 TSTP signal (SIGTSTP); 通常导致进程挂起(suspend)
		// syscall.SIGTSTP,

		// Ctrl-C 发送 INT signal (SIGINT)，通常导致进程结束
		syscall.SIGINT,

		// Ctrl-\ 发送 QUIT signal (SIGQUIT); 通常导致进程结束 和 dump core.
		syscall.SIGQUIT,

		// 结束程序
		syscall.SIGTERM,

		// 终端控制进程结束
		syscall.SIGHUP,
	)

	s := <-ch

	log.I("ExitGrace", "singal", s.String())

	switch s {

	// 程序正常退出
	case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
		for _, fn := range callback {
			fn()
		}

	// 终端连接断开
	case syscall.SIGHUP:
	}

}
