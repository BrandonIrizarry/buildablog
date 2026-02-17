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
	// Slug is the post's given filename minus the file
	// extension. The various templates use this to generate the
	// links to the posts.
	Slug string `json:"slug"`

	// Title is (by convention) the proper title of the post.
	Title string `json:"title"`

	// Summary is (by convention) a blurb summary of the post,
	// used in the archives listing.
	Summary string `json:"summary"`
}
