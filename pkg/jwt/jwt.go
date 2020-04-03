package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"

	// Required for SHA256 hashing
	_ "crypto/sha256"
)

// JWT holds configuration for signing and verifying
type JWT struct {
	algorithm string
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
}

// Token represents a JWT
type Token struct {
	header    map[string]string
	body      map[string]interface{}
	signature []uint8
}

// CreateToken prepares a new Token
func (jwt *JWT) CreateToken() *Token {
	return &Token{
		header: map[string]string{
			"typ": "jwt",
			"alg": jwt.algorithm,
		},
		body:      map[string]interface{}{},
		signature: []uint8{},
	}
}

// Decode will decode the payload without verifying the signature
func (jwt *JWT) Decode(jwtStr string) *Token {
	// To implement
	return nil
}

// Verify checks whether the given signature is valid
func (jwt *JWT) Verify(jwtStr string) *Token {
	token := jwt.Decode(jwtStr)
	return token
}

// Sign will add a signature to the given payload
func (jwt *JWT) Sign(t *Token) string {
	var r []byte
	// Encode header
	r, _ = json.Marshal(t.header)
	headerB64 := base64.RawURLEncoding.EncodeToString(r)

	// Encode body
	r, _ = json.Marshal(t.body)
	bodyB64 := base64.RawURLEncoding.EncodeToString(r)

	// Create signature
	pssh := crypto.SHA256.New()
	pssh.Write([]byte(headerB64 + "." + bodyB64))
	hashed := pssh.Sum(nil)
	sig, _ := rsa.SignPKCS1v15(rand.Reader, jwt.signKey, crypto.SHA256, hashed)
	sigB64 := base64.RawURLEncoding.EncodeToString(sig)

	return headerB64 + "." + bodyB64 + "." + sigB64
}
