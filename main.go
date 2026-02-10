package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /posts/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		fmt.Fprintf(w, "Post: %s", slug)
	})

	if err := http.ListenAndServe(":3030", mux); err != nil {
		log.Fatal(err)
	}
}
