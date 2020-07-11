package clagent_test

import (
	"errors"
	"github.com/cpapidas/clagent"
	"github.com/cpapidas/clagent/test/mock"
	"testing"
	"time"
)

func TestListenToPidShouldNotReturnAnyErrorForValidPort(t *testing.T) {
	conf := clagent.Config{
		Token: "token",
		Port:  123,
	}
	stop := make(chan bool, 1)
	proc := mock.Process{}
	l := mock.Log{}

	// We need to send a stop signal in order to stop
	// the infinity loop.
	go func(){
		time.Sleep(500 * time.Millisecond)
		stop <- true
	}()
	err := clagent.ListenToPid(conf, proc, l, stop)
	if err != nil {
		t.Fatal("expected error not to nil")
	}
}

func TestListenToPidShouldNotReturnAnyErrorForValidPid(t *testing.T) {
	conf := clagent.Config{
		Token: "token",
		Pid:  123,
	}
	stop := make(chan bool, 1)
	proc := mock.Process{}
	l := mock.Log{}
	stop <- true
	err := clagent.ListenToPid(conf, proc, l, stop)
	if err != nil {
		t.Fatal("expected error not to nil")
	}
}

func TestListenToPidShouldReturnAnErrorForInvalidPID(t *testing.T) {
	conf := clagent.Config{
		Token: "token",
		Port:  9999,
	}
	stop := make(chan bool, 1)
	proc := mock.Process{
		FindPIDByGivenPortNumberFnc: func(port int32) (int32, error) {
			return 0, errors.New("invalid pid")
		},
	}
	l := mock.Log{}
	err := clagent.ListenToPid(conf, proc, l, stop)
	if err == nil {
		t.Fatal("expected error not to be nil")
	}
}
