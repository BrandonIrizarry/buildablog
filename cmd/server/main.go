package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

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
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Set up the server.
	mux := http.NewServeMux()

	// Serve a blog post.
	contentPattern := fmt.Sprintf("GET /%s/{slug}", constants.PostsLabel)
	mux.HandleFunc(contentPattern, func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		templateContent, err := readers.ReadPage(slug, "posts")
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if err := feedTemplate(w, "posts", templateContent); err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
	})

	// Serve the site's front page.
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		templateContent, err := readers.ReadPage("index", "index")
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		err = feedTemplate(w, "index", templateContent)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
	})

	// Serve the archives page.
	mux.HandleFunc("GET /archives", func(w http.ResponseWriter, r *http.Request) {
		templateContent, err := readers.ReadPage("index", "archives")
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		rawJSON, err := readPublishedJSON("published.json")
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		type publishDataFormatted struct {
			types.PublishData
			CreatedHumanReadable string
			UpdatedHumanReadable string
			DidUpdate            bool
		}

		var publishedContent []publishDataFormatted

		if err := json.Unmarshal(rawJSON, &publishedContent); err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Prepare the human-readable time formats for display
		// on the archives page.
		const humanReadableFormat = "2006-1-2 (3:4 PM)"

		for i := range publishedContent {
			pc := &publishedContent[i]
			created := pc.Created
			updated := pc.Updated

			(*pc).CreatedHumanReadable = time.Unix(created, 0).Format(humanReadableFormat)
			(*pc).UpdatedHumanReadable = time.Unix(updated, 0).Format(humanReadableFormat)
			(*pc).DidUpdate = (updated > created)
		}

		err = feedTemplate(w, "archives", struct {
			Main      types.PostData
			Published []publishDataFormatted
		}{
			Main:      templateContent,
			Published: publishedContent,
		})

		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
	})

	// Serve the site's static assets (CSS files etc.)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	if err := http.ListenAndServe(":3030", mux); err != nil {
		log.Fatal(err)
	}
}

// feedTemplate loads and executes a template found under label using
// the given data parameter (used by [template.Template.Execute] to
// fill in the template.)
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

// readPublishedJSON returns the contents of filename as a byte slice
// (NOTE: such a slice should be marshallable as JSON.)
func readPublishedJSON(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("can't open %s: %w", filename, err)
	}
	defer f.Close()

	rawJSON, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("can't read %s: %w", filename, err)
	}

	return rawJSON, nil
}
