package readers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/BrandonIrizarry/buildablog/internal/types"
)

func ReadPublishingFile(filename string) ([]types.PublishData, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("can't open %s: %w", filename, err)
	}
	defer f.Close()

	rawJSON, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("can't read %s: %w", filename, err)
	}

	var publishedList []types.PublishData
	if err := json.Unmarshal(rawJSON, &publishedList); err != nil {
		return nil, err
	}

	return publishedList, nil
}
