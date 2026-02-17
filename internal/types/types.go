package types

import "html/template"

type FrontmatterData struct {
	Title   string `toml:"title"`
	Summary string `toml:"summary"`
	Publish bool   `toml:"publish"`
}

type PostData struct {
	FrontmatterData
	Content template.HTML
}

// PublishData represents a blog post in the archive file
// (published.json).
type PublishData struct {
	Date    string `json:"date"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
}
