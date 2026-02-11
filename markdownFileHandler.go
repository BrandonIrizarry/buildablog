package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

// markdownFileHandler returns the blog text found at the given slug
// path.
func markdownFileHandler(slug, label string) (string, error) {
	if slug == "" {
		slug = "index"
	}

	filename := fmt.Sprintf("%s/%s/%s.md", contentDirName, label, slug)
	log.Printf("Will attempt to open %s", filename)

	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
