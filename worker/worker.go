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
	mu          sync.Mutex
	id          int
	nodes       []string
	myAddress   string
	deliverFn   func(int, interface{})
	deliverArgs interface{}
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
	if request[0] == '2' {
		w.deliver(strings.Split(request, ";")[1])
	}
}

func (w *Worker) deliver(request string) {
	requestNum, _ := strconv.Atoi(request)
	log.Println("[Worker][Deliver] Delivering request:", requestNum, request)
	w.deliverFn(requestNum, w.deliverArgs)
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

func (w *Worker) BMulticastRequest(messageType int, message int, reference string) {
	finalMessage := strconv.Itoa(messageType) + ";" + w.myAddress + ";" + reference + ";" + strconv.Itoa(message) + "\n"
	w.mu.Lock()
	for _, node := range w.nodes {
		if node == w.myAddress {
			continue
		}
		go network.SendMessage(finalMessage, node)
	}
	w.mu.Unlock()
}

func (w *Worker) BMulticast(message int) {
	finalMessage := "2;" + strconv.Itoa(message)
	for _, node := range w.nodes {
		if node == w.myAddress {
			continue
		}
		go network.SendMessage(finalMessage, node)
	}
}

// func (w *Worker) BMulticastAgree(messageType int, message int, reference string, agreedNumber string) {
// 	finalMessage := strconv.Itoa(messageType) + ";" + reference + ";" + agreedNumber + ";" + strconv.Itoa(message) + "\n"
// 	w.mu.Lock()
// 	for _, node := range w.nodes {
// 		if node == w.myAddress {
// 			continue
// 		}
// 		go network.SendMessage(finalMessage, node)
// 	}
// 	w.mu.Unlock()
// }

func (w *Worker) StartWorker(managerAddress string, deliverFn func(int, interface{}), deliverArgs interface{}) {
	w.mu.Lock()
	var s network.Server
	s.Init(w.requestHandler)
	w.myAddress = s.GetAddress()
	w.deliverFn = deliverFn
	w.deliverArgs = deliverArgs
	request := "0;" + w.myAddress + "\n"
	go func() {
		time.Sleep(time.Second)
		network.SendMessage(request, managerAddress)
	}()
	s.Serve()
}
