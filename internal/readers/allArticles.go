package readers

import (
	"fmt"
	"os"
	"slices"

	"github.com/BrandonIrizarry/buildablog/internal/types"
)

func AllArticles[F types.Frontmatter](blogDir string, numPosts *int) ([]types.Article[F], error) {
	genre := (*new(F)).Genre()
	genreDir := fmt.Sprintf("%s/%s", blogDir, genre)

	entries, err := os.ReadDir(genreDir)
	if err != nil {
		return nil, fmt.Errorf("can't read %s: %w", genreDir, err)
	}

	if numPosts == nil {
		numPosts = new(len(entries))
	}

	// Accumulate the return value into this list.
	var articles []types.Article[F]

	// Increment i only when the article is published, that is, it
	// has a date field. Hence we can't use the i that comes with
	// the for-loop here.
	var i int

	for _, e := range entries {
		if i >= *numPosts {
			break
		}

		article, err := ReadArticle[F](blogDir, e.Name())
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
