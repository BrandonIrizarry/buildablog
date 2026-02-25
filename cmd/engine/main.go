package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"time"

	"github.com/BrandonIrizarry/buildablog/internal/constants"
	"github.com/BrandonIrizarry/buildablog/internal/readers"
	"github.com/BrandonIrizarry/buildablog/internal/types"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// For all drafts that have a date field, generate a symbolic
	// link in the appropriate format (YYYY-MM-DD) if it doesn't
	// exist already.
	fileData, err := allDrafts()
	if err != nil {
		log.Fatal(err)
	}

	for draftName, postData := range fileData {
		// If the date field is missing, the post is
		// considered unfinished.
		if postData.Date.IsZero() {
			continue
		}

		date := postData.Date.Format(time.DateOnly)
		publishedName := fmt.Sprintf("content/%s/%s", constants.PostsLabel, date)
		if err := os.Symlink("../drafts/"+draftName, publishedName); err != nil {
			if errors.Is(err, fs.ErrExist) {
				log.Printf("already published: %s → %s", draftName, publishedName)
				continue
			} else {
				log.Fatal(err)
			}
		} else {
			log.Printf("creating: %s → %s", draftName, publishedName)
		}
	}
}

func allDrafts() (map[string]types.PostData, error) {
	// FIXME: replace with constant
	entries, err := os.ReadDir("content/drafts")
	if err != nil {
		return nil, fmt.Errorf("can't read %s directory: %w", "content/drafts", err)
	}

	// Here we need to keep the association between filenames and post
	// data, in order to ultimately create the desired symbolic links.
	var fileData = make(map[string]types.PostData)

	for _, e := range entries {
		postData, err := readers.ReadMarkdown("drafts", e.Name())
		if err != nil {
			return nil, fmt.Errorf("can't read markdown file %s: %w", e.Name(), err)
		}

		fileData[e.Name()] = postData
	}

	return fileData, nil
}
