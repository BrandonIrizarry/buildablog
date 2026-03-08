package main

import (
	"log"
	"net/http"

	"github.com/BrandonIrizarry/buildablog/internal/genres/projects"
)

func (cfg config) getProjectsDate(w http.ResponseWriter, r *http.Request) {
	var err error
	article, err := findArticle[projects.Frontmatter](cfg.BlogDir, r.PathValue("date"))
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return

	}

	if err := tpl.ExecuteTemplate(w, "project", article); err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
