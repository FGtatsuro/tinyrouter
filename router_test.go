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

type route struct {
	path  string
	write string
}

type testcase struct {
	path string
	want string
}

func TestHandle(t *testing.T) {
	routes := []route{
		{"/", "root"},
		{"//", "root_multiple"},
		{"/next", "next"},
		{"/next/follow", "next/follow"},
	}
	var testcases []testcase
	for _, route := range routes {
		testcases = append(testcases, testcase{route.path, route.write})
	}

	for i, route := range routes {
		// FYI: https://github.com/golang/go/wiki/CommonMistakes
		write := route.write
		t.Run(route.path, func(t *testing.T) {
			router := tinyrouter.New()
			router.Handle(
				route.path,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(write))
				}))

			s := httptest.NewServer(router)
			defer s.Close()
			resp, _ := s.Client().Get(s.URL + testcases[i].path)

			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			got := string(body)
			if testcases[i].want != got {
				t.Errorf("Handler binding to '%v' must be called: want: %v/got %v", testcases[i].path, testcases[i].want, got)
			}
		})
	}
}

func TestHandleFunc(t *testing.T) {
	routes := []route{
		{"/", "root"},
		{"//", "root_multiple"},
		{"/next", "next"},
		{"/next/follow", "next/follow"},
	}
	var testcases []testcase
	for _, route := range routes {
		testcases = append(testcases, testcase{route.path, route.write})
	}

	for i, route := range routes {
		// FYI: https://github.com/golang/go/wiki/CommonMistakes
		write := route.write
		t.Run(route.path, func(t *testing.T) {
			router := tinyrouter.New()
			router.HandleFunc(
				route.path,
				func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(write))
				})

			s := httptest.NewServer(router)
			defer s.Close()
			resp, _ := s.Client().Get(s.URL + testcases[i].path)

			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			got := string(body)
			if testcases[i].want != got {
				t.Errorf("Handler binding to '%v' must be called: want: %v/got %v", testcases[i].path, testcases[i].want, got)
			}
		})
	}
}

func TestMultipleRoutes(t *testing.T) {
	routes := []route{
		{"/", "root"},
		{"//", "root_multiple"},
		{"/next", "next"},
		{"/next/follow", "next/follow"},
	}
	testcases := []testcase{
		// Same path overwrites previous registered route.
		{"/", "root_multiple"},
		{"//", "root_multiple"},
		{"/next", "next"},
		{"/next/follow", "next/follow"},
	}

	router := tinyrouter.New()
	for _, route := range routes {
		// FYI: https://github.com/golang/go/wiki/CommonMistakes
		write := route.write
		router.HandleFunc(
			route.path,
			func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(write))
			})
	}
	for _, tc := range testcases {
		t.Run(tc.path, func(t *testing.T) {
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

func TestRegexpRoutes(t *testing.T) {
	routes := []route{
		{"/regex/{[0-9]+}", "regex"},
		{"/regex/{[a-z][A-Z]}", "smallbig"},
	}
	testcases := []testcase{
		{"/regex/12345", "regex"},
		{"/regex/aB", "smallbig"},
	}

	router := tinyrouter.New()
	for _, route := range routes {
		// FYI: https://github.com/golang/go/wiki/CommonMistakes
		write := route.write
		router.HandleFunc(
			route.path,
			func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(write))
			})
	}
	for _, tc := range testcases {
		t.Run(tc.path, func(t *testing.T) {
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
