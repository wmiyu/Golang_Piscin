package main

import (
	"fmt"
	"io"
	"flag"
	"os/exec"
	"log"
)

func scan_arg() string {
	var c1 int
	var err1 error
	var s string
	var strings []string
	var result string

	line_cnt := 0
	for c1, err1 = fmt.Scan(&s); c1 != 0 && err1 == nil; c1, err1 = fmt.Scan(&s){

		strings = append(strings, s)
		line_cnt++
	}
	if err1 != io.EOF {
		fmt.Println("Error on line:", line_cnt, ":", err1)
		return ""
	}
	for _, str1 := range(strings) {
		result += str1 + " "
	}
	return result
}

func execute(cmdstr string, argstr string){

	cmd := exec.Command(cmdstr, argstr)

	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}

func main(){

    argstr := ""

	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Error: Too few params")
	} else {
		cmdstr := string(flag.Args()[0])
		for i := 1; i < len(flag.Args()); i++ {
			argstr += flag.Args()[i] + " "
		} 
		argstr += scan_arg()
		fmt.Println("Executing: ", cmdstr, " with: ", argstr)
		execute(cmdstr, argstr)
	}	

}
