package userserver

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/timvosch/togo/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

// UserServer ...
type UserServer struct {
	httpServer *http.Server
	router     *mux.Router
	repo       UserRepository
	jwt        *jwt.JWT
}

// NewServer creates a new server
func NewServer(addr, privKeyPath string) *UserServer {
	// Set up
	router := mux.NewRouter()
	httpServer := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	repo := NewUserMemoryRepository()
	j := jwt.NewJWT(privKeyPath)

	if j == nil {
		log.Fatalln("Could not create JWT")
	}

	// Build UserServer struct
	s := &UserServer{
		httpServer,
		router,
		repo,
		j,
	}

	setRoutes(s)

	return s
}

// StartServer creates and starts a UserServer
func (us *UserServer) StartServer() error {
	return us.httpServer.ListenAndServe()
}

// Shutdown the server
func (us *UserServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := us.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
}

func (us *UserServer) loginUser(email, password string) (string, error) {
	user := us.repo.GetUserByEmail(email)
	if user == nil {
		return "", errors.New("")
	}

	// Match password against hash
	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		return "", errors.New("")
	}

	// Create token with subject as user ID
	token := us.jwt.CreateToken()
	token.Body = map[string]interface{}{
		"sub": user.ID,
	}
	return us.jwt.Sign(token), nil
}
