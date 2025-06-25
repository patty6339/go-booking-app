package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const eventTickets int = 200

var eventName = "Scrabble National Championship"
var remainingTickets uint = 200
var bookings = make([]UserData, 0)

type UserData struct {
	firstName       string
	lastName        string
	email           string
	numberOfTickets uint
}

var wg = sync.WaitGroup{}
var db *sql.DB

func main() {
	// Command line flags
	webMode := flag.Bool("web", false, "Run in web mode")
	flag.Parse()

	// Initialize database
	db = initDB()
	defer db.Close()

	// Initialize events for multi-event support
	initializeEvents()

	if *webMode {
		startWebMode()
	} else {
		startCLIMode()
	}
}

func startWebMode() {
	// Set up routes - Remove auth middleware from login/register and add public home
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/book", authMiddleware(bookHandler))
	http.HandleFunc("/bookings", authMiddleware(bookingsHandler))
	http.HandleFunc("/events", eventsListHandler)
	http.HandleFunc("/book-event/", authMiddleware(bookEventHandler))
	http.HandleFunc("/payment", authMiddleware(paymentPageHandler))
	http.HandleFunc("/create-payment-intent", authMiddleware(paymentHandler))
	
	// Auth routes
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	
	fmt.Println("Enhanced booking application starting on http://localhost:8080")
	fmt.Println("Features available:")
	fmt.Println("- User authentication")
	fmt.Println("- Multiple events")
	fmt.Println("- Payment processing")
	fmt.Println("- Database storage")
	fmt.Println("- Email notifications")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func startCLIMode() {
	fmt.Println("Starting CLI mode...")
	
	// Load existing bookings from database
	loadedBookings, err := getBookingsFromDB(db)
	if err != nil {
		fmt.Printf("Error loading bookings: %v\n", err)
	} else {
		bookings = loadedBookings
		// Update remaining tickets based on loaded bookings
		for _, booking := range bookings {
			remainingTickets -= booking.numberOfTickets
		}
	}

	greetUsers()

	firstName, lastName, email, userTickets := getUserInput()

	isValidName, isValidEmail, isValidTicketNumber := ValidateUserInput(firstName, lastName, email, userTickets, remainingTickets)

	if isValidName && isValidEmail && isValidTicketNumber {
		fmt.Println("Thank you for your booking!")
		
		// Book the ticket
		bookTicket(userTickets, firstName, lastName, email)
		
		// Save to database
		userData := UserData{
			firstName:       firstName,
			lastName:        lastName,
			email:           email,
			numberOfTickets: userTickets,
		}
		
		err := saveBookingToDB(db, userData)
		if err != nil {
			fmt.Printf("Error saving booking to database: %v\n", err)
		} else {
			fmt.Println("Booking saved to database successfully!")
		}

		wg.Add(1)
		go sendTicketEnhanced(userTickets, firstName, lastName, email)

		firstNames := getFirstNames()
		fmt.Printf("The first names of the bookings are: %v\n", firstNames)

		wg.Wait()
	} else {
		if !isValidName {
			fmt.Println("Invalid name. Please enter a valid first and last name.")
		}
		if !isValidEmail {
			fmt.Println("Invalid email. Please enter a valid email address.")
		}
		if !isValidTicketNumber {
			fmt.Println("Invalid ticket number. Please enter a number greater than 0 and less than or equal to the remaining tickets.")
		}
	}
}

func sendTicketEnhanced(userTickets uint, firstName, lastName, email string) {
	defer wg.Done()
	
	time.Sleep(10 * time.Second)
	
	fmt.Println("###############")
	fmt.Printf("Sending ticket:\n %v tickets for %v %v\n", userTickets, firstName, lastName)
	
	// Try to send real email first, fall back to simulation
	sendTicketConfirmation(userTickets, firstName, lastName, email)
	
	fmt.Println("###############")
}

// Keep original functions for backward compatibility
func greetUsers() {
	fmt.Printf("Welcome to %v booking application\n", eventName)
	fmt.Printf("We have a total of %v tickets and %v are still available\n", eventTickets, remainingTickets)
	fmt.Println("Get Your Tickets Here to Attend!")
}

func getFirstNames() []string {
	firstNames := []string{}
	for _, booking := range bookings {
		firstNames = append(firstNames, booking.firstName)
	}
	return firstNames
}

func getUserInput() (string, string, string, uint) {
	var firstName string
	var lastName string
	var email string
	var userTickets uint

	fmt.Println("Enter your first name:")
	fmt.Scan(&firstName)
	fmt.Println("Enter your last name:")
	fmt.Scan(&lastName)
	fmt.Println("Enter your email:")
	fmt.Scan(&email)
	fmt.Println("Enter number of tickets:")
	fmt.Scan(&userTickets)

	return firstName, lastName, email, userTickets
}

func bookTicket(userTickets uint, firstName string, lastName string, email string) {
	remainingTickets = remainingTickets - userTickets

	var userData = UserData{
		firstName:       firstName,
		lastName:        lastName,
		email:           email,
		numberOfTickets: userTickets,
	}

	bookings = append(bookings, userData)
	fmt.Printf("List of bookings is %v\n", bookings)
	fmt.Printf("Thank you %v %v for booking %v tickets. You will receive a confirmation email at %v\n", firstName, lastName, userTickets, email)
	fmt.Printf("%v tickets remaining for %v\n", remainingTickets, firstName)
}

func sendTicket(userTickets uint, firstName string, lastName string, email string) {
	time.Sleep(10 * time.Second)
	fmt.Println("###############")
	fmt.Printf("Sending ticket:\n %v tickets for %v %v\n", userTickets, firstName, lastName)
	fmt.Printf("Sending ticket confirmation to %v %v at %v for %v tickets.\n", firstName, lastName, email, userTickets)
	fmt.Println("###############")
	wg.Done()
}
