// Package defers handles program-wide defers.
//
// Defers are executed when defers.Exit() is called or when an interrupt signal
// is caught, whichever happens first.
//
// If an interrupt signal is caught, the program will exit with a status of
// 128 plus the signal number. In the event the signal number cannot be
// determined, the program will exit with exit status 1.
package defers

import "os"

type defunc struct {
	fn    func()
	added chan struct{}
}

var (
	defers []func()

	exit = os.Exit

	exitc  = make(chan int)
	deferc = make(chan defunc)
	clearc = make(chan chan struct{})
	sigc   = make(chan os.Signal, 1)
)

func init() {
	go func() {
		for {
			select {
			case code := <-exitc:
				for _, f := range defers {
					f()
				}
				exit(code)
			case d := <-deferc:
				defers = append([]func(){d.fn}, defers...)
				d.added <- struct{}{}
			case c := <-clearc:
				defers = []func(){}
				c <- struct{}{}
			}
		}
	}()
	go func() { exitc <- signalCode(<-sigc) }()
	notify(sigc)
}

// Add adds a function to the front of the defer list.
func Add(f func()) {
	d := defunc{f, make(chan struct{})}
	deferc <- d
	<-d.added
}

// Exit runs all defers, then exits the program.
func Exit(code int) {
	exitc <- code
	select {}
}
