package main

import ("fmt"
		"sort")

/*It accepts a necessary amount and a sorted 
slice of unique denominations of coins. 
It may be something like [1,5,10,50,100,500,1000] 
or something exotic, like [1,3,4,7,13,15]. 
The output is supposed to be a slice of coins 
of minimal size that can be used to express 
the value (e.g. for 13 and [1,5,10] 
	it should give you [10,1,1,1]).*/

func minCoins(val int, coins []int) []int {
    res := make([]int, 0)
    i := len(coins) - 1
    for i >= 0 {
        for val >= coins[i] {
            val -= coins[i]
            res = append(res, coins[i])
        }
        i -= 1
    }
    return res
}
func makeResult(storage *[]int, sum int) []int {

	result := make([]int, 0)
	sort.Ints(*storage)
	i := len(*storage) - 1
	if (i < 0) {
		return []int{}
	}
	for (sum > 0 && i > -1) {
		if (sum >= (*storage)[i]) {
			sum -= (*storage)[i]
			result = append(result, (*storage)[i])
		} else {
			i -= 1
		}
	}
	if (sum != 0) {
		return []int{}
	}
	return result
}

func minCoins2(val int, coins []int) []int {
	mapValue := make(map[int]int)
	storage := make([]int, 0)

	for i := 0; i < len(coins); i++ {
		if (coins[i] <= 0) {
			return []int{}
		}
		if mapValue[coins[i]] == 0 {
			mapValue[coins[i]] = coins[i]
			storage = append(storage, coins[i])
		}
	}

	return makeResult(&storage, val)
}

func min(x,y int) int {
	if x < y {
		return x
	}
	return y
}

func main() {
	fmt.Println("==========================================")
	mc1 := []int{7,4,5,10,2} 
	res1 := minCoins(13, mc1)
	fmt.Println(res1)

	mc1 = []int{1,5,10,50,100,500,1000} 
	res1 = minCoins(13, mc1)
	fmt.Println(res1)

	mc1 = []int{1,3,4,7,13,15} 
	res1 = minCoins(13, mc1)
	fmt.Println(res1)
	fmt.Println("==========================================")
/*==========================================*/
	mc2 := []int{7,4,5,10,2} 
	res2 := minCoins2(13, mc2)
	fmt.Println(res2)

	mc2 = []int{1,5,10,50,100,500,1000} 
	res2 = minCoins2(13, mc2)
	fmt.Println(res2)

	mc2 = []int{1,3,4,7,13,15} 
	res2 = minCoins2(13, mc2)
	fmt.Println(res2)
	fmt.Println("==========================================")
}