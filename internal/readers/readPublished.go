package readers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"

	"github.com/BrandonIrizarry/buildablog/internal/constants"
	"github.com/BrandonIrizarry/buildablog/internal/types"
)

func ReadPublished(slug, label string) (any, template.HTML, error) {
	// Read anything the root file contains. Note that this file
	// should contain no frontmatter.
	if slug == "" {
		slug = "index"
	}

	filename := fmt.Sprintf("%s/%s/%s.md", constants.ContentDirName, label, slug)
	log.Printf("Will attempt to open %s", filename)

	f, err := os.Open(filename)
	if err != nil {
		return types.Metadata{}, "", err
	}
	defer f.Close()

	// There should be no frontmatter in this context.
	content, err := io.ReadAll(f)
	if err != nil {
		return types.Metadata{}, "", err
	}

	// Read what's currently published. We load each line
	// of 'published' into a slice of [types.PublishData],
	// which is then tossed into the archives.gohtml
	// template.
	f, err = os.Open("published")
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	var archiveData []types.PublishData

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var data types.PublishData
		entry := scanner.Text()

		if err := json.Unmarshal([]byte(entry), &data); err != nil {
			return nil, "", err
		}

		archiveData = append(archiveData, data)
	}

	return archiveData, template.HTML(content), nil
}
