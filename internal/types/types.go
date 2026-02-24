package types

import (
	"html/template"
	"time"
)

// FrontmatterData is used for unmarshalling blog post frontmatter.
type FrontmatterData struct {
	// Title is the blog post's title.
	Title string `toml:"title"`

	// Summary is a short summary of the blog post.
	Summary string `toml:"summary"`

	// Tags is the list of tags categorizing the blog post.
	Tags []string `toml:"tags"`

	// Publish flags whether to publish the blog post into the
	// publishing file.
	//
	// DEPRECATED: when generating the /archives page, only list
	// posts with the proper time format (symbolic links.) Hence
	// final inclusion of the date field will ideally signal the
	// publish (idea: a 'make publish' action can generate the
	// appropriate symlinks for posts.)
	Publish bool `toml:"publish"`

	// Date is the date of the post. It looks like TOML will
	// accept time.DateOnly as a value, so let's try it.
	Date time.Time `toml:"date"`
}

type PostData struct {
	FrontmatterData
	Content template.HTML
}
