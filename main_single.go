package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Global variables
var (
	eventTickets     = 100
	remainingTickets = uint(eventTickets)
	bookings         []UserData
	simpleMutex      sync.Mutex
	wg               sync.WaitGroup
	db               *sql.DB // Used in enhanced mode, initialized in initDB()
)

// UserData holds booking information
type UserData struct {
	FirstName       string
	LastName        string
	Email           string
	NumberOfTickets uint
}

// ValidateUserInputSimple validates user input for booking (simple mode)
func ValidateUserInputSimple(firstName, lastName, email string, userTickets, remainingTickets uint) (bool, bool, bool) {
	isValidName := len(firstName) >= 2 && len(lastName) >= 2
	isValidEmail := len(email) >= 5 && containsAt(email)
	isValidTicketNumber := userTickets > 0 && userTickets <= remainingTickets
	return isValidName, isValidEmail, isValidTicketNumber
}

// containsAt checks if an email contains '@'
func containsAt(email string) bool {
	for _, c := range email {
		if c == '@' {
			return true
		}
	}
	return false
}

// Dummy implementations for enhanced mode (to avoid errors if not present)
func initDB() *sql.DB                                           { return &sql.DB{} }
func initializeEvents()                                         {}
func authMiddleware(h http.HandlerFunc) http.HandlerFunc        { return h }
func bookingsHandler(w http.ResponseWriter, r *http.Request)    {}
func eventsListHandler(w http.ResponseWriter, r *http.Request)  {}
func bookEventHandler(w http.ResponseWriter, r *http.Request)   {}
func paymentPageHandler(w http.ResponseWriter, r *http.Request) {}
func paymentHandler(w http.ResponseWriter, r *http.Request)     {}
func loginHandler(w http.ResponseWriter, r *http.Request)       {}
func registerHandler(w http.ResponseWriter, r *http.Request)    {}

func main() {
	// Check for simple-web argument
	if len(os.Args) > 1 && os.Args[1] == "simple-web" {
		startSimpleWebSingle()
		return
	}

	// Parse flags for enhanced mode
	webMode := flag.Bool("web", false, "Run in web mode")
	flag.Parse()

	if *webMode {
		// Initialize database for enhanced mode
		db = initDB()
		defer db.Close()
		initializeEvents()
		startEnhancedWeb()
	} else {
		// CLI mode
		startCLISimple()
	}
}

func startSimpleWebSingle() {
	http.HandleFunc("/", simpleHomeHandlerMain)
	http.HandleFunc("/simple-book", simpleBookHandlerSingle)
	http.HandleFunc("/simple-bookings", simpleBookingsHandlerSingle)

	fmt.Println("🚀 Simple Booking App starting on http://localhost:8080")
	fmt.Println("📝 Features:")
	fmt.Println("   - Web-based ticket booking")
	fmt.Println("   - Input validation")
	fmt.Println("   - Booking management")
	fmt.Println("   - Real-time ticket tracking")
	fmt.Println()
	fmt.Println("🌐 Open your browser and go to: http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func startEnhancedWeb() {
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

func startCLISimple() {
	greetUsers()

	firstName, lastName, email, userTickets := getUserInput()

	isValidName, isValidEmail, isValidTicketNumber := ValidateUserInputSimple(firstName, lastName, email, userTickets, remainingTickets)

	if isValidName && isValidEmail && isValidTicketNumber {
		fmt.Println("Thank you for your booking!")
		bookTicketSimple(userTickets, firstName, lastName, email)
		wg.Add(1)
		go sendTicketSimple(userTickets, firstName, lastName, email)
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

// Simple web handlers
func simpleHomeHandlerMain(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>{{.EventName}} - Booking</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; }
        .form-group { margin-bottom: 15px; }
        label { display: block; margin-bottom: 5px; }
        input { width: 100%; padding: 8px; margin-bottom: 10px; box-sizing: border-box; }
        button { background-color: #4CAF50; color: white; padding: 10px 20px; border: none; cursor: pointer; }
        .error { color: red; }
        .success { color: green; }
        .info { background-color: #f0f0f0; padding: 15px; margin-bottom: 20px; border-radius: 5px; }
    </style>
</head>
<body>
    <h1>{{.EventName}}</h1>
    <div class="info">
        <p><strong>Total Tickets:</strong> {{.TotalTickets}}</p>
        <p><strong>Remaining:</strong> {{.RemainingTickets}}</p>
    </div>
    
    {{if .Error}}<p class="error">{{.Error}}</p>{{end}}
    {{if .Message}}<p class="success">{{.Message}}</p>{{end}}
    
    <form method="POST" action="/simple-book">
        <div class="form-group">
            <label>First Name:</label>
            <input type="text" name="firstName" required minlength="2">
        </div>
        <div class="form-group">
            <label>Last Name:</label>
            <input type="text" name="lastName" required minlength="2">
        </div>
        <div class="form-group">
            <label>Email:</label>
            <input type="email" name="email" required>
        </div>
        <div class="form-group">
            <label>Number of Tickets:</label>
            <input type="number" name="tickets" min="1" max="{{.RemainingTickets}}" required>
        </div>
        <button type="submit">Book Tickets</button>
    </form>
    
    <p><a href="/simple-bookings">View All Bookings ({{len .Bookings}} total)</a></p>
</body>
</html>`

	t, err := template.New("home").Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	// Get URL parameters for messages
	message := r.URL.Query().Get("message")
	errorMsg := r.URL.Query().Get("error")

	data := struct {
		EventName        string
		TotalTickets     int
		RemainingTickets uint
		Bookings         []UserData
		Message          string
		Error            string
	}{
		EventName:        eventName,
		TotalTickets:     eventTickets,
		RemainingTickets: remainingTickets,
		Bookings:         bookings,
		Message:          message,
		Error:            errorMsg,
	}

	t.Execute(w, data)
}

func simpleBookHandlerSingle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	email := r.FormValue("email")
	ticketsStr := r.FormValue("tickets")

	tickets, err := strconv.ParseUint(ticketsStr, 10, 32)
	if err != nil {
		http.Redirect(w, r, "/?error=Invalid+ticket+number", http.StatusSeeOther)
		return
	}

	userTickets := uint(tickets)

	// Thread-safe booking
	simpleMutex.Lock()
	defer simpleMutex.Unlock()

	// Validate input
	isValidName, isValidEmail, isValidTicketNumber := ValidateUserInputSimple(firstName, lastName, email, userTickets, remainingTickets)

	if !isValidName {
		http.Redirect(w, r, "/?error=Invalid+name.+Both+first+and+last+names+must+be+at+least+2+characters.", http.StatusSeeOther)
		return
	}

	if !isValidEmail {
		http.Redirect(w, r, "/?error=Invalid+email+address.", http.StatusSeeOther)
		return
	}

	if !isValidTicketNumber {
		http.Redirect(w, r, "/?error=Invalid+ticket+number.+Must+be+between+1+and+available+tickets.", http.StatusSeeOther)
		return
	}

	// Process booking
	remainingTickets -= userTickets

	userData := UserData{
		FirstName:       firstName,
		LastName:        lastName,
		Email:           email,
		NumberOfTickets: userTickets,
	}

	bookings = append(bookings, userData)

	// Simulate async ticket sending
	go func() {
		fmt.Printf("Sending %d tickets to %s %s at %s\n", userTickets, firstName, lastName, email)
	}()

	http.Redirect(w, r, "/?message=Booking+successful!+Thank+you+for+your+purchase.", http.StatusSeeOther)
}

func simpleBookingsHandlerSingle(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>All Bookings - {{.EventName}}</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        table { width: 100%; border-collapse: collapse; margin-top: 20px; }
        th, td { border: 1px solid #ddd; padding: 12px; text-align: left; }
        th { background-color: #f2f2f2; font-weight: bold; }
        tr:nth-child(even) { background-color: #f9f9f9; }
        .summary { background-color: #e7f3ff; padding: 15px; border-radius: 5px; margin-bottom: 20px; }
    </style>
</head>
<body>
    <h1>All Bookings - {{.EventName}}</h1>
    
    <div class="summary">
        <p><strong>Total Bookings:</strong> {{len .Bookings}}</p>
        <p><strong>Tickets Sold:</strong> {{.TicketsSold}}</p>
        <p><strong>Remaining Tickets:</strong> {{.RemainingTickets}}</p>
    </div>
    
    {{if .Bookings}}
    <table>
        <tr>
            <th>#</th>
            <th>First Name</th>
            <th>Last Name</th>
            <th>Email</th>
            <th>Tickets</th>
        </tr>
		{{range $index, $booking := .Bookings}}
		<tr>
			<td>{{add $index 1}}</td>
			<td>{{$booking.FirstName}}</td>
			<td>{{$booking.LastName}}</td>
			<td>{{$booking.Email}}</td>
			<td>{{$booking.NumberOfTickets}}</td>
		</tr>
		{{end}}
    </table>
    {{else}}
    <p>No bookings yet.</p>
    {{end}}
    
    <p><a href="/">Back to Booking</a></p>
</body>
</html>`

	// Template function to add numbers
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	t, err := template.New("bookings").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	ticketsSold := eventTickets - int(remainingTickets)

	data := struct {
		EventName        string
		Bookings         []UserData
		RemainingTickets uint
		TicketsSold      int
	}{
		EventName:        eventName,
		Bookings:         bookings,
		RemainingTickets: remainingTickets,
		TicketsSold:      ticketsSold,
	}

	t.Execute(w, data)
}

// CLI functions
func greetUsers() {
	fmt.Printf("Welcome to %v booking application\n", eventName)
	fmt.Printf("We have a total of %v tickets and %v are still available\n", eventTickets, remainingTickets)
	fmt.Println("Get Your Tickets Here to Attend!")
}

func getFirstNames() []string {
	firstNames := []string{}
	for _, booking := range bookings {
		firstNames = append(firstNames, booking.FirstName)
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

func bookTicketSimple(userTickets uint, firstName string, lastName string, email string) {
	remainingTickets = remainingTickets - userTickets

	var userData = UserData{
		FirstName:       firstName,
		LastName:        lastName,
		Email:           email,
		NumberOfTickets: userTickets,
	}

	bookings = append(bookings, userData)
	fmt.Printf("List of bookings is %v\n", bookings)
	fmt.Printf("Thank you %v %v for booking %v tickets. You will receive a confirmation email at %v\n", firstName, lastName, userTickets, email)
	fmt.Printf("%v tickets remaining for %v\n", remainingTickets, firstName)
}

func sendTicketSimple(userTickets uint, firstName string, lastName string, email string) {
	time.Sleep(10 * time.Second)
	fmt.Println("###############")
	fmt.Printf("Sending ticket:\n %v tickets for %v %v\n", userTickets, firstName, lastName)
	fmt.Printf("Sending ticket confirmation to %v %v at %v for %v tickets.\n", firstName, lastName, email, userTickets)
	fmt.Println("###############")
	wg.Done()
}
