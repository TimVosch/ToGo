package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	_ "crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"log"
	"testing"
	"time"
)

func createJWT() *JWT {
	// Gen keys
	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	pubKey := &privKey.PublicKey
	pemKey := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(pubKey),
	})
	log.Println("Got pubkey:\n", string(pemKey))

	jwt := &JWT{
		Algorithm:  "RS256",
		PrivateKey: privKey,
		PublicKey:  pubKey,
	}

	return jwt
}

func TestSignDoesVerify(t *testing.T) {
	jwt := createJWT()
	token := jwt.CreateToken()
	token.Body = map[string]interface{}{
		"wow": "hello",
	}

	jwtStr := jwt.Sign(token)
	t.Log("Got jwt:\n", jwtStr)

	_, err := jwt.Verify(jwtStr)

	if err != nil {
		t.Fatal("Failed to verify our own signed token")
	}
}

func TestErrorInvalidSignature(t *testing.T) {
	jwt1 := createJWT()
	jwt2 := createJWT()
	token := jwt1.CreateToken()
	token.Body = map[string]interface{}{
		"wow": "hello",
	}

	jwtStr := jwt1.Sign(token)
	t.Log("Got jwt:\n", jwtStr)

	_, err := jwt2.Verify(jwtStr)

	if err == nil {
		t.Fatal("Incorrect signature did not throw")
	}
}

func TestErrorIncorrectJwtFormat(t *testing.T) {
	jwt := createJWT()
	jwtStr := "incorrectJWT"

	_, err := jwt.Verify(jwtStr)

	if err == nil {
		t.Fatal("Incorrect JWT format did not throw")
	}
}

func TestExpiredToken(t *testing.T) {
	jwt := createJWT()
	token := jwt.CreateToken()
	token.Body = map[string]interface{}{
		"exp": time.Now().Unix() - 10,
	}

	jwtStr := jwt.Sign(token)
	log.Println("Got jwt:\n", jwtStr)

	_, err := jwt.Verify(jwtStr)

	if err == nil {
		t.Fatal("Expired JWT was verified")
	}
}
