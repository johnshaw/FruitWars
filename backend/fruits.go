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

func MakePlayers() []Player {
	players := make([]Player, 4)
	players[0] = Player{10, "banana", "Dick Face", 100}
	players[1] = Player{11, "apple", "Steve's Mom", 100}
	players[2] = Player{12, "grape", "Pete Chen", 100}
	players[3] = Player{13, "watermelon", "Oh Hai", 100}
	return players
}

func MakeDudes() []Dude {
	types := []string{"seedling", "fruit"}
	dudes := make([]Dude, 10)
	for i := 0; i<10; i++ {
		pid := rand.Intn(4) + 10
		t := types[rand.Intn(2)]
		pos := Location{rand.Intn(32), rand.Intn(18)}
		dir := Direction{rand.Float64()*2-1.0, rand.Float64()*2-1.0, 1.0}
		dudes[i] = Dude{i, pid, t, pos, dir, rand.Intn(50)+50}
	}
	return dudes
}

func PickBases() []Base {
	bases := make([]Base, 4)
	loc1 := Location{rand.Intn(28), rand.Intn(14)}
	bases[0] = Base{1, 10, loc1, []Tower{}}
	loc2 := Location{rand.Intn(28), rand.Intn(14)}
	bases[1] = Base{2, 11, loc2, []Tower{}}
	loc3 := Location{rand.Intn(28), rand.Intn(14)}
	bases[2] = Base{3, 12, loc3, []Tower{}}
	loc4 := Location{rand.Intn(28), rand.Intn(14)}
	bases[3] = Base{4, 13, loc4, []Tower{}}

	return bases
}

type BuyTowerParams struct {
	// Location of tower in base
	Pos int
}

type BuyDudeParams struct {
	Type string
}

// Action sent by a controller
type Action struct {
	PlayerId int
	// Name of action
	//   BuyTower
	//	 BuyDude
	//   DeployDude
	Action string
	// Params depends on the action
	Params interface{}
}

//func ProcessActions(state *State, actionChan) {
//}

func GameLoop(stateChan chan State) {
	bases := PickBases()
	players := MakePlayers()
	dudes := MakeDudes()
	state := State{players, bases, dudes}

	for {
		// Process Events from controllers

		// Update Positions

		// Award resources
		for i := 0; i<len(players); i++ {
			players[i].Water += 5
		}

		// Push state
		stateChan <- state

		time.Sleep(3 * time.Second)
	}
}

func ScreenHandler(bcast *StateBroadcaster, ws *websocket.Conn) {
	fmt.Println("Screen connected")

	// This channel will receive state data and publish it to the client
	id, stateChan := bcast.GetChan()
	// Clean up the channel when done
	defer bcast.DelChan(id)

	for {
		state := <-stateChan
		err := websocket.JSON.Send(ws, state)
		if err != nil {
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

func ControlHandler(actionChan chan Action, ws *websocket.Conn) {
	fmt.Println("Controller connected")
	for {
		var action Action
		err := websocket.JSON.Receive(ws, &action)
		if err != nil {
			return
		}
		actionChan <- action
	}
}

func MakeControlHandler(actionChan chan Action) func (*websocket.Conn) {
	handler := func (ws *websocket.Conn) {
		ControlHandler(actionChan, ws)
	}
	return handler
}

func main() {
	stateChan := make(chan State, 10)
	actChan := make(chan Action, 10)

	bcast := MakeStateBroadcaster(stateChan)

	go GameLoop(stateChan)

	screenHandler := MakeScreenHandler(bcast)
	controlHandler := MakeControlHandler(actChan)

	http.Handle("/screen", websocket.Handler(screenHandler))
	http.Handle("/control", websocket.Handler(controlHandler))
	fmt.Println("Starting server")
	if err := http.ListenAndServe("192.168.2.13:8080", nil); err != nil {
		log.Fatal("Server failed:", err)
	}
	fmt.Println("exiting server")
}
