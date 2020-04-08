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
	"strings"
	"time"

	// Required for SHA256 hashing
	_ "crypto/sha256"
)

// Signer is used to sign tokens
type Signer struct {
	PrivateKey *rsa.PrivateKey
}

// Verifier is used to verify token
type Verifier struct {
	PublicKey *rsa.PublicKey
}

// Token represents a JWT
type Token struct {
	header    map[string]string
	Body      map[string]interface{}
	signature []uint8
}

// NewSigner ...
func NewSigner(pemData []byte) (*Signer, error) {
	b, _ := pem.Decode(pemData)

	if b.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("Verifier requires an RSA private key")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(b.Bytes)
	if err != nil {
		return nil, err
	}

	return &Signer{
		PrivateKey: privKey,
	}, nil
}

// NewVerifier ...
func NewVerifier(pemData []byte) (*Verifier, error) {
	b, _ := pem.Decode(pemData)
	if b.Type == "RSA PRIVATE KEY" {
		return newVerifierFromPrivatekey(b)
	} else if b.Type == "RSA PUBLIC KEY" {
		return newVerifierFromPublicKey(b)
	}
	return nil, errors.New("Verifier requires an RSA key")
}

func newVerifierFromPrivatekey(b *pem.Block) (*Verifier, error) {
	privKey, err := x509.ParsePKCS1PrivateKey(b.Bytes)
	if err != nil {
		return nil, err
	}
	return &Verifier{
		PublicKey: &privKey.PublicKey,
	}, nil
}

func newVerifierFromPublicKey(b *pem.Block) (*Verifier, error) {
	pubKey, err := x509.ParsePKCS1PublicKey(b.Bytes)
	if err != nil {
		return nil, err
	}
	return &Verifier{
		PublicKey: pubKey,
	}, nil
}

// CreateToken prepares a new Token
func CreateToken() *Token {
	return &Token{
		header: map[string]string{
			"typ": "jwt",
			"alg": "RS256",
		},
		Body:      map[string]interface{}{},
		signature: []byte{},
	}
}

// Decode will decode the payload without verifying the signature
func (v *Verifier) Decode(jwtStr string) (*Token, error) {
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

// Verify decodes a token and checks whether the given signature is valid
func (v *Verifier) Verify(jwtStr string) (*Token, error) {
	var err error

	token, err := v.Decode(jwtStr)
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
	err = rsa.VerifyPKCS1v15(v.PublicKey, crypto.SHA256, hashed, token.signature)

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
func (s *Signer) Sign(t *Token) string {
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
	sig, _ := rsa.SignPKCS1v15(rand.Reader, s.PrivateKey, crypto.SHA256, hashed)
	sigB64 := base64.RawURLEncoding.EncodeToString(sig)

	return headerB64 + "." + bodyB64 + "." + sigB64
}
