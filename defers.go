// Package defers handles program-wide defers.
//
// It executes them when defers.Exit() is called or if an os.Interrupt signal
// is sent.
package defers

import (
	"os"
	"os/signal"
)

var (
	defers []func()

	exit = make(chan int)
	sig  = make(chan os.Signal, 1)
)

func init() {
	go func() {
		code := <-exit
		for _, f := range defers {
			f()
		}
		os.Exit(code)
	}()
	signal.Notify(sig, os.Interrupt)
	go func() { exit <- exitCode(<-sig) }()
}

// Add adds a function to the front of the defer list.
func Add(f func()) {
	defers = append([]func(){f}, defers...)
}

// Exit runs all defers, then exits the program.
func Exit(code int) {
	exit <- code
	select {}
}
