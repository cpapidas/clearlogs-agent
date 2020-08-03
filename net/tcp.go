package net

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

// TCP is responsible to describe the tcp client in order
// to send data to the external server.
type TCPClient struct {
	Conn net.Conn
}

// NewTCPClient is responsible to crate a new TCPClient object. If something
// occurred the function will return an error.
func NewTCPClient(token string) (*TCPClient, error) {
	url := os.Getenv("CL_BASEURL")
	if url == "" {
		url = ""
	}
	httpClient := NewHTTPClient(url, &http.Client{Timeout: 5 * time.Second})
	addres, err := httpClient.GetTCPAddress(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get TCP address with error: %v", err)
	}
	tcpc := TCPClient{}
	err = tcpc.Connect(addres)
	if err != nil {
		return nil, fmt.Errorf("failed to connect with error: %v", err)
	}

	return &tcpc, nil
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