package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/janislaus/figure10/internal/db"
	"github.com/janislaus/figure10/internal/handlers"
	"github.com/janislaus/figure10/internal/llm"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize database
	database, err := sql.Open("sqlite3", "./figure10.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	// Initialize database schema
	if err := db.InitDB(database); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Update the generator initialization in main.go
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println("Warning: GEMINI_API_KEY environment variable not set. Using fallback text generation.")
	} else {
		fmt.Printf("Using Gemini API with key: %s...\n", apiKey[:5]+"...") // Only show first few chars
	}
	generator := llm.NewTextGenerator(apiKey)

	// Create handler with dependencies
	h := handlers.NewHandler(database, generator)

	// Set up static file server
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Set up routes
	http.HandleFunc("/", h.HandleHome)
	http.HandleFunc("/generate-text", h.HandleGenerateText)
	http.HandleFunc("/start-session", h.HandleStartSession)
	http.HandleFunc("/submit-result", h.HandleSubmitResult)
	http.HandleFunc("/check-typing", h.HandleCheckTyping)
	http.HandleFunc("/history", h.HandleHistory)
	http.HandleFunc("/generate-practice", h.HandleGeneratePractice)

	// Create server
	port := "8081"
	server := &http.Server{
		Addr:    ":" + port,
		Handler: nil, // Use default ServeMux
	}

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		fmt.Printf("Starting Figure10 server on port %s...\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	fmt.Println("\nShutting down server...")

	// Create a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server gracefully stopped")
}
