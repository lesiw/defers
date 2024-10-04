package defers

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDefers(t *testing.T) {
	origExitFn := exit
	t.Cleanup(func() { exit = origExitFn })
	exitcode := make(chan int)
	exit = func(code int) { exitcode <- code }
	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("%d defers", i), func(t *testing.T) {
			wantcode := rand.Intn(127)
			a := []int{}
			for j := 0; j < i; j++ {
				j := j
				Add(func() { a = append(a, j) })
			}
			go Exit(wantcode)
			if got, want := <-exitcode, wantcode; got != want {
				t.Errorf("got exit(%d) call, want %d", got, want)
			}
			if got, want := len(a), i; got != want {
				t.Errorf("len(a) = %d, want %d", got, want)
			}
			for j := 0; j < len(a); j++ {
				if got, want := a[j], i-j-1; got != want {
					t.Errorf("a[%d] = %d, want %d", j, got, want)
				}
			}
			c := make(chan struct{})
			runc <- c
			<-c
		})
	}
}

func TestProc(t *testing.T) {
	if os.Getenv("DEFER_TEST_PROC") == "1" {
		Add(func() { fmt.Println("world") })
		Add(func() { fmt.Printf("hello ") })
		Exit(42)
		fmt.Println("you should not see this!")
	}
	buf := new(bytes.Buffer)
	cmd := exec.Command(os.Args[0], "-test.run=TestProc")
	cmd.Env = append(os.Environ(), "DEFER_TEST_PROC=1")
	cmd.Stdout = buf

	err := cmd.Run()

	if ee := new(exec.ExitError); errors.As(err, &ee) {
		if got, want := ee.ProcessState.ExitCode(), 42; got != want {
			t.Errorf("Exit(%[2]d) = %[1]d, want %[2]d", got, want)
		}
	} else if err != nil {
		t.Errorf("cmd.Run() = %q", err)
	}
	if got, want := buf.String(), "hello world\n"; got != want {
		t.Errorf("proc output -want +got\n%s", cmp.Diff(want, got))
	}
}

func TestRecoverPanic(t *testing.T) {
	if os.Getenv("DEFER_TEST_PROC") == "1" {
		func() {
			defer Recover()
			Add(func() { fmt.Println("world") })
			Add(func() { fmt.Printf("hello ") })
			panic("panicking")
		}()
		os.Exit(0) // Should not be reachable.
	}
	buf := new(bytes.Buffer)
	cmd := exec.Command(os.Args[0], "-test.run=TestRecoverPanic")
	cmd.Env = append(os.Environ(), "DEFER_TEST_PROC=1")
	cmd.Stdout = buf

	err := cmd.Run()

	if ee := new(exec.ExitError); errors.As(err, &ee) {
		if got, want := ee.ProcessState.ExitCode(), 2; got != want {
			t.Errorf("Exit(%[2]d) = %[1]d, want %[2]d", got, want)
		}
	} else if err != nil {
		t.Errorf("cmd.Run() = %q", err)
	}
	// The subprocess will contain additional test-related output,
	// so only check the first line.
	line := strings.Split(buf.String(), "\n")[0]
	if got, want := line, "hello world"; got != want {
		t.Errorf("proc output -want +got\n%s", cmp.Diff(want, got))
	}
}

func TestRecoverNoPanic(t *testing.T) {
	if os.Getenv("DEFER_TEST_PROC") == "1" {
		func() {
			defer Recover()
			Add(func() { fmt.Println("world") })
			Add(func() { fmt.Printf("hello ") })
		}()
		os.Exit(0)
	}
	buf := new(bytes.Buffer)
	cmd := exec.Command(os.Args[0], "-test.run=TestRecoverNoPanic")
	cmd.Env = append(os.Environ(), "DEFER_TEST_PROC=1")
	cmd.Stdout = buf

	err := cmd.Run()

	if err != nil {
		t.Errorf("cmd.Run() = %q, want <nil>", err)
	}
	if got, want := buf.String(), ""; got != want {
		t.Errorf("proc output -want +got\n%s", cmp.Diff(want, got))
	}
}
