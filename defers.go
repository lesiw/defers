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
	"os/signal"
	"syscall"
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
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
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

func exitCode(s os.Signal) (code int) {
	defer func() {
		if code == 0 {
			code = 1
		}
	}()
	sig, ok := s.(syscall.Signal)
	if !ok {
		return
	}
	code = (128 + int(sig)) % 255
	return
}
