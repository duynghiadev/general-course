package main

import (
	"log"
	"net/http"
	"os"
	"web-server/handlers"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Create a new ServeMux for routing
	mux := http.NewServeMux()

	// Serve static files from the 'public' directory
	fs := http.FileServer(http.Dir("public"))
	mux.Handle("/", fs)

	// Register API endpoint
	mux.HandleFunc("/api", handlers.APIHandler)

	// Register 404 handler
	mux.HandleFunc("/404", handlers.NotFoundHandler)

	// Get port from environment variable or default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Start the server
	log.Printf("Server running at http://localhost:%s", port)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
