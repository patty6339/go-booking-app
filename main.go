// package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// const eventTickets int = 201

// var eventName = "Scrabble National Championship"
// var remainingTickets uint = 201
// var bookings = make([]UserData, 0)

// type UserData struct {
// 	firstName       string
// 	lastName        string
// 	email           string
// 	numberOfTickets uint
// }

// var wg = sync.WaitGroup{}

// func main() {

// 	greetUsers()

// 	firstName, lastName, email, userTickets := getUserInput()

// 	isValidName, isValidEmail, isValidTicketNumber := ValidateUserInput(firstName, lastName, email, userTickets, remainingTickets)

// 	if isValidName && isValidEmail && isValidTicketNumber {
// 		fmt.Println("Thank you for your booking!")
// 	} else {
// 		if !isValidName {
// 			fmt.Println("Invalid name. Please enter a valid first and last name.")
// 		}
// 		if !isValidEmail {
// 			fmt.Println("Invalid email. Please enter a valid email address.")
// 		}
// 		if !isValidTicketNumber {
// 			fmt.Println("Invalid ticket number. Please enter a number greater than 0 and less than or equal to the remaining tickets.")
// 		}
// 		// continue
// 	}

// 	if userTickets > remainingTickets {
// 		fmt.Printf("Sorry we only have %v tickets remaining, you cannot book %v tickets\n", remainingTickets, userTickets)
// 		// continue
// 	}

// 	bookTicket(userTickets, firstName, lastName, email)

// 	wg.Add(1) // Add a goroutine to the wait group
// 	go sendTicket(userTickets, firstName, lastName, email)

// 	// call function print first names
// 	firstNames := getFirstNames()
// 	fmt.Printf("The first names of the bookings are: %v\n", firstNames)

// 	wg.Wait() // Wait for all goroutines to finish
// }

// func greetUsers() {
// 	fmt.Printf("Welcome to %v booking application\n", eventName)
// 	fmt.Printf("We have a total of %v tickets and %v are still available\n", eventTickets, remainingTickets)
// 	fmt.Println("Get Your Tickets Here to Attend!")
// }

// // This function greets users with the conference name, total tickets, and remaining tickets.
// // It is used to display the initial information about the event.
// // It takes three parameters: confName (string), confTickets (int), and remainingTickets (int).
// // It prints a welcome message and the ticket information to the console.

// func getFirstNames() []string {
// 	firstNames := []string{}
// 	for _, booking := range bookings {
// 		// Append the first name from each booking to the firstNames slice
// 		firstNames = append(firstNames, booking.firstName)
// 	}
// 	return firstNames

// }
// func getUserInput() (string, string, string, uint) {
// 	var firstName string
// 	var lastName string
// 	var email string
// 	var userTickets uint
// 	// ask for user input

// 	fmt.Println("Enter your first name:")
// 	fmt.Scan(&firstName)
// 	fmt.Println("Enter your last name:")
// 	fmt.Scan(&lastName)
// 	fmt.Println("Enter your email:")
// 	fmt.Scan(&email)
// 	fmt.Println("Enter number of tickets:")
// 	fmt.Scan(&userTickets)

// 	return firstName, lastName, email, userTickets
// }

// // This function takes user input for first name, last name, email, and number of tickets.

// func bookTicket(userTickets uint, firstName string, lastName string, email string) {
// 	remainingTickets = remainingTickets - userTickets

// 	// create a struct for a user

// 	var userData = UserData{
// 		firstName:       firstName,
// 		lastName:        lastName,
// 		email:           email,
// 		numberOfTickets: userTickets,
// 	}

// 	bookings = append(bookings, userData)
// 	fmt.Printf("List of bookings is %v\n", bookings)
// 	// Print a confirmation message
// 	fmt.Printf("Thank you %v %v for booking %v tickets. You will receive a confirmation email at %v\n", firstName, lastName, userTickets, email)
// 	fmt.Printf("%v tickets remaining for %v\n", remainingTickets, firstName)
// }

// func sendTicket(userTickets uint, firstName string, lastName string, email string) {
// 	time.Sleep(10 * time.Second) // Simulate a delay for sending the ticket
// 	// This function is a placeholder for sending tickets via email.
// 	// In a real application, you would implement email sending logic here.
// 	fmt.Println("###############")
// 	fmt.Printf("Sending ticket:\n %v tickets for %v %v\n", userTickets, firstName, lastName)
// 	fmt.Printf("Sending ticket confirmation to %v %v at %v for %v tickets.\n", firstName, lastName, email, userTickets)
// 	fmt.Println("###############")
// 	wg.Done() // Mark the goroutine as done
// }
