package worker

import (
	"TOMS/network"
	"TOMS/pqueue"
	"strconv"
	"sync"
	"time"
)

type Worker struct {
	mu              sync.Mutex
	id              int
	nodes           []string
	myAddress       string
	deliverFn       func(int, interface{})
	deliverArgs     interface{}
	messageCounter  int
	next            int
	messages        pqueue.PQueue
	messageLock     sync.Mutex
	waitingFor      int
	messagePriority string
}

func (w *Worker) SendReliable(message int) {
	w.messageLock.Lock()
	w.mu.Lock()
	finalMessage := string(REQUEST) + ";" + w.myAddress + ";" + w.makeReference() + ";" + strconv.Itoa(message) + "\n"
	w.waitingFor = len(w.nodes)
	w.messagePriority = ""

	for _, node := range w.nodes {
		go network.SendMessage(finalMessage, node)
	}
	w.mu.Unlock()
}

func (w *Worker) BMulticast(message int) {
	finalMessage := string(DELIVER_WITHOUT_ISIS) + ";" + strconv.Itoa(message) + "\n"
	for _, node := range w.nodes {
		go network.SendMessage(finalMessage, node)
	}
}

func (w *Worker) IsFree() bool {
	w.mu.Lock()
	ans := w.waitingFor == 0
	w.mu.Unlock()
	return ans
}

func (w *Worker) StartWorker(managerAddress string, deliverFn func(int, interface{}), deliverArgs interface{}) {
	w.mu.Lock()
	var s network.Server
	s.Init(w.requestHandler)
	w.myAddress = s.GetAddress()
	w.deliverFn = deliverFn
	w.deliverArgs = deliverArgs
	w.messagePriority = ""
	request := string(INITIAL_STATE) + ";" + w.myAddress + "\n"
	go func() {
		time.Sleep(time.Second)
		network.SendMessage(request, managerAddress)
	}()
	s.Serve()
}
