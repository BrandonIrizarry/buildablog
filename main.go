package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/yuin/goldmark"
	hl "github.com/yuin/goldmark-highlighting/v2"
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

		// Enable syntax highlighting in blog posts.
		mdRenderer := goldmark.New(
			goldmark.WithExtensions(hl.NewHighlighting(
				hl.WithStyle("dracula"),
			)),
		)

		// Render Markdown as HTML.
		var buf bytes.Buffer
		if err := mdRenderer.Convert([]byte(blogText), &buf); err != nil {
			http.Error(w, "Error converting Markdown", http.StatusInternalServerError)
			return
		}

		// Use the template.
		tpl, err := template.ParseFiles("post.gohtml")
		if err != nil {
			log.Printf("error parsing template: %v", err)
			http.Error(w, "error parsing template", http.StatusInternalServerError)
			return
		}

		err = tpl.Execute(w, postData{
			Title:   "My First Post",
			Content: template.HTML(buf.String()),
			Author:  "Brandon Irizarry",
		})
		if err != nil {
			log.Printf("error executing template: %v", err)
			http.Error(w, "error executing template", http.StatusInternalServerError)
			return
		}
	}
}
