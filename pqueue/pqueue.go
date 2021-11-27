package pqueue

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Reference struct {
	num int
	id  int
}

func makeReference(s string) *Reference {
	var r Reference
	split := strings.Split(s, ".")
	r.num, _ = strconv.Atoi(split[0])
	r.id, _ = strconv.Atoi(split[1])
	return &r
}

func (r *Reference) compare(r2 *Reference) int {
	if r.num < r2.num {
		return -1
	}
	if r.num > r2.num {
		return 1
	}
	if r.id < r2.id {
		return -1
	}
	if r.id > r2.id {
		return 1
	}
	return 0
}

type pqNode struct {
	next              *pqNode
	prev              *pqNode
	uniqueReference   Reference
	priorityReference Reference
	ready             bool
	message           string
}

type PQueue struct {
	head *pqNode
}

func (pq *PQueue) AddNode(reference string, priority string, message string) {
	var pqn = &pqNode{nil, nil, *makeReference(reference), *makeReference(priority), false, message}
	if pq.head == nil {
		pq.head = pqn
		return
	}
	if pq.head.priorityReference.compare(&pqn.priorityReference) > 0 {
		pq.head.prev = pqn
		pqn.next = pq.head
		pq.head = pqn
		log.Printf("******** pq head changed while adding a node. THIS SHOULDN'T HAPPEN ********")
		return
	}
	pq.insertNode(pq.head, pqn)
}

func (pq *PQueue) insertNode(startNode *pqNode, newNode *pqNode) {
	for !(startNode.next == nil || startNode.next.priorityReference.compare(&newNode.priorityReference) > 0) {
		startNode = startNode.next
	}
	newNode.next = startNode.next
	newNode.prev = startNode
	if startNode.next != nil {
		startNode.next.prev = newNode
	}
	startNode.next = newNode
}

func (pq *PQueue) UpdateNode(reference string, priority string) {
	ref := makeReference(reference)
	cur := pq.head
	for cur != nil && cur.uniqueReference.compare(ref) != 0 {
		cur = cur.next
	}
	if cur == nil {
		log.Fatal("There's no reference in priority queue")
	}

	cur.priorityReference = *makeReference(priority)
	cur.ready = true

	if cur.next == nil || cur.next.priorityReference.compare(&cur.priorityReference) > 0 {
		return
	}
	cur.next.prev = cur.prev
	if cur.prev != nil {
		cur.prev.next = cur.next
		pq.insertNode(cur.prev, cur)
	} else {
		pq.head = cur.next
		pq.insertNode(pq.head, cur)
	}
}

func (pq *PQueue) Print() {
	cur := pq.head
	for cur != nil {
		fmt.Println("( node: ", cur, " pr: ", cur.priorityReference, " ref: ", cur.uniqueReference)
		cur = cur.next
	}
}

func (pq *PQueue) Popable() bool {
	if pq.head == nil {
		return false
	}
	return pq.head.ready
}

func (pq *PQueue) Pop() string {
	if pq.head == nil {
		log.Fatal("Nothing to pop from queue")
	}
	cur := pq.head
	pq.head = cur.next
	if pq.head != nil {
		pq.head.prev = nil
	}
	return cur.message
}
