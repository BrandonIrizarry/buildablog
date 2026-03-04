package main

import (
	"log"
	"net/http"

	"github.com/BrandonIrizarry/buildablog/internal/projects"
	"github.com/BrandonIrizarry/buildablog/internal/readers"
)

func (cfg config) getProjectsDate(w http.ResponseWriter, r *http.Request) {
	date := r.PathValue("date")
	postData, err := readers.ReadArticle[projects.Frontmatter](cfg.PublishedDir("projects"), date)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tpl.ExecuteTemplate(w, "project", postData); err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
