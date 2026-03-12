package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"slices"

	"github.com/BrandonIrizarry/buildablog/v2/internal/types"
	"github.com/adrg/frontmatter"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/yuin/goldmark"
	hl "github.com/yuin/goldmark-highlighting/v2"
)

// allArticles returns all [types.Article] from the blog repo.
func allArticles[F types.Frontmatter](repo string) ([]types.Article[F], error) {
	fs := memfs.New()
	genre := (*new(F)).Genre()

	_, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: repo,
	})
	if err != nil {
		return nil, fmt.Errorf("can't clone repository %s: %w", repo, err)
	}

	log.Printf("Successfully cloned repository %s", repo)

	entries, err := fs.ReadDir("./" + genre)
	if err != nil {
		return nil, err
	}

	log.Printf("Successfully fetched genre entries for '%s'", genre)

	articles, err := entriesToArticles[F](fs, genre, entries)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

// entriesToArticles converts [os.FileInfo] entries into
// [types.Article], returning the slice of these along with an error.
func entriesToArticles[F types.Frontmatter](fs billy.Filesystem, genre string, entries []os.FileInfo) ([]types.Article[F], error) {
	// Accumulate the return value into this list.
	var articles []types.Article[F]

	for _, e := range entries {
		path := fmt.Sprintf("%s/%s", genre, e.Name())
		log.Printf("path: %s", path)
		f, err := fs.Open(path)
		if err != nil {
			log.Printf("can't read article %s: %v", path, err)
			return nil, err
		}
		defer f.Close()

		log.Printf("Successfully opened %s", path)

		article, err := readArticle[F](f)
		if err != nil {
			return nil, fmt.Errorf("can't read markdown file %s: %w", e.Name(), err)
		}

		log.Printf("Successfully fetched article: %v", article)

		// An article is published whenever the date field is
		// filled out.
		if !article.Frontmatter.GetDate().IsZero() {
			log.Printf("Adding %s to published articles", e.Name())
			articles = append(articles, article)
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
