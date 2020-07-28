package tinyrouter

import (
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strings"
)

type Router struct {
	root *node
}

type node struct {
	handler        *http.Handler
	exp            *regexp.Regexp
	staticChildren map[string]*node
	regexpChildren map[string]*node
}

func New() *Router {
	return &Router{
		&node{
			staticChildren: map[string]*node{},
			regexpChildren: map[string]*node{},
		},
	}
}

func (router *Router) Handle(pattern string, handler http.Handler) {
	current := router.root
	for _, segment := range strings.Split(path.Clean(pattern), "/")[1:] {
		if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
			// TODO: Now, we can't support path including Japanese
			segment = fmt.Sprintf("^%v$", segment[1:len(segment)-1])
			if n, ok := current.regexpChildren[segment]; ok {
				current = n
			} else {
				// TODO: Error handling
				// Include interface consideration. Return error or not?
				exp, _ := regexp.Compile(segment)
				newNode := &node{
					exp:            exp,
					staticChildren: map[string]*node{},
					regexpChildren: map[string]*node{},
				}
				current.regexpChildren[segment] = newNode
				current = newNode
			}
		} else {
			if n, ok := current.staticChildren[segment]; ok {
				current = n
			} else {
				newNode := &node{
					staticChildren: map[string]*node{},
					regexpChildren: map[string]*node{},
				}
				current.staticChildren[segment] = newNode
				current = newNode
			}
		}
	}
	current.handler = &handler
}

func (router *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	router.Handle(pattern, http.HandlerFunc(handler))
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Search order: static -> regexp
	current := router.root
Loop:
	for _, segment := range strings.Split(path.Clean(r.URL.Path), "/")[1:] {
		if n, ok := current.staticChildren[segment]; ok {
			current = n
			continue Loop
		}
		// TODO: Now, search order of regexp isn't fixed
		// TODO: This action is O(N)(N=regex segment num) order
		for _, n := range current.regexpChildren {
			if n.exp.Match([]byte(segment)) {
				current = n
				continue Loop
			}
		}

		http.NotFound(w, r)
		return
	}
	(*current.handler).ServeHTTP(w, r)
}
