package readers

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/BrandonIrizarry/buildablog/internal/types"
)

func ReadPublished() ([]types.PublishData, error) {
	// Read what's currently published. We load each line
	// of 'published' into a slice of [types.PublishData],
	// which is then tossed into the archives.gohtml
	// template.
	f, err := os.Open("published")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var archiveData []types.PublishData

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var data types.PublishData
		entry := scanner.Text()

		if err := json.Unmarshal([]byte(entry), &data); err != nil {
			return nil, err
		}

		archiveData = append(archiveData, data)
	}

	return archiveData, nil
}
