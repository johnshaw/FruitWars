package main

type NewChan struct {
	id int
	c chan interface{}
}

type Broadcaster struct {
	newChan chan NewChan
	delChan chan int
}

func (b *Broadcaster) GetChan() (int, chan interface{}) {
	// This channel will receive state data and publish it to the client
	c := make(chan interface{})
	id := rand.Int()
	b.newChan <- NewClientChan{id, c}
	return id, c
}

func (b *Broadcaster) DelChan(id int) {
	b.delChan <- id
}

func (b *Broadcaster) broadcastLoop(in chan interface{}) {
	out := map[int]chan interface{}
	for {
		select {
		case cnew := <-b.newChan:
			out[cnew.id] = cnew.c
		case cdel := <-b.delChan:
			delete(out, id)
		case d := <-in:
			for _, c := range out {
				c <- d
			}
		}
	}
}

func MakeBroadcaster(c chan interface{}) *Broadcaster {
	b := new(Broadcaster)
	b.newChan = make(chan NewClientChan)
	b.delChan = make(chan int)
	go b.broadcastLoop(c)
	return b
}
