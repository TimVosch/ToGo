package userserver

import (
	"context"
	"crypto/rsa"
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
	Router     *mux.Router
	repo       UserRepository
	signer     *jwt.Signer
	verifier   *jwt.Verifier
	jwks       *jose.JSONWebKeySet
}

func readKey(path string) *rsa.PrivateKey {
	data, err := ioutil.ReadFile("./private.pem")
	if err != nil {
		log.Fatalln("Could not read private key file: ", err)
	}
	b, _ := pem.Decode(data)
	key, err := x509.ParsePKCS1PrivateKey(b.Bytes)
	if err != nil {
		log.Fatalln("Could not create private key from file: ", err)
	}
	return key
}

func createJWT(key *rsa.PrivateKey) (*jwt.Signer, *jwt.Verifier) {
	signer := jwt.NewSigner(key)

	verifier, err := jwt.NewVerifier(key)
	if err != nil {
		log.Fatalln("Could not create verifier: ", err)
	}

	return signer, verifier
}

// NewServer creates a new server
func NewServer(addr, privKeyPath, mongoURI string) *UserServer {
	// Set up
	router := mux.NewRouter()
	httpServer := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	repo := NewUserMongoRepository(mongoURI, "togo", "users")

	// JWT and JWK
	key := readKey(privKeyPath)
	signer, verifier := createJWT(key)
	jwks := jwt.CreateJWKS(key)

	// Build UserServer struct
	s := &UserServer{
		httpServer,
		router,
		repo,
		signer,
		verifier,
		jwks,
	}

	setRoutes(s)

	return s
}

// NewServerless creates a UserServer without http server
func NewServerless(key *rsa.PrivateKey, mongoURI, prefix string) *UserServer {
	router := mux.NewRouter().PathPrefix(prefix).Subrouter()
	repo := NewUserMongoRepository(mongoURI, "togo", "users")
	signer, verifier := createJWT(key)
	jwks := jwt.CreateJWKS(key)

	// Build UserServer struct
	s := &UserServer{
		nil,
		router,
		repo,
		signer,
		verifier,
		jwks,
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
