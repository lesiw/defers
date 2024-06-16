//go:build unix
// +build unix

package defers

import (
	"os"

	"golang.org/x/sys/unix"
)

func exitCode(s os.Signal) (code int) {
	code = (128 + int(unix.SignalNum(s.String()))) % 255
	if code == 0 {
		code = 1
	}
	return
}
