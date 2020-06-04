package tinyrouter

import (
	"net/http"
	"path"
)

type Router struct {
	handlerMap map[string]*http.Handler
}

func New() *Router {
	return &Router{make(map[string]*http.Handler)}
}

func (router *Router) Handle(pattern string, handler http.Handler) {
	router.handlerMap[path.Clean(pattern)] = &handler
}

func (router *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	router.Handle(pattern, http.HandlerFunc(handler))
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	(*router.handlerMap[path.Clean(r.URL.Path)]).ServeHTTP(w, r)
}
