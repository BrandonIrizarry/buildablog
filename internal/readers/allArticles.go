package readers

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"slices"

	"github.com/BrandonIrizarry/buildablog/internal/types"
	"github.com/adrg/frontmatter"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	hl "github.com/yuin/goldmark-highlighting/v2"
)

// AllArticles returns all [types.Article] from the blog directory,
// which is read directly from the local filesystem.
func AllArticles[F types.Frontmatter](blogDir string) ([]types.Article[F], error) {
	genre := (*new(F)).Genre()
	genreDir := fmt.Sprintf("%s/%s", blogDir, genre)

	dir, err := os.Open(genreDir)
	if err != nil {
		return nil, fmt.Errorf("can't read %s: %w", genreDir, err)
	}
	defer dir.Close()

	entries, err := dir.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("can't read %s: %w", genreDir, err)
	}

	// Accumulate the return value into this list.
	var articles []types.Article[F]

	// Increment i only when the article is published, that is, it
	// has a date field. Hence we can't use the i that comes with
	// the for-loop here.
	var i int

	for _, e := range entries {
		readingPath := fmt.Sprintf("%s/%s/%s", blogDir, genre, e.Name())

		f, err := os.Open(readingPath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		article, err := readArticle[F](f)
		if err != nil {
			return nil, fmt.Errorf("can't read markdown file %s: %w", e.Name(), err)
		}

		// An article is published whenever the date field is
		// filled out.
		if !article.Frontmatter.GetDate().IsZero() {
			articles = append(articles, article)
			i++
		}
	}

	// Sort the articles by date, in reverse order (most recent
	// posts first.)
	slices.SortFunc(articles, func(a1, a2 types.Article[F]) int {
		date1 := a1.Frontmatter.GetDate()
		date2 := a2.Frontmatter.GetDate()

		return date2.Compare(date1)
	})

	return articles, nil
}

// readArticle reads the Markdown file 'basename' under blogDir/genre,
// where genre belongs to the [types.Frontmatter] set by the caller.
func readArticle[F types.Frontmatter](r io.Reader) (types.Article[F], error) {
	var zero types.Article[F]
	var fmdata F

	content, err := frontmatter.Parse(r, &fmdata)
	if err != nil {
		return zero, err
	}

	htmlContent, err := convertToHTML(content)
	if err != nil {
		return zero, err
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
