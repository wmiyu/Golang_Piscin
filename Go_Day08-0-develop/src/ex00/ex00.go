package main

import (
	"fmt"
	"unsafe"
)

func getElement(arr []int, idx int) (int, error) {
	if len(arr) == 0 {
		return 0, fmt.Errorf("EMPTY SLICE")
	}
	if idx < 0 {
		return 0, fmt.Errorf("IDX IS NEGATIVE")
	}
	start := unsafe.Pointer(&arr[0])
	size := unsafe.Sizeof(int(0))
	for i := 0; i < len(arr); i++ {
		item := *(*int)(unsafe.Pointer(uintptr(start) + size * uintptr(i)))
		if i == idx {
			return item, nil
		}
	}
	return 0, fmt.Errorf("IDX OUT OF BOUNDS")
}

func testCase(arr1 []int, index int) () {
	e, err := getElement(arr1, index)
	if err != nil {
		fmt.Printf("But some error occurs : %s\n", err.Error())
	} else {
		fmt.Printf(" :) YOUR VALUE: (%v)\n", e)
	}
}

func main() {
	arr1 := []int{11, 23, 44, 66, 88}

	testCase(arr1, 0)
	testCase(arr1, 1)
	testCase(arr1, 2)
	testCase(arr1, 3)
	testCase(arr1, 4)
	testCase(arr1, 5)
	testCase(arr1, -6)
	arr2 := []int{}
	testCase(arr2, 0)
}