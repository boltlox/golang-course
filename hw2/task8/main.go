package main

import "fmt"

func IsEqualArrays[T comparable](arr1, arr2 []T) bool {

	if len(arr1) != len(arr2) {
		return false
	}

	freqMap := make(map[T]int)

	for _, elem := range arr1 {
		freqMap[elem]++
	}

	for _, elem := range arr2 {
		freqMap[elem]--
	}

	for _, freq := range freqMap {
		if freq != 0 {
			return false
		}
	}

	return true
}

func main() {

	arr1 := []int{1, 2, 3, 4}
	arr2 := []int{3, 4, 1, 2}
	arr3 := []int{1, 2, 3, 4, 5}

	isEqual := IsEqualArrays(arr1, arr2)
	fmt.Println(isEqual)

	isEqual = IsEqualArrays(arr1, arr3)
	fmt.Println(isEqual)
}
