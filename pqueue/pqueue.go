package pqueue

import (
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
	cur := pq.head
	insertNode(cur, pqn)
}

func insertNode(startNode *pqNode, newNode *pqNode) {
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
	if cur.next == nil || cur.next.priorityReference.compare(&cur.priorityReference) > 0 {
		return
	}
	cur.next.prev = cur.prev
	if cur.prev != nil {
		cur.prev.next = cur.next
		insertNode(pq.head, cur)
	} else {
		pq.head = cur.next
		insertNode(cur.prev, cur)
	}
}
