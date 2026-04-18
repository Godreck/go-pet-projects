package main

import (
	"fmt"
	"strings"
)

func GetUserInput() UserData {
	var userName string
	var city string
	var userTicket int

	fmt.Println("Type city:")
	fmt.Scan(&city)
	fmt.Printf("%v selected.", city)

	fmt.Println("Input your name:")
	fmt.Scan(&userName)
	userName = strings.Replace(userName, "-", " ", -1)
	fmt.Printf("Ok, %v, nice to meet you!", userName)

	fmt.Println("Enter number of tickets:")
	fmt.Scan(&userTicket)

	var userData = UserData{
		userName:   userName,
		city:       city,
		userTicket: userTicket,
	}

	return userData
}
