package main

import (
	"container/heap"
	"fmt"
)

type TreeNode struct {
	Val    int
	HasToy bool
	Left   *TreeNode
	Right  *TreeNode
}

type Item struct {
	node     *TreeNode
	priority int
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {

	return false
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
func (pq *PriorityQueue) update(item *Item, node *TreeNode, priority int) {
	item.node = node
	item.priority = priority
	heap.Fix(pq, item.index)
}

func unrollGarland(root *TreeNode) {
	item := &Item{root, 0, 0}
	cs := PriorityQueue{item}
	ns := PriorityQueue{}

	heap.Init(&cs)
	heap.Init(&ns)
	for cs.Len() > 0 || ns.Len() > 0 {
		for cs.Len() > 0 {
			item = heap.Remove(&cs, cs.Len()-1).(*Item)
			fmt.Println("(", item.node.Val, "\t", item.node.HasToy, ")")
			if item.node.Right != nil {
				heap.Push(&ns, &Item{item.node.Right, 0, 0})
			}
			if item.node.Left != nil {
				heap.Push(&ns, &Item{item.node.Left, 0, 0})
			}
		}
		for ns.Len() > 0 {
			item = heap.Remove(&ns, ns.Len()-1).(*Item)
			fmt.Println("(", item.node.Val, "\t", item.node.HasToy, ")")
			if item.node.Left != nil {
				heap.Push(&cs, &Item{item.node.Left, 0, 0})
			}
			if item.node.Right != nil {
				heap.Push(&cs, &Item{item.node.Right, 0, 0})
			}
		}
	}
}

func sampleTree1() *TreeNode {
	t := &TreeNode{10, true,
		&TreeNode{5, true,
			&TreeNode{4, true, nil, nil},
			&TreeNode{6, false, nil, nil}},
		&TreeNode{20, false,
			&TreeNode{15, true, nil, nil},
			&TreeNode{25, true, nil, nil}}}
	fmt.Println(`
      1 0
    /     \
   5      20
 /  \    /  \
4    6  15  25 `)
	return t
}

func main() {

	st1 := sampleTree1()
	fmt.Println("\n...: unrollGarland :...")
	fmt.Println()
	unrollGarland(st1)
}
