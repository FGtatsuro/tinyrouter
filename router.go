package tinyrouter

import (
	"net/http"
)

type Router struct {
}

func New() *Router {
	return &Router{}
}

func (r *Router) Handle(pattern string, handler http.Handler) {
}

func (r *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
}
