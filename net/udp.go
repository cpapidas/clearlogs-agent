package net

import (
	"fmt"
	"net"
)

// UDP is responsible to describe the UDP functionalities.
type UDP struct {
	// PC defined the UDP connection interface.
	PC net.PacketConn
	// MessageSize defines the message size that we can receive.
	MessageSize int
}

// NewUDP initializes and returns a new UDP object. If something
// occurred the function will return an error.
func NewUDP(address string, ) (*UDP, error) {
	pc, err := net.ListenPacket("udp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to set up the UDP server, %v", err)
	}

	return &UDP{PC: pc}, nil
}

// GetLogFromProcessType describes the function which will get the logs
// from UDP. This implementation does not need the first argument, "input".
// The second argument is the "message chan" which is the logs per line
//and the third the "error channel" in case something happen.
func (u UDP) GetLogFromProcess(input string, message chan<- string, errChan chan<- error) {
	for {
		buf := make([]byte, u.MessageSize)
		_, _, err := u.PC.ReadFrom(buf)
		mes := string(buf)
		if mes != "" {
			message <- mes
		}
		if err != nil {
			err = fmt.Errorf("error on udp message, %v\n", err)
			errChan <- err
			continue
		}
	}
}
