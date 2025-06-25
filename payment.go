package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
)

const ticketPrice = 5000 // $50.00 in cents

func init() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY") // Set this in your environment
}

type PaymentRequest struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Tickets  int    `json:"tickets"`
}

type PaymentResponse struct {
	ClientSecret string `json:"client_secret"`
	Error        string `json:"error,omitempty"`
}

func createPaymentIntent(tickets int) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(tickets * ticketPrice)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Metadata: map[string]string{
			"tickets": strconv.Itoa(tickets),
		},
	}

	return paymentintent.New(params)
}

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	pi, err := createPaymentIntent(req.Tickets)
	if err != nil {
		response := PaymentResponse{Error: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := PaymentResponse{ClientSecret: pi.ClientSecret}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func paymentPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Payment - {{.EventName}}</title>
    <script src="https://js.stripe.com/v3/"></script>
    <style>
        body { font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; }
        .form-group { margin-bottom: 15px; }
        #card-element { padding: 10px; border: 1px solid #ccc; }
        button { background-color: #4CAF50; color: white; padding: 10px 20px; border: none; cursor: pointer; }
        .error { color: red; }
    </style>
</head>
<body>
    <h1>Payment for {{.EventName}}</h1>
    <div id="payment-form">
        <div class="form-group">
            <label>Number of Tickets:</label>
            <input type="number" id="tickets" min="1" max="10" value="1">
            <p>Price per ticket: $50.00</p>
            <p>Total: $<span id="total">50.00</span></p>
        </div>
        
        <div class="form-group">
            <label>Card Details:</label>
            <div id="card-element"></div>
            <div id="card-errors" class="error"></div>
        </div>
        
        <button id="submit-payment">Pay Now</button>
    </div>

    <script>
        const stripe = Stripe('pk_test_your_publishable_key_here'); // Replace with your publishable key
        const elements = stripe.elements();
        const cardElement = elements.create('card');
        cardElement.mount('#card-element');

        const ticketsInput = document.getElementById('tickets');
        const totalSpan = document.getElementById('total');
        
        ticketsInput.addEventListener('input', function() {
            const tickets = parseInt(this.value) || 1;
            const total = (tickets * 50).toFixed(2);
            totalSpan.textContent = total;
        });

        document.getElementById('submit-payment').addEventListener('click', async function() {
            const tickets = parseInt(ticketsInput.value) || 1;
            
            // Create payment intent
            const response = await fetch('/create-payment-intent', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    amount: tickets * 5000,
                    currency: 'usd',
                    tickets: tickets
                })
            });
            
            const { client_secret, error } = await response.json();
            
            if (error) {
                document.getElementById('card-errors').textContent = error;
                return;
            }
            
            // Confirm payment
            const result = await stripe.confirmCardPayment(client_secret, {
                payment_method: {
                    card: cardElement
                }
            });
            
            if (result.error) {
                document.getElementById('card-errors').textContent = result.error.message;
            } else {
                alert('Payment successful! Redirecting to booking confirmation...');
                window.location.href = '/booking-success?payment_intent=' + result.paymentIntent.id;
            }
        });
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}
