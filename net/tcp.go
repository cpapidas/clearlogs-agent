package net

import (
	"fmt"
	"net"
)

// TCP is responsible to describe the tcp client in order
// to send data to the external server.
type TCPClient struct {
	Conn net.Conn
}

// Connect function is responsible to set up the TCP connection
// with the server. It gets the server address (e.g. localhost:3000)
// as a string and will return an error if something occurred.
func (t *TCPClient) Connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect to the server using tcp connection, address: %s, error: %v", address, err)
	}
	t.Conn = conn
	return nil
}

// Send is responsible to send a message to the client. It gets the
// message as a string and in any case it returns an error.
func (t TCPClient) Send(message string) error {
	_, err := t.Conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("faield to write to server with error: %v", err)
	}
	return nil
}