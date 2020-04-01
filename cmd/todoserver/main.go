package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/timvosch/togo/pkg/todoserver"
)

func main() {
	// Capture system signals
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan)

	// Create server
	s := todoserver.NewServer()

	// Start the server in a new goroutine
	go func() {
		if err := s.StartServer(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen: %v\n", err)
		}
	}()
	log.Println("Server is ready to handle requests")

	<-sigChan

	// Shutdown the server
	log.Println("Shutting down...")
	s.Shutdown()
}
