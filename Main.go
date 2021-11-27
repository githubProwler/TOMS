package main

import (
	"TOMS/colouredCircle"
	"TOMS/manager"
	"TOMS/worker"
	"flag"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func setColor(inputString string, args interface{}) {
	cc, ok := args.(*colouredCircle.ColouredCircle)
	if !ok {
		log.Fatal("Passed wrong type of args to setColor callback")
	}

	inputColor := strings.TrimSpace(inputString)
	inputNumber, _ := strconv.Atoi(inputColor)
	inputNumber = inputNumber % 256

	cc.AddColor(inputNumber)
}

func cbk(inputString string, args interface{}) {
	inputColor := strings.TrimSpace(inputString)
	inputNumber, _ := strconv.Atoi(inputColor)
	inputNumber = inputNumber % 256

	w, ok := args.(*worker.Worker)
	if !ok {
		log.Fatal("Passed wrong type of args to ckb callback")
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

func startManager() {
	var m manager.Manager
	m.StartManager()
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	var cc colouredCircle.ColouredCircle
	server := flag.Bool("server", false, "Set to run program as a server")
	managerAddress := flag.String("mAddr", "", "Manager server address")
	flag.Parse()

	if *server {
		log.Println("Starting a server")
		startManager()
	} else {
		if len(*managerAddress) > 9 {
			var w worker.Worker
			go cc.Main("Server", cbk, &w)
			log.Println("Starting a client, manager: ", *managerAddress)
			w.StartWorker(*managerAddress, deliver, &cc)
		}
	}
}
