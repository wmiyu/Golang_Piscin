package main

import (
	"fmt"
	"time"
)

func sleepSort(slice []int) chan string {
	message := make(chan string)

	for i := 0; i < len(slice); i++ {
		go sleepPrint(slice[i], &message)
	}
	return message
}

func sleepPrint(delay int, message *chan string) {
	time.Sleep(time.Millisecond *1000 * time.Duration(delay))
	*message <- fmt.Sprintf("#:_ value = delay: %d", delay)
}

func main() {
	slice := []int{5,7,3,9,4,8,2,6}
	fmt.Println(slice)
	message := sleepSort(slice)
	for i := 0; i < len(slice); i++{
		msg, open := <-message
		if !open {
			break
		}
		fmt.Println(msg)
	}
}
