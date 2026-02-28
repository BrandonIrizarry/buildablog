package main

import (
	"log"
	"net/http"
	"slices"

	"github.com/BrandonIrizarry/buildablog/internal/posts"
	"github.com/BrandonIrizarry/buildablog/internal/readers"
)

func getPosts(w http.ResponseWriter, r *http.Request) {
	ps, err := readers.AllPosts(nil)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter posts by tag.
	tag := r.FormValue("tag")

	if tag != "" {
		ps = slices.DeleteFunc(ps, func(p posts.Post) bool {
			return !slices.Contains(p.Tags, tag)
		})
	}

	if err := feedTemplate(w, "posts", ps); err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
