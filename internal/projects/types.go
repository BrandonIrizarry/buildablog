package projects

import (
	"html/template"
)

// Frontmatter is used for unmarshalling project post frontmatter.
type Frontmatter struct {
	// Name is the name of the project.
	Name string `toml:"title"`

	// HostURL is the project repo's Web URL (i.e. where it's
	// hosted for viewing, contributing, etc.)
	HostURL string `toml:"host_url"`

	// Synopsis is a short blurb summarizing what the project is
	// about.
	Synopsis string `toml:"synopsis"`

	// Stack is the lists of languages, frameworks, etc. used
	// to make the project.
	Stack []string `toml:"stack"`

	// Thumbnail is the path to the thumbnail used for display
	// with the project.
	Thumbnail string `toml:"thumbnail"`
}

type Post struct {
	Frontmatter
	Content template.HTML
}
