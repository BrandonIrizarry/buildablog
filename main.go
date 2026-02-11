package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/adrg/frontmatter"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	hl "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	contentDirName = "content"
	postsPrefix    = "posts"
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
	contentPattern := fmt.Sprintf("GET /%s/{slug}", postsPrefix)
	mux.HandleFunc(contentPattern, postHandler(markdownFileHandler))

	// Serve the site's front page.
	mux.HandleFunc("GET /{$}", postHandler(markdownFileHandler))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	if err := http.ListenAndServe(":3030", mux); err != nil {
		log.Fatal(err)
	}
}

// postHandler "decorates" the given reader using an
// [http.HandlerFunc].
func postHandler(reader reader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var post postData
		post.Slug = r.PathValue("slug")
		log.Printf("Slug: %s", post.Slug)

		whole, err := reader(post.Slug)
		if err != nil {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		blogContentBytes, err := frontmatter.Parse(strings.NewReader(whole), &post)
		if err != nil {
			http.Error(w, "Error parsing frontmatter", http.StatusInternalServerError)
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
		if err := mdRenderer.Convert(blogContentBytes, &buf); err != nil {
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

		post.Content = template.HTML(buf.String())

		if err := tpl.Execute(w, post); err != nil {
			log.Printf("error executing template: %v", err)
			http.Error(w, "error executing template", http.StatusInternalServerError)
			return
		}
	}
}
