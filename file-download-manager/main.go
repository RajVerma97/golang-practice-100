package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func downloadFile(url string, wg *sync.WaitGroup) error {

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to get the current working dir with err %s", err)
	}

	defer wg.Done()
	splitUrlSlice := strings.Split(url, "/")
	fileName := splitUrlSlice[len(splitUrlSlice)-1]

	dirPath := currentDir + "/downloads/"
	filePath := dirPath + fileName
	fmt.Println(filePath)
	err = os.MkdirAll(dirPath, 0755)

	if err != nil {
		return fmt.Errorf("not able to create the dir /downloads with err %s", err)
	}

	file, err := os.Create(filePath)

	if err != nil {
		return fmt.Errorf("not able to create the file with filepath %s with err %s", filePath, err)
	}

	defer file.Close()

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("invalid response with errr  %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("invalid body  %s", err)
	}
	fmt.Printf("File Downloaded from the url:%s\n", url)

	err = os.WriteFile(filePath, body, 0644)

	if err != nil {
		return fmt.Errorf("failed to write file: %s", err)
	}

	return nil

}
func main() {

	var wg sync.WaitGroup

	start := time.Now()

	urls := []string{
		"https://www.w3.org/TR/png/iso_8859-1.txt",
		"https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
		"https://jsonplaceholder.typicode.com/posts.json",
		"https://www.youtube.com/shorts/Mv0_WQ9uBjY",
	}

	for _, url := range urls {
		wg.Add(1)
		go downloadFile(url, &wg)
	}

	wg.Wait()

	duration := time.Since(start)
	fmt.Println(duration)
}
