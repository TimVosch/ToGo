package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
	"strings"
	"time"

	// Required for SHA256 hashing
	_ "crypto/sha256"
)

// JWT holds configuration for signing and verifying
type JWT struct {
	Algorithm  string
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// Token represents a JWT
type Token struct {
	header    map[string]string
	Body      map[string]interface{}
	signature []uint8
}

// NewJWT ...
func NewJWT(privKeyPath string) *JWT {
	var dat []byte
	dat, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		return nil
	}

	b, _ := pem.Decode(dat)

	privKey, err := x509.ParsePKCS1PrivateKey(b.Bytes)
	if err != nil {
		return nil
	}

	log.Println("Loaded private key")

	return &JWT{
		Algorithm:  "RS256",
		PrivateKey: privKey,
		PublicKey:  &privKey.PublicKey,
	}
}

// CreateToken prepares a new Token
func (jwt *JWT) CreateToken() *Token {
	return &Token{
		header: map[string]string{
			"typ": "jwt",
			"alg": jwt.Algorithm,
		},
		Body:      map[string]interface{}{},
		signature: []byte{},
	}
}

// Decode will decode the payload without verifying the signature
func (jwt *JWT) Decode(jwtStr string) (*Token, error) {
	var token Token

	strs := strings.Split(jwtStr, ".")
	if len(strs) != 3 {
		return nil, errors.New("Invalid JWT String")
	}

	headerRaw, _ := base64.RawURLEncoding.DecodeString(strs[0])
	bodyRaw, _ := base64.RawURLEncoding.DecodeString(strs[1])
	sig, _ := base64.RawURLEncoding.DecodeString(strs[2])

	json.Unmarshal(headerRaw, &token.header)
	json.Unmarshal(bodyRaw, &token.Body)
	token.signature = sig

	return &token, nil
}

// Verify checks whether the given signature is valid
func (jwt *JWT) Verify(jwtStr string) (*Token, error) {
	var err error

	token, err := jwt.Decode(jwtStr)
	if err != nil {
		return nil, err
	}

	// Extract headerB64 || . || payloadB64
	strs := strings.Split(jwtStr, ".")
	hashInput := strs[0] + "." + strs[1]

	// Verify signature
	sha256 := crypto.SHA256.New()
	sha256.Write([]byte(hashInput))
	hashed := sha256.Sum(nil)
	err = rsa.VerifyPKCS1v15(jwt.PublicKey, crypto.SHA256, hashed, token.signature)

	if err != nil {
		return nil, errors.New("Invalid JWT signature")
	}

	// Verify exp
	if exp, ok := token.Body["exp"].(float64); ok == true {
		if time.Now().Unix() >= int64(exp) {
			return nil, errors.New("JWT has expired")
		}
	}

	return token, nil
}

// Sign will add a signature to the given payload
func (jwt *JWT) Sign(t *Token) string {
	var r []byte
	// Encode header
	r, _ = json.Marshal(t.header)
	headerB64 := base64.RawURLEncoding.EncodeToString(r)

	// Encode body
	r, _ = json.Marshal(t.Body)
	bodyB64 := base64.RawURLEncoding.EncodeToString(r)

	// Create signature
	sha256 := crypto.SHA256.New()
	sha256.Write([]byte(headerB64 + "." + bodyB64))
	hashed := sha256.Sum(nil)
	sig, _ := rsa.SignPKCS1v15(rand.Reader, jwt.PrivateKey, crypto.SHA256, hashed)
	sigB64 := base64.RawURLEncoding.EncodeToString(sig)

	return headerB64 + "." + bodyB64 + "." + sigB64
}
