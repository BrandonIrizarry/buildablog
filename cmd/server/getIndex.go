package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BrandonIrizarry/buildablog/v2/internal/genres/index"
	"github.com/BrandonIrizarry/buildablog/v2/internal/genres/posts"
	"github.com/BrandonIrizarry/buildablog/v2/internal/genres/projects"
	"github.com/BrandonIrizarry/buildablog/v2/internal/types"
)

func (cfg config) getIndex(w http.ResponseWriter, r *http.Request) {
	frontPages, err := allArticles[index.Frontmatter](cfg.BlogDir)

	// We expect there for now, by suitable convention, to be only
	// one front page; but let's still guard for any degenerate
	// cases.
	if len(frontPages) == 0 {
		err := fmt.Errorf("no front page")
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	frontPage := frontPages[0]

	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the top three most recent posts.
	recentPosts, err := allArticles[posts.Frontmatter](cfg.BlogDir)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the top three most recent projects.
	recentProjects, err := allArticles[projects.Frontmatter](cfg.BlogDir)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type payload struct {
		Intro    types.Article[index.Frontmatter]
		Posts    []types.Article[posts.Frontmatter]
		Projects []types.Article[projects.Frontmatter]
	}

	pld := payload{
		Intro:    frontPage,
		Posts:    recentPosts,
		Projects: recentProjects,
	}

	if err := tpl.ExecuteTemplate(w, "index", pld); err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
