package main

import (
	"fmt"
	"math/rand"
	"code.google.com/p/go.net/websocket"
)

type ControlRequest struct {
	// Controller Id. This is not the player id, but instead a unique
	// identifier of the controller.
	ControlId int
	MsgType string
	Data string
	ws *websocket.Conn
}

type ControlResponse struct {
	MsgType string
	Data string
}

type BuyTower struct {
	// Location of tower in base
	Pos int
}

type BuyDude struct {
	// tower, dude, base
	Type string
	Dir Direction
	Count int
}

type BuyReject struct {
	Reason string
}

type DeployReject struct {
	Reason string
}

type SelectPlayer struct {
	PlayerId string
}

// Sent to controller on connect
type PlayerList struct {
	PlayersIds []string
}

// Confirms a successful player pick
type ConfirmPlayer struct {
	PlayerId string
}

// Reject a player pick
type RejectPlayer struct {
	PlayerId string
}

type ControlUpdate struct {
}

func ControlHandler(cmchan chan ControlRequest, ws *websocket.Conn) {
	fmt.Println("Controller connected")
	id := rand.Int()
	for {
		var req ControlRequest

		err := websocket.JSON.Receive(ws, &req)
		if err != nil {
			fmt.Printf("Control receive failed: %s\n", err)
			req.ControlId = id
			req.MsgType = "ControlDisconnect"
			cmchan <- req
			return
		}
		fmt.Println("Control Request!")
		fmt.Printf("%+v\n", req)

		req.ControlId = id
		req.ws = ws

		cmchan <- req
	}
}

func MakeControlHandler(cmchan chan ControlRequest) func (*websocket.Conn) {
	handler := func (ws *websocket.Conn) {
		ControlHandler(cmchan, ws)
	}
	return handler
}
