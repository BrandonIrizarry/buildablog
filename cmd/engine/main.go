package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/BrandonIrizarry/buildablog/internal/readers"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Publish the given candidate.
	//
	// This marshals the title and description into a file, which
	// the /archives endpoint reads and uses to define its
	// content.
	var candidate string
	flag.StringVar(&candidate, "publish", "", "Publish content/posts/<post>")
	flag.Parse()

	fmData, _, err := readers.ReadMarkdownFile(candidate, "posts")
	if err != nil {
		log.Fatal("couldn't read content/posts/" + candidate)
	}

	if !fmData.Publish {
		log.Printf("'%s' is a draft, and so will not be published", candidate)
		os.Exit(0)
	}

	// Create this type on the fly, so that we don't have to
	// marshal the entire FrontmatterData struct. The idea is that
	// each entry in the archive will only use these fields.
	type publishData struct {
		Date    string `json:"date"`
		Title   string `json:"title"`
		Summary string `json:"summary"`
	}

	b, err := json.Marshal(publishData{
		Date:    time.Now().Format(time.DateOnly),
		Title:   fmData.Title,
		Summary: fmData.Summary,
	})
	if err != nil {
		log.Fatal("couldn't marshal title and summary")
	}

	f, err := os.OpenFile("published", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("couldn't open 'published' file")
	}
	defer f.Close()

	// Append a newline for human-readability purposes, as well as
	// for the 'GET /archives' endpoint to be able to read this
	// file line by line.
	b = append(b, '\n')

	_, err = f.Write(b)
	if err != nil {
		log.Fatal("couldn't write to 'published' file")
	}
}
