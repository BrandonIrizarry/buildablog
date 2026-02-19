package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"slices"
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
		fmdata, content, err := readers.ReadPage(slug, "posts")
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		data := struct {
			Title   string
			Content template.HTML
		}{
			Title:   fmdata.Title,
			Content: content,
		}

		if err := feedTemplate(w, "posts", data); err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
	})

	// Serve the site's front page.
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmdata, content, err := readers.ReadPage("index", "index")
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		data := struct {
			Title   string
			Content template.HTML
		}{
			Title:   fmdata.Title,
			Content: content,
		}

		err = feedTemplate(w, "index", data)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
	})

	// Serve the archives page.
	mux.HandleFunc("GET /archives", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("GET /archives")

		pdataList, err := readers.ReadPublishingFile("published.json")
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Published []types.PublishData
			Tag       string
		}{
			Published: pdataList,
			Tag:       r.FormValue("tag"),
		}

		if err := feedTemplate(w, "archives", data); err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
	})

	mux.HandleFunc("GET /tags", func(w http.ResponseWriter, r *http.Request) {
		publishedList, err := readers.ReadPublishingFile("published.json")
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Since different blog posts will (as is desirable)
		// intersect across various tags, we need to
		// deduplicate them using a set, which is in turn fed
		// into the tags.gohtml template.
		tagSet := make(map[string]struct{})
		for _, p := range publishedList {
			for _, t := range p.Tags {
				tagSet[t] = struct{}{}
			}
		}

		log.Printf("Tag set: %v", tagSet)

		if err := feedTemplate(w, "tags", tagSet); err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
	templateName := label + ".gohtml"
	funcMap := template.FuncMap{
		"dec": func(value int) int {
			return value - 1
		},
		"humanReadable": func(timestamp int64) string {
			const humanReadableFormat = "2006-1-2 (3:04 PM)"
			return time.Unix(timestamp, 0).Format(humanReadableFormat)
		},
		"hasTag": func(tag string, tags []string) bool {
			return slices.Contains(tags, tag)
		},
	}

	tpl, err := template.New(templateName).Funcs(funcMap).ParseFiles("gohtml/"+templateName, "html/nav.html")
	if err != nil {
		return err
	}

	// Use the template.
	if err := tpl.Execute(w, data); err != nil {
		return err
	}

	return nil
}
