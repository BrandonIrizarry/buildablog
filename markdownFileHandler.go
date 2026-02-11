package main

import (
	"fmt"
	"log"
	"os"

	"github.com/adrg/frontmatter"
)

// readMarkdownFile returns the blog text found at the given slug
// path.
func readMarkdownFile(slug, label string) (frontmatterData, []byte, error) {
	var data frontmatterData

	if slug == "" {
		slug = "index"
	}

	filename := fmt.Sprintf("%s/%s/%s.md", contentDirName, label, slug)
	log.Printf("Will attempt to open %s", filename)

	f, err := os.Open(filename)
	if err != nil {
		return frontmatterData{}, nil, err
	}
	defer f.Close()

	blogContent, err := frontmatter.Parse(f, &data)
	if err != nil {
		return frontmatterData{}, nil, err
	}

	return data, blogContent, nil
}
