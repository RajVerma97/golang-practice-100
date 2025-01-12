package main

import (
	"fmt"
	"sync"
)

type Person struct {
	Name  string
	Age   int
	Email string
}
type BankAccount struct {
	AccountNumber int
	Person
	Amount int
}

type BankAccounts []BankAccount
type Persons []Person

var (
	bankAccounts []BankAccount
	mu           sync.Mutex
)

func (b *BankAccount) Deposit(amount int) {
	mu.Lock()
	defer mu.Unlock()

	b.Amount += amount
	fmt.Printf("Deposited Amount %d in Bank Account Holder Name %s\n", amount, b.Name)
}

func (b *BankAccount) Withdraw(amount int) error {
	mu.Lock()
	defer mu.Unlock()
	if b.Amount < amount {
		return fmt.Errorf("insufficient funds in Account Holder Name %s Currrent Balance:%d", b.Name, b.Amount)

	}
	b.Amount -= amount
	fmt.Printf("Withdrawn Amount %d in Bank Account Holder Name %s\n", amount, b.Name)
	return nil
}

func CreateAccount(person Person, wg *sync.WaitGroup) {
	defer wg.Done()

	bankAccount := BankAccount{
		AccountNumber: len(bankAccounts) + 1,
		Person:        person,
		Amount:        0,
	}
	bankAccounts = append(bankAccounts, bankAccount)
	// fmt.Printf("Account Created with Holder Name %s and AccountNumber %d", person.Name, bankAccount.AccountNumber)
}

func main() {

	var wg sync.WaitGroup

	persons := Persons{
		{Name: "Alice", Age: 30, Email: "alice@example.com"},
		{Name: "Bob", Age: 25, Email: "bob@example.com"},
		{Name: "Charlie", Age: 35, Email: "charlie@example.com"},
		{Name: "David", Age: 40, Email: "david@example.com"},
		{Name: "Eva", Age: 28, Email: "eva@example.com"},
		{Name: "Frank", Age: 32, Email: "frank@example.com"},
		{Name: "Grace", Age: 29, Email: "grace@example.com"},
		{Name: "Helen", Age: 38, Email: "helen@example.com"},
		{Name: "Isaac", Age: 24, Email: "isaac@example.com"},
		{Name: "Jack", Age: 27, Email: "jack@example.com"},
		{Name: "Karen", Age: 31, Email: "karen@example.com"},
		{Name: "Leo", Age: 33, Email: "leo@example.com"},
		{Name: "Mona", Age: 26, Email: "mona@example.com"},
		{Name: "Nina", Age: 34, Email: "nina@example.com"},
		{Name: "Oscar", Age: 39, Email: "oscar@example.com"},
		{Name: "Paul", Age: 30, Email: "paul@example.com"},
		{Name: "Quinn", Age: 22, Email: "quinn@example.com"},
		{Name: "Rita", Age: 36, Email: "rita@example.com"},
		{Name: "Steve", Age: 28, Email: "steve@example.com"},
		{Name: "Tina", Age: 41, Email: "tina@example.com"},
	}

	for _, person := range persons {
		wg.Add(1)
		go CreateAccount(person, &wg)

	}

	wg.Wait()
	for i := range bankAccounts {
		wg.Add(2)
		go func(bankAccount *BankAccount) {
			defer wg.Done()
			bankAccount.Deposit(100)
		}(&bankAccounts[i])

		go func(bankAccount *BankAccount) {
			defer wg.Done()
			err := bankAccount.Withdraw(50)
			if err != nil {
				fmt.Printf("err %s", err)
			}
		}(&bankAccounts[i])

	}

	wg.Wait()

}
