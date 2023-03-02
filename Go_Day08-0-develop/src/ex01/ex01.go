package main

import (
	"fmt"
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

func main() {
	up1 := UnknownPlant{FlowerType: "uft1", LeafType: "lanceolate", Color: 10}
	up2 := AnotherUnknownPlant{FlowerColor: 155, LeafType: "cannabis", Height: 33}
	fmt.Println("==============================")
	describePlant(up1)

	fmt.Println("==============================")
	describePlant(up2)
}