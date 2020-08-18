package main
//
//import (
//	"github.com/cpapidas/clagent/net"
//	"log"
//)
//
//func main() {
//	tcpClient := net.TCPClient{}
//
//
//	err := tcpClient.Connect("wss://0vqkck5wa9.execute-api.eu-central-1.amazonaws.com/production/")
//	if err != nil {
//		log.Fatalf("failed to connect to socket with error: %v", err)
//	}
//	err = tcpClient.Send("testMessage")
//	if err != nil {
//		log.Fatalf("failed to send a message with error: %v", err)
//	}
//
//}

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
)

var addr = flag.String("addr", "0vqkck5wa9.execute-api.eu-central-1.amazonaws.com", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: *addr, Path: "/production"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	err = c.WriteMessage(websocket.TextMessage, []byte(`{"action": "send-message", "data": "f"}`))
	if err != nil {
		log.Println("write:", err)
		return
	}

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		}
	}
}