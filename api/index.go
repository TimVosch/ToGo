package api

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"net/http"

	"github.com/timvosch/togo/pkg/jwt"
	"github.com/timvosch/togo/pkg/userserver"
)

func createJWT(key *rsa.PrivateKey) (*jwt.Signer, *jwt.Verifier) {
	signer := jwt.NewSigner(key)

	verifier, err := jwt.NewVerifier(key)
	if err != nil {
		log.Fatalln("Could not create verifier: ", err)
	}

	return signer, verifier
}

// Users is the handler for the Now serverless function
func Users(w http.ResponseWriter, r *http.Request) {
	// Set up
	b, _ := pem.Decode([]byte(
		"",
	))
	key, err := x509.ParsePKCS1PrivateKey(b.Bytes)
	if err != nil {
		log.Fatalln("Could not read private key!")
	}

	us := userserver.NewServerless()
	us.Router.ServeHTTP(w, r)
}
