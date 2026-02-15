package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/BrandonIrizarry/buildablog/internal/constants"
	"github.com/BrandonIrizarry/buildablog/internal/readers"
	"github.com/BrandonIrizarry/buildablog/internal/types"
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

	// Serve a blog post.
	contentPattern := fmt.Sprintf("GET /%s/{slug}", constants.PostsLabel)
	mux.HandleFunc(contentPattern, func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		templateContent, err := readers.ReadPage(slug, "posts")
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Error loading post", http.StatusNotFound)
			return
		}

		if err := feedTemplate(w, "posts", templateContent); err != nil {
			log.Printf("%v", err)
			http.Error(w, "Templating error", http.StatusNotFound)
			return
		}
	})

	// Serve the site's front page.
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		templateContent, err := readers.ReadPage("index", "index")
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Error loading post", http.StatusNotFound)
			return
		}

		err = feedTemplate(w, "index", templateContent)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Templating error", http.StatusNotFound)
			return
		}
	})

	// Serve the archives page.
	mux.HandleFunc("GET /archives", func(w http.ResponseWriter, r *http.Request) {
		templateContent, err := readers.ReadPage("index", "archives")
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Error loading post", http.StatusNotFound)
			return
		}

		// FIXME: make a separate function out of this, to
		// simplify error handling.
		f, err := os.Open("published.json")
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Error loading post", http.StatusNotFound)
			return
		}
		defer f.Close()

		var publishedContent []types.PublishData
		rawJSON, err := io.ReadAll(f)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Error loading post", http.StatusNotFound)
			return
		}

		json.Unmarshal(rawJSON, &publishedContent)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Error loading post", http.StatusNotFound)
			return
		}

		err = feedTemplate(w, "archives", struct {
			Main      types.PostData
			Published []types.PublishData
		}{
			Main:      templateContent,
			Published: publishedContent,
		})

		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Templating error", http.StatusNotFound)
			return
		}
	})

	// Serve the site's static assets (CSS files etc.)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	if err := http.ListenAndServe(":3030", mux); err != nil {
		log.Fatal(err)
	}
}

func feedTemplate(w http.ResponseWriter, label string, data any) error {
	// Load the template.
	tpl, err := template.ParseFiles("gohtml/"+label+".gohtml", "html/nav.html")
	if err != nil {
		return err
	}

	// Use the template.
	if err := tpl.Execute(w, data); err != nil {
		return err
	}

	return nil
}
