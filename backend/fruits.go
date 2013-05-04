package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
	//"encoding/json"
	"code.google.com/p/go.net/websocket"
)

type Location struct {
	X int
	Y int
}

type Base struct {
	Player       string
	BaseLocation Location
}

type State struct {
	Bases []Base
}

type NewClientChan struct {
	ClientId  int
	stateChan chan State
}

type DelClientChan struct {
	ClientId int
}

func PickBases() []Base {
	bases := make([]Base, 4)
	loc1 := Location{rand.Intn(28), rand.Intn(14)}
	bases[0] = Base{"Dick Face", loc1}
	loc2 := Location{rand.Intn(28), rand.Intn(14)}
	bases[1] = Base{"Steve's Mom", loc2}
	loc3 := Location{rand.Intn(28), rand.Intn(14)}
	bases[2] = Base{"Nigga Please", loc3}
	loc4 := Location{rand.Intn(28), rand.Intn(14)}
	bases[3] = Base{"Oh Hai", loc4}

	return bases
}

func Broadcast(
	newChan chan NewClientChan,
	delChan chan DelClientChan,
	stateChan chan State) {
	out := map[int]chan State{}
	for {
		select {
		case cnew := <-newChan:
			out[cnew.ClientId] = cnew.stateChan
		case cdel := <-delChan:
			delete(out, cdel.ClientId)
		case state := <-stateChan:
			for _, c := range out {
				c <- state
			}
		}
	}
}

func GameLoop(stateChan chan State) {
	bases := PickBases()
	state := State{bases}

	for {
		stateChan <- state

		time.Sleep(3 * time.Second)
	}
}

func ScreenHandler(newChan chan NewClientChan, delChan chan DelClientChan, ws *websocket.Conn) {
	fmt.Println("Screen connected")

	// This channel will receive state data and publish it to the client
	stateChan := make(chan State)
	id := rand.Int()
	newChan <- NewClientChan{id, stateChan}

	for {
		state := <-stateChan
		err := websocket.JSON.Send(ws, state)
		if err != nil {
			// Clean up the channel when done
			delChan <- DelClientChan{id}
			return
		}
	}
}

func MakeScreenHandler(newChan chan NewClientChan, delChan chan DelClientChan) func (*websocket.Conn) {
	handler := func (ws *websocket.Conn) {
		ScreenHandler(newChan, delChan, ws)
	}
	return handler
}

func ControlHandler(ws *websocket.Conn) {
	fmt.Println("Controller connected")
}

func main() {
	newChan := make(chan NewClientChan)
	delChan := make(chan DelClientChan)
	stateChan := make(chan State, 10)

	go Broadcast(newChan, delChan, stateChan)
	go GameLoop(stateChan)

	screenHandler := MakeScreenHandler(newChan, delChan)

	http.Handle("/screen", websocket.Handler(screenHandler))
	http.Handle("/control", websocket.Handler(ControlHandler))
	fmt.Println("Starting server")
	if err := http.ListenAndServe("192.168.2.13:8080", nil); err != nil {
		log.Fatal("Server failed:", err)
	}
	fmt.Println("exiting server")
}
