package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func initializeDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./bookings.db")
	if err != nil {

		log.Fatal(err)
	}

	// Create bookings table
	createTable := `
	CREATE TABLE IF NOT EXISTS bookings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT NOT NULL,
		number_of_tickets INTEGER NOT NULL,
		booking_date DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func saveBookingToDB(db *sql.DB, userData UserData) error {
	query := `INSERT INTO bookings (first_name, last_name, email, number_of_tickets) 
			  VALUES (?, ?, ?, ?)`

	_, err := db.Exec(query, userData.firstName, userData.lastName, userData.email, userData.numberOfTickets)
	return err
}

func getBookingsFromDB(db *sql.DB) ([]UserData, error) {
	rows, err := db.Query("SELECT first_name, last_name, email, number_of_tickets FROM bookings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []UserData
	for rows.Next() {
		var userData UserData
		err := rows.Scan(&userData.firstName, &userData.lastName, &userData.email, &userData.numberOfTickets)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, userData)
	}

	return bookings, nil
}
