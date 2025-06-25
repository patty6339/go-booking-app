package main

import (
	"fmt"
	"net/smtp"
	"os"
)

// EmailConfig holds SMTP server configuration and sender credentials.
type EmailConfig struct {
	SMTPHost    string
	SMTPPort    string
	SenderEmail string
	SenderPass  string
}

// getEmailConfig retrieves email configuration from environment variables.
func getEmailConfig() EmailConfig {
	return EmailConfig{
		SMTPHost:    "smtp.gmail.com",
		SMTPPort:    "587",
		SenderEmail: os.Getenv("SENDER_EMAIL"), // Set this in your environment (must be a valid Gmail address)
		SenderPass:  os.Getenv("SENDER_PASS"),  // Use Gmail App Password, not your regular password
	}
}

// sendRealEmail sends an email using the SMTP configuration.
func sendRealEmail(recipientEmail, subject, body string) error {
	config := getEmailConfig()

	if config.SenderEmail == "" || config.SenderPass == "" {
		return fmt.Errorf("email credentials not configured")
	}
	if config.SMTPHost == "" || config.SMTPPort == "" {
		return fmt.Errorf("SMTP host or port not configured")
	}

	// Set up authentication information.
	auth := smtp.PlainAuth("", config.SenderEmail, config.SenderPass, config.SMTPHost)

	// Compose the message.
	message := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", recipientEmail, subject, body))

	// Send the email.
	err := smtp.SendMail(config.SMTPHost+":"+config.SMTPPort, auth, config.SenderEmail, []string{recipientEmail}, message)
	return err
}

// ticketConfirmationTemplate is the email template for ticket confirmations.
const ticketConfirmationTemplate = `Dear %s %s,

Thank you for your booking!

Booking Details:
- Event: Scrabble National Championship
- Number of Tickets: %d
- Email: %s

Your tickets will be sent to you shortly.

Best regards,
Booking Team
`

// TicketConfirmationParams holds the parameters for sending a ticket confirmation email.
type TicketConfirmationParams struct {
	UserTickets uint
	FirstName   string
	LastName    string
	Email       string
}

// sendTicketConfirmation sends a ticket confirmation email, or simulates sending if it fails.
func sendTicketConfirmation(params TicketConfirmationParams) {
	subject := "Ticket Confirmation - Scrabble National Championship"
	body := fmt.Sprintf(ticketConfirmationTemplate, params.FirstName, params.LastName, params.UserTickets, params.Email)

	err := sendRealEmail(params.Email, subject, body)
	if err != nil {
		fmt.Printf("Failed to send email: %v\n", err)
		fmt.Println("Falling back to simulation...")
		// Fall back to the original simulation.
		fmt.Printf("Sending ticket confirmation to %s %s at %s for %d tickets.\n", params.FirstName, params.LastName, params.Email, params.UserTickets)
	} else {
		fmt.Printf("Email sent successfully to %s\n", params.Email)
	}
}

// main function to demonstrate sending a ticket confirmation email.
