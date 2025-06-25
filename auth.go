package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	ID       int
	Username string
	Email    string
	Password string // This should be hashed in production
	Created  time.Time
}

type Session struct {
	Token   string
	UserID  int
	Expires time.Time
}

var users = make(map[string]User)       // username -> User
var usersByID = make(map[int]User)      // userID -> User
var sessions = make(map[string]Session) // token -> Session

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func generateSessionToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// fallback to time-based token if random fails (not recommended for production)
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

func registerUser(username, email, password string) error {
	if _, exists := users[username]; exists {
		return fmt.Errorf("user already exists")
	}
	user := User{
		ID:       len(users) + 1,
		Username: username,
		Email:    email,
		Password: hashPassword(password),
		Created:  time.Now(),
	}

	users[username] = user
	usersByID[user.ID] = user
	return nil
}

func loginUser(username, password string) (string, error) {
	user, exists := users[username]
	if !exists {
		return "", fmt.Errorf("user not found")
	}

	if user.Password != hashPassword(password) {
		return "", fmt.Errorf("invalid password")
	}

	// Create session
	token := generateSessionToken()
	session := Session{
		Token:   token,
		UserID:  user.ID,
		Expires: time.Now().Add(24 * time.Hour), // 24 hour session
	}

	sessions[token] = session
	return token, nil
}

func validateSession(token string) (User, bool) {
	session, exists := sessions[token]
	if !exists || time.Now().After(session.Expires) {
		return User{}, false
	}
	// Find user by ID efficiently
	user, exists := usersByID[session.UserID]
	if exists {
		return user, true
	}
	return User{}, false
}

func requireAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		_, valid := validateSession(cookie.Value)
		if !valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next(w, r)
	}
}

var loginTemplate = `
<!DOCTYPE html>
<html>
<head><title>Login</title></head>
<body>
	<h2>Login</h2>
	<form method="POST">
		<div>
			<label>Username:</label>
			<input type="text" name="username" required>
		</div>
		<div>
			<label>Password:</label>
			<input type="password" name="password" required>
		</div>
		<button type="submit">Login</button>
	</form>
	<p><a href="/register">Register</a></p>
</body>
</html>`

func authLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		token, err := loginUser(username, password)
		if err != nil {
			// Redirect to a fixed, safe relative path to prevent SSRF
			http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
			return
		}

		// Set session cookie
		cookie := &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   true,                    // Ensure cookie is sent only over HTTPS
			SameSite: http.SameSiteStrictMode, // Optional: helps prevent CSRF
		}
		http.SetCookie(w, cookie)

		// Redirect to a fixed, safe relative path to prevent SSRF
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Show login form
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(loginTemplate))
}

func authRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		err := registerUser(username, email, password)
		if err != nil {
			http.Redirect(w, r, "/register?error=Registration failed", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/login?message=Registration successful", http.StatusSeeOther)
		return
	}

	// Show registration form
	tmpl := `
<!DOCTYPE html>
<html>
<head><title>Register</title></head>
<body>
	<h2>Register</h2>
	<form method="POST">
		<div>
			<label>Username:</label>
			<input type="text" name="username" required>
		</div>
		<div>
			<label>Email:</label>
			<input type="email" name="email" required>
		</div>
		<div>
			<label>Password:</label>
			<input type="password" name="password" required>
		</div>
		<button type="submit">Register</button>
	</form>
	<p><a href="/login">Login</a></p>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}
