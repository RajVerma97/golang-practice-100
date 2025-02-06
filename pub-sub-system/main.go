package main

import (
	"fmt"
	"sync"
)

type PubSub struct {
	mu          sync.RWMutex
	subscribers map[string][]chan interface{}
	closed      bool // A flag to prevent further publishing after the PubSub system is closed.
}

// initializes and returns a new PubSub instance with an empty subscribers map.
func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[string][]chan interface{}),
	}
}

// Adds a new subscriber to a given topic.
func (ps *PubSub) Subscribe(topic string) chan interface{} {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan interface{}, 1)
	ps.subscribers[topic] = append(ps.subscribers[topic], ch)
	return ch
}

func (ps *PubSub) Publish(topic string, message interface{}) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return
	}

	// Get subscribers for this topic
	subscribers, exists := ps.subscribers[topic]
	if !exists {
		return
	}

	// Publish message to all subscribers
	for _, ch := range subscribers {
		select {
		case ch <- message:
		default:
			// Skip if subscriber's channel is full
		}
	}
}

//Removes a subscriber from a topic.
func (ps *PubSub) Unsubscribe(topic string, ch chan interface{}) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	subscribers, exists := ps.subscribers[topic]
	if !exists {
		return
	}

	// Find and remove the subscriber
	for i, subscriber := range subscribers {
		if subscriber == ch {
			ps.subscribers[topic] = append(subscribers[:i], subscribers[i+1:]...)
			close(ch)
			break
		}
	}
}

func (ps *PubSub) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.closed {
		ps.closed = true
		// Close all subscriber channels
		for _, subscribers := range ps.subscribers {
			for _, ch := range subscribers {
				close(ch)
			}
		}
		// Clear the subscribers map
		ps.subscribers = make(map[string][]chan interface{})
	}
}

func main() {
	ps := NewPubSub()
	defer ps.Close()

	// Create subscribers
	sub1 := ps.Subscribe("topic1")
	sub2 := ps.Subscribe("topic1")
	sub3 := ps.Subscribe("topic2")

	// Start subscriber goroutines
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		msg := <-sub1
		fmt.Printf("Subscriber 1 received: %v\n", msg)
	}()

	go func() {
		defer wg.Done()
		msg := <-sub2
		fmt.Printf("Subscriber 2 received: %v\n", msg)
	}()

	go func() {
		defer wg.Done()
		msg := <-sub3
		fmt.Printf("Subscriber 3 received: %v\n", msg)
	}()

	// Publish messages
	ps.Publish("topic1", "Hello topic 1")
	ps.Publish("topic2", "Hello topic 2")

	wg.Wait()
}
