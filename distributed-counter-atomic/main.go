package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Counter struct {
	value int64
}

func NewCounter() *Counter {
	return &Counter{value: 0}
}

// Increment atomically increases the counter by 1
func (c *Counter) Increment() int64 {
	return atomic.AddInt64(&c.value, 1)
}

// Decrement atomically decreases the counter by 1
func (c *Counter) Decrement() int64 {
	return atomic.AddInt64(&c.value, -1)
}

// Get returns the current value of the counter
func (c *Counter) Get() int64 {
	return atomic.LoadInt64(&c.value)
}

// AddValue atomically adds a value to the counter
func (c *Counter) AddValue(val int64) int64 {
	return atomic.AddInt64(&c.value, val)
}

// Reset sets the counter back to 0
func (c *Counter) Reset() {
	atomic.StoreInt64(&c.value, 0)
}

func main() {
	counter := NewCounter()
	var wg sync.WaitGroup

	// Test concurrent increments
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	// Test concurrent decrements
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Decrement()
		}()
	}

	// Wait for all operations to complete
	wg.Wait()

	fmt.Printf("Final counter value: %d\n", counter.Get())

	// Test concurrent reads
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			value := counter.Get()
			fmt.Printf("Reader %d got value: %d\n", id, value)
		}(i)
	}

	wg.Wait()

	// Test adding specific values
	counter.AddValue(100)
	fmt.Printf("After adding 100: %d\n", counter.Get())

	// Test reset
	counter.Reset()
	fmt.Printf("After reset: %d\n", counter.Get())
}
