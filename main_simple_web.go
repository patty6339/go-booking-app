package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Simple web version without authentication for testing
const simpleEventTickets int = 200
var simpleEventName = "Scrabble National Championship"
var simpleRemainingTickets uint = 200
var simpleBookings = make([]UserData, 0)
var simpleMutex = sync.Mutex{}

func simpleHomeHandler(w http.ResponseWriter, r *http.Request) {
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
		EventName:        simpleEventName,
		TotalTickets:     simpleEventTickets,
		RemainingTickets: simpleRemainingTickets,
		Bookings:         simpleBookings,
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
	isValidName, isValidEmail, isValidTicketNumber := ValidateUserInput(firstName, lastName, email, userTickets, simpleRemainingTickets)
	
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
	simpleRemainingTickets -= userTickets
	
	userData := UserData{
		firstName:       firstName,
		lastName:        lastName,
		email:           email,
		numberOfTickets: userTickets,
	}
	
	simpleBookings = append(simpleBookings, userData)
	
	// Simulate async ticket sending
	go func() {
		fmt.Printf("Sending %d tickets to %s %s at %s\n", userTickets, firstName, lastName, email)
	}()
	
	http.Redirect(w, r, "/?message=Booking+successful!+Thank+you+for+your+purchase.", http.StatusSeeOther)
}

func simpleBookingsHandler(w http.ResponseWriter, r *http.Request) {
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
            <td>{{$booking.firstName}}</td>
            <td>{{$booking.lastName}}</td>
            <td>{{$booking.email}}</td>
            <td>{{$booking.numberOfTickets}}</td>
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

	ticketsSold := simpleEventTickets - int(simpleRemainingTickets)

	data := struct {
		EventName        string
		Bookings         []UserData
		RemainingTickets uint
		TicketsSold      int
	}{
		EventName:        simpleEventName,
		Bookings:         simpleBookings,
		RemainingTickets: simpleRemainingTickets,
		TicketsSold:      ticketsSold,
	}
	
	t.Execute(w, data)
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
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Simple main function for testing
func runSimpleWeb() {
	fmt.Println("Starting Simple Web Booking Application...")
	startSimpleWeb()
}
