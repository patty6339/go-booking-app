package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type PageData struct {
	EventName        string
	TotalTickets     int
	RemainingTickets uint
	Bookings         []UserData
	Message          string
	Error            string
}

func startWebServer() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/book", bookHandler)
	http.HandleFunc("/bookings", bookingsHandler)
	
	fmt.Println("Web server starting on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>{{.EventName}} - Booking</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; }
        .form-group { margin-bottom: 15px; }
        label { display: block; margin-bottom: 5px; }
        input { width: 100%; padding: 8px; margin-bottom: 10px; }
        button { background-color: #4CAF50; color: white; padding: 10px 20px; border: none; cursor: pointer; }
        .error { color: red; }
        .success { color: green; }
    </style>
</head>
<body>
    <h1>{{.EventName}}</h1>
    <p>Total Tickets: {{.TotalTickets}} | Remaining: {{.RemainingTickets}}</p>
    
    {{if .Error}}<p class="error">{{.Error}}</p>{{end}}
    {{if .Message}}<p class="success">{{.Message}}</p>{{end}}
    
    <form method="POST" action="/book">
        <div class="form-group">
            <label>First Name:</label>
            <input type="text" name="firstName" required>
        </div>
        <div class="form-group">
            <label>Last Name:</label>
            <input type="text" name="lastName" required>
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
    
    <p><a href="/bookings">View All Bookings</a></p>
</body>
</html>`

	t, _ := template.New("home").Parse(tmpl)
	data := PageData{
		EventName:        eventName,
		TotalTickets:     eventTickets,
		RemainingTickets: remainingTickets,
	}
	t.Execute(w, data)
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		firstName := r.FormValue("firstName")
		lastName := r.FormValue("lastName")
		email := r.FormValue("email")
		ticketsStr := r.FormValue("tickets")
		
		tickets, err := strconv.ParseUint(ticketsStr, 10, 32)
		if err != nil {
			http.Redirect(w, r, "/?error=Invalid ticket number", http.StatusSeeOther)
			return
		}
		
		userTickets := uint(tickets)
		isValidName, isValidEmail, isValidTicketNumber := ValidateUserInput(firstName, lastName, email, userTickets, remainingTickets)
		
		if !isValidName || !isValidEmail || !isValidTicketNumber {
			http.Redirect(w, r, "/?error=Invalid input data", http.StatusSeeOther)
			return
		}
		
		bookTicket(userTickets, firstName, lastName, email)
		go sendTicket(userTickets, firstName, lastName, email)
		
		http.Redirect(w, r, "/?message=Booking successful!", http.StatusSeeOther)
		return
	}
	
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func bookingsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>All Bookings</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        table { width: 100%; border-collapse: collapse; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
    <h1>All Bookings</h1>
    <table>
        <tr>
            <th>First Name</th>
            <th>Last Name</th>
            <th>Email</th>
            <th>Tickets</th>
        </tr>
        {{range .Bookings}}
        <tr>
            <td>{{.firstName}}</td>
            <td>{{.lastName}}</td>
            <td>{{.email}}</td>
            <td>{{.numberOfTickets}}</td>
        </tr>
        {{end}}
    </table>
    <p><a href="/">Back to Booking</a></p>
</body>
</html>`

	t, _ := template.New("bookings").Parse(tmpl)
	data := PageData{
		Bookings: bookings,
	}
	t.Execute(w, data)
}
