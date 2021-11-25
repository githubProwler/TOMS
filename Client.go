package main

import (
	"TOMS/myRPC"
	"log"
)

func startClient() {
	var reply AddColorResponse
	var client myRPC.Client

	client.InitClient("localhost:1234", 5000, true)
	// Create a TCP connection to localhost on port 1234
	// client, err := rpc.Dial("tcp", "localhost:1234")
	// if err != nil {
	// log.Fatal("Connection error: ", err)
	// }
	var rqs AddColorRequest

	rqs.RedAmount = 255
	client.Call("Message.AddColor", rqs, &reply)
	log.Println("AddColor", rqs, reply)

	rqs.RedAmount = 0
	client.Call("Message.AddColor", rqs, &reply)
	log.Println("AddColor", rqs, reply)

	rqs.RedAmount = 0
	client.Call("Message.AddColor", rqs, &reply)
	log.Println("AddColor", rqs, reply)

	rqs.RedAmount = 100
	client.Call("Message.AddColor", rqs, &reply)
	log.Println("AddColor", rqs, reply)

	rqs.RedAmount = 200
	client.Call("Message.AddColor", rqs, &reply)
	log.Println("AddColor", rqs, reply)

	rqs.RedAmount = 255
	client.Call("Message.AddColor", rqs, &reply)
	log.Println("AddColor", rqs, reply)

	rqs.RedAmount = 0
	client.Call("Message.AddColor", rqs, &reply)
	log.Println("AddColor", rqs, reply)

	rqs.RedAmount = 250
	client.Call("Message.AddColor", rqs, &reply)
	log.Println("AddColor", rqs, reply)

	rqs.RedAmount = 0
	client.Call("Message.AddColor", rqs, &reply)
	log.Println("AddColor", rqs, reply)

	rqs.RedAmount = 0
	client.Call("Message.AddColor", rqs, &reply)
	log.Println("AddColor", rqs, reply)

	rqs.RedAmount = 200
	client.Call("Message.AddColor", rqs, &reply)
	log.Println("AddColor", rqs, reply)

	rqs.RedAmount = 0
	client.Call("Message.AddColor", rqs, &reply)
	log.Println("AddColor", rqs, reply)

	rqs.RedAmount = 0
	client.Call("Message.AddColor", rqs, &reply)
	log.Println("AddColor", rqs, reply)

	rqs.RedAmount = 0
	client.Call("Message.AddColor", rqs, &reply)
	log.Println("AddColor", rqs, reply)
}
