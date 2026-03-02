package readers

import (
	"fmt"

	"github.com/BrandonIrizarry/buildablog/internal/constants"
)

func GenrePublished(genre string) string {
	return fmt.Sprintf("%s/published/%s", constants.BlogDir, genre)
}

func GenreDrafts(genre string) string {
	return fmt.Sprintf("%s/drafts/%s", constants.BlogDir, genre)
}
