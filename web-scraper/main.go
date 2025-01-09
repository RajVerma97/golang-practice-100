package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func fetchDataFromUrl(url string, wg *sync.WaitGroup, resultChan chan string) {
	resp, err := http.Get(url)

	defer wg.Done()

	if err != nil {
		fmt.Printf("Error fetching from %s", url)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading the response body %s \n", err)
	}

	resultChan <- string(body)
}
func main() {

	var wg sync.WaitGroup
	start := time.Now()

	urls := []string{
		"https://www.bbc.com/news",
		"https://www.bbc.com/news",
		"https://www.bbc.com/news",
		"https://www.bbc.com/news",
	}
	resultChan := make(chan string, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go fetchDataFromUrl(url, &wg, resultChan)
	}

	go func() {
		wg.Wait()
		close(resultChan)

	}()

	for result := range resultChan {
		fmt.Println(result)
	}
	fmt.Println(time.Since(start))

}
