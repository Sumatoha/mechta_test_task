package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

type Data struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {
	filename := flag.String("file", "data.json", "JSON file for calculations")
	numGoroutines := flag.Int("goroutines", 4, "Number of goroutines")
	flag.Parse()

	data, err := readFile(*filename)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	totalSum := sumData(data, *numGoroutines)
	fmt.Printf("Total Sum: %d\n", totalSum)
}

func readFile(filename string) ([]Data, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var data []Data
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return data, nil
}

func worker(data []Data, resultChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	localSum := 0
	for _, d := range data {
		localSum += d.A + d.B
	}
	resultChan <- localSum
}

func sumData(data []Data, numGoroutines int) int {
	dataLen := len(data)
	chunkSize := (dataLen + numGoroutines - 1) / numGoroutines

	resultChan := make(chan int, numGoroutines)
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > dataLen {
			end = dataLen
		}
		wg.Add(1)
		go worker(data[start:end], resultChan, &wg)
	}

	wg.Wait()
	close(resultChan)

	totalSum := 0
	for sum := range resultChan {
		totalSum += sum
	}

	return totalSum
}
