//go:build !plan9
// +build !plan9

package defers

import (
	"os"
	"os/signal"
	"syscall"
)

func notify(c chan<- os.Signal) {
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
}

func signalCode(s os.Signal) (code int) {
	defer func() {
		if code == 0 {
			code = 1
		}
	}()
	sig, ok := s.(syscall.Signal)
	if !ok {
		return
	}
	code = (128 + int(sig)) % 256
	return
}
