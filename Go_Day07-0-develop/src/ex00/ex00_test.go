package main

import ("testing"
		"reflect")

func TestMinCoins(t *testing.T) {
	testTable := []struct {
		coins []int
		amount int
		expected []int
	}{
		{
			coins: []int{7,4,5,10,2},
			amount: 13,
			expected: []int{},
		},
		{
			coins: []int{1,5,10,50,100,500,1000},
			amount: 13,
			expected: []int{10,1,1,1},
		},
		{
			coins: []int{1,3,4,7,13,15},
			amount: 13,
			expected: []int{13},
		},
	}

	for _, testCase := range testTable {
		result := minCoins(testCase.amount, testCase.coins)
	
		if !reflect.DeepEqual(result, testCase.expected) {
			t.Error("Incorrect result. Expected:", testCase.expected,"got:", result)
		}
	}
}

func TestMinCoins2(t *testing.T) {
	testTable := []struct {
		coins []int
		amount int
		expected []int
	}{
		{
			coins: []int{7,4,5,10,2},
			amount: 13,
			expected: []int{},
		},
		{
			coins: []int{1,5,10,50,100,500,1000},
			amount: 13,
			expected: []int{10,1,1,1},
		},
		{
			coins: []int{1,3,4,7,13,15},
			amount: 13,
			expected: []int{13},
		},
	}

	for _, testCase := range testTable {
		result := minCoins2(testCase.amount, testCase.coins)
	
		if !reflect.DeepEqual(result, testCase.expected) {
			t.Error("Incorrect result. Expected:", testCase.expected,"got:", result)
		}
	}
}