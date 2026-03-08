package main

import (
	"log"
	"net/http"

	"github.com/BrandonIrizarry/buildablog/internal/genres/index"
	"github.com/BrandonIrizarry/buildablog/internal/genres/posts"
	"github.com/BrandonIrizarry/buildablog/internal/genres/projects"
	"github.com/BrandonIrizarry/buildablog/internal/readers"
	"github.com/BrandonIrizarry/buildablog/internal/types"
)

func (cfg config) getIndex(w http.ResponseWriter, r *http.Request) {
	// FIXME: For now, parse the index.md page as if it were a post.
	frontPage, err := readers.ReadArticle[index.Frontmatter](cfg.BlogDir, "index.md")
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the top three most recent posts.
	//
	// FIXME: make the argument to AllPosts here
	// configurable somehow.
	recentPosts, err := readers.AllArticles[posts.Frontmatter](cfg.BlogDir, new(3))
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the top three most recent projects.
	recentProjects, err := readers.AllArticles[projects.Frontmatter](cfg.BlogDir, new(3))
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
