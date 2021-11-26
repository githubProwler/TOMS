package main

import (
	"TOMS/colouredCircle"
	"TOMS/manager"
	"TOMS/worker"
	"flag"
	"log"
	"strconv"
	"strings"
)

type Message struct {
	cc *colouredCircle.ColouredCircle
}

type AddColorRequest struct {
	RedAmount int
}

func cbk(inputString string, args interface{}) {
	inputColor := strings.TrimSpace(inputString)
	inputNumber, _ := strconv.Atoi(inputColor)
	inputNumber = inputNumber % 256

	w, ok := args.(*worker.Worker)
	if !ok {
		log.Fatal("HELLO MOTHERFUCKER")
	}

	w.BMulticast(inputNumber)
}

func deliver(inputNumber int, args interface{}) {
	cc, ok := args.(*colouredCircle.ColouredCircle)
	if !ok {
		log.Fatal("Argument was not colouredCircle type")
	}

	cc.AddColor(inputNumber)
}

type AddColorResponse struct {
	Success bool
}

func (msg *Message) AddColor(rqs AddColorRequest, rsp *AddColorResponse) error {
	msg.cc.AddColor(rqs.RedAmount)
	return nil
}

func main() {
	// var client myRPC.Client
	msg := new(Message)
	msg.cc = new(colouredCircle.ColouredCircle)
	server := flag.Bool("server", false, "Set to run program as a server")
	managerAddress := flag.String("mAddr", "", "Manager server address")
	flag.Parse()

	if *server {
		log.Println("Starting a server")
		// myRPC.StartServer(msg, ":1234")
		startManager()
	} else {
		if len(*managerAddress) > 9 {
			var w worker.Worker
			go msg.cc.Main("Server", cbk, &w)
			log.Println("Starting a client, manager: ", *managerAddress)
			// startClient()
			testClient(*managerAddress)
			w.StartWorker(*managerAddress, deliver, msg.cc)
		}
	}
}

func testServer() {
	// network.StartServer()
}

func startManager() {
	var m manager.Manager
	m.StartManager()
}

func testClient(manager string) {
	// var w worker.Worker
	// w.StartWorker(manager, deliver)
}
