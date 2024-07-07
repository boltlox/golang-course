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

func Filter[T any](s []T, f func(T) bool) []T {
	var r []T
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
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

	filteredStudents := Filter(data.Students, func(student Student) bool {
		for _, result := range data.Results {
			if result.StudentID == student.ID && result.Result != 5 {
				return false
			}
		}
		return true
	})

	fmt.Println("Students with all grades equal to 5:")
	for _, student := range filteredStudents {
		fmt.Println(student.Name)
	}
}
