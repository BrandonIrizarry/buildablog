package readers

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/BrandonIrizarry/buildablog/internal/constants"
	"github.com/BrandonIrizarry/buildablog/internal/types"
)

func AllArticles[F types.Frontmatter](numPosts *int) ([]types.Article[F], error) {
	publishedDir := fmt.Sprintf("%s/published/%s", constants.BlogDir, (*new(F)).Genre())

	entries, err := os.ReadDir(publishedDir)
	if err != nil {
		return nil, fmt.Errorf("can't read %s: %w", publishedDir, err)
	}

	// The nice thing is that, thanks to [time.DateOnly], posts
	// are already sorted on the filesystem in order from oldest
	// to newest. However, for display in Posts, RSS, etc., posts
	// should appear from newest to oldest.
	slices.Reverse(entries)

	if numPosts == nil {
		numPosts = new(len(entries))
	}

	// Accumulate the return value into this list.
	var articles []types.Article[F]

	for i, e := range entries {
		if i >= *numPosts {
			break
		}

		filename := e.Name()

		filenameDate, err := time.ParseInLocation(time.DateOnly, filename, constants.TZOffset)
		if err != nil {
			return nil, fmt.Errorf("%s isn't in YYYY-MM-DD format: %w", filename, err)
		}

		article, err := ReadArticle[F](publishedDir, filename)
		if err != nil {
			return nil, fmt.Errorf("can't read markdown file %s: %w", filename, err)
		}

		//  The frontmatter date and the symlink name should
		//  correspond by definition.
		date := article.Frontmatter.GetDate()
		if !date.Equal(filenameDate) {
			return nil, fmt.Errorf("filename %s doesn't match frontmatter date %s", filenameDate, date)
		}

		articles = append(articles, article)
	}

	return articles, nil
}
