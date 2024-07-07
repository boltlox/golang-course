package main

import (
	"fmt"
	"testing"
)

func BenchmarkInefficientLogger(b *testing.B) {
	logger := NewInefficientLogger("benchmark_inefficient.log")
	for i := 0; i < b.N; i++ {
		logger.Info(fmt.Sprintf("Inefficient log message %d", i))
	}
}

func BenchmarkEfficientLogger(b *testing.B) {
	logger, err := NewEfficientLogger("benchmark_efficient.log")
	if err != nil {
		b.Fatalf("Error creating efficient logger: %v", err)
	}
	defer logger.Close()

	for i := 0; i < b.N; i++ {
		logger.Info(fmt.Sprintf("Efficient log message %d", i))
	}
}
