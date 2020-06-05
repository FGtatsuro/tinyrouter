package tinyrouter_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/FGtatsuro/tinyrouter"

	"net/http/httptest"
)

func TestNew(t *testing.T) {
	var router *tinyrouter.Router = tinyrouter.New()
	if router == nil {
		t.Errorf("Router must be created")
	}
}

var paths = []string{
	"/",
	"//",
	"/next",
	"/next/follow",
}

func TestHandle(t *testing.T) {
	wants := []string{
		"root",
		"root_multiple",
		"next",
		"next/follow",
	}
	for i, path := range paths {
		t.Run(path, func(t *testing.T) {
			router := tinyrouter.New()
			router.Handle(
				path,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(wants[i]))
				}))

			s := httptest.NewServer(router)
			defer s.Close()
			resp, _ := s.Client().Get(s.URL + path)

			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			got := string(body)
			if wants[i] != got {
				t.Errorf("Handler binding to '%v' must be called: want: %v/got %v", path, wants[i], got)
			}
		})
	}
}

func TestHandleFunc(t *testing.T) {
	wants := []string{
		"root",
		"root_multiple",
		"next",
		"next/follow",
	}
	for i, path := range paths {
		t.Run(path, func(t *testing.T) {
			router := tinyrouter.New()
			router.HandleFunc(
				path,
				func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(wants[i]))
				})

			s := httptest.NewServer(router)
			defer s.Close()
			resp, _ := s.Client().Get(s.URL + path)

			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			got := string(body)
			if wants[i] != got {
				t.Errorf("Handler binding to '%v' must be called: want: %v/got %v", path, wants[i], got)
			}
		})
	}

}
