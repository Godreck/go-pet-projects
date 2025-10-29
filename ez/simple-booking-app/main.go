package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var confName = "SALUKI concert"

const confTotalTickets = 50

var confRemainingTickets = 50

var wg = sync.WaitGroup{}

type UserData struct {
	userName   string
	city       string
	userTicket int
}

func main() {
	var bookings = make([]UserData, 0)

	greetUser()

	for confRemainingTickets > 0 && len(bookings) < 50 {
		userData := GetUserInput()
		userTicket := userData.userTicket

		if userTicket > confRemainingTickets {
			fmt.Printf("Now we have only %v tickets, you can't book ticket %v\n", confRemainingTickets, userTicket)
			continue
		}
		confRemainingTickets = confTotalTickets - userTicket
		bookings = append(bookings, userData)

		fmt.Printf("There are %v tickets available\n", confRemainingTickets)
		wg.Add(1)
		go sendTicket(userData.userName, userData.userTicket)

		firstNames := getFirstNames(bookings)
		fmt.Printf("All bookings: %v\n", bookings)
		fmt.Printf("All bookings: %v\n", firstNames)

		if confRemainingTickets == 0 {
			fmt.Println("Our ivent is out of tickets, come beack in next ivent.")
			break
		}
	}
	wg.Wait()
}

func greetUser() {
	fmt.Printf("Wellcome to %v booking application!\n", confName)
	fmt.Printf("There's %v tickets currently avaliable from %v total number of tickets.\n", confRemainingTickets, confTotalTickets)
	fmt.Println("You can get your ticket here!")
}

func getFirstNames(bookings []UserData) []string {
	var firstNames = []string{}
	for _, booking := range bookings {
		var names = strings.Fields(booking.userName)
		firstNames = append(firstNames, names[0])
	}
	return firstNames
}

func sendTicket(userName string, userTicket int) {
	time.Sleep(50 * time.Second)
	fmt.Printf("%v ticket for %v\n", userTicket, userName)
	wg.Done()
}
