package linux

import (
	"os"
	"testing"
	"time"
)

func TestGetLogFromProcessShouldReturnAMessageForAValidPid(t *testing.T) {
	go panicOnTimeout(10 * time.Second)

	stdErrPath := "../test/stderr"
	stdOutPath := "../test/stdout"

	expectedMessage := "test write \n"
	writeTextToANamedPipe(expectedMessage, stdErrPath, t)
	writeTextToANamedPipe(expectedMessage, stdOutPath, t)

	l := Log{StdErr: stdErrPath, StdOut: stdOutPath}
	message := make(chan string, 1)
	errCh := make(chan error, 1)
	go l.GetLogFromProcess(123, message, errCh)
	mess := <- message
	if mess != "test write \n" {
		t.Errorf("expected message to be: %s but got: %s", expectedMessage, mess)
	}
}

func TestGetLogFromProcessShouldReturnAnErrorForInvalidFilePath(t *testing.T) {
	go panicOnTimeout(10 * time.Second)

	stdErrPath := "../test/invalid"
	stdOutPath := "../test/invalid"

	l := Log{StdErr: stdErrPath, StdOut: stdOutPath}
	message := make(chan string, 1)
	errCh := make(chan error, 1)
	go l.GetLogFromProcess(999999, message, errCh)
	err := <- errCh
	if err == nil {
		t.Errorf("expected error to be nil but got: %v", err)
	}
}

// panicOnTimeout with panic after a period of time.
// This function used as a go routine in tests:
// e.g. go panicOnTimeout(10 * time.Second)
func panicOnTimeout(d time.Duration) {
	<-time.After(d)
	panic("Test timed out")
}

// writeTextToANamedPipe is responsible to write a dummy text to
// the given named pipe
func writeTextToANamedPipe(message, pipeFile string, t *testing.T) {
	f, err := os.OpenFile(pipeFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		t.Fatalf("error opening file: %v", err)
	}
	_, err = f.WriteString(message)
	if err != nil {
		t.Fatalf("faield to write to pipe: %v", err)
	}
	time.Sleep(time.Second)
}