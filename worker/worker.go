package worker

import (
	"TOMS/network"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Worker struct {
	mu        sync.Mutex
	id        int
	nodes     []string
	myAddress string
}

func (w *Worker) requestHandler(request string) {
	if len(request) == 0 {
		log.Fatal("[Manager][Handler] Error: Empty request")
	}

	if request[0] == '0' {
		w.init(strings.Split(request, ";")[1:])
	}
	if request[0] == '1' {
		w.addNode(strings.Split(request, ";")[1])
	}
}

func (w *Worker) addNode(node string) {
	w.mu.Lock()
	w.nodes = append(w.nodes, node)
	log.Println("[Worker][AddNode][0] ", node)
	log.Println("[Worker][AddNode][1] ", w.nodes)
	w.mu.Unlock()
}

func (w *Worker) init(workers []string) {
	w.id, _ = strconv.Atoi(workers[0])
	w.nodes = append(w.nodes, workers[1:]...)
	log.Println("[Worker][AddNode][1] ", w.id, " ", w.nodes)
	w.mu.Unlock()
}

func (w *Worker) StartWorker(managerAddress string) {
	w.mu.Lock()
	var s network.Server
	s.Init(w.requestHandler)
	w.myAddress = s.GetAddress()
	request := "0;" + w.myAddress + "\n"
	go func() {
		time.Sleep(time.Second)
		network.SendMessage(request, managerAddress)
	}()
	s.Serve()
}
