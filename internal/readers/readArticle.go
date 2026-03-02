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
)

// ReadArticle returns the article found at pathPrefix/basename.
//
// Currently, only Markdown posts with TOML frontmatter are
// recognized for reading.
func ReadArticle[F types.Frontmatter](pathPrefix, basename string) (types.Article[F], error) {
	var zeroArticle types.Article[F]

	path := fmt.Sprintf("%s/%s", pathPrefix, basename)
	f, err := os.Open(path)
	if err != nil {
		return zeroArticle, err
	}
	defer f.Close()

	var fmdata F
	content, err := frontmatter.Parse(f, &fmdata)
	if err != nil {
		return zeroArticle, err
	}

	htmlContent, err := convertToHTML(content)
	if err != nil {
		return zeroArticle, err
	}

	article := types.Article[F]{
		Frontmatter: fmdata,
		Content:     htmlContent,
	}

	return article, nil
}

// convertToHTML converts the given article's content, a byte slice,
// to [template.HTML], which is then returned with an error.
//
// Syntax highlighting is also added.
//
// For available styles, see https://github.com/alecthomas/chroma/tree/master/styles
//
// See also https://xyproto.github.io/splash/docs/ for a list
// of canonical themes (though some may not be available here;
// try 'go get -u' to update chroma etc. to fetch the latest
// styles.)
func convertToHTML(content []byte) (template.HTML, error) {
	syntaxStyle := "catppuccin-macchiato"

	mdRenderer := goldmark.New(
		goldmark.WithExtensions(hl.NewHighlighting(
			hl.WithStyle(syntaxStyle),
			hl.WithFormatOptions(
				chromahtml.WithLineNumbers(true),
				chromahtml.ClassPrefix("content"),
			),
		)),
	)

	// Render Markdown as HTML.
	var buf bytes.Buffer
	if err := mdRenderer.Convert(content, &buf); err != nil {
		return "", err
	}

	return template.HTML(buf.String()), nil
}
