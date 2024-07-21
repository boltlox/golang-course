package main

import (
	"fmt"
)

func validateInputs(matrix [][]int, userAnswer []int) error {
	n := len(matrix)

	// Проверка, что матрица квадратная
	for i := range matrix {
		if len(matrix[i]) != n {
			return fmt.Errorf("матрица должна быть квадратной")
		}
	}

	// Проверка, что нет петель
	for i := 0; i < n; i++ {
		if matrix[i][i] != 0 {
			return fmt.Errorf("в графе не может быть петель (матрица стоимости должна иметь нули на диагонали)")
		}
	}

	// Проверка, что ответы пользователя в пределах диапазона матрицы
	for _, ans := range userAnswer {
		if ans < 0 || ans >= n {
			return fmt.Errorf("ответы пользователя не должны выходить за диапазон матрицы")
		}
	}

	// Проверка уникальности элементов в слайсе ответов пользователя
	answerSet := make(map[int]bool)
	for _, ans := range userAnswer {
		if _, exists := answerSet[ans]; exists {
			return fmt.Errorf("элементы в слайсе ответов пользователя должны быть уникальными")
		}
		answerSet[ans] = true
	}

	return nil
}

func calcMaxGrade(matrix [][]int) int {
	maxGrade := 0
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix); j++ {
			if i != j {
				maxGrade += matrix[i][j]
			}
		}
	}
	return maxGrade
}

func calcUserGrade(matrix [][]int, userAnswer []int) int {
	userGrade := 0
	for i := 0; i < len(userAnswer)-1; i++ {
		userGrade += matrix[userAnswer[i]][userAnswer[i+1]]
	}
	return userGrade
}

func EvalSequence(matrix [][]int, userAnswer []int) int {
	if err := validateInputs(matrix, userAnswer); err != nil {
		fmt.Println("Ошибка:", err)
		return 0
	}

	maxGrade := calcMaxGrade(matrix)
	userGrade := calcUserGrade(matrix, userAnswer)

	percent := userGrade * 100 / maxGrade
	return percent
}

func main() {
	matrix := [][]int{
		{0, 2, 3},
		{2, 0, 1},
		{3, 1, 0},
	}
	userAnswer := []int{0, 1, 2}

	percent := EvalSequence(matrix, userAnswer)
	fmt.Printf("Процент корректности: %d%%\n", percent)
}
