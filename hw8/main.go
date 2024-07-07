package main

import (
	"fmt"
	"os"
	"time"
)

// InefficientLogger - неэффективный логгер
type InefficientLogger struct {
	filePath string
}

// NewInefficientLogger - конструктор для неэффективного логгера
func NewInefficientLogger(filePath string) *InefficientLogger {
	return &InefficientLogger{filePath: filePath}
}

// Info - метод для записи информационного сообщения с датой
func (l *InefficientLogger) Info(message string) {
	f, err := os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	logMessage := fmt.Sprintf("%s: %s\n", time.Now().Format(time.RFC3339), message)
	if _, err := f.WriteString(logMessage); err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

// EfficientLogger - эффективный логгер
type EfficientLogger struct {
	file *os.File
}

// NewEfficientLogger - конструктор для эффективного логгера
func NewEfficientLogger(filePath string) (*EfficientLogger, error) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &EfficientLogger{file: f}, nil
}

// Close - метод для закрытия файла логгера
func (l *EfficientLogger) Close() {
	l.file.Close()
}

// Info - метод для записи информационного сообщения с датой
func (l *EfficientLogger) Info(message string) {
	logMessage := fmt.Sprintf("%s: %s\n", time.Now().Format(time.RFC3339), message)
	if _, err := l.file.WriteString(logMessage); err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func main() {
	// Пример использования неэффективного логгера
	inefficientLogger := NewInefficientLogger("inefficient.log")
	for i := 0; i < 10; i++ {
		inefficientLogger.Info(fmt.Sprintf("Inefficient log message %d", i))
	}

	// Пример использования эффективного логгера
	efficientLogger, err := NewEfficientLogger("efficient.log")
	if err != nil {
		fmt.Println("Error creating efficient logger:", err)
		return
	}
	defer efficientLogger.Close()

	for i := 0; i < 10; i++ {
		efficientLogger.Info(fmt.Sprintf("Efficient log message %d", i))
	}
}
