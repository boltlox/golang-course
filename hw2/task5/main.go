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

type Cache[K comparable, V any] struct {
	m map[K]V
}

func (c *Cache[K, V]) Init() {
	c.m = make(map[K]V)
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.m[key] = value
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	value, ok := c.m[key]
	return value, ok
}

var studentCache Cache[int, Student]
var objectCache Cache[int, Object]

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

	studentCache.Init()
	for _, student := range data.Students {
		studentCache.Set(student.ID, student)
	}

	objectCache.Init()
	for _, object := range data.Objects {
		objectCache.Set(object.ID, object)
	}

	fmt.Println("__________________________________________")
	fmt.Println("Student name  | Grade | Object    | Result")
	fmt.Println("__________________________________________")

	for _, result := range data.Results {
		student, found := studentCache.Get(result.StudentID)
		if !found {
			log.Fatalf("Не удалось найти студента с ID: %d", result.StudentID)
		}
		object, found := objectCache.Get(result.ObjectID)
		if !found {
			log.Fatalf("Не удалось найти предмет с ID: %d", result.ObjectID)
		}
		fmt.Printf("%-13s | %-5d | %-9s | %-6d\n", student.Name, student.Grade, object.Name, result.Result)
	}
}
