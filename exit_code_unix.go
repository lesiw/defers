//go:build unix
// +build unix

package defers

import (
	"os"
	"syscall"
)

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
