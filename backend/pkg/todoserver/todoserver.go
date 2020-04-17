package todoserver

import (
	"context"
	"crypto/rsa"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/timvosch/togo/pkg/common/keys"
	"github.com/timvosch/togo/pkg/jwt"
)

// TodoServer contains relevant todo Server objects
type TodoServer struct {
	httpServer *http.Server
	Router     *mux.Router
	db         TodoRepository
	jwt        *jwt.Verifier
}

func createJWT(jwksURL string) *jwt.Verifier {
	jwks, err := keys.FetchPublicKeys(jwksURL)
	if err != nil {
		log.Fatalln("Could not read key file: ", err)
	}

	// Decode JWK into PubKey
	pubKey, ok := jwks.Keys[0].Key.(*rsa.PublicKey)
	if !ok {
		log.Fatalln("Could not create RSA key from JWK.")
	}

	// Create verifier with pubkey
	v, _ := jwt.NewVerifier(pubKey)

	return v
}

// NewServer creates a new server
func NewServer(addr, jwksURL, mongoConnection string) *TodoServer {
	// Set up
	router := mux.NewRouter()
	httpServer := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// db := NewTodoMemoryRepository()
	db := NewMongoRepository(mongoConnection, "togo", "todos")
	jwt := createJWT(jwksURL)

	// Build TodoServer struct
	s := &TodoServer{
		httpServer,
		router,
		db,
		jwt,
	}

	setRoutes(s)

	return s
}

// NewServerless creates a TodoServer without http server
func NewServerless(key *rsa.PrivateKey, mongoURI, jwksURL, prefix string) *TodoServer {
	router := mux.NewRouter().PathPrefix(prefix).Subrouter()
	db := NewMongoRepository(mongoURI, "togo", "todos")
	jwt := createJWT(jwksURL)

	// Build UserServer struct
	s := &TodoServer{
		nil,
		router,
		db,
		jwt,
	}

	setRoutes(s)
	return s
}

// StartServer creates and starts a TodoServer
func (s *TodoServer) StartServer() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown the server
func (s *TodoServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}

}
