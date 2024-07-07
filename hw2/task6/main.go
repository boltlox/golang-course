package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Grade int    `json:"grade"`
}

type Object struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Result struct {
	ObjectID  int `json:"object_id"`
	StudentID int `json:"student_id"`
	Result    int `json:"result"`
}

type Data struct {
	Students []Student `json:"students"`
	Objects  []Object  `json:"objects"`
	Results  []Result  `json:"results"`
}

func main() {

	file, err := os.Open("dz3.json")
	if err != nil {
		log.Fatalf("Не удалось открыть файл: %s", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Не удалось прочитать файл: %s", err)
	}

	var data Data
	if err := json.Unmarshal(byteValue, &data); err != nil {
		log.Fatalf("Ошибка при парсинге JSON: %s", err)
	}

	studentMap := mapByID(data.Students, func(s Student) int { return s.ID })
	objectMap := mapByID(data.Objects, func(o Object) int { return o.ID })

	type GradeResult struct {
		Total float64
		Count int
	}

	subjectResults := make(map[string]map[int]GradeResult)

	for _, result := range data.Results {
		student := studentMap[result.StudentID]
		object := objectMap[result.ObjectID]

		if _, ok := subjectResults[object.Name]; !ok {
			subjectResults[object.Name] = make(map[int]GradeResult)
		}

		gradeResult := subjectResults[object.Name][student.Grade]
		gradeResult.Total += float64(result.Result)
		gradeResult.Count++
		subjectResults[object.Name][student.Grade] = gradeResult
	}

	for subject, grades := range subjectResults {
		fmt.Println("________________")
		fmt.Printf("%-9s | Mean\n", subject)
		fmt.Println("________________")

		var totalSum float64
		var totalCount int

		for grade, result := range grades {
			mean := result.Total / float64(result.Count)
			fmt.Printf("%-2d grade  | %.1f\n", grade, mean)
			totalSum += result.Total
			totalCount += result.Count
		}

		overallMean := totalSum / float64(totalCount)
		fmt.Println("________________")
		fmt.Printf("mean      | %.1f\n", overallMean)
		fmt.Println("________________")
	}
}

func mapByID[T any](items []T, keySelector func(T) int) map[int]T {
	result := make(map[int]T)
	for _, item := range items {
		key := keySelector(item)
		result[key] = item
	}
	return result
}
