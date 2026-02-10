package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/yuin/goldmark"
)

func main() {
	// Set up logging.
	logFile, err := os.OpenFile("buildablog.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("can't open logfile: %v", err)
		os.Exit(1)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Set up the server.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /posts/{slug}", postHandler(markdownFileHandler))

	if err := http.ListenAndServe(":3030", mux); err != nil {
		log.Fatal(err)
	}
}

// postHandler "decorates" the given reader using an
// [http.HandlerFunc].
func postHandler(reader reader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Slug: %s", r.PathValue("slug"))

		slug := r.PathValue("slug")
		blogText, err := reader(slug)
		if err != nil {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		var buf bytes.Buffer
		if err := goldmark.Convert([]byte(blogText), &buf); err != nil {
			http.Error(w, "Error converting Markdown", http.StatusInternalServerError)
			return
		}

		io.Copy(w, &buf)
	}
}
