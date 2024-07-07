package main

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

func setUp() {
	messageChan = make(chan Message, 1000)
	cache = Cache{messages: make(map[string][]Message)}
}

func tearDown() {
	close(messageChan)
	for k := range cache.messages {
		os.Remove(k)
	}
}

// TestSuccessfulWrite проверяет успешную запись сообщения в файл
func TestSuccessfulWrite(t *testing.T) {
	setUp()
	defer tearDown()

	messageChan <- Message{Token: "token1", FileID: "testfile.txt", Content: "Test message"}

	go cacheMessages()
	time.Sleep(100 * time.Millisecond)

	cache.mu.Lock()
	if len(cache.messages["testfile.txt"]) != 1 {
		t.Errorf("Expected 1 message in cache, got %d", len(cache.messages["testfile.txt"]))
	}
	cache.mu.Unlock()

	go worker()
	time.Sleep(2 * time.Second)

	data, err := os.ReadFile("testfile.txt")
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "Test message\n" {
		t.Errorf("Expected 'Test message', got '%s'", string(data))
	}
}

// TestInvalidToken проверяет обработку сообщения с недействительным токеном
func TestInvalidToken(t *testing.T) {
	setUp()
	defer tearDown()

	messageChan <- Message{Token: "invalid", FileID: "testfile.txt", Content: "Test message"}

	go cacheMessages()
	time.Sleep(100 * time.Millisecond)

	cache.mu.Lock()
	if len(cache.messages["testfile.txt"]) != 0 {
		t.Errorf("Expected 0 messages in cache, got %d", len(cache.messages["testfile.txt"]))
	}
	cache.mu.Unlock()
}

// TestGracefulShutdown проверяет корректное завершение работы и запись сообщений в файл
func TestGracefulShutdown(t *testing.T) {
	setUp()
	defer tearDown()

	cache.messages["testfile.txt"] = []Message{{Token: "token1", FileID: "testfile.txt", Content: "Test message"}}

	gracefulShutdown()

	data, err := os.ReadFile("testfile.txt")
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "Test message\n" {
		t.Errorf("Expected 'Test message', got '%s'", string(data))
	}
}

// TestHighLoad проверяет обработку большого количества сообщений
func TestHighLoad(t *testing.T) {
	setUp()
	defer tearDown()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			messageChan <- Message{Token: "token1", FileID: "testfile.txt", Content: "Test message"}
		}(i)
	}
	wg.Wait()

	go cacheMessages()
	time.Sleep(5 * time.Second)

	cache.mu.Lock()
	msgCount := len(cache.messages["testfile.txt"])
	cache.mu.Unlock()
	fmt.Printf("Messages in cache: %d\n", msgCount)

	if msgCount != 100 {
		t.Errorf("Expected 100 messages in cache, got %d", msgCount)
	}

	go worker()
	time.Sleep(2 * time.Second)

	data, err := os.ReadFile("testfile.txt")
	if err != nil {
		t.Fatal(err)
	}

	expectedContent := ""
	for i := 0; i < 100; i++ {
		expectedContent += "Test message\n"
	}

	if string(data) != expectedContent {
		t.Errorf("Expected '%s', got '%s'", expectedContent, string(data))
	}
}

// TestSimultaneousWrite проверяет одновременную запись сообщений в файл
func TestSimultaneousWrite(t *testing.T) {
	setUp()
	defer tearDown()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				messageChan <- Message{Token: "token1", FileID: "testfile.txt", Content: fmt.Sprintf("Message %d from user %d", j, i)}
			}
		}(i)
	}
	wg.Wait()

	go cacheMessages()
	time.Sleep(1 * time.Second)

	go worker()
	time.Sleep(2 * time.Second)

	data, err := os.ReadFile("testfile.txt")
	if err != nil {
		t.Fatal(err)
	}

	expectedCount := 50
	actualCount := len(string(data))
	if actualCount != expectedCount*len("Message X from user X\n") {
		t.Errorf("Expected %d messages, got %d", expectedCount, actualCount)
	}
}

// TestWorkerRetry проверяет поведение воркера при сбоях записи
func TestWorkerRetry(t *testing.T) {
	setUp()
	defer tearDown()

	errFile := "errorfile.txt"
	cache.messages[errFile] = []Message{{Token: "token1", FileID: errFile, Content: "Test message with error"}}

	originalWriteToFile := writeToFileFn
	writeToFileFn = func(fileID string, messages []Message) error {
		if fileID == errFile {
			return errors.New("simulated write error")
		}
		return originalWriteToFile(fileID, messages)
	}

	go worker()
	time.Sleep(5 * time.Second)

	if _, err := os.Stat(errFile); err == nil {
		t.Fatal("Expected write error, but file was created")
	}

	writeToFileFn = originalWriteToFile
}
