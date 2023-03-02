package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"unicode/utf8"
)

func countMetrics(flags map[string]bool, filename string, waitGroup *sync.WaitGroup) {
	switch {
	case flags["W"]:
		file, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		splitFile := strings.Split(string(file), " ")
		count := len(splitFile)
		fmt.Printf("\t%d\t%s\n", count, filename)
	case flags["L"]:
		file, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		count := strings.Count(string(file), "\n")
		fmt.Printf("\t%d\t%s\n", count, filename)
	case flags["M"]:
		file, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		count := utf8.RuneCountInString(string(file))
		fmt.Printf("\t%d\t%s\n", count, filename)
	}

	waitGroup.Done()
}

func main() {
	flagW := flag.Bool("w", false, "Count words")
	flagL := flag.Bool("l", false, "Count lines")
	flagM := flag.Bool("m", false, "Count character")

	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Error: Invalid filename")
		os.Exit(0)
	}

	flags := map[string]bool{
		"W": *flagW,
		"L": *flagL,
		"M": *flagM,
	}

	var waitGroup sync.WaitGroup

	for _, file := range flag.Args() {
		waitGroup.Add(1)
		go countMetrics(flags, file, &waitGroup)
	}
	waitGroup.Wait()
}
