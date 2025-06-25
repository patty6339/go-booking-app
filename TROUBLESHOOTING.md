# Troubleshooting Guide

## ✅ RESOLVED: Your App is Ready!

Your Go booking application is now **ready to run**! All major issues have been resolved.

## Quick Start Commands

### Option 1: Use the Runner Script (Recommended)
```bash
./run.sh
```

### Option 2: Direct Commands

**Simple Web Mode (Recommended for testing):**
```bash
export PATH=$PATH:/usr/local/go/bin
go run main_single.go helper.go simple-web
```

**Enhanced Web Mode (Full features):**
```bash
export PATH=$PATH:/usr/local/go/bin
go run *.go -web
```

**CLI Mode:**
```bash
export PATH=$PATH:/usr/local/go/bin
go run main_single.go helper.go
```

## What Was Fixed

### ✅ 1. Go Installation
- **Issue**: `bash: go: command not found`
- **Solution**: Installed Go 1.24.4 and added to PATH

### ✅ 2. Dependencies
- **Issue**: Missing SQLite and Stripe packages
- **Solution**: Ran `go mod tidy` to install all dependencies

### ✅ 3. Multiple Main Functions
- **Issue**: Conflicting main functions in different files
- **Solution**: Created `main_single.go` as single entry point

### ✅ 4. Missing Imports
- **Issue**: Missing template and fmt imports in web files
- **Solution**: Added all required imports

### ✅ 5. Authentication Issues
- **Issue**: Auth middleware blocking public pages
- **Solution**: Removed auth from public routes in simple mode

## Current Status: ✅ WORKING

Your application now has:

- ✅ **Simple Web Interface**: Basic booking system without authentication
- ✅ **Enhanced Web Interface**: Full features with auth, payments, etc.
- ✅ **CLI Interface**: Original command-line version
- ✅ **Database Integration**: SQLite for data persistence
- ✅ **Input Validation**: Comprehensive form validation
- ✅ **Concurrent Processing**: Goroutines for async operations
- ✅ **Email System**: SMTP integration (needs configuration)
- ✅ **Payment Processing**: Stripe integration (needs API keys)

## Environment Variables (Optional)

For full functionality, set these environment variables:

```bash
# For email functionality
export SENDER_EMAIL="your-email@gmail.com"
export SENDER_PASS="your-app-password"

# For Stripe payments
export STRIPE_SECRET_KEY="sk_test_your_secret_key"
export STRIPE_PUBLISHABLE_KEY="pk_test_your_publishable_key"
```

## Testing the Application

1. **Start the simple web server:**
   ```bash
   ./run.sh
   # Choose option 1
   ```

2. **Open your browser to:** `http://localhost:8080`

3. **Test booking tickets:**
   - Fill out the form with valid data
   - Submit the booking
   - Check the bookings page

## Common Issues (If Any)

### Port Already in Use
```bash
# Kill any process using port 8080
sudo lsof -ti:8080 | xargs kill -9
```

### Permission Issues
```bash
# Make sure the script is executable
chmod +x run.sh
```

### Go Path Issues
```bash
# Ensure Go is in PATH
export PATH=$PATH:/usr/local/go/bin
go version
```

## Success! 🎉

Your Go booking application is now fully functional and ready for use!
