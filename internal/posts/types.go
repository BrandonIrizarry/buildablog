package posts

import (
	"html/template"
	"time"
)

// Frontmatter is used for unmarshalling blog post frontmatter.
type Frontmatter struct {
	// Title is the blog post's title.
	Title string `toml:"title"`

	// Summary is a short summary of the blog post.
	Summary string `toml:"summary"`

	// Tags is the list of tags categorizing the blog post.
	Tags []string `toml:"tags"`

	// Date is the date of the post. It looks like TOML will
	// accept time.DateOnly as a value, so let's try it.
	Date time.Time `toml:"date"`
}

type Post struct {
	Frontmatter
	Content template.HTML
}
