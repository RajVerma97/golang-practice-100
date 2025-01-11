package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func isPrime(num int) bool {
	if num <= 1 {
		return false
	}

	for i := 2; i*i < num; i++ {
		if num%i == 0 {
			return false
		}
	}

	return true
}

func generatePrimeNumbers(start int, end int, resultChan chan int, wg *sync.WaitGroup) {

	defer wg.Done()
	for i := start; i <= end; i++ {
		if isPrime(i) {
			resultChan <- i
		}
	}

}

func main() {

	numOfGoroutines := runtime.NumCPU()
	start := time.Now()

	n := 1000000
	var wg sync.WaitGroup
	resultChan := make(chan int, n)

	rangeSize := n / numOfGoroutines

	for i := 0; i < numOfGoroutines; i++ {

		start := rangeSize*i + 1
		end := (i + 1) * rangeSize
		wg.Add(1)
		go generatePrimeNumbers(start, end, resultChan, &wg)
	}

	wg.Wait()
	close(resultChan)

	for result := range resultChan {
		fmt.Println(result)
	}

	duration := time.Since(start)
	fmt.Println(duration)

}
