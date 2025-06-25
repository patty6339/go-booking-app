#!/bin/bash

# Go Booking App Runner Script

export PATH=$PATH:/usr/local/go/bin

echo "🎫 Go Booking App Runner"
echo "========================"
echo ""
echo "Choose an option:"
echo "1. 🌐 Simple Web Mode (Recommended - Works perfectly!)"
echo "2. 💻 CLI Mode (Command line interface)"
echo "3. 🚀 Enhanced Web Mode (Full features - requires setup)"
echo ""
read -p "Enter your choice (1-3): " choice

case $choice in
    1)
        echo ""
        echo "🚀 Starting Simple Web Mode..."
        echo "📍 Open your browser to: http://localhost:8080"
        echo "⏹️  Press Ctrl+C to stop the server"
        echo ""
        go run app.go simple-web
        ;;
    2)
        echo ""
        echo "🚀 Starting CLI Mode..."
        echo ""
        go run app.go
        ;;
    3)
        echo ""
        echo "🚀 Starting Enhanced Web Mode..."
        echo "📍 Open your browser to: http://localhost:8080"
        echo "⏹️  Press Ctrl+C to stop the server"
        echo "⚠️  Note: This requires database and other dependencies"
        echo ""
        go run *.go -web
        ;;
    *)
        echo "Invalid choice. Please run the script again."
        exit 1
        ;;
esac
