package jwt_test

import (
	"crypto/rand"
	"crypto/rsa"
	_ "crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"log"
	"testing"
	"time"

	"github.com/timvosch/togo/pkg/jwt"
)

func createPEM() []byte {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalln("Error occured while creating rsa key: ", err)
	}
	// pubKey := &privKey.PublicKey
	pemKey := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})
	return pemKey
}

func createSet() (*jwt.Signer, *jwt.Verifier, []byte) {
	pem := createPEM()
	signer, err := jwt.NewSigner(pem)
	if err != nil {
		log.Fatalln("Error occured while creating signer: ", err)
	}
	verifier, err := jwt.NewVerifier(pem)
	if err != nil {
		log.Fatalln("Error occured while creating verifier: ", err)
	}
	return signer, verifier, pem
}

func TestSignDoesVerify(t *testing.T) {
	signer, verifier, _ := createSet()
	token := jwt.CreateToken()
	token.Body = map[string]interface{}{
		"wow": "hello",
	}

	jwtStr := signer.Sign(token)
	t.Log("Got jwt:\n", jwtStr)

	_, err := verifier.Verify(jwtStr)

	if err != nil {
		t.Fatal("Failed to verify our own signed token")
	}
}

func TestErrorInvalidSignature(t *testing.T) {
	// signer and verifier do not share the PEM key, therefore are different
	_, verifier, _ := createSet()
	signer, _, _ := createSet()
	token := jwt.CreateToken()
	token.Body = map[string]interface{}{
		"wow": "hello",
	}

	jwtStr := signer.Sign(token)
	t.Log("Got jwt:\n", jwtStr)

	_, err := verifier.Verify(jwtStr)

	if err == nil {
		t.Fatal("Incorrect signature did not throw")
	}
}

func TestErrorIncorrectJwtFormat(t *testing.T) {
	_, verifier, _ := createSet()
	jwtStr := "incorrectJWT"

	_, err := verifier.Verify(jwtStr)

	if err == nil {
		t.Fatal("Incorrect JWT format did not throw")
	}
}

func TestExpiredToken(t *testing.T) {
	signer, verifier, _ := createSet()
	token := jwt.CreateToken()
	token.Body = map[string]interface{}{
		"exp": time.Now().Unix() - 10,
	}

	jwtStr := signer.Sign(token)
	log.Println("Got jwt:\n", jwtStr)

	_, err := verifier.Verify(jwtStr)

	if err == nil {
		t.Fatal("Expired JWT was verified")
	}
}
