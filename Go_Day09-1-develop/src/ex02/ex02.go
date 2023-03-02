package main

import (
	"fmt"
	"sync"
	"reflect"
	"strings"
)

type UnknownPlant struct {
    FlowerType  string
    LeafType    string
    Color       int `color_scheme:"rgb"`
}

type AnotherUnknownPlant struct {
    FlowerColor int
    LeafType    string
    Height      int `unit:"inches"`
}

func describePlant(x interface{}) () {
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)
	for i:=0; i < t.NumField(); i++ {
		fmt.Printf("%s", t.Field(i).Name)
		if t.Field(i).Tag != "" {
			tag := string(t.Field(i).Tag)
			tag = strings.Replace(tag, ":", "=", 1)
			tag = strings.Replace(tag, "\"", "", -1)
			fmt.Printf("(%s)", tag)
		}
		fmt.Printf(":%v", v.Field(i) )
		fmt.Println()
	}
}

func multiplex(in ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	var wg = &sync.WaitGroup{}

	for _, c := range in {
			wg.Add(1)
			go func(c <-chan interface{}) {
					for d := range c {
							out <- d
					}
					wg.Done()
			}(c)
	}
	go func() {
			wg.Wait()
			close(out)
	}()
	return out
}

func main() {

	fmt.Printf(" === Starting multiplex ==== \n") 
	inChannel1 := make(chan interface{})
	inChannel2 := make(chan interface{})

	unknownPlant := UnknownPlant{FlowerType: "uft1", LeafType: "lanceolate", Color: 10}
	anotherUnknownPlant := AnotherUnknownPlant{FlowerColor: 155, LeafType: "cannabis", Height: 33}

	go func() {
		inChannel1 <- unknownPlant
		close(inChannel1)
	}()
	go func() {
		inChannel2 <- anotherUnknownPlant
		close(inChannel2)
	}()	
	outChannel := multiplex(inChannel1, inChannel2)

	for {
		outIface, open := <-outChannel
		if !open{
			break
		}
	fmt.Println("==============================")
	describePlant(outIface)
	}
	fmt.Println("==============================")
	fmt.Printf(" === All Done. ===\n");
}