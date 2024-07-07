package main

import (
	"fmt"
	"sort"
)

func intersection(slices ...[]int) []int {
	if len(slices) == 0 {
		return []int{}
	}

	makeSet := func(slice []int) map[int]struct{} {
		set := make(map[int]struct{})
		for _, v := range slice {
			set[v] = struct{}{}
		}
		return set
	}

	resultSet := makeSet(slices[0])

	for _, slice := range slices[1:] {
		currentSet := makeSet(slice)
		for k := range resultSet {
			if _, found := currentSet[k]; !found {
				delete(resultSet, k)
			}
		}
	}

	result := make([]int, 0, len(resultSet))
	for k := range resultSet {
		result = append(result, k)
	}
	sort.Ints(result)
	return result

}

func main() {
	fmt.Println(intersection([]int{1, 2, 3, 2}))
	fmt.Println(intersection([]int{1, 2, 3, 2}, []int{3, 2}))
	fmt.Println(intersection([]int{1, 2, 3, 2}, []int{3, 2}, []int{}))
}
