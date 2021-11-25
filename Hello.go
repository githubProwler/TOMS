package main

import (
	"TOMS/colouredCircle"
	"TOMS/myRPC"
	"flag"
	"log"
	"strconv"
	"strings"
)

type ToDo struct {
	Title  string
	Status string
}

type EditToDo struct {
	Title     string
	NewTitle  string
	NewStatus string
}

type Task int

var todoSlice []ToDo

// GetToDo takes a string type and returns a ToDo
func (t *Task) GetToDo(title string, reply *ToDo) error {
	log.Print("[GetToDo]", title)
	var found ToDo
	for _, v := range todoSlice {
		if v.Title == title {
			found = v
		}
	}

	*reply = found
	return nil
}

// MakeToDo takes a ToDo type and appends to the todoArray
func (t *Task) MakeToDo(todo ToDo, reply *ToDo) error {
	log.Print("[MakeToDo]", todo)
	todoSlice = append(todoSlice, todo)
	*reply = todo
	return nil
}

// EditToDo takes a string type and a ToDo type and edits an item in the todoArray
func (t *Task) EditToDo(todo EditToDo, reply *ToDo) error {
	log.Print("[EditToDo]", todo)
	var edited ToDo
	// 'i' is the index in the array and 'v' the value
	for i, v := range todoSlice {
		if v.Title == todo.Title {
			todoSlice[i] = ToDo{todo.NewTitle, todo.NewStatus}
			edited = ToDo{todo.NewTitle, todo.NewStatus}
		}
	}
	// edited will be the edited ToDo or a zeroed ToDo
	*reply = edited
	return nil
}

func (t *Task) DeleteToDo(todo ToDo, reply *ToDo) error {
	log.Print("[DeleteToDo]", todo)
	var deleted ToDo
	for i, v := range todoSlice {
		if v.Title == todo.Title && v.Status == todo.Status {
			todoSlice = append(todoSlice[:i], todoSlice[i+1:]...)
			deleted = todo
			break
		}
	}
	*reply = deleted
	return nil
}

func cbk(cc *colouredCircle.ColouredCircle, inputString string) {
	inputColor := strings.TrimSpace(inputString)
	inputNumber, _ := strconv.Atoi(inputColor)
	inputNumber = inputNumber % 256

	cc.AddColor(inputNumber)
}

type Message struct {
	cc *colouredCircle.ColouredCircle
}

type AddColorRequest struct {
	RedAmount int
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
	flag.Parse()

	if *server {
		go msg.cc.Main("Server", cbk)
		log.Println("Starting a server")
		myRPC.StartServer(msg, ":1234")
	} else {
		log.Println("Starting a client")
		startClient()
	}
}
