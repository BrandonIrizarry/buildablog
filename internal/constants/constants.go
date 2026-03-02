package constants

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	// TZOffset is used for localizing time stamps we create.
	TZOffset *time.Location

	// BlogDir is the filesystem location of the site's content.
	BlogDir string
)

func init() {
	var err error

	// Configure TZOffset.
	TZOffset, err = time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatal(err)
	}

	// Configure BlogDir.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	BlogDir = fmt.Sprintf("%s/blog", homeDir)
}
