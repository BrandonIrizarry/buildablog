package readers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BrandonIrizarry/buildablog/internal/constants"
	"github.com/BrandonIrizarry/buildablog/internal/types"
)

func AllPosts() ([]types.PostData, error) {
	// FIXME: move TZ string to its own constant. Also, replace
	// log.Fatal with a returned error. Also, if possible, we
	// should move this definition outside this function.
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatal(err)
	}

	entries, err := os.ReadDir("content/" + constants.PostsLabel)
	if err != nil {
		return nil, fmt.Errorf("can't read content directory: %w", err)
	}

	// Accumulate the return value into this list.
	var postDataList []types.PostData

	for _, e := range entries {
		filename := e.Name()

		filenameDate, err := time.ParseInLocation(time.DateOnly, filename, location)
		if err != nil {
			// By convention, if the filename doesn't
			// parse according to the given layout, it's
			// considered a draft, so skip it.
			continue
		}

		fmdata, content, err := ReadMarkdown(constants.PostsLabel, filename)
		if err != nil {
			return nil, fmt.Errorf("can't read markdown file %s: %w", filename, err)
		}

		//  These should naturally correspond. In the
		//  future there will be a mechanism to
		//  automatically generate the needed symbolic
		//  links.
		if !fmdata.Date.Equal(filenameDate) {
			return nil, fmt.Errorf("filename %s doesn't match frontmatter date %s", filenameDate, fmdata.Date)
		}

		postData := types.PostData{
			FrontmatterData: fmdata,
			Content:         content,
		}

		postDataList = append(postDataList, postData)
	}

	return postDataList, nil
}
