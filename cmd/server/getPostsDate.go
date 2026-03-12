package main

import (
	"log"
	"net/http"

	"github.com/BrandonIrizarry/buildablog/v2/internal/genres/posts"
)

func (cfg config) getPostsDate(w http.ResponseWriter, r *http.Request) {
	var err error
	article, err := findArticle[posts.Frontmatter](cfg.BlogDir, r.PathValue("date"))
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return

	}

	if err := tpl.ExecuteTemplate(w, "post", article); err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
