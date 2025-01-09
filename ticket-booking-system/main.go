package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
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

type UserContactDetails struct {
	MobileNumber string
	Address      string
	City         string
	Country      string
	Pincode      string
}

type User struct {
	ID            int
	Username      string
	BookedTickets []Ticket
}

var (
	tickets map[int]*Ticket
	users   map[int]*User
	mu      sync.Mutex
	wg      sync.WaitGroup
)

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
	mu.Lock()
	defer mu.Unlock()

	ticket, err := fetchTicketById(ticketId)
	if err != nil {
		return fmt.Errorf("ticket not found with ticket %d ", ticketId)
	}

	if ticket.Status != AVAILABLE {
		return fmt.Errorf("ticket  BOOKED ALREADY with ticket %d ", ticketId)

	}

	user, err := fetchUserById(userId)
	if err != nil {
		return fmt.Errorf("user not found with user %d ", userId)
	}

	ticket.Status = BOOKED
	userBookedTickets := user.BookedTickets
	user.BookedTickets = append(userBookedTickets, *ticket)

	return nil
}

func simulateUserBooking(userId int) {
	defer wg.Done()

	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 3; i++ {
		ticketId := random.Intn(len(tickets)) + 1
		err := BookTicket(userId, ticketId)
		if err != nil {
			fmt.Printf("User %d failed to book ticket %d: %s\n", userId, ticketId, err)
		} else {
			fmt.Printf("User %d successfully booked ticket %d\n", userId, ticketId)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {

	start := time.Now()

	ticketCount := 1000000
	userCount := 200000

	CreateTickets(ticketCount)
	CreateUsers(userCount)

	for _, user := range users {
		wg.Add(1)
		go simulateUserBooking(user.ID)
	}
	wg.Wait()

	fmt.Println(time.Since(start))

}
