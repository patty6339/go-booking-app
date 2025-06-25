package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Global variables
const eventTickets int = 200
var eventName = "Scrabble National Championship"
var remainingTickets uint = 200
var bookings = make([]UserData, 0)
var wg = sync.WaitGroup{}
var simpleMutex = sync.Mutex{}

type UserData struct {
	firstName       string
	lastName        string
	email           string
	numberOfTickets uint
}

// Validation function
func ValidateUserInput(firstName string, lastName string, email string, userTickets uint, remainingTickets uint) (bool, bool, bool) {
	isValidName := len(firstName) >= 2 && len(lastName) >= 2
	isValidEmail := strings.Contains(email, "@")
	isValidTicketNumber := userTickets > 0 && userTickets <= remainingTickets
	return isValidName, isValidEmail, isValidTicketNumber
}

func main() {
	// Check for simple-web argument
	if len(os.Args) > 1 && os.Args[1] == "simple-web" {
		startSimpleWeb()
		return
	}

	// Default to CLI mode
	startCLI()
}

func startSimpleWeb() {
	http.HandleFunc("/", simpleHomeHandler)
	http.HandleFunc("/simple-book", simpleBookHandler)
	http.HandleFunc("/simple-bookings", simpleBookingsHandler)
	
	fmt.Println("🚀 Simple Booking App starting on http://localhost:8080")
	fmt.Println("📝 Features:")
	fmt.Println("   - Web-based ticket booking")
	fmt.Println("   - Input validation")
	fmt.Println("   - Booking management")
	fmt.Println("   - Real-time ticket tracking")
	fmt.Println()
	fmt.Println("🌐 Open your browser and go to: http://localhost:8080")
	fmt.Println("⏹️  Press Ctrl+C to stop the server")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func startCLI() {
	greetUsers()

	firstName, lastName, email, userTickets := getUserInput()

	isValidName, isValidEmail, isValidTicketNumber := ValidateUserInput(firstName, lastName, email, userTickets, remainingTickets)

	if isValidName && isValidEmail && isValidTicketNumber {
		fmt.Println("Thank you for your booking!")
		bookTicket(userTickets, firstName, lastName, email)
		wg.Add(1)
		go sendTicket(userTickets, firstName, lastName, email)
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
func simpleHomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>{{.EventName}} - Booking</title>
    <style>
        body { 
            font-family: Arial, sans-serif; 
            max-width: 600px; 
            margin: 0 auto; 
            padding: 20px; 
            background-color: #f5f5f5;
        }
        .container {
            background-color: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .form-group { margin-bottom: 15px; }
        label { display: block; margin-bottom: 5px; font-weight: bold; }
        input { 
            width: 100%; 
            padding: 10px; 
            margin-bottom: 10px; 
            box-sizing: border-box; 
            border: 1px solid #ddd;
            border-radius: 5px;
        }
        button { 
            background-color: #4CAF50; 
            color: white; 
            padding: 12px 24px; 
            border: none; 
            cursor: pointer; 
            border-radius: 5px;
            font-size: 16px;
            width: 100%;
        }
        button:hover { background-color: #45a049; }
        .error { color: red; padding: 10px; background-color: #ffe6e6; border-radius: 5px; }
        .success { color: green; padding: 10px; background-color: #e6ffe6; border-radius: 5px; }
        .info { 
            background-color: #e7f3ff; 
            padding: 15px; 
            margin-bottom: 20px; 
            border-radius: 5px; 
            border-left: 4px solid #2196F3;
        }
        .nav-link {
            display: inline-block;
            margin-top: 20px;
            padding: 10px 15px;
            background-color: #2196F3;
            color: white;
            text-decoration: none;
            border-radius: 5px;
        }
        .nav-link:hover { background-color: #1976D2; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎫 {{.EventName}}</h1>
        <div class="info">
            <p><strong>📊 Total Tickets:</strong> {{.TotalTickets}}</p>
            <p><strong>🎟️ Remaining:</strong> {{.RemainingTickets}}</p>
            <p><strong>✅ Sold:</strong> {{.TicketsSold}}</p>
        </div>
        
        {{if .Error}}<div class="error">❌ {{.Error}}</div>{{end}}
        {{if .Message}}<div class="success">✅ {{.Message}}</div>{{end}}
        
        <form method="POST" action="/simple-book">
            <div class="form-group">
                <label>👤 First Name:</label>
                <input type="text" name="firstName" required minlength="2" placeholder="Enter your first name">
            </div>
            <div class="form-group">
                <label>👤 Last Name:</label>
                <input type="text" name="lastName" required minlength="2" placeholder="Enter your last name">
            </div>
            <div class="form-group">
                <label>📧 Email:</label>
                <input type="email" name="email" required placeholder="Enter your email address">
            </div>
            <div class="form-group">
                <label>🎫 Number of Tickets:</label>
                <input type="number" name="tickets" min="1" max="{{.RemainingTickets}}" required placeholder="How many tickets?">
            </div>
            <button type="submit">🎯 Book Tickets Now</button>
        </form>
        
        <a href="/simple-bookings" class="nav-link">📋 View All Bookings ({{len .Bookings}} total)</a>
    </div>
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

	ticketsSold := eventTickets - int(remainingTickets)

	data := struct {
		EventName        string
		TotalTickets     int
		RemainingTickets uint
		TicketsSold      int
		Bookings         []UserData
		Message          string
		Error            string
	}{
		EventName:        eventName,
		TotalTickets:     eventTickets,
		RemainingTickets: remainingTickets,
		TicketsSold:      ticketsSold,
		Bookings:         bookings,
		Message:          message,
		Error:            errorMsg,
	}
	
	t.Execute(w, data)
}

func simpleBookHandler(w http.ResponseWriter, r *http.Request) {
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
	isValidName, isValidEmail, isValidTicketNumber := ValidateUserInput(firstName, lastName, email, userTickets, remainingTickets)
	
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
		firstName:       firstName,
		lastName:        lastName,
		email:           email,
		numberOfTickets: userTickets,
	}
	
	bookings = append(bookings, userData)
	
	// Simulate async ticket sending
	go func() {
		fmt.Printf("📧 Sending %d tickets to %s %s at %s\n", userTickets, firstName, lastName, email)
	}()
	
	successMsg := fmt.Sprintf("Booking+successful!+%d+tickets+booked+for+%s+%s", userTickets, firstName, lastName)
	http.Redirect(w, r, "/?message="+successMsg, http.StatusSeeOther)
}

func simpleBookingsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>All Bookings - {{.EventName}}</title>
    <style>
        body { 
            font-family: Arial, sans-serif; 
            max-width: 1000px; 
            margin: 0 auto; 
            padding: 20px; 
            background-color: #f5f5f5;
        }
        .container {
            background-color: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        table { 
            width: 100%; 
            border-collapse: collapse; 
            margin-top: 20px; 
        }
        th, td { 
            border: 1px solid #ddd; 
            padding: 12px; 
            text-align: left; 
        }
        th { 
            background-color: #4CAF50; 
            color: white;
            font-weight: bold; 
        }
        tr:nth-child(even) { background-color: #f9f9f9; }
        tr:hover { background-color: #f5f5f5; }
        .summary { 
            background-color: #e7f3ff; 
            padding: 20px; 
            border-radius: 5px; 
            margin-bottom: 20px; 
            border-left: 4px solid #2196F3;
        }
        .nav-link {
            display: inline-block;
            margin-top: 20px;
            padding: 10px 15px;
            background-color: #2196F3;
            color: white;
            text-decoration: none;
            border-radius: 5px;
        }
        .nav-link:hover { background-color: #1976D2; }
        .no-bookings {
            text-align: center;
            padding: 40px;
            color: #666;
            font-style: italic;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>📋 All Bookings - {{.EventName}}</h1>
        
        <div class="summary">
            <h3>📊 Booking Summary</h3>
            <p><strong>🎫 Total Bookings:</strong> {{len .Bookings}}</p>
            <p><strong>🎟️ Tickets Sold:</strong> {{.TicketsSold}}</p>
            <p><strong>📈 Remaining Tickets:</strong> {{.RemainingTickets}}</p>
            <p><strong>💰 Revenue:</strong> ${{.Revenue}} (estimated at $50/ticket)</p>
        </div>
        
        {{if .Bookings}}
        <table>
            <tr>
                <th>#</th>
                <th>👤 First Name</th>
                <th>👤 Last Name</th>
                <th>📧 Email</th>
                <th>🎫 Tickets</th>
                <th>💰 Value</th>
            </tr>
            {{range $index, $booking := .Bookings}}
            <tr>
                <td>{{add $index 1}}</td>
                <td>{{$booking.firstName}}</td>
                <td>{{$booking.lastName}}</td>
                <td>{{$booking.email}}</td>
                <td>{{$booking.numberOfTickets}}</td>
                <td>${{multiply $booking.numberOfTickets 50}}</td>
            </tr>
            {{end}}
        </table>
        {{else}}
        <div class="no-bookings">
            <h3>🎭 No bookings yet!</h3>
            <p>Be the first to book tickets for this amazing event.</p>
        </div>
        {{end}}
        
        <a href="/" class="nav-link">🎯 Back to Booking</a>
    </div>
</body>
</html>`

	// Template functions
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"multiply": func(a uint, b int) int {
			return int(a) * b
		},
	}

	t, err := template.New("bookings").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	ticketsSold := eventTickets - int(remainingTickets)
	revenue := ticketsSold * 50 // $50 per ticket

	data := struct {
		EventName        string
		Bookings         []UserData
		RemainingTickets uint
		TicketsSold      int
		Revenue          int
	}{
		EventName:        eventName,
		Bookings:         bookings,
		RemainingTickets: remainingTickets,
		TicketsSold:      ticketsSold,
		Revenue:          revenue,
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
