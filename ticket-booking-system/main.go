package main

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type TicketStatus int32

const (
	AVAILABLE TicketStatus = iota
	BOOKED
)

type Ticket struct {
	ID     int
	Price  int
	Status TicketStatus
}

type Metrics struct {
	SuccessfulBookings int32
	FailedBookings     int32
}

type User struct {
	ID            int
	Username      string
	BookedTickets []Ticket
}

var (
	tickets    map[int]*Ticket
	users      map[int]*User
	ticketLock sync.RWMutex
	metrics    Metrics
	wg         sync.WaitGroup
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func CreateTickets(count int) {
	tickets = make(map[int]*Ticket, count)
	for i := 1; i <= count; i++ {
		tickets[i] = &Ticket{
			ID:     i,
			Status: AVAILABLE,
		}
	}
}

func CreateUsers(count int) {
	users = make(map[int]*User, count)
	for i := 1; i <= count; i++ {
		users[i] = &User{
			ID:            i,
			Username:      fmt.Sprintf("User %d", i),
			BookedTickets: []Ticket{},
		}
	}
}

func fetchTicketById(ticketId int) (*Ticket, error) {
	if ticket, exists := tickets[ticketId]; exists {
		return ticket, nil
	}
	return nil, errors.New("ticket not found")
}

func fetchUserById(userId int) (*User, error) {
	if user, exists := users[userId]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func BookTicket(userId int, ticketId int) error {
	ticketLock.Lock()
	defer ticketLock.Unlock()

	ticket, err := fetchTicketById(ticketId)
	if err != nil {
		return fmt.Errorf("ticket not found with ticket %d ", ticketId)
	}

	if ticket.Status != AVAILABLE {
		return fmt.Errorf("ticket BOOKED ALREADY with ticket %d ", ticketId)
	}

	user, err := fetchUserById(userId)
	if err != nil {
		return fmt.Errorf("user not found with user %d ", userId)
	}

	ticket.Status = BOOKED
	user.BookedTickets = append(user.BookedTickets, *ticket)

	return nil
}

func BatchBooking(userId int, ticketIds []int) {
	defer wg.Done()

	for _, ticketId := range ticketIds {
		err := BookTicket(userId, ticketId)
		if err != nil {
			fmt.Printf("User %d failed to book ticket %d: %s\n", userId, ticketId, err)
			atomic.AddInt32(&metrics.FailedBookings, 1)
		} else {
			fmt.Printf("User %d successfully booked ticket %d\n", userId, ticketId)
			atomic.AddInt32(&metrics.SuccessfulBookings, 1)
		}
	}
}

func worker(jobs <-chan []int) {
	for ticketIds := range jobs {
		userId := random.Intn(len(users)) + 1 // Randomly assign a user ID
		BatchBooking(userId, ticketIds)
	}
}

func simulateUserBooking(batchSize int, jobs chan<- []int) {
	var ticketIdsBatch []int

	for i := 1; i <= 3; i++ {
		ticketId := random.Intn(len(tickets)) + 1
		ticketIdsBatch = append(ticketIdsBatch, ticketId)

		if len(ticketIdsBatch) == batchSize {
			wg.Add(1)
			jobs <- ticketIdsBatch
			ticketIdsBatch = []int{}
		}
	}

	if len(ticketIdsBatch) > 0 {
		wg.Add(1)
		jobs <- ticketIdsBatch
	}
}

func PrintMetrics() {
	fmt.Println("Successful Bookings:", metrics.SuccessfulBookings)
	fmt.Println("Failed Bookings:", metrics.FailedBookings)
}

func main() {
	start := time.Now()

	ticketCount := 1000000
	userCount := 200000
	batchSize := 1000
	workerCount := runtime.NumCPU() * 2

	CreateTickets(ticketCount)
	CreateUsers(userCount)

	// Create a buffered channel for jobs
	jobs := make(chan []int, workerCount)

	// Start worker pool
	for w := 1; w <= workerCount; w++ {
		go worker(jobs)
	}

	//send tasks  to the channels (jobs)
	for range users {
		go simulateUserBooking(batchSize, jobs)
	}

	// Wait for all bookings to complete
	wg.Wait()

	// Close the jobs channel
	close(jobs)

	PrintMetrics()
	fmt.Println(time.Since(start))
}
