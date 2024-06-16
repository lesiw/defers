//go:build !unix
// +build !unix

package defers

import "os"

func exitCode(os.Signal) int {
	return 1
}
