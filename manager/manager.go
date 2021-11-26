package manager

import (
	"TOMS/network"
	"log"
	"strconv"
	"strings"
	"sync"
)

type Manager struct {
	mu      sync.Mutex
	workers []string
	n       int
	length  int
}

func (m *Manager) addNewWorker(ip string) {
	m.mu.Lock()
	m.n++
	m.workers = append(m.workers, ip)
	m.length += len(ip) + 1
	go sendInitialState(m.workers, m.n, m.length)
	go updateAllWorkers(m.n, ip, m.workers[:m.n-1])
	log.Println("[Manager][addNewWorker][0] ", ip)
	m.mu.Unlock()
}

func sendInitialState(workers []string, n int, length int) {
	var request strings.Builder
	reqLength := length + 3 + len(strconv.Itoa(n))
	request.Grow(reqLength)
	request.WriteString("0;")
	request.WriteString(strconv.Itoa(n))
	for _, address := range workers {
		request.WriteString(";" + address)
	}
	request.WriteRune('\n')
	network.SendMessage(request.String(), workers[n-1])
}

func updateAllWorkers(n int, ip string, workers []string) {
	request := "1;" + ip + "\n"

	for _, addr := range workers[:n-1] {
		log.Println("[Manager][updateWorker][0] Rcvr: ", addr, " Msg: ", ip)
		go network.SendMessage(request, addr)
	}
}

func (m *Manager) requestHandler(request string) {
	if len(request) == 0 {
		log.Fatal("[Manager][Handler] Error: Empty request")
	}

	if request[0] == '0' {
		newIp := strings.Split(request, ";")[1]

		log.Printf("[Manager][Handler] Request: \"%s\"\n", request)

		m.addNewWorker(newIp)
		return
	}

	log.Printf("[Manager][Handler] Unrecognized request: \"%s\"\n", request)
}

func (m *Manager) StartManager() {
	m.n = 0
	m.length = 0

	var s network.Server
	s.Init(m.requestHandler)
	log.Println("Manager Address is " + s.GetAddress())
	s.Serve()
}
