// Package defers handles program-wide defers.
//
// Defers are executed when defers.Exit() is called or when an interrupt signal
// is caught, whichever happens first.
//
// If an interrupt signal is caught, the program will exit with a status of
// 128 plus the signal number. In the event the signal number cannot be
// determined, the program will exit with exit status 1.
package defers

import (
	"os"
)

type empty = struct{}

type defunc struct {
	fn    func()
	added chan empty
}

var (
	defers []func()

	exit  = make(chan int)
	queue = make(chan defunc)
	sig   = make(chan os.Signal, 1)
)

func init() {
	go func() {
		code := <-exit
		for _, f := range defers {
			f()
		}
		os.Exit(code)
	}()
	go func() {
		for d := range queue {
			defers = append([]func(){d.fn}, defers...)
			d.added <- empty{}
		}
	}()
	notify(sig)
	go func() { exit <- signalCode(<-sig) }()
}

// Add adds a function to the front of the defer list.
func Add(f func()) {
	d := defunc{f, make(chan empty)}
	queue <- d
	<-d.added
}

// Exit runs all defers, then exits the program.
func Exit(code int) {
	exit <- code
	select {}
}
