package worker

import (
	"TOMS/network"
	"log"
	"strconv"
	"strings"
)

func (w *Worker) requestHandler(request string) {
	if len(request) == 0 {
		log.Fatal("[Manager][Handler] Error: Empty request")
	}
	args := strings.Split(request, ";")[1:]
	switch request[0] {
	case INITIAL_STATE:
		w.init(args)
	case ADD_NEW_WORKER:
		w.addNode(args[0])
	case DELIVER_WITHOUT_ISIS:
		w.deliver(args[0])
	case REQUEST:
		w.propose(args)
	case PROPOSAL:
		w.receiveProposal(args[0])
	case AGREED:
		w.receiveAgreed(args)
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

func (w *Worker) propose(args []string) {
	from, reference, message := args[0], args[1], args[2]
	w.mu.Lock()
	proposal := w.makeProposal()
	w.messages.AddNode(reference, proposal, message)

	finalMessage := string(PROPOSAL) + ";" + proposal + "\n"
	network.SendMessage(finalMessage, from)
	w.mu.Unlock()
}

func (w *Worker) receiveAgreed(args []string) {
	reference, agreed := args[0], args[1]
	w.mu.Lock()
	agreedNumber := getProposalNumber(agreed)
	if agreedNumber >= w.next {
		w.next = agreedNumber + 1
	}
	w.messages.UpdateNode(reference, agreed)
	for w.messages.Poppable() {
		message := w.messages.Pop()
		w.deliver(message)
	}
	w.mu.Unlock()
}

func (w *Worker) receiveProposal(proposal string) {
	w.mu.Lock()
	w.waitingFor--

	if w.messagePriority == "" || comparePriority(w.messagePriority, proposal) {
		w.messagePriority = proposal
	}

	if w.waitingFor > 0 {
		w.mu.Unlock()
		return
	}

	finalMessage := string(AGREED) + ";" + w.makeReference() + ";" + w.messagePriority + "\n"
	for _, node := range w.nodes {
		go network.SendMessage(finalMessage, node)
	}
	w.mu.Unlock()
	w.messageLock.Unlock()
}
