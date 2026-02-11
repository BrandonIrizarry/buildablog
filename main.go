package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	hl "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	contentDirName = "content"
	postsLabel     = "posts"
	indexLabel     = "index"
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
	contentPattern := fmt.Sprintf("GET /%s/{slug}", postsLabel)
	mux.HandleFunc(contentPattern, postHandler(markdownFileHandler, postsLabel))

	// Serve the site's front page.
	mux.HandleFunc("GET /{$}", postHandler(markdownFileHandler, indexLabel))

	// Serve the site's static assets (CSS files etc.)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	if err := http.ListenAndServe(":3030", mux); err != nil {
		log.Fatal(err)
	}
}

// postHandler "decorates" the given reader using an
// [http.HandlerFunc].
func postHandler(reader reader, label string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		log.Printf("Slug: %s", slug)

		fmData, blogContent, err := reader(slug, label)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Error loading post", http.StatusNotFound)
			return
		}

		// Enable syntax highlighting in blog posts.
		//
		// For available styles, see https://xyproto.github.io/splash/docs/
		mdRenderer := goldmark.New(
			goldmark.WithExtensions(hl.NewHighlighting(
				hl.WithStyle("gruvbox"),
				hl.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
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
		tpl, err := template.ParseFiles(label + ".gohtml")
		if err != nil {
			log.Printf("error parsing template: %v", err)
			http.Error(w, "error parsing template", http.StatusInternalServerError)
			return
		}

		// Use the template.
		post := postData{
			frontmatterData: fmData,
			Content:         template.HTML(buf.String()),
		}

		if err := tpl.Execute(w, post); err != nil {
			log.Printf("error executing template: %v", err)
			http.Error(w, "error executing template", http.StatusInternalServerError)
			return
		}
	}
}
