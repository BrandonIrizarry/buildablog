package main

import (
	"log"
	"net/http"

	"github.com/BrandonIrizarry/buildablog/internal/posts"
	"github.com/BrandonIrizarry/buildablog/internal/readers"
)

func (cfg config) getTags(w http.ResponseWriter, r *http.Request) {
	genre := (*new(posts.Frontmatter)).Genre()

	posts, err := readers.AllArticles[posts.Frontmatter](cfg.PublishedDir(genre), nil)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Dump all tags into a list, then deduplicate using a
	// set.
	var tagList []string
	for _, p := range posts {
		tagList = append(tagList, p.Frontmatter.Tags...)
	}

	tagSet := make(map[string]struct{})
	for _, t := range tagList {
		tagSet[t] = struct{}{}
	}

	log.Printf("Tag set: %v", tagSet)

	if err := feedTemplate(w, "tags", tagSet); err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
