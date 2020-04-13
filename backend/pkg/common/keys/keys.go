package keys

import (
	"encoding/json"
	"net/http"

	"gopkg.in/square/go-jose.v2"
)

// FetchPublicKeys retrieves the public key from an URL offering a JWKS
func FetchPublicKeys(url string) (*jose.JSONWebKeySet, error) {
	var jwks jose.JSONWebKeySet

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	json.NewDecoder(resp.Body).Decode(&jwks)

	return &jwks, nil
}
