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
	// Чтение JSON файла
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

	studentMap := make(map[int]Student)
	for _, student := range data.Students {
		studentMap[student.ID] = student
	}

	objectMap := make(map[int]Object)
	for _, object := range data.Objects {
		objectMap[object.ID] = object
	}

	// Вывод таблицы
	fmt.Println("__________________________________________")
	fmt.Println("Student name  | Grade | Object    | Result")
	fmt.Println("__________________________________________")

	for _, result := range data.Results {
		student := studentMap[result.StudentID]
		object := objectMap[result.ObjectID]
		fmt.Printf("%-13s | %-5d | %-9s | %-6d\n", student.Name, student.Grade, object.Name, result.Result)
	}
}
