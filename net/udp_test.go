package net

import (
	"errors"
	"github.com/cpapidas/clagent/test/mock"
	"net"
	"strings"
	"testing"
	"time"
)

func TestUDPGetLogFromProcessShouldSendAMessageForValidArguments(t *testing.T) {
	want := "m"
	packetConn := mock.UDP{ReadFromFun: func(p []byte) (n int, addr net.Addr, err error) {
		rp := []byte(want)
		_p0 := &rp[0]
		p[0] = *_p0
		return 0, nil, nil
	}}
	messageCh := make(chan string, 1)
	errCh := make(chan error, 1)

	udp := &UDP{PC: packetConn, MessageSize: 1024}
	go udp.GetLogFromProcess("", messageCh, errCh)

	go func() {
		time.Sleep(1 * time.Second)
		t.Fatal("timeout")
	}()

	got := <-messageCh

	if got[0:1] != want[0:1] {
		t.Errorf("expected: %s; got:%s", want, got)
	}
}

func TestUDPGetLogFromProcessShouldReturnErrorForInvalidMessage(t *testing.T) {
	want := errors.New("test_error")
	packetConn := mock.UDP{ReadFromFun: func(p []byte) (n int, addr net.Addr, err error) {
		return 0, nil, want
	}}
	messageCh := make(chan string, 1)
	errCh := make(chan error, 1)

	udp := &UDP{PC: packetConn, MessageSize: 1024}
	go udp.GetLogFromProcess("", messageCh, errCh)

	go func() {
		time.Sleep(1 * time.Second)
		t.Fatal("timeout")
	}()

	got := <-errCh

	if !strings.Contains(got.Error(), want.Error()) {
		t.Errorf("expected: %s to contain; got: %s", got, want)
	}
}
