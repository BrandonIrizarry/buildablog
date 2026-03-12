package main

import (
	"log"
	"net/http"

	"github.com/BrandonIrizarry/buildablog/internal/genres/projects"
)

func (cfg config) getProjects(w http.ResponseWriter, r *http.Request) {
	ps, err := allArticles[projects.Frontmatter](cfg.BlogDir, cfg.IsRepo)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tpl.ExecuteTemplate(w, "projects", ps); err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
