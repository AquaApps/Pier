package common

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitingSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
