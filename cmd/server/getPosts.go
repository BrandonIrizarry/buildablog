package main

import (
	"log"
	"net/http"
	"slices"

	"github.com/BrandonIrizarry/buildablog/internal/genres/posts"
	"github.com/BrandonIrizarry/buildablog/internal/readers"
	"github.com/BrandonIrizarry/buildablog/internal/types"
)

func (cfg config) getPosts(w http.ResponseWriter, r *http.Request) {
	ps, err := readers.AllArticles[posts.Frontmatter](cfg.BlogDir, nil)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter posts by tag.
	tag := r.FormValue("tag")

	if tag != "" {
		ps = slices.DeleteFunc(ps, func(p types.Article[posts.Frontmatter]) bool {
			return !slices.Contains(p.Frontmatter.Tags, tag)
		})
	}

	if err := tpl.ExecuteTemplate(w, "posts", ps); err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
