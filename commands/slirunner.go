package commands

import (
	"os"
	"os/signal"
	"syscall"
)

var SLIRunner struct {
	Once  onceCommand  `command:"once" description:"performs a single run of the SLIs suite"`
	Start startCommand `command:"start" description:"initiates the periodic run of the SLIs suite"`
}

func onTerminationSignal(f func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	f()
}
