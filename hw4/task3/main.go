package main

import (
	"fmt"
	"sync"
)

func mergeChannels(ch1, ch2 <-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	output := func(ch <-chan int) {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic:", r)
			}
		}()
		for val := range ch {
			out <- val
		}
	}

	wg.Add(2)
	go output(ch1)
	go output(ch2)

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		defer close(ch1)
		for i := 1; i <= 5; i++ {
			ch1 <- i
		}
	}()

	go func() {
		defer close(ch2)
		for i := 6; i <= 10; i++ {
			ch2 <- i
		}
	}()

	merged := mergeChannels(ch1, ch2)

	for val := range merged {
		fmt.Println(val)
	}
}
