package main

import (
	"fmt"
	"log"
	"code.google.com/p/go.net/websocket"
)

var origin = "http://192.168.2.13"
var url = "ws://192.168.2.13:8080/screen"

func main() {
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal("Dial failure:", err)
	}

	msg := make([]byte, 1024)
	for {
		n, err := ws.Read(msg)
		if n > 0 {
			fmt.Printf("msg: %s\n", msg[:n])
		}
		if err != nil {
			log.Fatal("Failed to read:", err)
		}
	}
}

