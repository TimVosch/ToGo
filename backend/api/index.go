package api

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"log"
	"net/http"
	"os"

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

// Handler is the handler for the Now serverless function
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling url: ", r.URL.Path)
	// Set up
	pemB64, _ := base64.StdEncoding.DecodeString("LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBMksvY3VUYzlRMTladmZuOU9zdnpzMTFSOUVPWjlrbVJrZXdpSDQ2dElMakJlaTNSClBsZ016UVdLYkh6RzJGNzFpaUx0RnZKSFk4MnQwZlhsTTVjaFIxNWtQd1JodHNFMU9CUHdLa0M0VWhPSVBrV0kKSEd1VGdrc2RteEZWbUNLMXlsZDZBbFUyVDk5VnVkUE9GS0ZwN2szMUhJSmI1cUxvaVQ2VDNUN09HQWFTeWZjSQpwbkl1WW1LTE04S2ZTWlVOV2FCSUxvRGV6V1poTzFHNktVWStZNzd2VThpMU9JT2ZlWjVEL1FQdSsxOEluUzQ3CkJKL0xkeDFXMVVhS1BJY1pFMklzcURUWjYwNVRaNmdVYzNFQlVRVXE1c0VyU0xEQkpXeHRLK3YrTElYK3Z0cloKcEZWOXJzTjQybFlSTS8vVlErUTZXb2F0ZXl2Rm5HSmNXdElQTHdJREFRQUJBb0lCQVFEVEo4ejMwSlNxcXFoUgpNT05NQUtPakRqVm15dG1sMzFzenorQVEwSUIxZXBWUDhvWU5NdENHbWZlaWNKVjFGRlJDSUhiWi9ZOEQrdEovCjRCZFNodHV2S0pTWHREVmtXakw2U3JPbStScWxJTk9MbTBaZ0s1UzdTMmUyVE5ZVVF2N1VCeHFtVzFOcDBrRS8Kck44Tnk1M20wNkVmL3doL1lCRXFiUWk2ZVJGczVtL3BTTCtQUXVObGxxVHdOMUNpelpnamh1K1J2R2F5ZHkxbApHcDIzTmdONnJIVTU2UjVxWjgrZ1Z2Vy96d1RmR2d0KzJJR3gxSFlxWjAxZnR0ZWoyK29oMkhtMmxiRWIvZXlPCnhzMW1lS1lXRVJIT0dHcE9jVXh0dTZUUTBnYmwzN1BOalpNdUhFbWtzMkxNcVJ5aU0vMkNEUG13NVZKSTNWRHkKOUgvbmpJbkJBb0dCQU8xVjVNVGJyRjl3aXh3Y3JoVlgvbzdKMDVobDVLSUd3c0lJS0trTFpsTjNybWdlNkF1SQpBNnZubGhnNXVGb3g3QVAxUjFoYU1mYUEvNTNYM09RU3ZhSndQd3dBTTl4V2ZabW1rdXVoNUJPdk52d295cHBLCjlDMFZHWkgxclVVMlpXYnFoUFRQM3E2SEl6WkIya2VyTnNoN2hPaVBqeWtFUktHQldoamVUTytsQW9HQkFPbTYKUkJBeHhBaXlrQnF3L1lLd0ZHb0ltcDJwbGFFdmt4dDZrMzFmRmhvZnFaYjltVWNoUWk0WUhVKy9qbVIzdm5iUQpuTVlWYW03RHgzNEZRS0x5RWdKMm5iM3NGa0RCMy96aTFrSFBnWWFhTFROMW04ZWp6WmhWTEhYa0J1U2JTS3ZxClE3MklIS3FvUTJTT1BQWG5RVzlXQWVQVzl0MERDdE5iWHdkUGNrdERBb0dBVnp0eTBraExtdWlxdUxKeDZiWm0KQWVWOHVFNzdNZko2TXdiOHF1VmR2dUFHWW82NWkwTjYxZnhRMXFhZ2M4WlZrVDdkOGtOMGliM3dOZnZaWEpybQp5SVdwSnFnTVo3Z0NnaThQWVR3bnNIUitLVUIwOXpFRmZteDY4WUx6SkxWUm4vb2kxRGh6Q0lMekZrWXVESm1KCmtUYVZLMFZZd1NLb2R4UXNJV2ZUcjJFQ2dZQW40N3RVRERwSnhiZmtaa3FOOEdFN2k0WmYzQjZHYU9reGFtVWIKbzR2Ukg1QkJEYjBJTDd2c3cvN1VxbnV4MStId3d6L09hcjlFY2pOczVaYVhlTHJzSXJSZlFwaTFxcUVBdHZJYwpQejc2NnZ0RjZnK1JMZnFid2dXWmhUWkw4OWllUnBnVEU5VFlwMmtCRTJtQ0NsclhscFV4L25FWlhUaU93K3hmCkFnY2Vvd0tCZ1FEUlpNSjZDc21JQ0JsSEFKamI0WDU2OCtna3g3UWtDRWhpaUJkRVl4K2JDN2JQUVFDc2E4MnUKd00wbE9xUlVHbUhqNDB4cmxlTDRPWkFUSDYxNElnOUovZXJaNUxOWG40amtLd0YzK0NFb1hCeFo5SjBUcW9ybQpKL1pNVzJycDhzSjc1ZVNvR1cyZWFubDVwSUtmMEQvUFJXcVhFSGtZTS90M3BKZGZtR0hpbnc9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=")
	b, _ := pem.Decode(pemB64)
	key, err := x509.ParsePKCS1PrivateKey(b.Bytes)
	if err != nil {
		log.Fatalln("Could not read private key!")
	}

	us := userserver.NewServerless(key, os.Getenv("MONGO_URI"), "/api/")
	us.Router.ServeHTTP(w, r)
}
