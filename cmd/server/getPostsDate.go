package main

import (
	"log"
	"net/http"

	"github.com/BrandonIrizarry/buildablog/internal/constants"
	"github.com/BrandonIrizarry/buildablog/internal/readers"
)

func getPostsDate(w http.ResponseWriter, r *http.Request) {
	date := r.PathValue("date")
	postData, err := readers.ReadMarkdown(constants.GenrePublished("posts"), date)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := feedTemplate(w, "post", postData); err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
