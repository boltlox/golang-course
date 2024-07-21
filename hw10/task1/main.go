package main

import (
	"fmt"
	"math"
)

// BFS функция для графа, представленного в виде матрицы стоимости
func BFS(cost [][]int, start int) {
	n := len(cost)
	visited := make([]bool, n)
	queue := []int{start}
	visited[start] = true

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		fmt.Printf("Visited %d\n", node)

		for neighbor, weight := range cost[node] {
			if weight > 0 && !visited[neighbor] {
				queue = append(queue, neighbor)
				visited[neighbor] = true
			}
		}
	}
}

func main() {
	// Пример графа в виде матрицы стоимости
	inf := math.MaxInt32
	cost := [][]int{
		{0, 1, 4, inf, inf, inf},
		{1, 0, 4, 2, 7, inf},
		{4, 4, 0, 3, 5, inf},
		{inf, 2, 3, 0, 4, 6},
		{inf, 7, 5, 4, 0, 7},
		{inf, inf, inf, 6, 7, 0},
	}

	start := 0
	fmt.Printf("Starting BFS from node %d\n", start)
	BFS(cost, start)
}
