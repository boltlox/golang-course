package main

import (
	"fmt"
	"math"
)

// Функция проверки, является ли число простым
func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// Функция для разделения чисел на простые и составные
func splitNumbers(numbers []int, primeCh chan<- int, compositeCh chan<- int) {
	defer close(primeCh)
	defer close(compositeCh)

	for _, number := range numbers {
		if isPrime(number) {
			primeCh <- number
		} else {
			compositeCh <- number
		}
	}
}

// Функция владельца канала для чтения данных из канала и записи в слайс
func channelOwner(ch <-chan int, result *[]int, done chan<- bool) {
	for number := range ch {
		*result = append(*result, number)
	}
	done <- true
}

func main() {
	numbers := []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	primeCh := make(chan int)
	compositeCh := make(chan int)

	var primes []int
	var composites []int

	done := make(chan bool)

	go channelOwner(primeCh, &primes, done)
	go channelOwner(compositeCh, &composites, done)

	go splitNumbers(numbers, primeCh, compositeCh)

	<-done
	<-done

	fmt.Println("Простые числа:", primes)
	fmt.Println("Составные числа:", composites)
}
