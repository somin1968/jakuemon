package main

import (
	"fmt"
	"google.golang.org/appengine"
	"net/http"
)

func main() {
	http.HandleFunc("/api/articles/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})
	appengine.Main()
}
