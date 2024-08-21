//go:build !plan9 && !windows
// +build !plan9,!windows

package defers

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

func TestSignal(t *testing.T) {
	if os.Getenv("DEFER_TEST_PROC") == "1" {
		Add(func() { fmt.Println("world") })
		Add(func() { fmt.Printf("hello ") })
		fmt.Println("READY") // Signal that the process is ready.
		select {}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	cmd := exec.CommandContext(ctx, os.Args[0],
		"-test.v", "-test.run=TestSignal")
	cmd.Env = append(os.Environ(), "DEFER_TEST_PROC=1")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("failed to attach StdoutPipe() on test process: %s", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start test process: %s", err)
	}
	buf := bufio.NewReader(stdout)
	for {
		text, err := buf.ReadString('\n')
		if text == "READY\n" {
			break
		} else if err != nil {
			t.Fatalf("error waiting for test process: "+
				"failed to read stdout: %s", err)
		}
	}
	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		t.Fatalf("failed to send signal to test process: %s", err)
	}
	text, err := buf.ReadString('\n')
	if err != nil {
		t.Errorf("failed to read test process stdout: %s", err)
	}
	if got, want := text, "hello world\n"; got != want {
		t.Errorf("test process stdout = %q, want %q", got, want)
	}
	err = cmd.Wait()
	cancel()

	if ee := new(exec.ExitError); errors.As(err, &ee) {
		if got, want := ee.ProcessState.ExitCode(), 130; got != want {
			t.Errorf("kill -SIGINT = %d, want %d", got, want)
			t.Error(ee.Error())
		}
	} else if err != nil {
		t.Errorf("cmd.Wait() = %q", err)
	}
}
