package main

import (
	"math"
	"math/rand"
)

type Location struct {
	X float64
	Y float64
}

type Direction struct {
	X float64
	Y float64
	// Velocity
	V float64
}

func (d *Direction) Normalize() {
	l := math.Hypot(d.X, d.Y)
	d.X = d.X / l
	d.Y = d.Y / l
}

func (l *Location) Advance(dir Direction) Location {
	x := l.X + dir.X * dir.V
	y := l.Y + dir.Y * dir.V

	if x > BOARD_WIDTH {
		x -= BOARD_WIDTH 
	} else if x < 0 {
		x += BOARD_WIDTH
	}

	if y > BOARD_HEIGHT {
		y -= BOARD_HEIGHT;
	} else if y < 0 {
		y += BOARD_HEIGHT
	}

	return Location{x, y}
}

func Wrap(x float64, v float64) float64 {
	if x > v {
		x -= v
	} else if x < 0 {
		x += v
	}
	return x
}

func WrapWidth(x float64) float64 {
	return Wrap(x, BOARD_WIDTH)
}

func WrapHeight(y float64) float64 {
	return Wrap(y, BOARD_HEIGHT)
}

// Get neighbor
// +-+-+-+
// |0|1|2|
// +-+-+-+
// |7| |3|
// +-+-+-+
// |6|5|4|
// +-+-+-+
func (l *Location) Neighbor(n int) Location {
	switch n {
	case 0:
		return Location{WrapWidth(l.X - 1), WrapHeight(l.Y - 1)}
	case 1:
		return Location{WrapWidth(l.X), WrapHeight(l.Y - 1)}
	case 2:
		return Location{WrapWidth(l.X + 1), WrapHeight(l.Y - 1)}
	case 3:
		return Location{WrapWidth(l.X + 1), WrapHeight(l.Y)}
	case 4:
		return Location{WrapWidth(l.X + 1), WrapHeight(l.Y + 1)}
	case 5:
		return Location{WrapWidth(l.X), WrapHeight(l.Y + 1)}
	case 6:
		return Location{WrapWidth(l.X - 1), WrapHeight(l.Y + 1)}
	case 7:
		return Location{WrapWidth(l.X - 1), WrapHeight(l.Y)}
	}
	return Location{l.X, l.Y}
}

type Player struct {
	Id string
	Name string
	Water int
	assigned bool
}

type Tower struct {
	// Position is 0->11, CW from top-left corner
	Pos int
	Type string
	Health int
}

type Base struct {
	// Player id
	Id string
	Pos Location
	Towers []Tower
}

type Dude struct {
	Id int
	// Player id "banana", "grape", "watermelon", "apple"
	PlayerId string
	// "antitower", "antidude", "antibase"
	Type string
	Pos Location
	Dir Direction
	Health int
	Alive bool
}

type State struct {
	// Each string key is actually an int, but the json encoding requires
	// string keys
	Players map[string]*Player
	Bases map[string]*Base
	Dudes map[string]*Dude
	Map [][]string

	controlToPlayer map[int]string
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
