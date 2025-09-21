package main

import "container/heap"

type PriorityQueue []HuffBaseNode

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Freq() < pq[j].Freq()
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	*pq = append(*pq, x.(HuffBaseNode))
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[0]
	old[0] = old[n-1]
	*pq = old[0 : n-1]
	heap.Fix(pq, 0)
	return item
}
