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

var testcases = []struct {
	path string
	want string
}{
	{"/", "root"},
	{"//", "root"},
	{"/next", "next"},
	{"/next/follow", "next/follow"},
}

func TestHandle(t *testing.T) {
	for _, tc := range testcases {
		t.Run(tc.path, func(t *testing.T) {
			router := tinyrouter.New()
			router.Handle(
				tc.path,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(tc.want))
				}))

			s := httptest.NewServer(router)
			defer s.Close()
			resp, _ := s.Client().Get(s.URL + tc.path)

			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			got := string(body)
			if tc.want != got {
				t.Errorf("Handler binding to '%v' must be called: want: %v/got %v", tc.path, tc.want, got)
			}
		})
	}
}

func TestHandleFunc(t *testing.T) {
	for _, tc := range testcases {
		t.Run(tc.path, func(t *testing.T) {
			router := tinyrouter.New()
			router.HandleFunc(
				tc.path,
				func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(tc.want))
				})

			s := httptest.NewServer(router)
			defer s.Close()
			resp, _ := s.Client().Get(s.URL + tc.path)

			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			got := string(body)
			if tc.want != got {
				t.Errorf("Handler binding to '%v' must be called: want: %v/got %v", tc.path, tc.want, got)
			}
		})
	}
}
