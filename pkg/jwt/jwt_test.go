package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	_ "crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

func TestSign(t *testing.T) {
	// Gen keys
	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	pubKey := &privKey.PublicKey
	pemKey := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(pubKey),
	})
	// pemKeyB64 := base64.StdEncoding.EncodeToString(pemKey)
	t.Log("Got pubkey:\n", string(pemKey))

	jwt := &JWT{
		algorithm: "RS256",
		signKey:   privKey,
		verifyKey: pubKey,
	}

	token := jwt.CreateToken()
	token.body = map[string]interface{}{
		"wow": "hello",
	}

	jwtStr := jwt.Sign(token)
	t.Log("Got jwt: ", jwtStr)

}
