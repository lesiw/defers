package main

import (
	"os"

	"labs.lesiw.io/ops/golib"
	"lesiw.io/ops"
)

type Ops struct{ golib.Ops }

func main() {
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "build")
	}
	ops.Handle(Ops{})
}
