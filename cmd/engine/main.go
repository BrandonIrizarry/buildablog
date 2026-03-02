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
	"github.com/BrandonIrizarry/buildablog/internal/types"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := publish[posts.Frontmatter](); err != nil {
		log.Fatal(err)
	}
}

// publish generates a symbolic link for all dated posts in the
// appropriate format (YYYY-MM-DD) if it doesn't exist already. A post
// is considered dated if its frontmatter's 'date' field is non-zero
// (i.e. the user has signed it off with a publishing date.)
func publish[F types.Frontmatter]() error {
	fileData, err := allDrafts[F]()
	if err != nil {
		return err
	}

	for originalFilename, fmdata := range fileData {
		// If the date field is missing, the post is
		// considered unfinished.
		if fmdata.GetDate().IsZero() {
			continue
		}

		date := fmdata.GetDate().Format(time.DateOnly)
		symlinkTarget := fmt.Sprintf("%s/drafts/%s/%s", constants.BlogDir, fmdata.Genre(), originalFilename)
		publishedName := fmt.Sprintf("%s/published/%s/%s", constants.BlogDir, fmdata.Genre(), date)

		if err := os.Symlink(symlinkTarget, publishedName); err != nil {
			if errors.Is(err, fs.ErrExist) {
				log.Printf("already published: %s → %s", originalFilename, publishedName)
				continue
			} else {
				return err
			}
		} else {
			log.Printf("creating: %s → %s", originalFilename, publishedName)
		}
	}

	return nil
}

func allDrafts[F types.Frontmatter]() (map[string]F, error) {
	draftsDir := fmt.Sprintf("%s/drafts/%s", constants.BlogDir, (*new(F)).Genre())
	entries, err := os.ReadDir(draftsDir)
	if err != nil {
		return nil, fmt.Errorf("can't read %s directory: %w", draftsDir, err)
	}

	// Here we need to keep the association between filenames and post
	// data, in order to ultimately create the desired symbolic links.
	var fileData = make(map[string]F)

	for _, e := range entries {
		fmdata, _, err := readers.ReadFrontmatter[F](draftsDir, e.Name())
		if err != nil {
			return nil, fmt.Errorf("can't read markdown file %s: %w", e.Name(), err)
		}

		fileData[e.Name()] = fmdata
	}

	return fileData, nil
}
