package types

import (
	"html/template"
)

type FrontmatterData struct {
	Title   string `toml:"title"`
	Summary string `toml:"summary"`
	Publish bool   `toml:"publish"`
	Author  author `toml:"author"`
	Style   style  `toml:"style"`
}

type author struct {
	Name  string `toml:"name"`
	Email string `toml:"email"`
}

type style struct {
	Syntax string `toml:"syntax"`
}

type PostData struct {
	FrontmatterData
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
