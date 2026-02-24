package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/BrandonIrizarry/buildablog/internal/constants"
	"github.com/BrandonIrizarry/buildablog/internal/readers"
	"github.com/BrandonIrizarry/buildablog/internal/types"
)

// tpls maps template basenames to actual templates. This is so that
// we can parse all our templates up front, as opposed to parsing them
// on each request.
var tpls = make(map[string]*template.Template)

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

	// Flags: if -flagLocal is specified, set RSS siteURL to
	// localhost:PORT.
	var flagLocal bool
	flag.BoolVar(&flagLocal, "local", false, "Whether we're serving from localhost")
	flag.Parse()

	// Parse the templates up front.
	funcMap := template.FuncMap{
		"dec": func(value int) int {
			return value - 1
		},
		"humanReadable": func(t time.Time) string {
			return t.Format(time.DateOnly)
		},
		"hasTag": func(tag string, tags []string) bool {
			return slices.Contains(tags, tag)
		},
	}

	gohtmlFiles, err := filepath.Glob("gohtml/*")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range gohtmlFiles {
		log.Printf("Parsing template file '%s'", file)

		name := filepath.Base(file)
		tpl, err := template.New(name).Funcs(funcMap).ParseFiles(file, "html/nav.html")
		if err != nil {
			log.Fatal(err)
		}

		tpls[name] = tpl
	}

	// Set up the server.
	mux := http.NewServeMux()

	// Serve a blog post.
	contentPattern := fmt.Sprintf("GET /%s/{date}", constants.PostsLabel)
	mux.HandleFunc(contentPattern, func(w http.ResponseWriter, r *http.Request) {
		date := r.PathValue("date")
		postData, err := readers.ReadMarkdown(constants.PostsLabel, date)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := feedTemplate(w, constants.PostsLabel, postData); err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Serve the site's front page.
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		postData, err := readers.ReadMarkdown(constants.IndexLabel, "index.md")
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		err = feedTemplate(w, "index", postData)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Serve the archives.
	mux.HandleFunc("GET /archives", func(w http.ResponseWriter, r *http.Request) {
		posts, err := readers.AllPosts()
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// The nice thing is that, because of the file naming
		// convention, posts are already sorted on the
		// filesystem. However, for display in Archives, the
		// most recent post should come first.
		slices.Reverse(posts)

		payload := struct {
			Posts []types.PostData
			Tag   string
		}{
			Posts: posts,
			Tag:   r.FormValue("tag"),
		}

		if err := feedTemplate(w, "archives", payload); err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	//  Serve the tags page.
	mux.HandleFunc("GET /tags", func(w http.ResponseWriter, r *http.Request) {
		posts, err := readers.AllPosts()
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Dump all tags into a list, then deduplicate using a
		// set.
		var tagList []string
		for _, p := range posts {
			tagList = append(tagList, p.Tags...)
		}

		tagSet := make(map[string]struct{})
		for _, t := range tagList {
			tagSet[t] = struct{}{}
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

	log.Print("Killing any previous server instance; starting server on port 3030")

	if err := http.ListenAndServe(":3030", mux); err != nil {
		log.Fatal(err)
	}
}

// feedTemplate loads and executes a template found under label using
// the given data parameter (used by [template.Template.Execute] to
// fill in the template.)
func feedTemplate(w http.ResponseWriter, label string, data any) error {
	// Load the template.
	t, ok := tpls[label+".gohtml"]
	if !ok {
		return fmt.Errorf("no template under label '%s'", label)
	}

	// Use the template.
	if err := t.Execute(w, data); err != nil {
		return fmt.Errorf("can't execute template under label '%s': %w", label, err)
	}

	return nil
}
