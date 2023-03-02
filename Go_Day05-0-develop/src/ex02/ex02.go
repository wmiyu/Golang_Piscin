package main

import (
	"container/heap"
	"fmt"

)

type Present struct {
    Value int
    Size int
}

type Item struct {
	node     *Present
	priority int
	index    int
}

// An Item is something we manage in a priority queue.
// type Item struct {
// 	value    string // The value of the item; arbitrary.
// 	priority int    // The priority of the item in the queue.
// 	// The index is needed by update and is maintained by the heap.Interface methods.
// 	index int // The index of the item in the heap.
// }

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.

	if pq[i].priority == pq[j].priority {
		return pq[i].node.Size < pq[j].node.Size
	}
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, node *Present, priority int) {
	item.node = node
	item.priority = priority
	heap.Fix(pq, item.index)
}

func getNCoolestPresents(ph []*Present, N int) []*Present {
	if N <= 0 || len(ph) < N {
		return nil
	}
	pq := PriorityQueue{}
	heap.Init(&pq)

	for i:=0; i < len(ph); i++ {
		heap.Push(&pq, &Item{ph[i], ph[i].Value, 0})
	}
	result := []*Present{}
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		if len(result) < N {
			result = append(result, item.node)
		}
		fmt.Println(item.index, item.priority, item.node)
	}
	return result
}

func presentHeap1() []*Present {
	p := []*Present{&Present{5,10}, 
					&Present{4,5},
					&Present{3,1},
					&Present{5,2},
					&Present{5,5},}
	return p
}

func main() {

	ph1 := presentHeap1()

	for i:=0; i < len(ph1); i++ {
		fmt.Println(ph1[i])
	}
	fmt.Println("...: getNCoolestPresents :...")
	ph_coolest := getNCoolestPresents(ph1, 5)
	if ph_coolest == nil {
		fmt.Println("Error N of presents")
		return
	}
	fmt.Println("...: Coolest Presents :...")
	for _, item := range ph_coolest {
		fmt.Println(item)
	}
}
