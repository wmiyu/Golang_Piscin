package main

import (
	"gococoa"
)

func main() {
	
	gococoa.InitApplication()
	window := gococoa.NewWindow(10, 20, 300, 200, "Day08 ex02: School 21")
	window.MakeKeyAndOrderFront()
	gococoa.RunApplication()
}