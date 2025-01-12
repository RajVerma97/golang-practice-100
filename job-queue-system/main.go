package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, ch chan string, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	for c := range ch {
		mu.Lock()
		fmt.Printf("Worker %d started job: %s\n", id, c)
		time.Sleep(2 * time.Second)
		fmt.Printf("Worker %d finished job: %s\n", id, c)
		mu.Unlock()
	}

}
func main() {

	var wg sync.WaitGroup
	var mu sync.Mutex
	urls := []string{
		"first",
		"second",
		"third",
	}

	ch := make(chan string, len(urls))

	//enqueue
	for _, url := range urls {
		ch <- url
	}

	close(ch)
	numOfWorkers := 4

	for i := 1; i <= numOfWorkers; i++ {
		wg.Add(1)
		go worker(i, ch, &wg, &mu)
	}

	wg.Wait()
	fmt.Println("All jobs finished")

}
