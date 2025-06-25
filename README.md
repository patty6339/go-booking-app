# Go Booking App

A comprehensive ticket booking application built with Go that demonstrates modern web development concepts including concurrent programming, database integration, user authentication, payment processing, and multi-event management.

## Overview

This application provides both command-line and web interfaces for booking tickets to various events. It started as a simple CLI application for the "Scrabble National Championship" and has evolved into a full-featured booking system with enterprise-level capabilities.

## Features

### Core Features
- **Interactive CLI**: Original command-line interface for user input
- **Web Interface**: Modern HTML-based booking system with responsive design
- **Input Validation**: Comprehensive validation for user data and business rules
- **Concurrent Processing**: Uses goroutines for asynchronous operations
- **Real-time Inventory**: Tracks remaining tickets and prevents overbooking

### Enhanced Features
- **Database Integration**: SQLite database for persistent data storage
- **User Authentication**: Secure login/registration system with session management
- **Real Email Notifications**: SMTP integration for ticket confirmations
- **Payment Processing**: Stripe integration for secure credit card payments
- **Multiple Event Support**: Manage and book tickets for various events
- **Booking Management**: Complete CRUD operations for bookings
- **Admin Panel**: Event management and booking oversight

## Project Structure

```
go-booking-app/
├── main.go              # Original CLI application
├── main_enhanced.go     # Enhanced application with web support
├── helper.go           # Input validation helper functions
├── database.go         # Database operations and SQLite integration
├── email.go            # Email sending functionality (SMTP)
├── web.go              # Web interface handlers and templates
├── auth.go             # User authentication and session management
├── payment.go          # Stripe payment processing
├── events.go           # Multiple event management system
├── go.mod              # Go module definition with dependencies
├── go.sum              # Dependency checksums
├── bookings.db         # SQLite database (created at runtime)
└── README.md           # Project documentation
```

## Requirements

- Go 1.24.4 or higher
- Terminal/Command prompt
- SQLite3 (for database functionality)
- SMTP email account (for email notifications)
- Stripe account (for payment processing)

## Dependencies

The application uses the following Go packages:
- `github.com/mattn/go-sqlite3` - SQLite database driver
- `github.com/stripe/stripe-go/v72` - Stripe payment processing
- Standard library packages for HTTP, templates, crypto, etc.

## Installation

1. **Install Go** (if not already installed):
   ```bash
   # Download and install Go
   wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
   source ~/.bashrc
   ```

2. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd go-booking-app
   ```

3. **Install dependencies**:
   ```bash
   go mod tidy
   ```

4. **Set up environment variables** (create a `.env` file or export):
   ```bash
   # For email functionality
   export SENDER_EMAIL="your-email@gmail.com"
   export SENDER_PASS="your-app-password"
   
   # For Stripe payments
   export STRIPE_SECRET_KEY="sk_test_your_secret_key"
   export STRIPE_PUBLISHABLE_KEY="pk_test_your_publishable_key"
   ```

## Usage

### CLI Mode (Original)
1. Run the original CLI application:
   ```bash
   go run main.go helper.go
   ```

### Enhanced CLI Mode
1. Run the enhanced CLI application with database support:
   ```bash
   go run *.go
   ```

### Web Mode
1. Run the web application:
   ```bash
   go run *.go -web
   ```
   
   Or build and run:
   ```bash
   go build -o booking-app
   ./booking-app -web
   ```

2. Open your browser and navigate to:
   - `http://localhost:8080` - Main booking interface
   - `http://localhost:8080/events` - Browse all events
   - `http://localhost:8080/login` - User login
   - `http://localhost:8080/register` - User registration

### CLI Usage Steps
1. Follow the interactive prompts:
   - Enter your first name (minimum 2 characters)
   - Enter your last name (minimum 2 characters)
   - Enter your email address (must contain @)
   - Enter the number of tickets you want to book

2. The application will:
   - Validate your input
   - Process your booking if valid
   - Save booking to database
   - Display booking confirmation
   - Send email confirmation (if configured)
   - Show updated ticket availability

## Code Structure

### Main Components

- **UserData Struct**: Stores user booking information
- **Global Variables**: 
  - `eventTickets`: Total available tickets (200)
  - `remainingTickets`: Current available tickets
  - `bookings`: Slice storing all booking records
  - `wg`: WaitGroup for goroutine synchronization

### Key Functions

- `main()`: Application entry point and flow control
- `greetUsers()`: Displays welcome message and ticket availability
- `getUserInput()`: Collects user input from command line
- `ValidateUserInput()`: Validates user input (in helper.go)
- `bookTicket()`: Processes valid bookings and updates inventory
- `sendTicket()`: Simulates asynchronous ticket delivery
- `getFirstNames()`: Extracts first names from all bookings

### Enhanced Functions

#### Database Operations (database.go)
- `initDB()`: Initialize SQLite database and create tables
- `saveBookingToDB()`: Save booking information to database
- `getBookingsFromDB()`: Retrieve all bookings from database

#### Email Functionality (email.go)
- `sendRealEmail()`: Send actual emails via SMTP
- `sendTicketConfirmation()`: Send booking confirmation emails
- `getEmailConfig()`: Configure email settings from environment

#### Web Interface (web.go)
- `startWebServer()`: Initialize HTTP server and routes
- `homeHandler()`: Handle main booking page
- `bookHandler()`: Process web-based bookings
- `bookingsHandler()`: Display all bookings in web interface

#### Authentication (auth.go)
- `registerUser()`: User registration with password hashing
- `loginUser()`: User authentication and session creation
- `validateSession()`: Session validation middleware
- `authMiddleware()`: Protect routes requiring authentication

#### Payment Processing (payment.go)
- `createPaymentIntent()`: Create Stripe payment intent
- `paymentHandler()`: Handle payment processing requests
- `paymentPageHandler()`: Display payment form

#### Event Management (events.go)
- `createEvent()`: Create new events
- `bookEventTicket()`: Book tickets for specific events
- `eventsListHandler()`: Display available events
- `bookEventHandler()`: Handle event-specific bookings

## Validation Rules

- **Name**: Both first and last names must be at least 2 characters long
- **Email**: Must contain the "@" symbol
- **Tickets**: Must be greater than 0 and not exceed remaining tickets

## Concurrency Features

The application demonstrates Go's concurrency features:
- **Goroutines**: Ticket sending runs asynchronously
- **WaitGroups**: Ensures main function waits for goroutine completion
- **Synchronization**: Prevents race conditions in ticket booking

## Example Output

```
Welcome to Scrabble National Championship booking application
We have a total of 200 tickets and 200 are still available
Get Your Tickets Here to Attend!

Enter your first name:
John
Enter your last name:
Doe
Enter your email:
john.doe@email.com
Enter number of tickets:
2

Thank you for your booking!
Thank you John Doe for booking 2 tickets. You will receive a confirmation email at john.doe@email.com
198 tickets remaining for John
The first names of the bookings are: [John]

###############
Sending ticket:
 2 tickets for John Doe
Sending ticket confirmation to John Doe at john.doe@email.com for 2 tickets.
###############
```

## Future Enhancements

- Database integration for persistent storage ✅ **IMPLEMENTED**
- Real email sending functionality ✅ **IMPLEMENTED**
- Web interface ✅ **IMPLEMENTED**
- Payment processing ✅ **IMPLEMENTED**
- Multiple event support ✅ **IMPLEMENTED**
- User authentication ✅ **IMPLEMENTED**
- Booking cancellation
- Ticket transfer functionality
- Admin dashboard with analytics
- Mobile responsive design improvements
- API endpoints for third-party integrations
- Advanced reporting and analytics
- Notification system (SMS, push notifications)
- QR code generation for tickets
- Seat selection for events
- Discount codes and promotional pricing

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is open source and available under the [MIT License](LICENSE).

## Contact

For questions or suggestions, please open an issue in the repository.
