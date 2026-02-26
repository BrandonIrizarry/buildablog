package readers

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/BrandonIrizarry/buildablog/internal/posts"
	"github.com/adrg/frontmatter"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	hl "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/renderer/html"
)

// ReadMarkdown returns blog post data as two separate parts: frontmatter
// (as a [posts.FrontmatterData] struct) and content (as a
// [template.HTML] string.)
func ReadMarkdown(label, basename string) (posts.PostData, error) {
	var fmdata posts.FrontmatterData

	path := fmt.Sprintf("content/%s/%s", label, basename)
	f, err := os.Open(path)
	if err != nil {
		return posts.PostData{}, err
	}
	defer f.Close()

	content, err := frontmatter.Parse(f, &fmdata)
	if err != nil {
		return posts.PostData{}, err
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
	if err := mdRenderer.Convert(content, &buf); err != nil {
		return posts.PostData{}, err
	}

	postData := posts.PostData{
		FrontmatterData: fmdata,
		Content:         template.HTML(buf.String()),
	}

	return postData, nil
}
