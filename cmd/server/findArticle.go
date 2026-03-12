package main

import (
	"fmt"
	"time"

	"github.com/BrandonIrizarry/buildablog/internal/types"
)

func findArticle[F types.Frontmatter](blogDir, date string, isRepo bool) (types.Article[F], error) {
	var zero types.Article[F]
	var err error

	articles, err := allArticles[F](blogDir, isRepo)
	if err != nil {
		return zero, err
	}

	for _, a := range articles {
		articleDate := a.Frontmatter.GetDate().Format(time.DateOnly)
		if articleDate == date {
			return a, nil
		}
	}

	return zero, fmt.Errorf("article under %s with date %s not found", blogDir, date)
}
