package constants

import (
	"log"
	"time"
)

const (
	PostsLabel = "posts"
	IndexLabel = "index"
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
