package main

import (
	"log"
	"net/http"

	"github.com/BrandonIrizarry/buildablog/internal/constants"
	"github.com/BrandonIrizarry/buildablog/internal/posts"
	"github.com/BrandonIrizarry/buildablog/internal/readers"
	"github.com/BrandonIrizarry/buildablog/internal/types"
)

func getIndex(w http.ResponseWriter, r *http.Request) {
	// For now, parse the index.md page as if it were a post.
	frontPage, err := readers.ReadArticle[posts.Frontmatter](constants.BlogDir, "index.md")
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the top three most recent posts.
	//
	// FIXME: make the argument to AllPosts here
	// configurable somehow.
	recentPosts, err := readers.AllArticles[posts.Frontmatter](new(3))
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// The template code will know how to interpret this ad-hoc
	// scheme.
	ps := append([]types.Article[posts.Frontmatter]{frontPage}, recentPosts...)

	if err := feedTemplate(w, "index", ps); err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
