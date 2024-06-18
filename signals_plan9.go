//go:build plan9
// +build plan9

package defers

import (
	"os"
	"os/signal"
)

func notify(c chan<- os.Signal) {
	signal.Notify(c, os.Interrupt)
}

func signalCode(os.Signal) (code int) {
	return 130 // Consistent with SIGINT (128 + 2).
}
