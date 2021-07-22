module github.com/timvosch/togo

go 1.13

require (
	github.com/gorilla/mux v1.7.4
	go.mongodb.org/mongo-driver v1.5.1
	golang.org/x/crypto v0.0.0-20200406173513-056763e48d71
	gopkg.in/square/go-jose.v2 v2.5.0
)

replace github.com/timvosch/togo => ./
