# tinyrouter

Minimum HTTP request router

## Usage

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/FGtatsuro/tinyrouter"
)

func main() {
	router := tinyrouter.New()
	// Same signature to http.Handle
	router.Handle("/handle", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("handle\n"))
	}))
	// Same signature to http.HandleFunc
	router.HandleFunc("/handlefunc", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("handlefunc\n"))
	})
	// Get PathVars via context
	router.HandleFunc("/users/{[0-9a-zA-Z]+}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := ctx.Value(tinyrouter.PathVarsContextKey).([]string)
		w.Write([]byte(fmt.Sprintf("user %v\n", vars[0])))
	})
	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```
