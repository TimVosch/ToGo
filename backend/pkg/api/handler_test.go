package api_test

import (
	"log"
	"testing"

	"github.com/timvosch/togo/pkg/api"
)

func TestHandlerShouldChainNext(t *testing.T) {
	called := map[int]bool{
		0: false,
		1: false,
		2: false,
	}
	// Every item in the chain sets
	// a `called` field to true
	httpHandler := api.Handler(
		func(ctx *api.CTX, next func()) {
			log.Println("I am the first handler")
			called[0] = true
			next()
		},
		func(ctx *api.CTX, next func()) {
			log.Println("I am the second handler")
			called[1] = true
			next()
		},
		func(ctx *api.CTX, next func()) {
			log.Println("I am the third handler")
			called[2] = true
		},
	)

	httpHandler(nil, nil)

	for _, v := range called {
		if v == false {
			t.Fatal("Chain did not call each handler")
		}
	}
}

func TestSilentlyHandleEndOfChain(t *testing.T) {
	httpHandler := api.Handler(
		func(ctx *api.CTX, next func()) {
			next()
			next()
			next()
		},
	)

	httpHandler(nil, nil)
}
