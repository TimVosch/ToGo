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
	addr     = flag.String("addr", ":3000", "Set the listening address")
	jwksURL  = flag.String("jwks", "http://127.0.0.1:3001/.well-known/jwks.json", "The URL for fetching public keys as JWKS")
	mongoURI = flag.String("mongoURI", "mongo://your-connection-string/", "This is the mongo connection URI")
)

func init() {
	flag.Parse()
}

func main() {
	// Capture system signals
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan)

	// Create server
	s := todoserver.NewServer(*addr, *jwksURL, *mongoURI)

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
