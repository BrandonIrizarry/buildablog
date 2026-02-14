package types

import (
	"html/template"
)

type Metadata struct {
	Title   string `toml:"title"`
	Summary string `toml:"summary"`
	Publish bool   `toml:"publish"`
}

type PostData struct {
	Metadata
	Content template.HTML
}

// PublishData lets us not have to marshal the entire FrontmatterData
// struct when publishing a blog post. The idea is that each entry in
// the archive will depend only on these fields.
type PublishData struct {
	Date    string `json:"date"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
}
