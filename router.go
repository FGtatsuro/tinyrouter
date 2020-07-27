package tinyrouter

import (
	"net/http"
	"path"
	"regexp"
	"strings"
)

type Router struct {
	root *node
}

type node struct {
	handler           *http.Handler
	exp               *regexp.Regexp
	static            string
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
	// TODO: More precise way to find regex
	// TODO: Now, repetitions like {2} can't be used
	current := router.root
	for _, segment := range strings.Split(path.Clean(pattern), "/")[1:] {
		if strings.Contains(segment, "{") && strings.Contains(segment, "}") {
			segment = strings.ReplaceAll(segment, "{", "")
			segment = strings.ReplaceAll(segment, "}", "")
			if n, ok := current.regexpChildren[segment]; ok {
				current = n
			} else {
				// TODO: Error handling
				// Include interface consideration. Return error or not?
				exp, _ := regexp.Compile(segment)
				newNode := &node{
					exp:               exp,
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
					// TODO: No need?
					static:            segment,
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
		// TODO: Handle the case path doesn't exist in handler map.
		// TODO: return error
		return
	}
	(*current.handler).ServeHTTP(w, r)
}
