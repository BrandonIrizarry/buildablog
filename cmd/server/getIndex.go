package main

import (
	"log"
	"net/http"

	"github.com/BrandonIrizarry/buildablog/internal/posts"
	"github.com/BrandonIrizarry/buildablog/internal/readers"
	"github.com/BrandonIrizarry/buildablog/internal/types"
)

func (cfg config) getIndex(w http.ResponseWriter, r *http.Request) {
	// For now, parse the index.md page as if it were a post.
	frontPage, err := readers.ReadArticle[posts.Frontmatter](cfg.BlogDir, "index.md")
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the top three most recent posts.
	//
	// FIXME: make the argument to AllPosts here
	// configurable somehow.
	genre := (*new(posts.Frontmatter)).Genre()
	recentPosts, err := readers.AllArticles[posts.Frontmatter](cfg.PublishedDir(genre), new(3), cfg.Timezone)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// The template code will know how to interpret this ad-hoc
	// scheme.
	ps := append([]types.Article[posts.Frontmatter]{frontPage}, recentPosts...)

	if err := tpl.ExecuteTemplate(w, "index", ps); err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
