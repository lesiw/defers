# lesiw.io/defers

[![Go Reference](https://pkg.go.dev/badge/lesiw.io/defers.svg)](https://pkg.go.dev/lesiw.io/defers)

``` go
// Package defers handles program-wide defers.
//
// Defers are executed when defers.Exit() is called or when an interrupt signal
// is caught, whichever happens first.
//
// If an interrupt signal is caught, the program will exit with a status of
// 128 plus the signal number. In the event the signal number cannot be
// determined, the program will exit with exit status 1.
```

## Example

``` go
package main

import (
    "fmt"
    "os"

    "lesiw.io/defers"
)

// Set stop to true to halt the program.
// This forces the Go Playground to send an os.Interrupt.
// Global defers will still run before the program ends.
var stop = false

var success bool

func main() {
    defers.Add(func() {
        if success {
            fmt.Fprintln(os.Stderr, "The program executed successfully.")
        } else {
            fmt.Fprintln(os.Stderr, "The program was interrupted.")
        }
    })
    fmt.Println("Preparing to send a greeting...")
    if stop {
        select {}
    }
    fmt.Println("Hello world!")
    success = true
    defers.Exit(0)
}
```

[▶️ Run this example on the Go Playground](https://go.dev/play/p/amY5VkD51QF)
