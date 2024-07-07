package main

import (
	"fmt"
	"sync"
	"time"
)

func RunProcessor(wg *sync.WaitGroup, prices <-chan map[string]float64) {
	go func() {
		defer wg.Done()
		for price := range prices {
			for key, value := range price {
				price[key] = value + 1
			}
			fmt.Println(price)
		}
	}()
}

func RunWriter() <-chan map[string]float64 {
	prices := make(chan map[string]float64)
	go func() {
		defer close(prices)
		var currentPrice = map[string]float64{
			"inst1": 1.1,
			"inst2": 2.1,
			"inst3": 3.1,
			"inst4": 4.1,
		}
		for i := 1; i < 5; i++ {
			newPrice := make(map[string]float64)
			for key, value := range currentPrice {
				newPrice[key] = value + 1
			}
			prices <- newPrice
			for key := range currentPrice {
				currentPrice[key] = newPrice[key]
			}
			time.Sleep(time.Second)
		}
	}()
	return prices
}

func main() {
	p := RunWriter()
	var wg sync.WaitGroup

	wg.Add(3)
	for i := 0; i < 3; i++ {
		RunProcessor(&wg, p)
	}
	wg.Wait()
}
