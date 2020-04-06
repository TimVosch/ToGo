package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/timvosch/togo/pkg/todoserver"
)

var (
	addr = flag.String("addr", ":3000", "Set the listening address")
)

func init() {
	flag.Parse()
}

func main() {
	// Capture system signals
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan)

	// Create server
	s := todoserver.NewServer(*addr)

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
