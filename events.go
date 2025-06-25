package main

import (
	"fmt"
	"time"
)

type Event struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Date             time.Time `json:"date"`
	Location         string    `json:"location"`
	TotalTickets     int       `json:"total_tickets"`
	RemainingTickets int       `json:"remaining_tickets"`
	TicketPrice      float64   `json:"ticket_price"`
	Active           bool      `json:"active"`
}

type EventBooking struct {
	ID              int       `json:"id"`
	EventID         int       `json:"event_id"`
	UserID          int       `json:"user_id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email"`
	NumberOfTickets int       `json:"number_of_tickets"`
	TotalAmount     float64   `json:"total_amount"`
	BookingDate     time.Time `json:"booking_date"`
	Status          string    `json:"status"` // pending, confirmed, cancelled
}

var events = make(map[int]Event)
var eventBookings = make([]EventBooking, 0)
var nextEventID = 1
var nextBookingID = 1

/* Removed duplicate initializeEvents function to resolve redeclaration error.
   The implementation should exist in only one file in the package. */

func createEvent(name, description, location string, date time.Time, totalTickets int, ticketPrice float64) Event {
	event := Event{
		ID:               nextEventID,
		Name:             name,
		Description:      description,
		Date:             date,
		Location:         location,
		TotalTickets:     totalTickets,
		RemainingTickets: totalTickets,
		TicketPrice:      ticketPrice,
		Active:           true,
	}

	events[nextEventID] = event
	nextEventID++
	return event
}

func bookEventTicket(eventID, userID int, firstName, lastName, email string, numberOfTickets int) (*EventBooking, error) {
	event, exists := events[eventID]
	if !exists {
		return nil, fmt.Errorf("event not found")
	}

	if !event.Active {
		return nil, fmt.Errorf("event is not active")
	}

	if numberOfTickets > event.RemainingTickets {
		return nil, fmt.Errorf("not enough tickets available")
	}

	// Update event tickets
	event.RemainingTickets -= numberOfTickets
	events[eventID] = event

	// Create booking
	booking := EventBooking{
		ID:              nextBookingID,
		EventID:         eventID,
		UserID:          userID,
		FirstName:       firstName,
		LastName:        lastName,
		Email:           email,
		NumberOfTickets: numberOfTickets,
		TotalAmount:     float64(numberOfTickets) * event.TicketPrice,
		BookingDate:     time.Now(),
		Status:          "confirmed",
	}

	eventBookings = append(eventBookings, booking)
	nextBookingID++

	return &booking, nil
}

/* Removed duplicate eventsListHandler to resolve redeclaration error.
   The implementation should exist in only one file in the package. */

/* Removed duplicate bookEventHandler to resolve redeclaration error.
   The implementation should exist in only one file in the package. */
