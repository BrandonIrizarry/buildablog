package readers

import (
	"fmt"
	"log"
	"os"

	"github.com/BrandonIrizarry/buildablog/internal/constants"
	"github.com/adrg/frontmatter"
)

// readMarkdownFile returns the blog text found at the given slug
// path.
func ReadMarkdownFile(slug, label string) (FrontmatterData, []byte, error) {
	var data FrontmatterData

	if slug == "" {
		slug = "index"
	}

	filename := fmt.Sprintf("%s/%s/%s.md", constants.ContentDirName, label, slug)
	log.Printf("Will attempt to open %s", filename)

	f, err := os.Open(filename)
	if err != nil {
		return FrontmatterData{}, nil, err
	}
	defer f.Close()

	blogContent, err := frontmatter.Parse(f, &data)
	if err != nil {
		return FrontmatterData{}, nil, err
	}

	return data, blogContent, nil
}
