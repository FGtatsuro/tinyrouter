package tinyrouter_test

import (
	"net/http"
	"testing"

	"github.com/FGtatsuro/tinyrouter"
)

func TestNew(t *testing.T) {
	var router *tinyrouter.Router = tinyrouter.New()
	if router == nil {
		t.Errorf("Router must be created")
	}
}

func TestHandle(t *testing.T) {
	router := tinyrouter.New()
	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
}

func TestHandleFunc(t *testing.T) {
	router := tinyrouter.New()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	})
}
