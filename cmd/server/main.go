package main

import (
	"bufio"
	"bytes"
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
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	hl "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/renderer/html"
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
	mux.HandleFunc(contentPattern, gohtmlHandler(constants.PostsLabel))

	// Serve the site's front page.
	mux.HandleFunc("GET /{$}", gohtmlHandler(constants.IndexLabel))

	// Serve the archives page.
	mux.HandleFunc("GET /archives", func(w http.ResponseWriter, r *http.Request) {
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
		tpl, err := template.ParseFiles("gohtml/archives.gohtml")
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
	})

	// Serve the site's static assets (CSS files etc.)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	if err := http.ListenAndServe(":3030", mux); err != nil {
		log.Fatal(err)
	}
}

// gohtmlHandler "decorates" the given reader using an
// [http.HandlerFunc].
func gohtmlHandler(label string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		log.Printf("Slug: %s", slug)

		fmData, blogContent, err := readers.ReadMarkdownFile(slug, label)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Error loading post", http.StatusNotFound)
			return
		}

		// Enable syntax highlighting in blog posts.
		//
		// For available styles, see https://github.com/alecthomas/chroma/tree/master/styles
		//
		// See also https://xyproto.github.io/splash/docs/ for
		// a list of canonical themes (though some may not be
		// available here; try 'go get -u' to update chroma
		// and friends.)
		syntaxStyle := fmData.Style.Syntax
		if syntaxStyle == "" {
			syntaxStyle = "gruvbox"
		}

		mdRenderer := goldmark.New(
			goldmark.WithExtensions(hl.NewHighlighting(
				hl.WithStyle(syntaxStyle),
				hl.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
					chromahtml.ClassPrefix("content"),
				),
			)),
			// This enables us to use raw HTML in our
			// files, such as anchor-tags (for TOC
			// destinations) and <br> (for adding extra
			// spaces.)
			//
			// I found this out on
			// https://deepwiki.com/yuin/goldmark/2.1-configuration-options
			// 😞
			goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
		)

		// Render Markdown as HTML.
		var buf bytes.Buffer
		if err := mdRenderer.Convert(blogContent, &buf); err != nil {
			http.Error(w, "Error converting Markdown", http.StatusInternalServerError)
			return
		}

		// Load the template.
		tpl, err := template.ParseFiles("gohtml/"+label+".gohtml", "html/nav.html")
		if err != nil {
			log.Printf("error parsing template: %v", err)
			http.Error(w, "error parsing template", http.StatusInternalServerError)
			return
		}

		// Use the template.
		post := types.PostData{
			FrontmatterData: fmData,
			Content:         template.HTML(buf.String()),
		}

		if err := tpl.Execute(w, post); err != nil {
			log.Printf("error executing template: %v", err)
			http.Error(w, "error executing template", http.StatusInternalServerError)
			return
		}
	}
}
