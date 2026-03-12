package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"slices"

	"github.com/BrandonIrizarry/buildablog/internal/types"
	"github.com/adrg/frontmatter"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/yuin/goldmark"
	hl "github.com/yuin/goldmark-highlighting/v2"
)

// allArticles returns all [types.Article] from the blog directory,
// which is read directly from the local filesystem.
func allArticles[F types.Frontmatter](blogDir string, isRepo bool) ([]types.Article[F], error) {
	var err error
	var entries []os.FileInfo

	genre := (*new(F)).Genre()
	genreDir := fmt.Sprintf("%s/%s", blogDir, genre)

	if isRepo {
		entries, err = getEntriesRepo(blogDir, genre)
	} else {
		entries, err = getEntriesDir(blogDir, genre)
	}

	if err != nil {
		return nil, err
	}

	articles, err := entriesToArticles[F](genreDir, entries)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func getEntriesDir(blogDir, genre string) ([]os.FileInfo, error) {
	genreDir := fmt.Sprintf("%s/%s", blogDir, genre)

	dir, err := os.Open(genreDir)
	if err != nil {
		return nil, fmt.Errorf("can't open directory %s: %w", genreDir, err)
	}
	defer dir.Close()

	entries, err := dir.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("can't read directory %s: %w", genreDir, err)
	}

	for _, e := range entries {
		log.Printf("Entry: %v", e.Name())
	}

	return entries, nil
}

func getEntriesRepo(blogDir, genre string) ([]os.FileInfo, error) {
	fs := memfs.New()

	// FIXME: for now, blogDir can refer to both an
	// ordinary directory, or else either a local or
	// remote Git repo.
	_, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: blogDir,
	})

	if err != nil {
		return nil, fmt.Errorf("can't clone repository %s: %w", blogDir, err)
	}

	log.Print("Successfully cloned repository")

	// Note that, here, the Git repo doesn't know anything about
	// my local home directory, even in the case where I simply
	// cloned from there. So the following ends up being the
	// correct way to read something from within the repo.
	entries, err := fs.ReadDir("./" + genre)
	if err != nil {
		return nil, fmt.Errorf("can't read repository: %w", err)
	}

	for _, e := range entries {
		log.Printf("Entry: %v", e.Name())
	}

	return entries, nil
}

// entriesToArticles converts [os.FileInfo] entries into
// [types.Article], returning the slice of these along with an error.
func entriesToArticles[F types.Frontmatter](genrePath string, entries []os.FileInfo) ([]types.Article[F], error) {
	// Accumulate the return value into this list.
	var articles []types.Article[F]

	for _, e := range entries {
		readingPath := fmt.Sprintf("%s/%s", genrePath, e.Name())
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
