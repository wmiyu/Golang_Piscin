package main

import "fmt"

type Present struct {
    size   int64
	value  int64
}

func (p Present) Weight() int64 {
	return (p.size)
}

func (p Present) Value() int64 {
	return (p.value)
}

func TestPresents1 () *[]Present {
	items := []Present {
		Present{3,5},
		Present{2,3},
		Present{1,4},
		}
		return &items
}

func grabPresents(presents *[]Present, capacity int64) []int64 {
	if presents == nil || capacity <= 0 || len(*presents) == 0{
		return nil
	}
	items := []Packable{}
	for _, present := range *presents {
		items = append(items, Packable(present))
	}
	indxs := Knapsack(items, capacity)
	return indxs
}


func main() {

	fmt.Println("..: TestPresents 1 (size, value) capacity :5 :..")
	somePresents := TestPresents1()
	for idx, present := range *somePresents {
		fmt.Println(idx + 1, present)
	}
	knaps1 := grabPresents(somePresents, 5)
	for _, idx := range knaps1 {
		fmt.Println("You'd better take #:", idx + 1)
	}
	fmt.Println("..: TestPresents 2 (size, value) capacity :3 :..")
	somePresents = TestPresents1()
	for idx, present := range *somePresents {
		fmt.Println(idx + 1, present)
	}
	knaps1 = grabPresents(somePresents, 3)
	for _, idx := range knaps1 {
		fmt.Println("You'd better take #:", idx + 1)
	}
	fmt.Println("..: TestPresents 3 (size, value) capacity :0 :..")
	somePresents = TestPresents1()
	for idx, present := range *somePresents {
		fmt.Println(idx + 1, present)
	}
	knaps1 = grabPresents(somePresents, 0)
	for _, idx := range knaps1 {
		fmt.Println("You'd better take #:", idx + 1)
	}
	fmt.Println("..: TestPresents 4 (size, value) NO PRESENTS :..")
	knaps1 = grabPresents(nil, 100)

}