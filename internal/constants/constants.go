package constants

import (
	"log"
	"time"
)

// Label definitions.
//
// Labels help associate various assets into common classes (e.g.,
// blog posts use a certain set of templates, directories, endpoints
// etc., as do projects, and so on.q)
const (
	PostsLabel = "posts"
	IndexLabel = "index"
)

const (
	PostDraftsDir = "drafts/posts"
)

// TZOffset is used for localizing time stamps we create.
var TZOffset *time.Location

func init() {
	var err error

	TZOffset, err = time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatal(err)
	}
}
