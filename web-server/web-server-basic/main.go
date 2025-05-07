package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	server *http.Server
	logger *log.Logger
}

func NewServer(port string, logger *log.Logger) *Server {
	// Create file server handler
	fs := http.FileServer(http.Dir("./static"))

	// Create mux and register routes
	mux := http.NewServeMux()
	mux.Handle("/", fs)

	// Configure server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		server: srv,
		logger: logger,
	}
}

func (s *Server) Start() error {
	// Start server in a goroutine
	go func() {
		s.logger.Printf("Starting server on %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Server error: %v", err)
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Println("Server is shutting down...")
	return s.server.Shutdown(ctx)
}

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "[SERVER] ", log.LstdFlags|log.Lshortfile)

	// Create and start server
	srv := NewServer("8081", logger)
	if err := srv.Start(); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server exited properly")
}
