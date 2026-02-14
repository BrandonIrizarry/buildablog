package types

import (
	"html/template"
)

type Reader func(string, string) (any, template.HTML, error)

type Metadata struct {
	Title   string `toml:"title"`
	Summary string `toml:"summary"`
	Publish bool   `toml:"publish"`
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
