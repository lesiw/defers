// Package defers handles program-wide defers.
//
// It executes them when defers.Exit() is called or if an os.Interrupt signal
// is sent.
package defers

import (
	"os"
	"os/signal"
)

type empty = struct{}

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
	signal.Notify(sig, os.Interrupt)
	go func() { exit <- exitCode(<-sig) }()
}

// Add adds a function to the front of the defer list.
func Add(f func()) {
	d := newDefunc(f)
	queue <- d
	<-d.added
}

// Exit runs all defers, then exits the program.
func Exit(code int) {
	exit <- code
	select {}
}

type defunc struct {
	fn    func()
	added chan empty
}

func newDefunc(f func()) defunc {
	return defunc{f, make(chan empty)}
}
