package myRPC

import (
	"log"
	"net"
	"net/rpc"
	"time"
)

type Client struct {
	rpcClient     *rpc.Client
	address       string
	jitter        int
	jitterEnabled bool
}

func StartServer(rcvr interface{}, address string) {
	err := rpc.Register(rcvr)
	if err != nil {
		log.Fatal("Format of service is incorrect. ", err)
	}

	// rpc.HandleHTTP()

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Listen error: ", err)
	}
	log.Printf("Serving RPC server on port %s.", address)

	// err = http.Serve(listener, nil)
	rpc.Accept(listener)
	if err != nil {
		log.Fatal("Error serving: ", err)
	}
}

func (client *Client) InitClient(address string, jitter int, jitterEnabled bool) {
	var err error
	client.rpcClient, err = rpc.Dial("tcp", address)
	// client.rpcClient, err = rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	client.address = address
	client.jitter = jitter
	client.jitterEnabled = jitterEnabled
}

func (client *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	if client.jitterEnabled {
		time.Sleep(time.Millisecond * time.Duration(client.jitter))
	}
	return client.rpcClient.Call(serviceMethod, args, reply)
}
