package net

import (
	"github.com/cpapidas/clagent/test"
	"log"
	"testing"
)

func init() {
	// Start the new server
	tcp, err := test.NewServer("tcp", ":1123")
	if err != nil {
		log.Println("error starting TCP server")
		return
	}
	// Run the servers in goroutines to stop blocking
	go func() {
		err := tcp.Run()
		if err != nil {
			log.Fatalf("failed to start tcp server with error: %v", err)
		}
	}()
}

func TestConnectionShouldConnectForAValidAddress(t *testing.T) {
	client := TCPClient{}
	err := client.Connect(":1123")
	if err != nil {
		t.Errorf("expected error to be nil but got: %v", err)
	}
	defer client.Conn.Close()
}

func TestConnectionShouldReturnAnErrorForAInvalidAddress(t *testing.T) {
	client := TCPClient{}
	err := client.Connect("")
	if err == nil {
		t.Fatal("expected error to not be nil")
	}
}

func TestSendShouldSendAMessageForAValidConnection(t *testing.T) {
	client := TCPClient{}
	err := client.Connect(":1123")
	if err != nil {
		t.Errorf("expected error to be nil but got: %v", err)
	}
	err = client.Send("simple message")
	if err != nil {
		t.Errorf("expected error to be nil but got: %v", err)
	}
	defer client.Conn.Close()
}

func TestSendShouldReturnAnErrorForInvalidConnection(t *testing.T) {
	client := TCPClient{}
	err := client.Connect(":1123")
	if err != nil {
		t.Errorf("expected error to be nil but got: %v", err)
	}
	client.Conn.Close()
	err = client.Send("simple message")
	if err == nil {
		t.Fatal("expected error to be nil")
	}
}
