package readers

import (
	"fmt"
	"os"
	"time"

	"github.com/BrandonIrizarry/buildablog/internal/constants"
	"github.com/BrandonIrizarry/buildablog/internal/posts"
)

func AllPosts() ([]posts.Post, error) {
	publishedDir := constants.GenrePublished("posts")

	entries, err := os.ReadDir(publishedDir)
	if err != nil {
		return nil, fmt.Errorf("can't read %s: %w", publishedDir, err)
	}

	// Accumulate the return value into this list.
	var postDataList []posts.Post

	for _, e := range entries {
		filename := e.Name()

		filenameDate, err := time.ParseInLocation(time.DateOnly, filename, constants.TZOffset)
		if err != nil {
			return nil, fmt.Errorf("%s isn't in YYYY-MM-DD format: %w", filename, err)
		}

		postData, err := ReadMarkdown(constants.GenrePublished("posts"), filename)
		if err != nil {
			return nil, fmt.Errorf("can't read markdown file %s: %w", filename, err)
		}

		//  These should naturally correspond. In the
		//  future there will be a mechanism to
		//  automatically generate the needed symbolic
		//  links.
		if !postData.Date.Equal(filenameDate) {
			return nil, fmt.Errorf("filename %s doesn't match frontmatter date %s", filenameDate, postData.Date)
		}

		postDataList = append(postDataList, postData)
	}

	return postDataList, nil
}
