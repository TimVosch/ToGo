package userserver

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/timvosch/togo/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/square/go-jose.v2"
)

// UserServer ...
type UserServer struct {
	httpServer *http.Server
	router     *mux.Router
	repo       UserRepository
	signer     *jwt.Signer
	verifier   *jwt.Verifier
	jwks       *jose.JSONWebKeySet
}

func createJWT(keyData []byte) (*jwt.Signer, *jwt.Verifier) {
	signer, err := jwt.NewSigner(keyData)
	if err != nil {
		log.Fatalln("Could not create signer: ", err)
	}

	verifier, err := jwt.NewVerifier(keyData)
	if err != nil {
		log.Fatalln("Could verifier: ", err)
	}

	return signer, verifier
}

func createJWK(keyData []byte) *jose.JSONWebKey {
	b, _ := pem.Decode(keyData)
	key, _ := x509.ParsePKCS1PrivateKey(b.Bytes)

	// Create jwk
	var jwk = &jose.JSONWebKey{
		Key:       key,
		Algorithm: "RS256",
		Use:       "sig",
		KeyID:     "0",
	}

	return jwk
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

	// JWT and JWK
	data, err := ioutil.ReadFile("./private.pem")
	if err != nil {
		log.Fatalln("Could not read private key file: ", err)
	}
	signer, verifier := createJWT(data)
	jwk := createJWK(data).Public()

	// Build UserServer struct
	s := &UserServer{
		httpServer,
		router,
		repo,
		signer,
		verifier,
		&jose.JSONWebKeySet{
			Keys: []jose.JSONWebKey{
				jwk,
			},
		},
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
	token := jwt.CreateToken()
	token.Body = map[string]interface{}{
		"sub": user.ID,
	}
	return us.signer.Sign(token), nil
}
