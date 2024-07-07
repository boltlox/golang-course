package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Message struct {
	Token   string
	FileID  string
	Content string
}

type Cache struct {
	messages map[string][]Message
	mu       sync.Mutex
}

var (
	validTokens    = map[string]bool{"token1": true, "token2": true, "token3": true}
	messageChan    = make(chan Message, 1000) // Увеличен размер буфера канала
	cache          = Cache{messages: make(map[string][]Message)}
	workerInterval = time.Second
	done           = make(chan bool)
	wg             sync.WaitGroup
	enableLogging  bool // Флаг для включения/отключения логирования
	numHandlers    = 10 // Количество обработчиков кеша
	numWorkers     = 5  // Количество воркеров
	maxRetries     = 3  // Максимальное количество попыток повторной записи
	retryDelay     = time.Second
	writeToFileFn  = writeToFile
)

func main() {
	flag.BoolVar(&enableLogging, "log", false, "Enable logging of cached messages")
	flag.Parse()

	wg.Add(1) // Добавляем воркеров и обработчиков кеша к группе ожидания
	go simulateUsers()

	for i := 0; i < numHandlers; i++ {
		wg.Add(1)
		go cacheMessages()
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker()
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	close(done) // Закрытие канала завершения
	wg.Wait()   // Ждем завершения всех горутин
	gracefulShutdown()
}

func simulateUsers() {
	defer wg.Done() // Уменьшаем счетчик горутин по завершению
	users := 50
	for i := 0; i < users; i++ {
		go func(user int) {
			for {
				select {
				case <-done:
					return
				default:
					messageChan <- Message{
						Token:   validTokensList()[user%len(validTokensList())],
						FileID:  fmt.Sprintf("file%d.txt", user),
						Content: fmt.Sprintf("Message from user %d", user),
					}
					time.Sleep(time.Millisecond * 50)
				}
			}
		}(i)
	}
}

func cacheMessages() {
	defer wg.Done() // Уменьшаем счетчик горутин по завершению
	for {
		select {
		case msg := <-messageChan:
			if validateToken(msg.Token) {
				cache.mu.Lock()
				cache.messages[msg.FileID] = append(cache.messages[msg.FileID], msg)
				if enableLogging {
					fmt.Printf("Cached message for file %s\n", msg.FileID)
				}
				cache.mu.Unlock()
			} else {
				if enableLogging {
					fmt.Printf("Invalid token for message to file %s\n", msg.FileID)
				}
			}
		case <-done:
			return
		}
	}
}

func validateToken(token string) bool {
	_, exists := validTokens[token]
	return exists
}

func worker() {
	defer wg.Done() // Уменьшаем счетчик горутин по завершению
	for {
		select {
		case <-time.After(workerInterval):
			cache.mu.Lock()
			for fileID, messages := range cache.messages {
				if len(messages) > 0 {
					if err := writeToFileWithRetry(fileID, messages); err != nil {
						fmt.Println("Failed to write to file after retries:", err)
					} else {
						cache.messages[fileID] = []Message{}
					}
				}
			}
			cache.mu.Unlock()
		case <-done:
			return
		}
	}
}

func writeToFileWithRetry(fileID string, messages []Message) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = writeToFileFn(fileID, messages)
		if err == nil {
			return nil
		}
		time.Sleep(retryDelay)
	}
	return err
}

func writeToFile(fileID string, messages []Message) error {
	file, err := os.OpenFile(fileID, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, msg := range messages {
		if _, err := file.WriteString(fmt.Sprintf("%s\n", msg.Content)); err != nil {
			return err
		}
	}
	return nil
}

func gracefulShutdown() {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	for fileID, messages := range cache.messages {
		if len(messages) > 0 {
			if err := writeToFileWithRetry(fileID, messages); err != nil {
				fmt.Println("Error writing to file during shutdown:", err)
			}
		}
	}
}

func validTokensList() []string {
	keys := make([]string, 0, len(validTokens))
	for k := range validTokens {
		keys = append(keys, k)
	}
	return keys
}
