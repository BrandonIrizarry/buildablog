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

	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatal(err)
	}

	// Set up the server.
	mux := http.NewServeMux()

	// Serve a blog post.
	contentPattern := fmt.Sprintf("GET /%s/{date}", constants.PostsLabel)
	mux.HandleFunc(contentPattern, func(w http.ResponseWriter, r *http.Request) {
		date := r.PathValue("date")
		fmdata, content, err := readers.ReadPost(constants.PostsLabel, date)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		payload := struct {
			Title   string
			Content template.HTML
		}{
			Title:   fmdata.Title,
			Content: content,
		}

		if err := feedTemplate(w, constants.PostsLabel, payload); err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Serve the site's front page.
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmdata, content, err := readers.ReadPost(constants.IndexLabel, "index.md")
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Serve the archives.
	mux.HandleFunc("GET /archives", func(w http.ResponseWriter, r *http.Request) {
		// Thanks to the magic of symbolic links, posts can
		// have local human-readable names (something I'm
		// adamant about), and also possess a straightforward
		// and (sensibly) immutable slug for publication
		// purposes.
		postDirEntries, err := os.ReadDir("content/" + constants.PostsLabel)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		payload := struct {
			Tag   string
			Posts []types.FrontmatterData
		}{
			Tag:   r.FormValue("tag"),
			Posts: make([]types.FrontmatterData, 0),
		}

		for _, p := range postDirEntries {
			filename := p.Name()

			// This has to do with the convention we use
			// for naming published posts.
			filenameDate, err := time.ParseInLocation(time.DateOnly, filename, location)
			if err != nil {
				// Post doesn't count as published, so skip.
				continue
			}

			// Read the post's frontmatter.
			fmdata, _, err := readers.ReadPost(constants.PostsLabel, filename)
			if err != nil {
				log.Printf("%v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//  These should naturally correspond. In the
			//  future there will be a mechanism to
			//  automatically generate the needed symbolic
			//  links.
			if !fmdata.Date.Equal(filenameDate) {
				err := fmt.Errorf("Filename %s doesn't match frontmatter date %s", filenameDate, fmdata.Date)
				log.Printf("%v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			payload.Posts = append(payload.Posts, fmdata)
		}

		if err := feedTemplate(w, "archives", payload); err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	//  Serve the tags page.
	mux.HandleFunc("GET /tags", func(w http.ResponseWriter, r *http.Request) {
		postEntries, err := os.ReadDir("content/" + constants.PostsLabel)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Accumulate all tags into a list. Naturally there
		// are duplicates, and so we remove these later using
		// a set.
		var tagList []string
		for _, p := range postEntries {
			filename := p.Name()

			filenameDate, err := time.ParseInLocation(time.DateOnly, filename, location)
			if err != nil {
				// Post doesn't count as published, so skip.
				continue
			}

			fmdata, _, err := readers.ReadPost(constants.PostsLabel, filename)
			if err != nil {
				log.Printf("%v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//  These should naturally correspond. In the
			//  future there will be a mechanism to
			//  automatically generate the needed symbolic
			//  links.
			if !fmdata.Date.Equal(filenameDate) {
				err := fmt.Errorf("Filename %s doesn't match frontmatter date %s", filenameDate, fmdata.Date)
				log.Printf("%v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tagList = append(tagList, fmdata.Tags...)
		}

		// Here we remove duplicates, as promised. This is
		// what gets fed into the appropriate template.
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
