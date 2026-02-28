package main

import (
	"log"
	"net/http"
)

func getProjects(w http.ResponseWriter, r *http.Request) {
	if err := feedTemplate(w, "projects", struct{}{}); err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
