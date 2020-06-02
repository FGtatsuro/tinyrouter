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
	router := tinyrouter.New()
	// Same signature to http.Handle
	router.Handle("/handle", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "handle\n")
	}))
	// Same signature to http.HandleFunc
	router.HandleFunc("/handlefunc", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "handlefunc\n")
	})
	router.HandleFunc("/group/{groupid}/users/{userid:[0-9a-zA-Z]}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		pathVars := ctx.Value(tinyrouter.PathVarsContextKey).(map[string]string)
		io.WriteString(w, fmt.Sprintf("group %v/user %v", pathVars["groupid"], pathVars["userid"]))
	})
	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```
