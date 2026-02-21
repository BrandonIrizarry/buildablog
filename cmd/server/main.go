package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
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
		"humanReadable": func(timestamp int64) string {
			const humanReadableFormat = "2006-1-2 (3:04 PM)"
			return time.Unix(timestamp, 0).Format(humanReadableFormat)
		},
		"hasTag": func(tag string, tags []string) bool {
			return slices.Contains(tags, tag)
		},
		"dateOnly": func(timeObj time.Time) string {
			return timeObj.Format(time.DateOnly)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("GET /archives2", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("GET /archives2")

		tag := r.FormValue("tag")
		postEntries, err := os.ReadDir("content/posts")
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var pdataList []types.PublishData

		for _, postEntry := range postEntries {
			postName := postEntry.Name()
			slug := strings.TrimSuffix(postName, ".md")
			fmdata, _, err := readers.ReadPage(slug, "posts")
			if err != nil {
				log.Printf("%v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			finfo, err := postEntry.Info()
			if err != nil {
				log.Printf("%v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if fmdata.Publish {
				updatedTime := finfo.ModTime().Unix()

				date, err := time.Parse(time.DateOnly, fmdata.Date)
				if err != nil {
					log.Printf("%v", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				publishData := types.PublishData{
					Title:   fmdata.Title,
					Summary: fmdata.Summary,
					Tags:    fmdata.Tags,
					Slug:    slug,
					Updated: updatedTime,
					Date:    date,
				}

				pdataList = append(pdataList, publishData)
			}
		}

		// Sort the posts by recency (as reported in the
		// frontmatter.) Newer posts should naturally come
		// before older posts.
		sort.Slice(pdataList, func(i int, j int) bool {
			pdata1 := pdataList[i]
			pdata2 := pdataList[j]

			if pdata1.Date.Equal(pdata2.Date) {
				return pdata2.Title < pdata1.Title
			}

			return pdata2.Date.Before(pdata1.Date)
		})

		data := struct {
			Published []types.PublishData
			Tag       string
			Date      string
		}{
			Published: pdataList,
			Tag:       tag,
		}

		if err := feedTemplate(w, "archives", data); err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("GET /tags", func(w http.ResponseWriter, r *http.Request) {
		publishedList, err := readers.ReadPublishingFile("published.json")
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

	mux.HandleFunc("GET /rss", func(w http.ResponseWriter, r *http.Request) {
		siteTitle := "Biome of Ideas"
		var siteURL string

		// Use flagLocal to set the correct siteURL for
		// purposes of testing the RSS feed locally with
		// something like newsboat.
		if flagLocal {
			siteURL = "http://localhost:3030"
		} else {
			siteURL = "https://brandonirizarry.xyz"
		}

		publishedList, err := readers.ReadPublishingFile("published.json")
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var items []types.RSSItem
		for _, p := range publishedList {
			link := fmt.Sprintf("posts/%s", p.Slug)
			pubDate := time.Unix(p.Created, 0).Format(time.RFC1123)

			item := types.RSSItem{
				Title:       p.Title,
				Link:        link,
				GUID:        link,
				PubDate:     pubDate,
				Description: p.Summary,
			}

			items = append(items, item)
		}

		image := types.RSSImage{
			Title:  siteTitle,
			Link:   siteURL,
			URL:    fmt.Sprintf("%s/static/bitmap.png", siteURL),
			Width:  300,
			Height: 300,
		}

		rssChannel := types.RSSChannel{
			Title:       siteTitle,
			Link:        siteURL,
			Description: "My personal website and blog",
			Language:    "en-us",
			Image:       image,
			Items:       items,
		}

		type rss struct {
			Channel types.RSSChannel `xml:"channel"`
			Version string           `xml:"version,attr"`
		}

		rssPayload := rss{
			Channel: rssChannel,
			Version: "2.0",
		}

		// Marshal the data to XML
		feed, err := xml.MarshalIndent(rssPayload, "", strings.Repeat(" ", 4))
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, xml.Header+string(feed))
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
