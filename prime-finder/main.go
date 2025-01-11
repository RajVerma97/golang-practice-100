package main

import (
	"fmt"
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

func generatePrimeNumbers(n int, resultChan chan int, wg *sync.WaitGroup) {

	defer wg.Done()
	for i := 1; i <= n; i++ {
		if isPrime(i) {
			resultChan <- i
		}
	}

}

func main() {

	start := time.Now()

	n := 1000000
	var wg sync.WaitGroup
	resultChan := make(chan int, n)
	wg.Add(1)

	go generatePrimeNumbers(n, resultChan, &wg)

	wg.Wait()
	close(resultChan)

	for result := range resultChan {
		fmt.Println(result)
	}

	duration := time.Since(start)
	fmt.Println(duration)

}
