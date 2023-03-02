package main

import "fmt"

type TreeNode struct {
	HasToy	bool
	Left	*TreeNode
	Right	*TreeNode
}

func getNodeSumm(tr *TreeNode) int {
	var sum int

	if tr == nil {
		return 0
	}

	switch tr.HasToy {
		case true : sum = 1
		case false: sum = 0 
	}

	if tr.Left == nil && tr.Right == nil {
		return sum
	}

	if tr.Left != nil {
		sum += getNodeSumm(tr.Left)
	}
	if tr.Right != nil {
		sum += getNodeSumm(tr.Right)
	}
	return sum
}

func areToysBalanced(tr *TreeNode) bool {
	leftsum := getNodeSumm(tr.Left)
	rightsum := getNodeSumm(tr.Right)
	fmt.Printf("TOYS BALANCE = LEFT: %d RIGHT: %d \n", leftsum, rightsum)
	return leftsum == rightsum
}

func main() {
	fmt.Println("HAPPY NEW YEAR!")
	
	st1 := sampleTree1()
	fmt.Println(areToysBalanced(st1))
	st2 := sampleTree2()
	fmt.Println(areToysBalanced(st2))
	st3 := sampleTree3()
	fmt.Println(areToysBalanced(st3))
	st4 := sampleTree4()
	fmt.Println(areToysBalanced(st4))

}