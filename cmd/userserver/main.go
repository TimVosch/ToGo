package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/timvosch/togo/pkg/userserver"
)

var (
	addr        = flag.String("addr", ":3000", "Set the listening address")
	privKeyPath = flag.String("privkey", "./private.pem", "Path to the private RSA key")
)

func init() {
	flag.Parse()
}

func main() {
	// Capture system signals
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan)

	// Create server
	us := userserver.NewServer(*addr, *privKeyPath)

	// Start the server in a new goroutine
	go func() {
		if err := us.StartServer(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen: %v\n", err)
		}
	}()
	log.Println("Server is ready to handle requests")

	<-sigChan

	// Shutdown the server
	log.Println("Shutting down...")
	us.Shutdown()
}
