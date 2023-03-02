package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"sort"
)

func print_round_fl(str string, f float64) {
	fmt.Println(str, math.Round(f*100)/100)
}

func scan_arr() []float64 {
	var arr []float64
	var n float64
	var c1 int
	var err1 error

	line_cnt := 0
	for c1, err1 = fmt.Scan(&n); c1 != 0; c1, err1 = fmt.Scan(&n) {
		if n < -100000 || n > 100000 {
			fmt.Println("Error on line:", line_cnt+1, ": OVERFLOW")
			return nil
		}
		arr = append(arr, n)
		line_cnt++
	}
	if err1 != io.EOF {
		fmt.Println("Error on line:", line_cnt, ":", err1)
		return nil
	}
	return arr
}

func mean_arr(arr []float64) {
	var size int
	var mean float64

	sort.Float64s(arr)
	size = len(arr)
	for i := 0; i < size; i++ {
		mean += arr[i]
	}
	mean /= float64(size)
	print_round_fl("Mean:", mean)
}

func median_arr(arr []float64) {
	var size int

	sort.Float64s(arr)
	size = len(arr)
	if size == 2 {
		print_round_fl("Median:", (arr[0]+arr[1])/2)
	} else if size%2 > 0 {
		print_round_fl("Median:", arr[size/2])
	} else {
		print_round_fl("Median:", (arr[size/2]+arr[size/2+1])/2)
	}
}

func mode_arr(arr []float64) {

	sort.Float64s(arr)
	size := len(arr)

	count := 1
	countMode :=1;

	number := arr[0]
	mode := number

	for i := 1; i < size; i++ {
		if arr[i] == number {
			count += 1
		} else {
			count = 1
			number = arr[i]
		}
		if (count > countMode) {
			countMode = count
			mode = number
		} else if count == countMode {
			if number < mode {
				mode = number
			}
		}
	}
	print_round_fl("Mode:", mode)
}

func sd_arr(arr []float64) {
	var size int
	var mean, sd float64

	sort.Float64s(arr)
	size = len(arr)
	for i := 0; i < size; i++ {
		mean += arr[i]
	}
	mean /= float64(size)
	for j := 0; j < size; j++ {
		sd += math.Pow(arr[j]-mean, 2)
	}
	sd = math.Sqrt(sd / float64(size-1))

	print_round_fl("SD:", sd)
}

func info_arr(arr []float64) {
	fmt.Println(arr)
	fmt.Println("(+) Size:", len(arr))
}

func main() {

	var (
		hidemean   bool
		hidemedian bool
		hidemode   bool
		hidesd     bool
		showinfo   bool
	)

	flag.BoolVar(&hidemean, "hidemean", false, "dont show Mean")
	flag.BoolVar(&hidemedian, "hidemedian", false, "dont show Median")
	flag.BoolVar(&hidemode, "hidemode", false, "dont show Mode")
	flag.BoolVar(&hidesd, "hidesd", false, "dont show SD")
	flag.BoolVar(&showinfo, "showinfo", false, "show show extra info")
	flag.Parse()

	var arr []float64 = scan_arr()

	if arr != nil {
		if !hidemean {
			mean_arr(arr)
		}
		if !hidemedian {
			median_arr(arr)
		}
		if !hidemode {
			mode_arr(arr)
		}
		if !hidesd {
			sd_arr(arr)
		}
	}
	if showinfo {
		info_arr(arr)
	}
}
