package server

import (
	"fmt"
	"code.google.com/p/go.net/websocket"
)

func ScreenHandler(bcast *StateBroadcaster, ws *websocket.Conn) {
	fmt.Println("Screen connected")

	// This channel will receive state data and publish it to the client
	id, schan := bcast.GetChan()
	// Clean up the channel when done
	defer bcast.DelChan(id)

	for {
		state := <-schan
		err := websocket.JSON.Send(ws, state)
		if err != nil {
			fmt.Printf("Screen Handler Error: %s\n", err)
			return
		}
	}
}

func MakeScreenHandler(bcast *StateBroadcaster) func (*websocket.Conn) {
	handler := func (ws *websocket.Conn) {
		ScreenHandler(bcast, ws)
	}
	return handler
}
