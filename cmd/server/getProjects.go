package main

import (
	"log"
	"net/http"

	"github.com/BrandonIrizarry/buildablog/internal/projects"
	"github.com/BrandonIrizarry/buildablog/internal/readers"
)

func (cfg config) getProjects(w http.ResponseWriter, r *http.Request) {
	ps, err := readers.AllArticles[projects.Frontmatter](nil)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := feedTemplate(w, "projects", ps); err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
