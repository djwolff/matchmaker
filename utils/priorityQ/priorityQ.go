package priorityQ

import (
	"container/heap"

	"github.com/djwolff/matchmaker/models/mm"
)

type PriorityQueue []*mm.Player

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].OfferingRole > pq[j].OfferingRole
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*mm.Player)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// func (pq PriorityQueue) Push(x any) {
// 	n := len(pq)
// 	item := x.(*mm.Player)
// 	item.Index = n
// 	pq = append(pq, item)
// }

// func (pq PriorityQueue) Pop() any {
// 	old := pq
// 	n := len(old)
// 	item := old[n-1]
// 	old[n-1] = nil  // avoid memory leak
// 	item.Index = -1 // for safety
// 	pq = old[0 : n-1]
// 	return item
// }

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *mm.Player, offeringRole string) {
	// item.Conn = conn
	item.OfferingRole = offeringRole
	heap.Fix(pq, item.Index)
}
