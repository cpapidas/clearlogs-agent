package mac

import (
	"os"
	"testing"
	"time"
)

func TestGetLogFromProcessShouldReturnAnErrorForInvalidPID(t *testing.T) {
	go panicOnTimeout(10 * time.Second)
	l := Log{}
	message := make(chan string, 1)
	errCh := make(chan error, 1)
	go l.GetLogFromProcess(999999, message, errCh)
	err := <- errCh
	if err.Error() != "EOF" {
		t.Error("expected error to not be empty")
	}
}

func TestGetLogFromProcessShouldReturnAMessageForAValidPid(t *testing.T) {
	go panicOnTimeout(10 * time.Second)
	l := Log{}
	message := make(chan string, 1)
	errCh := make(chan error, 1)
	go l.GetLogFromProcess(os.Getpid(), message, errCh)
	err := <- errCh
	if err == nil || err.Error() != "EOF" {
		t.Errorf("expected error to be empty but got: %v", err)
	}
}

// panicOnTimeout with panic after a period of time.
// This function used as a go routine in tests:
// e.g. go panicOnTimeout(10 * time.Second)
func panicOnTimeout(d time.Duration) {
	<-time.After(d)
	panic("Test timed out")
}
