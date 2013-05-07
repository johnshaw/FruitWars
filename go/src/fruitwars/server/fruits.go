package server

import (
	"fmt"
	"log"
	"math"
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

const TOWER_COST = 50
const DUDE_COST = 20

const CLOCK_TICK = 100

const STARTING_WATER = 100
const WATER_PER_ROUND = 1

const DUDE_HEALTH = 100
const TOWER_HEALTH = 100
const BASE_HEALTH = 100

const DUDE_VELOCITY = 0.4

func MakePlayers() map[string]*Player {
	players := map[string]*Player{}
	players["banana"] = &Player{"banana", "Dick Face", STARTING_WATER, false}
	players["apple"] = &Player{"apple", "Steve's Mom", STARTING_WATER, false}
	players["grape"] = &Player{"grape", "Pete Chen", STARTING_WATER, false}
	players["watermelon"] = &Player{"watermelon", "Oh Hai", STARTING_WATER, false}
	return players
}

func MakeDudes() map[string]*Dude {
	types := []string{"antitower", "antidude", "antitower"}
	teams := []string{"apple", "banana", "grape", "watermelon"}
	dudes := map[string]*Dude{}
	for i := 0; i<10; i++ {
		dude := MakeDude(types[rand.Intn(3)])
		dude.PlayerId = teams[rand.Intn(4)]
		dude.Pos = Location{float64(rand.Intn(BOARD_WIDTH)), float64(rand.Intn(BOARD_HEIGHT))}

		x := rand.Float64()*2 - 1.0
		y := rand.Float64()*2 - 1.0
		dir := Direction{x, y, DUDE_VELOCITY}
		dir.Normalize()
		dude.Dir = dir

		dudes[strconv.Itoa(dude.Id)] = &dude
	}
	return dudes
}

func PickBases(m [][]string) map[string]*Base {
	basepos := []Location{}
	for i, row := range m {
		for j, c := range row {
			if c == "B" {
				basepos = append(basepos, Location{float64(j), float64(i)})
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
		bases["apple"] = &Base{"apple", basepos[0], []Tower{}, BASE_HEALTH}
		if len(basepos) > 1 {
			bases["banana"] = &Base{"banana", basepos[1], []Tower{}, BASE_HEALTH}
			if len(basepos) > 2 {
				bases["grape"] = &Base{"grape", basepos[2], []Tower{}, BASE_HEALTH}
				if len(basepos) > 3 {
					bases["watermelon"] = &Base{"watermelon", basepos[3], []Tower{}, BASE_HEALTH}
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

func ProcessControlRequest(state *State, msg ControlRequest, reset chan bool) {
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
			if (player.Water >= DUDE_COST) {
				player.Water -= DUDE_COST
				fmt.Printf("Player %s Bought Dude, Water=%d\n", pid, player.Water)
			} else {
				// TODO: Send back notification that all dudes could not be bought?
				break
			}

			fmt.Printf("BuyDude: %+v\n", bd)

			d := MakeDude(bd.Type)
			d.PlayerId = pid
			d.Pos = Location{state.Bases[pid].Pos.X+1, state.Bases[pid].Pos.Y+1}
			d.Dir = bd.Dir
			d.Dir.V = DUDE_VELOCITY
			d.Dir.Normalize()

			state.Dudes[strconv.Itoa(d.Id)] = &d

			fmt.Printf("Added Dude: %+v\n", d)
		}

	case "BuyTower":
		fmt.Println("Buy Tower!")
		var bt BuyTower
		json.Unmarshal([]byte(msg.Data), &bt)

		pid, ok := state.controlToPlayer[msg.ControlId]
		if !ok {
			data, _ := json.Marshal(BuyReject{"Unregistered Player"})
			resp := ControlResponse{"BuyReject", string(data)}
			websocket.JSON.Send(msg.ws, resp)
			return
		}

		player := state.Players[pid]

		fmt.Printf("BuyTower: %+v\n", bt)

		// Don't buy already bought tower
		for _, tower := range state.Bases[pid].Towers {
			if tower.Pos == bt.Pos {
				fmt.Printf("Already bought tower %d", bt.Pos)
				data, _ := json.Marshal(BuyReject{"Already bought tower"})
				resp := ControlResponse{"BuyReject", string(data)}
				websocket.JSON.Send(msg.ws, resp)
				return
			}
		}

		// Deduct required moneys
		if (player.Water >= TOWER_COST) {
			player.Water -= TOWER_COST
			fmt.Printf("Player %s Bought Tower, Water=%d\n", pid, player.Water)
		} else {
			// TODO: Send back notification that tower could not be bought?
			data, _ := json.Marshal(BuyReject{"Not enough funds"})
			resp := ControlResponse{"BuyReject", string(data)}
			websocket.JSON.Send(msg.ws, resp)
			break
		}

		var tower Tower
		tower.PlayerId = pid
		tower.Pos = bt.Pos
		tower.Health = TOWER_HEALTH
		tower.Type = ""

		state.Bases[pid].Towers = append(state.Bases[pid].Towers, tower)

		fmt.Printf("Added Tower: %+v\n", tower)

		data, _ := json.Marshal(BuyTowerConfirm{bt.Pos})
		resp := ControlResponse{"BuyTowerConfirm", string(data)}
		websocket.JSON.Send(msg.ws, resp)

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

		fmt.Println("Checking player")
		player := state.Players[sp.PlayerId]
		if !player.assigned {
			fmt.Println("Assigned!")
			player.assigned = true
			state.controlToPlayer[msg.ControlId] = player.Id
			data, _ := json.Marshal(ConfirmPlayer{sp.PlayerId})
			resp := ControlResponse{"ConfirmPlayer", string(data)}
			websocket.JSON.Send(msg.ws, resp)
			return
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
			// If all controllers disconnect, reset game state
			if len(state.controlToPlayer) == 0 {
				reset <- true
			}
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
	d1.Health -= d2.Stats.DudeDmg
	d2.Health -= d1.Stats.DudeDmg
}

func DudeVsTower(d *Dude, t *Tower) {
	t.Health -= d.Stats.TowerDmg
	// Tower dmg is fixed
	d.Health -= 5
	if t.Health < 0 {
		fmt.Printf("Tower %s(%d) destroyed!", t.PlayerId, t.Pos)
		t.Health = 0
	}
}

func DudeVsBase(d *Dude, b *Base) {
	b.Health -= d.Stats.BaseDmg
	//fmt.Printf("Dude attacking %s base! health=%d\n", b.Id, b.Health)
	if b.Health < 0 {
		fmt.Printf("Base %s destroyed!\n", b.Id)
		b.Health = 0
	}
}

func DoIt(state *State) bool {
	changed := false

	// Dudes will be left in the map when their health drops to 0 and removed on
	// the next tick. This is the next tick.
	// Remove Dudes with 0 health
	for id, dude := range state.Dudes {
		if dude.Health == 0 {
			changed = true
			delete(state.Dudes, id)
		}
	}

	basepos := map[Location]*Base{}
	towerpos := map[Location]*Tower{}
	for _, base := range state.Bases {
		// All four base squares
		pos := base.Pos.Neighbor(4)
		basepos[pos] = base
		basepos[pos.Neighbor(3)] = base
		basepos[pos.Neighbor(4)] = base
		basepos[pos.Neighbor(5)] = base
		// Towers!
		for _, tower := range base.Towers {
			pos := tower.TowerLocation(base)
			towerpos[pos] = &tower
		}
	}

	//oldpos := GetDudePositionMap(state)
	for _, dude := range state.Dudes {
		changed = true
		pos := dude.Pos.Advance(dude.Dir)
		advance := true
		ipos := Location{math.Floor(pos.X), math.Floor(pos.Y)}

		//
		// Dude Vs. Tower
		//

		//fmt.Println("ipos:", ipos)

		// On tower
		if tower, ok := towerpos[ipos]; ok {
			if tower.Health > 0 {
				if dude.PlayerId != tower.PlayerId {
					DudeVsTower(dude, tower)
					advance = false
				}
			}
		}

		// Nearby tower
		for i := 0; i<8; i++ {
			if tower, ok := towerpos[ipos.Neighbor(i)]; ok {
				if tower.Health > 0 {
					if dude.PlayerId != tower.PlayerId {
						DudeVsTower(dude, tower)
					}
				}
			}
		}

		//
		// Dude Vs. Base
		//
		if base, ok := basepos[ipos]; ok {
			if base.Health > 0 {
				if dude.PlayerId != base.Id {
					DudeVsBase(dude, base)
					advance = false
				}
			}
		}

		// Dude Vs. Dude
		//if d1.PlayerId != d2.PlayerId {
		//}

		// Store Advanced Position
		if advance {
			dude.Pos = pos
		}

		// Subtract 1 health (rot)
		dude.Health -= 1

		// Clamp health at 0
		if dude.Health < 0 {
			dude.Health = 0
		}
	}

	return changed
}

func GameLoop(schan chan State, cmchan chan ControlRequest) {
	rand.Seed(time.Now().Unix())
	m := ReadMap("./map")
	reset := make(chan bool, 10)
	for {
		players := MakePlayers()
		//dudes := MakeDudes()
		dudes := map[string]*Dude{}
		bases := PickBases(m)
		state := State{players, bases, dudes, m, map[int]string{}, true}
		tick := time.Tick(CLOCK_TICK * time.Millisecond)

		// Initial update
		schan <- state

		gloop:
		for {
			select {
			case <-tick:
				// Update positions and fight
				state.changed = DoIt(&state)

				// Award resources
				for _, p := range players {
					p.Water += WATER_PER_ROUND
				}

				// Push state
				schan <- state
			case msg := <-cmchan:
				// Process Events from controllers
				ProcessControlRequest(&state, msg, reset)
			case <-reset:
				break gloop
			}
		}

		fmt.Println("Game Reset")
	}
}

func Serve() {
	schan := make(chan State, 100)
	cmchan := make(chan ControlRequest, 100)

	//cuchans := map[int]chan ControlUpdate{}

	bcast := MakeStateBroadcaster(schan)
	//bcast := MakeBroadcaster(schan)

	go GameLoop(schan, cmchan)

	screenHandler := MakeScreenHandler(bcast)
	controlHandler := MakeControlHandler(cmchan)

	http.Handle("/screen", websocket.Handler(screenHandler))
	http.Handle("/control", websocket.Handler(controlHandler))
	http.Handle("/", http.FileServer(http.Dir("./static")))
	fmt.Println("Starting server")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal("Server failed:", err)
	}
	fmt.Println("exiting server")
}
