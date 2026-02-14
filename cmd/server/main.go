package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"slices"

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

		publishedContent, err := readers.ReadPublished()
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

func archivesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read what's currently published. We load each line
		// of 'published' into a slice of [types.PublishData],
		// which is then tossed into the archives.gohtml
		// template.
		f, err := os.Open("published")
		if err != nil {
			log.Printf("error opening 'published' file: %v", err)
			http.Error(w, "it looks like nothing is published yet.", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		var archiveData []types.PublishData

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			var data types.PublishData
			entry := scanner.Text()

			if err := json.Unmarshal([]byte(entry), &data); err != nil {
				log.Printf("error unmarshaling entry: %v", err)
				http.Error(w, "oops, something happened", http.StatusInternalServerError)
				return
			}

			archiveData = append(archiveData, data)
		}
		if err := scanner.Err(); err != nil {
			log.Printf("error while scanning 'published' file")
			http.Error(w, "oops, something happened", http.StatusInternalServerError)
			return
		}

		// Load the template.
		tpl, err := template.ParseFiles("gohtml/archives.gohtml", "html/nav.html")
		if err != nil {
			log.Printf("error parsing template: %v", err)
			http.Error(w, "error parsing template", http.StatusInternalServerError)
			return
		}

		// Note that published items are appended to the
		// 'published' file, meaning that the file is in
		// chronological order. But we want to list archive
		// entries in reverse chronological order, as is the
		// custom.
		slices.Reverse(archiveData)

		if err := tpl.Execute(w, archiveData); err != nil {
			log.Printf("error executing template: %v", err)
			http.Error(w, "error executing template", http.StatusInternalServerError)
			return
		}
	}
}
