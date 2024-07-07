package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

func readInput(ctx context.Context, wg *sync.WaitGroup, inputChan chan<- string) {
	defer wg.Done()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Введите текст (для завершения введите 'exit'):")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Чтение данных завершено по контексту.")
			return
		default:
			if scanner.Scan() {
				text := scanner.Text()
				if text == "exit" {
					close(inputChan)
					return
				}
				inputChan <- text
			}
		}
	}
}

func writeToFile(ctx context.Context, wg *sync.WaitGroup, inputChan <-chan string, fileName string) {
	defer wg.Done()
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Запись данных завершена по контексту.")
			writer.Flush()
			return
		case text, ok := <-inputChan:
			if !ok {
				writer.Flush()
				return
			}
			_, err := writer.WriteString(text + "\n")
			if err != nil {
				fmt.Println("Ошибка при записи в файл:", err)
				return
			}
			writer.Flush()
		}
	}
}

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	inputChan := make(chan string)
	wg.Add(2)

	go readInput(ctx, &wg, inputChan)
	go writeToFile(ctx, &wg, inputChan, "output.txt")

	wg.Wait()
	fmt.Println("Программа завершена.")
}
