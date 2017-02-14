package signal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/wgyuuu/embryonic/helper/dlog"
	"github.com/wgyuuu/embryonic/helper/util"
)

func OpenMonitor(exitFunc func()) {
	go util.SafeFunc(func() { monitor(exitFunc) })
}

func monitor(exitFunc func()) {
	defer os.Exit(0)

	c := make(chan os.Signal)
	signal.Notify(c,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer signal.Stop(c)

	s := <-c
	dlog.Strace("ExitSignal", s)

	exitFunc()
}
