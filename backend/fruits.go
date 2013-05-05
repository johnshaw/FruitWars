package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
	"encoding/json"
	"strconv"
	"strings"
	"os"
	"bufio"
	"code.google.com/p/go.net/websocket"
)

const BOARD_WIDTH = 32
const BOARD_HEIGHT = 18

func MakePlayers() map[string]*Player {
	players := map[string]*Player{}
	players["banana"] = &Player{"banana", "Dick Face", 100, false}
	players["apple"] = &Player{"apple", "Steve's Mom", 100, false}
	players["grape"] = &Player{"grape", "Pete Chen", 100, false}
	players["watermelon"] = &Player{"watermelon", "Oh Hai", 100, false}
	return players
}

func MakeDudes() map[string]*Dude {
	types := []string{"seedling", "fruit"}
	teams := []string{"apple", "banana", "grape", "watermelon"}
	dudes := map[string]*Dude{}
	for i := 0; i<10; i++ {
		pid := teams[rand.Intn(4)]
		t := types[rand.Intn(2)]
		pos := Location{float64(rand.Intn(BOARD_WIDTH)), float64(rand.Intn(BOARD_HEIGHT))}

		x := rand.Float64()*2 - 1.0
		y := rand.Float64()*2 - 1.0
		dir := Direction{x, y, 1.0}
		dir.Normalize()

		dudes[strconv.Itoa(i)] = &Dude{i, pid, t, pos, dir, rand.Intn(50)+50, true}
	}
	return dudes
}

func PickBases(m [][]string) map[string]*Base {
	basepos := []Location{}
	for i, row := range m {
		for j, c := range row {
			if c == "B" {
				basepos = append(basepos, Location{float64(i), float64(j)})
			}
		}
	}

	// Randomly shuffle the basepos array
	for i := len(basepos) - 1; i > 0; i-- {
		if j := rand.Intn(i + 1); i != j {
			tmp := basepos[i]
			basepos[i] = basepos[j]
			basepos[j] = tmp
		}
	}

	rand.Intn(len(basepos))

	bases := map[string]*Base{}

	if len(basepos) > 0 {
		bases["apple"] = &Base{"apple", basepos[0], []Tower{}}
		if len(basepos) > 1 {
			bases["banana"] = &Base{"banana", basepos[1], []Tower{}}
			if len(basepos) > 2 {
				bases["grape"] = &Base{"grape", basepos[2], []Tower{}}
				if len(basepos) > 3 {
					bases["watermelon"] = &Base{"watermelon", basepos[3], []Tower{}}
				}
			}
		}
	}

	return bases
}

func ReadMap(path string) [][]string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Failed to read map file", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	lines := [][]string{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.Trim(line, "\n ")
		if len(line) > 0 {
			lines = append(lines, strings.Split(line, ""))
		}
	}
	return lines
}

func ProcessControlRequest(state *State, msg ControlRequest) {
	switch msg.MsgType {
	case "BuyDude":
		fmt.Println("Buy Dude!")
		var bd BuyDude
		json.Unmarshal([]byte(msg.Data), &bd)

		pid, ok := state.controlToPlayer[msg.ControlId]
		if !ok {
			data, _ := json.Marshal(BuyReject{"Unregistered Player"})
			resp := ControlResponse{"BuyReject", string(data)}
			websocket.JSON.Send(msg.ws, resp)
			return
		}

		player := state.Players[pid]

		for i := 0; i<bd.Count; i++ {
			// Deduct required moneys
			if (player.Water >= 20) {
				player.Water -= 20
				fmt.Printf("Player %s Bought Dude, Water=%d\n", pid, player.Water)
			} else {
				// TODO: Send back notification that all dudes could not be bought?
				break
			}
			state.Players[pid] = player

			fmt.Printf("BuyDude: %+v\n", bd)

			var d Dude
			d.Id = int(rand.Int31())
			d.PlayerId = pid
			d.Type = bd.Type
			d.Pos = Location{state.Bases[pid].Pos.X+1, state.Bases[pid].Pos.Y+1}
			d.Dir = bd.Dir
			d.Dir.V = 1.0
			d.Dir.Normalize()
			d.Health = 100

			state.Dudes[strconv.Itoa(d.Id)] = &d

			fmt.Printf("Added Dude: %+v\n", d)
		}

	case "SelectPlayer":
		fmt.Println("Select Player!")
		var sp SelectPlayer
		json.Unmarshal([]byte(msg.Data), &sp)
		pid, ok := state.controlToPlayer[msg.ControlId]
		if ok {
			if pid == sp.PlayerId {
				data, _ := json.Marshal(ConfirmPlayer{sp.PlayerId})
				resp := ControlResponse{"ConfirmPlayer", string(data)}
				websocket.JSON.Send(msg.ws, resp)
			} else {
				data, _ := json.Marshal(ConfirmPlayer{sp.PlayerId})
				resp := ControlResponse{"RejectPlayer", string(data)}
				websocket.JSON.Send(msg.ws, resp)
			}
			return
		}
		for _, player := range state.Players {
			fmt.Println("Checking player")
			if !player.assigned {
				fmt.Println("Assigned!")
				player.assigned = true
				state.controlToPlayer[msg.ControlId] = player.Id
				data, _ := json.Marshal(ConfirmPlayer{sp.PlayerId})
				resp := ControlResponse{"ConfirmPlayer", string(data)}
				websocket.JSON.Send(msg.ws, resp)
				return
			}
		}
		data, _ := json.Marshal(ConfirmPlayer{sp.PlayerId})
		resp := ControlResponse{"RejectPlayer", string(data)}
		websocket.JSON.Send(msg.ws, resp)
	case "ControlDisconnect":
		fmt.Println("ControlDisconnect!")
		pid, ok := state.controlToPlayer[msg.ControlId]
		if ok {
			state.Players[pid].assigned = false
			delete(state.controlToPlayer, msg.ControlId)
		}
	}
}

func GetDudePositionMap(state *State) map[Location]Dude {
	oldpos := map[Location]Dude{}
	for _, dude := range state.Dudes {
		oldpos[dude.Pos] = *dude
	}
	return oldpos
}

func DudeVsDude(d1 *Dude, d2 *Dude) {
}

func DudeVsTower(d *Dude, t *Tower) {
}

func DudeVsBase(d *Dude, b *Base) {
}

func DoIt(state *State) {
	//oldpos := GetDudePositionMap(state)
	for _, dude := range state.Dudes {
		pos := dude.Pos.Advance(dude.Dir)

		// Dude Vs. Dude
		//if dude.Pos
		// Dude Vs. Tower
		// Dude Vs. Base
		// Store Advanced Position
		dude.Pos = pos

		// Subtract 1 health (rot)
		dude.Health -= 1
	}
}

func GameLoop(schan chan State, cmchan chan ControlRequest) {
	rand.Seed(time.Now().Unix())
	players := MakePlayers()
	dudes := MakeDudes()
	m := ReadMap("./map")
	bases := PickBases(m)
	state := State{players, bases, dudes, m, map[int]string{}}
	tick := time.Tick(1 * time.Second)

	for {
		select {
		case <-tick:
			// Update positions and fight
			DoIt(&state)

			// Award resources
			for _, p := range players {
				p.Water += 5
			}

			// Push state
			schan <- state
		case msg := <-cmchan:
			// Process Events from controllers
			ProcessControlRequest(&state, msg)
		}
	}
}

func main() {
	schan := make(chan State, 10)
	cmchan := make(chan ControlRequest, 10)

	//cuchans := map[int]chan ControlUpdate{}

	bcast := MakeStateBroadcaster(schan)
	//bcast := MakeBroadcaster(schan)

	go GameLoop(schan, cmchan)

	screenHandler := MakeScreenHandler(bcast)
	controlHandler := MakeControlHandler(cmchan)

	http.Handle("/screen", websocket.Handler(screenHandler))
	http.Handle("/control", websocket.Handler(controlHandler))
	fmt.Println("Starting server")
	if err := http.ListenAndServe("192.168.2.13:8080", nil); err != nil {
		log.Fatal("Server failed:", err)
	}
	fmt.Println("exiting server")
}
