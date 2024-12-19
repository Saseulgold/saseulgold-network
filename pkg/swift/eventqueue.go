package swift

import (
	"sync"
)

type Priority int

const (
	LowPriority    Priority = 1
	NormalPriority Priority = 2
	HighPriority   Priority = 3
)

type PriorityMessage struct {
	message  *Message
	priority Priority
	index    int
}

type PriorityQueue struct {
	items []*PriorityMessage
	mu    sync.RWMutex
}

func (pq *PriorityQueue) Len() int { return len(pq.items) }

func (pq *PriorityQueue) Less(i, j int) bool {
	return pq.items[i].priority > pq.items[j].priority
}

func (pq *PriorityQueue) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].index = i
	pq.items[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(pq.items)
	item := x.(*PriorityMessage)
	item.index = n
	pq.items = append(pq.items, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := pq.items
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	pq.items = old[0 : n-1]
	return item
}
