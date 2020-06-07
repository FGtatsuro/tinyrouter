package tinyrouter

import (
	"net/http"
	"path"
	"regexp"
	"strings"
)

type Router struct {
	handlerMap       map[string]*http.Handler
	regexpHandlerMap map[string]*http.Handler
}

func New() *Router {
	return &Router{
		make(map[string]*http.Handler),
		make(map[string]*http.Handler),
	}
}

func (router *Router) Handle(pattern string, handler http.Handler) {
	// TODO: More precise way to find regex
	// TODO: Now, repetitions like {2} can't be used
	if strings.Contains(pattern, "{") && strings.Contains(pattern, "}") {
		pattern = strings.ReplaceAll(pattern, "{", "")
		pattern = strings.ReplaceAll(pattern, "}", "")
		router.regexpHandlerMap[path.Clean(pattern)] = &handler
	} else {
		router.handlerMap[path.Clean(pattern)] = &handler
	}
}

func (router *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	router.Handle(pattern, http.HandlerFunc(handler))
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Search order: regexp -> normal
	// TODO: Now, search order of regexp isn't fixed
	// TODO: Compile regexp only registration time
	for pattern, handler := range router.regexpHandlerMap {
		re, _ := regexp.Compile(pattern)
		if re.Match([]byte(path.Clean(r.URL.Path))) {
			(*handler).ServeHTTP(w, r)
			return
		}
	}
	// TODO: Handle the case path doesn't exist in handler map.
	(*router.handlerMap[path.Clean(r.URL.Path)]).ServeHTTP(w, r)
}
