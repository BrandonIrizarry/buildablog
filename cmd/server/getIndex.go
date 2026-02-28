package main

import (
	"log"
	"net/http"

	"github.com/BrandonIrizarry/buildablog/internal/constants"
	"github.com/BrandonIrizarry/buildablog/internal/posts"
	"github.com/BrandonIrizarry/buildablog/internal/readers"
)

func getIndex(w http.ResponseWriter, r *http.Request) {
	frontPage, err := readers.ReadMarkdown(constants.BlogDir, "index.md")
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the top three most recent posts.
	//
	// FIXME: make the argument to AllPosts here
	// configurable somehow.
	recentPosts, err := readers.AllPosts(new(3))
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// For simplicity, reuse the same [posts.Post] slice
	// datatype, just as we do for the GET /posts
	// endpoint. The template code will know how to
	// interpret this ad-hoc scheme.
	ps := append([]posts.Post{frontPage}, recentPosts...)

	if err := feedTemplate(w, "index", ps); err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
