package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"time"

	"github.com/BrandonIrizarry/buildablog/internal/constants"
	"github.com/BrandonIrizarry/buildablog/internal/posts"
	"github.com/BrandonIrizarry/buildablog/internal/readers"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := publish(); err != nil {
		log.Fatal(err)
	}
}

// publish generates a symbolic link for all dated posts in the
// appropriate format (YYYY-MM-DD) if it doesn't exist already. A post
// is considered dated if its frontmatter's 'date' field is non-zero
// (i.e. the user has signed it off with a publishing date.)
func publish() error {
	fileData, err := allDrafts()
	if err != nil {
		return err
	}

	for draftName, postData := range fileData {
		// If the date field is missing, the post is
		// considered unfinished.
		if postData.Date.IsZero() {
			continue
		}

		date := postData.Date.Format(time.DateOnly)
		publishedName := fmt.Sprintf("content/%s/%s", constants.PostsLabel, date)
		if err := os.Symlink("../"+constants.PostDraftsDir+"/"+draftName, publishedName); err != nil {
			if errors.Is(err, fs.ErrExist) {
				log.Printf("already published: %s → %s", draftName, publishedName)
				continue
			} else {
				return err
			}
		} else {
			log.Printf("creating: %s → %s", draftName, publishedName)
		}
	}

	return nil
}

func allDrafts() (map[string]posts.Post, error) {
	entries, err := os.ReadDir("content/" + constants.PostDraftsDir)
	if err != nil {
		return nil, fmt.Errorf("can't read %s directory: %w", "content/"+constants.PostDraftsDir, err)
	}

	// Here we need to keep the association between filenames and post
	// data, in order to ultimately create the desired symbolic links.
	var fileData = make(map[string]posts.Post)

	for _, e := range entries {
		postData, err := readers.ReadMarkdown(constants.PostDraftsDir, e.Name())
		if err != nil {
			return nil, fmt.Errorf("can't read markdown file %s: %w", e.Name(), err)
		}

		fileData[e.Name()] = postData
	}

	return fileData, nil
}
