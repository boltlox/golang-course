package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	numWorkers = 5
	timeout    = 10 * time.Second
)

type Task struct {
	url      string
	filename string
}

type Result struct {
	url string
	err error
}

type FileDownloader struct {
	tasks   chan Task
	results chan Result
	wg      sync.WaitGroup
}

func NewFileDownloader(taskBufferSize int) *FileDownloader {
	return &FileDownloader{
		tasks:   make(chan Task, taskBufferSize),
		results: make(chan Result, taskBufferSize),
	}
}

func (fd *FileDownloader) worker(id int) {
	for task := range fd.tasks {
		fmt.Printf("Worker %d начал загрузку %s\n", id, task.url)
		fd.results <- fd.downloadFile(task.url, task.filename)
		fmt.Printf("Worker %d закончил загрузку %s\n", id, task.url)
		fd.wg.Done()
	}
}

func (fd *FileDownloader) downloadFile(url, filename string) Result {
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return Result{url: url, err: err}
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return Result{url: url, err: err}
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return Result{url: url, err: err}
	}

	return Result{url: url, err: nil}
}

func (fd *FileDownloader) addTask(url, filename string) {
	fd.wg.Add(1)
	fd.tasks <- Task{url: url, filename: filename}
}

func (fd *FileDownloader) run(numWorkers int) {
	for i := 1; i <= numWorkers; i++ {
		go fd.worker(i)
	}

	go func() {
		fd.wg.Wait()
		close(fd.tasks)
		close(fd.results)
	}()
}

func main() {
	urls := []string{
		"https://example.com/file1.jpg",
		"https://example.com/file2.jpg",
		"https://example.com/file3.jpg",
	}

	fileDownloader := NewFileDownloader(len(urls))
	fileDownloader.run(numWorkers)

	for i, url := range urls {
		filename := fmt.Sprintf("file%d.jpg", i+1)
		fileDownloader.addTask(url, filename)
	}

	for result := range fileDownloader.results {
		if result.err != nil {
			fmt.Printf("Ошибка загрузки %s: %v\n", result.url, result.err)
		} else {
			fmt.Printf("Успешно загружено %s\n", result.url)
		}
	}
}
