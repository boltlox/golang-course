package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func readInput(dataChan chan<- string, doneChan chan struct{}) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Введите данные (для завершения работы нажмите Ctrl+C):")
	for {
		select {
		case <-doneChan:
			close(dataChan)
			return
		default:
			if scanner.Scan() {
				dataChan <- scanner.Text()
			} else {
				fmt.Println("Ошибка чтения данных.")
			}
		}
	}
}

func writeToFile(dataChan <-chan string, doneChan chan struct{}) {
	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		doneChan <- struct{}{}
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for data := range dataChan {
		_, err := writer.WriteString(data + "\n")
		if err != nil {
			fmt.Println("Ошибка записи в файл:", err)
			doneChan <- struct{}{}
			return
		}
		writer.Flush()
	}
}

func main() {
	dataChan := make(chan string)
	doneChan := make(chan struct{})

	go readInput(dataChan, doneChan)
	go writeToFile(dataChan, doneChan)

	// Обработка сигнала Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	close(doneChan)

	fmt.Println("Завершение работы программы.")
}
