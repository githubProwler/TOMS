package worker

import (
	"strconv"
	"strings"
)

func (w *Worker) makeReference() string {
	return strconv.Itoa(w.messageCounter) + "." + strconv.Itoa(w.id)
}

func (w *Worker) makeProposal() string {
	proposal := strconv.Itoa(w.next) + "." + strconv.Itoa(w.id)
	w.next++
	return proposal
}

func comparePriority(pref1 string, pref2 string) bool {

	split1 := strings.Split(pref1, ".")
	split2 := strings.Split(pref2, ".")
	num1, _ := strconv.Atoi(split1[0])
	num2, _ := strconv.Atoi(split2[0])
	if num1 > num2 {
		return false
	}
	if num1 < num2 {
		return true
	}
	id1, _ := strconv.Atoi(split1[1])
	id2, _ := strconv.Atoi(split2[1])
	return id1 < id2
}

func getProposalNumber(proposal string) int {
	num, _ := strconv.Atoi(strings.Split(proposal, ".")[0])
	return num
}
