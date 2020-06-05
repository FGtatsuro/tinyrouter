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
	writes := []string{
		"root",
		"root_multiple",
		"next",
		"next/follow",
	}
	wants := writes
	for i, path := range paths {
		// FYI: https://github.com/golang/go/wiki/CommonMistakes
		write := writes[i]
		t.Run(path, func(t *testing.T) {
			router := tinyrouter.New()
			router.Handle(
				path,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(write))
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
	writes := []string{
		"root",
		"root_multiple",
		"next",
		"next/follow",
	}
	wants := writes
	for i, path := range paths {
		// FYI: https://github.com/golang/go/wiki/CommonMistakes
		write := writes[i]
		t.Run(path, func(t *testing.T) {
			router := tinyrouter.New()
			router.HandleFunc(
				path,
				func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(write))
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

func TestMultipleRoutes(t *testing.T) {
	writes := []string{
		"root",
		"root_multiple",
		"next",
		"next/follow",
	}
	wants := []string{
		// Same path overwrites previous registered route.
		"root_multiple",
		"root_multiple",
		"next",
		"next/follow",
	}
	router := tinyrouter.New()
	for i, path := range paths {
		// FYI: https://github.com/golang/go/wiki/CommonMistakes
		write := writes[i]
		router.HandleFunc(
			path,
			func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(write))
			})
	}
	for i, path := range paths {
		t.Run(path, func(t *testing.T) {
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
