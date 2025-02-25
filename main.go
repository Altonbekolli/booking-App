package main

import (
	"fmt"
	"strings"
)

const ticket uint = 50

var remainingTickets uint = 50
var conferenceName = "Go Conference"
var bookings = []string{}

func main() {

	if remainingTickets == 0 {
		fmt.Printf("All tickets are sold out! Booking closed.\n")
		return
	}

	for {
		greetUser(conferenceName, ticket, remainingTickets)

		firstName, lastName, email, ticket := getUserInput()

		isValidEmail, isValidName, isValidTicket := validInput(firstName, lastName, email, ticket)

		if ticket > remainingTickets {
			fmt.Printf("Sorry, we only have %d tickets left.\n", remainingTickets)
			continue
		}

		if isValidEmail && isValidName && isValidTicket {
			processBooking(firstName, lastName, email, ticket)

			fmt.Printf("These are the last names of our booking list: %v\n", strings.Join(printLastName(), ", "))

			// Falls keine Tickets mehr übrig sind, die Schleife beenden
			if remainingTickets == 0 {
				fmt.Printf("All tickets are sold out! Booking closed.\n")
				break
			}
		} else {
			printValidationErrors(isValidName, isValidEmail, isValidTicket)
		}
	}
}

// Main zu Ende

func greetUser(confname string, ticket uint, ticketAvailable uint) {
	fmt.Printf("Welcome to our %v booking application.\n", confname)
	fmt.Printf("We have %d tickets and %d are still available.\n", ticket, ticketAvailable)
}

func getUserInput() (string, string, string, uint) {
	var firstName, lastName, email string
	var ticket uint

	fmt.Printf("What is your first name?\n")
	fmt.Scan(&firstName)

	fmt.Printf("What is your last name?\n")
	fmt.Scan(&lastName)

	fmt.Printf("What is your email address?\n")
	fmt.Scan(&email)

	fmt.Printf("How many tickets do you want to book?\n")
	fmt.Scan(&ticket)

	return firstName, lastName, email, ticket
}

func validInput(firstName, lastName, email string, ticket uint) (bool, bool, bool) {
	isValidEmail := strings.Contains(email, "@")
	isValidName := len(firstName) >= 2 && len(lastName) >= 2
	isValidTicket := ticket > 0
	return isValidEmail, isValidName, isValidTicket
}

func printLastName() []string {
	lastNames := []string{}
	for _, booking := range bookings {
		var names = strings.Fields(booking)
		if len(names) >= 2 {
			lastNames = append(lastNames, names[1])
		}
	}
	return lastNames
}

func processBooking(firstName, lastName, email string, ticket uint) {
	bookings = append(bookings, firstName+" "+lastName)
	remainingTickets -= ticket

	fmt.Printf("Thank you %v %v for booking %v tickets. We will send you an email at %v.\n", firstName, lastName, ticket, email)
	fmt.Printf("There are %d tickets available for the conference %s.\n", remainingTickets, conferenceName)
}

func printValidationErrors(isValidName, isValidEmail, isValidTicket bool) {
	if !isValidName {
		fmt.Printf("Invalid input! Make sure your first and last name has at least 2 letters.\n")
	}
	if !isValidEmail {
		fmt.Printf("Invalid input! Your email must contain '@'.\n")
	}
	if !isValidTicket {
		fmt.Printf("You need to buy at least one ticket.\n")
	}
}
