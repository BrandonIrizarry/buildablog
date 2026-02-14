package readers

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/BrandonIrizarry/buildablog/internal/types"
	"github.com/adrg/frontmatter"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	hl "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/renderer/html"
)

func ReadPage(slug, label string) (types.PostData, error) {
	var data types.PostData

	filename := fmt.Sprintf("content/%s/%s.md", label, slug)
	f, err := os.Open(filename)
	if err != nil {
		return types.PostData{}, err
	}
	defer f.Close()

	// This works because [types.FrontmatterData] is embedded
	// inside [types.PostData].
	blogContent, err := frontmatter.Parse(f, &data)
	if err != nil {
		return types.PostData{}, err
	}

	// Enable syntax highlighting in blog posts.
	//
	// For available styles, see https://github.com/alecthomas/chroma/tree/master/styles
	//
	// See also https://xyproto.github.io/splash/docs/ for
	// a list of canonical themes (though some may not be
	// available here; try 'go get -u' to update chroma
	// and friends.)
	syntaxStyle := "catppuccin-macchiato"

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
		return types.PostData{}, err
	}

	data.Content = template.HTML(buf.String())

	return data, nil
}
