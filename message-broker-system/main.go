package main

import (
	"fmt"
	"sync"
	"time"
)

type Message struct {
	ID        int
	Content   string
	CreatedAt time.Time
}

func main() {

	var wg sync.WaitGroup
	var mu sync.Mutex

	messages := []Message{}

	messageCount := 10

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= messageCount; i++ {
			mu.Lock()

			message := Message{

				ID:        len(messages) + 1,
				Content:   fmt.Sprintf("message %d", i),
				CreatedAt: time.Now(),
			}
			messages = append(messages, message)
			mu.Unlock()
		}
	}()

	wg.Wait()

	ch := make(chan Message, len(messages))

	wg.Add(1)
	go func() {

		defer wg.Done()
		for _, message := range messages {

			ch <- message

		}

		close(ch)
	}()

	for c := range ch {
		fmt.Println(c)
	}

	fmt.Println("All messages processed")
}
