# tinyrouter

Minimum HTTP request router

## Usage

```go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/FGtatsuro/tinyrouter"
)

func main() {
	r := tinyrouter.New()
	// Same signature to http.Handle
	r.Handle("/handle", HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "handle\n")
	}))
	// Same signature to http.HandleFunc
	r.HandleFunc("/handlefunc", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "handlefunc\n")
	})
	r.HandleFunc("/group/{groupid}/users/{userid:[0-9a-zA-Z]}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		pathVars := ctx.Value(tinyrouter.PathVarsContextKey).(map[string]string)
		io.WriteString(w, fmt.Sprintf("group %v/user %v", pathVars["groupid"], pathVars["userid"]))
	})
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```
