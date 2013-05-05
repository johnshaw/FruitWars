package main

import (
	"math/rand"
)

type Location struct {
	X int
	Y int
}

type Direction struct {
	X float64
	Y float64
	// Velocity
	V float64
}

type Player struct {
	Id int
	Team string
	Name string
	Water int
}

type Tower struct {
	// Position is 0->11, CW from top-left corner
	Pos int
	Type string
	Health int
}

type Base struct {
	Id int
	// Player id
	PlayerId int
	BaseLocation Location
	Towers []Tower
}

type Dude struct {
	Id int
	// Player id
	PlayerId int
	Type string
	Pos Location
	Dir Direction
	Health int
}

type State struct {
	Players []Player
	Bases []Base
	Dudes []Dude
}

type NewClientChan struct {
	ClientId  int
	StateChan chan State
}

type DelClientChan struct {
	ClientId int
}

type StateBroadcaster struct {
	newChan chan NewClientChan
	delChan chan DelClientChan
}

func (b *StateBroadcaster) GetChan() (int, chan State) {
	// This channel will receive state data and publish it to the client
	c := make(chan State)
	id := rand.Int()
	b.newChan <- NewClientChan{id, c}
	return id, c
}

func (b *StateBroadcaster) DelChan(id int) {
	b.delChan <- DelClientChan{id}
}

func (b *StateBroadcaster) broadcastLoop(stateChan chan State) {
	out := map[int]chan State{}
	for {
		select {
		case cnew := <-b.newChan:
			out[cnew.ClientId] = cnew.StateChan
		case cdel := <-b.delChan:
			delete(out, cdel.ClientId)
		case state := <-stateChan:
			for _, c := range out {
				c <- state
			}
		}
	}
}

func MakeStateBroadcaster(stateChan chan State) *StateBroadcaster {
	b := new(StateBroadcaster)
	b.newChan = make(chan NewClientChan)
	b.delChan = make(chan DelClientChan)
	go b.broadcastLoop(stateChan)
	return b
}
