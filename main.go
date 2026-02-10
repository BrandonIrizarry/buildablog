package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

// markdownFileHandler returns the blog text found at the given slug
// path.
func markdownFileHandler(slug string) (string, error) {
	filename := fmt.Sprintf("posts/%s.md", slug)
	log.Printf("Will attempt to open %s", filename)

	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
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

		fmt.Fprint(w, blogText)
	}
}
